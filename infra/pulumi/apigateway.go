package main

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/servicediscovery"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// APIGatewayResources holds the API Gateway resources created by createAPIGateway.
type APIGatewayResources struct {
	API         *apigatewayv2.Api
	Stage       *apigatewayv2.Stage
	InvokeRoles map[string]*iam.Role
}

// createAPIGateway creates an HTTP API with IAM auth, VPC Link to Cloud Map,
// and a cross-account invoke role per caller account.
func createAPIGateway(ctx *pulumi.Context, appEnv string, callerAccounts map[string]string, net *NetworkResources, cmService *servicediscovery.Service) (*APIGatewayResources, error) {
	api, err := apigatewayv2.NewApi(ctx, "api", &apigatewayv2.ApiArgs{
		Name:         pulumi.Sprintf("subject-data-%s", appEnv),
		ProtocolType: pulumi.String("HTTP"),
		Tags:         tags("api", appEnv),
	})
	if err != nil {
		return nil, err
	}

	vpcLink, err := apigatewayv2.NewVpcLink(ctx, "vpc-link", &apigatewayv2.VpcLinkArgs{
		Name:             pulumi.Sprintf("subject-data-%s", appEnv),
		SubnetIds:        net.SubnetIDs,
		SecurityGroupIds: pulumi.StringArray{net.ECSSG.ID()},
		Tags:             tags("vpc-link", appEnv),
	})
	if err != nil {
		return nil, err
	}

	integration, err := apigatewayv2.NewIntegration(ctx, "integration", &apigatewayv2.IntegrationArgs{
		ApiId:                api.ID(),
		IntegrationType:      pulumi.String("HTTP_PROXY"),
		IntegrationMethod:    pulumi.String("ANY"),
		ConnectionType:       pulumi.String("VPC_LINK"),
		ConnectionId:         vpcLink.ID(),
		IntegrationUri:       cmService.Arn,
		PayloadFormatVersion: pulumi.String("1.0"),
	})
	if err != nil {
		return nil, err
	}

	_, err = apigatewayv2.NewRoute(ctx, "route", &apigatewayv2.RouteArgs{
		ApiId:             api.ID(),
		RouteKey:          pulumi.String("ANY /{proxy+}"),
		Target:            pulumi.Sprintf("integrations/%s", integration.ID()),
		AuthorizationType: pulumi.String("AWS_IAM"),
	})
	if err != nil {
		return nil, err
	}

	stage, err := apigatewayv2.NewStage(ctx, "stage", &apigatewayv2.StageArgs{
		ApiId:      api.ID(),
		Name:       pulumi.String("$default"),
		AutoDeploy: pulumi.Bool(true),
		Tags:       tags("stage", appEnv),
	})
	if err != nil {
		return nil, err
	}

	// --- Cross-account invoke roles ---
	//
	// Each caller account gets its own role + policy. The caller's ECS task
	// role assumes this role to call the API Gateway. The trust policy scopes
	// to that account; the permission grants execute-api:Invoke.

	// Sort caller names for deterministic resource ordering.
	callerNames := make([]string, 0, len(callerAccounts))
	for name := range callerAccounts {
		callerNames = append(callerNames, name)
	}
	sort.Strings(callerNames)

	invokeRoles := make(map[string]*iam.Role, len(callerAccounts))
	for _, name := range callerNames {
		accountID := callerAccounts[name]

		role, err := iam.NewRole(ctx, fmt.Sprintf("apigw-invoke-role-%s", name), &iam.RoleArgs{
			Description: pulumi.Sprintf("Allows %s (%s) to invoke subject-data API Gateway", name, accountID),
			AssumeRolePolicy: pulumi.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Principal": {"AWS": "arn:aws:iam::%s:root"},
				"Action": "sts:AssumeRole",
				"Condition": {
					"StringEquals": {
						"sts:ExternalId": "subject-data-%s"
					}
				}
			}]
		}`, accountID, appEnv),
			Tags: tags(fmt.Sprintf("apigw-invoke-role-%s", name), appEnv),
		})
		if err != nil {
			return nil, err
		}

		_, err = iam.NewRolePolicy(ctx, fmt.Sprintf("apigw-invoke-policy-%s", name), &iam.RolePolicyArgs{
			Role: role.Name,
			Policy: api.ExecutionArn.ApplyT(func(execArn string) string {
				policy, _ := json.Marshal(map[string]any{
					"Version": "2012-10-17",
					"Statement": []map[string]any{
						{
							"Effect":   "Allow",
							"Action":   "execute-api:Invoke",
							"Resource": fmt.Sprintf("%s/*", execArn),
						},
					},
				})
				return string(policy)
			}).(pulumi.StringOutput),
		})
		if err != nil {
			return nil, err
		}

		invokeRoles[name] = role
	}

	return &APIGatewayResources{
		API:         api,
		Stage:       stage,
		InvokeRoles: invokeRoles,
	}, nil
}
