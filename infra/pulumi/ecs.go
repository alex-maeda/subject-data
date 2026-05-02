package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/cloudwatch"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ecs"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/secretsmanager"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/servicediscovery"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ECSResources holds the ECS resources created by createECS.
type ECSResources struct {
	Cluster         *ecs.Cluster
	CloudMapService *servicediscovery.Service
}

// createECS creates the ECS cluster, task definition, service, and supporting
// resources (log group, IAM roles, Cloud Map service discovery).
//
// The task definition includes two containers:
//  1. subject-data — the app, port 8380
//  2. tailscale        — sidecar for backend connectivity via Tailscale
//
// Fargate tasks share a network namespace, so subject-data can reach
// backend services via the Tailscale interface without any proxy config.
func createECS(ctx *pulumi.Context, appEnv string, net *NetworkResources, repoURL string, rdsRes *RDSResources, s3Res *S3BucketResources) (*ECSResources, error) {
	// --- CloudWatch log group ---

	_, err := cloudwatch.NewLogGroup(ctx, "log-group", &cloudwatch.LogGroupArgs{
		Name:            pulumi.Sprintf("/ecs/subject-data-%s", appEnv),
		RetentionInDays: pulumi.Int(30),
		Tags:            tags("logs", appEnv),
	})
	if err != nil {
		return nil, err
	}

	// --- IAM: execution role (pull ECR, write logs) ---

	execRole, err := iam.NewRole(ctx, "ecs-exec-role", &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(`{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Principal": {"Service": "ecs-tasks.amazonaws.com"},
				"Action": "sts:AssumeRole"
			}]
		}`),
		Tags: tags("ecs-exec-role", appEnv),
	})
	if err != nil {
		return nil, err
	}

	_, err = iam.NewRolePolicyAttachment(ctx, "exec-ecr", &iam.RolePolicyAttachmentArgs{
		Role:      execRole.Name,
		PolicyArn: pulumi.String("arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"),
	})
	if err != nil {
		return nil, err
	}

	// --- Secrets Manager ---
	//
	// Pulumi creates the secret resources. Values are set out-of-band via CLI:
	//   aws secretsmanager put-secret-value --secret-id <name> --secret-string <value>

	tsAuthKey, err := secretsmanager.NewSecret(ctx, "ts-authkey", &secretsmanager.SecretArgs{
		Name: pulumi.Sprintf("subject-data/%s/ts-authkey", appEnv),
		Tags: tags("ts-authkey", appEnv),
	})
	if err != nil {
		return nil, err
	}

	authTokens, err := secretsmanager.NewSecret(ctx, "auth-tokens", &secretsmanager.SecretArgs{
		Name: pulumi.Sprintf("subject-data/%s/api-auth-tokens", appEnv),
		Tags: tags("auth-tokens", appEnv),
	})
	if err != nil {
		return nil, err
	}

	// Store the DATABASE_URL as a secret so ECS can inject it at runtime.
	dbURLSecret, err := secretsmanager.NewSecret(ctx, "db-url", &secretsmanager.SecretArgs{
		Name: pulumi.Sprintf("subject-data/%s/database-url", appEnv),
		Tags: tags("db-url", appEnv),
	})
	if err != nil {
		return nil, err
	}

	_, err = secretsmanager.NewSecretVersion(ctx, "db-url-version", &secretsmanager.SecretVersionArgs{
		SecretId:     dbURLSecret.ID(),
		SecretString: rdsRes.ConnStringFn,
	})
	if err != nil {
		return nil, err
	}

	// Grant the execution role access to pull secrets at task start.
	_, err = iam.NewRolePolicy(ctx, "exec-secrets", &iam.RolePolicyArgs{
		Role: execRole.Name,
		Policy: pulumi.All(tsAuthKey.Arn, authTokens.Arn, dbURLSecret.Arn).ApplyT(
			func(args []any) string {
				policy, _ := json.Marshal(map[string]any{
					"Version": "2012-10-17",
					"Statement": []map[string]any{
						{
							"Effect": "Allow",
							"Action": "secretsmanager:GetSecretValue",
							"Resource": []string{
								args[0].(string),
								args[1].(string),
								args[2].(string),
							},
						},
					},
				})
				return string(policy)
			},
		).(pulumi.StringOutput),
	})
	if err != nil {
		return nil, err
	}

	// --- IAM: task role (what the running container can do) ---

	taskRole, err := iam.NewRole(ctx, "ecs-task-role", &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(`{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Principal": {"Service": "ecs-tasks.amazonaws.com"},
				"Action": "sts:AssumeRole"
			}]
		}`),
		Tags: tags("ecs-task-role", appEnv),
	})
	if err != nil {
		return nil, err
	}

	// Grant the task role S3 read/write for the ingestion bucket.
	// Read: pull parquet files and manifests during batch import jobs.
	// Write: publish ingestion.json completion marker back to the same prefix.
	_, err = iam.NewRolePolicy(ctx, "task-s3-access", &iam.RolePolicyArgs{
		Role: taskRole.Name,
		Policy: s3Res.Bucket.Bucket.ApplyT(
			func(bucket string) string {
				bucketArn := fmt.Sprintf("arn:aws:s3:::%s", bucket)
				policy, _ := json.Marshal(map[string]any{
					"Version": "2012-10-17",
					"Statement": []map[string]any{
						{
							"Sid":    "S3ReadIngestionBucket",
							"Effect": "Allow",
							"Action": []string{
								"s3:GetObject",
								"s3:ListBucket",
							},
							"Resource": []string{
								bucketArn,
								bucketArn + "/*",
							},
						},
						{
							"Sid":    "S3WriteIngestionJson",
							"Effect": "Allow",
							"Action": []string{
								"s3:PutObject",
							},
							"Resource": bucketArn + "/published/*/ingestion.json",
						},
					},
				})
				return string(policy)
			},
		).(pulumi.StringOutput),
	})
	if err != nil {
		return nil, err
	}

	// --- Cluster ---

	cluster, err := ecs.NewCluster(ctx, "cluster", &ecs.ClusterArgs{
		Name: pulumi.Sprintf("subject-data-%s", appEnv),
		Settings: ecs.ClusterSettingArray{
			&ecs.ClusterSettingArgs{
				Name:  pulumi.String("containerInsights"),
				Value: pulumi.String("enabled"),
			},
		},
		Tags: tags("cluster", appEnv),
	})
	if err != nil {
		return nil, err
	}

	// --- Task definition ---

	imageTag := os.Getenv("IMAGE_TAG")
	if imageTag == "" {
		imageTag = "latest"
	}

	containerDefs := pulumi.All(
		tsAuthKey.Arn, authTokens.Arn, dbURLSecret.Arn,
	).ApplyT(func(args []any) (string, error) {
		tsAuthKeyArn := args[0].(string)
		authTokensArn := args[1].(string)
		dbURLArn := args[2].(string)

		defs := []map[string]any{
			{
				"name":      "subject-data",
				"image":     fmt.Sprintf("%s:%s", repoURL, imageTag),
				"essential": true,
				"portMappings": []map[string]any{
					{"containerPort": 8380, "protocol": "tcp"},
				},
				"environment": []map[string]string{
					{"name": "PORT", "value": "8380"},
					{"name": "APP_ENV", "value": appEnv},
					{"name": "DB_DRIVER", "value": "postgres"},
					{"name": "TS_SOCKS5_PROXY", "value": "localhost:1055"},
				},
				"secrets": []map[string]string{
					{"name": "AUTH_TOKENS", "valueFrom": authTokensArn},
					{"name": "DATABASE_URL", "valueFrom": dbURLArn},
				},
				"logConfiguration": map[string]any{
					"logDriver": "awslogs",
					"options": map[string]string{
						"awslogs-group":         fmt.Sprintf("/ecs/subject-data-%s", appEnv),
						"awslogs-region":        "us-east-2",
						"awslogs-stream-prefix": "app",
					},
				},
				"dependsOn": []map[string]string{
					{"containerName": "tailscale", "condition": "HEALTHY"},
				},
			},
			{
				"name":      "tailscale",
				"image":     "ghcr.io/tailscale/tailscale:latest",
				"essential": true,
				"environment": []map[string]string{
					{"name": "TS_USERSPACE", "value": "true"},
					{"name": "TS_SOCKS5_SERVER", "value": ":1055"},
					{"name": "TS_HOSTNAME", "value": fmt.Sprintf("subject-data-%s", appEnv)},
					{"name": "TS_EXTRA_ARGS", "value": "--advertise-tags=tag:server"},
				},
				"secrets": []map[string]string{
					{"name": "TS_AUTHKEY", "valueFrom": tsAuthKeyArn},
				},
				"healthCheck": map[string]any{
					"command":     []string{"CMD-SHELL", "tailscale status --peers=false"},
					"interval":    10,
					"timeout":     5,
					"retries":     3,
					"startPeriod": 30,
				},
				"logConfiguration": map[string]any{
					"logDriver": "awslogs",
					"options": map[string]string{
						"awslogs-group":         fmt.Sprintf("/ecs/subject-data-%s", appEnv),
						"awslogs-region":        "us-east-2",
						"awslogs-stream-prefix": "tailscale",
					},
				},
			},
		}
		b, err := json.Marshal(defs)
		return string(b), err
	}).(pulumi.StringOutput)

	taskDef, err := ecs.NewTaskDefinition(ctx, "task-def", &ecs.TaskDefinitionArgs{
		Family:                  pulumi.Sprintf("subject-data-%s", appEnv),
		Cpu:                     pulumi.String("512"),
		Memory:                  pulumi.String("1024"),
		NetworkMode:             pulumi.String("awsvpc"),
		RequiresCompatibilities: pulumi.StringArray{pulumi.String("FARGATE")},
		ExecutionRoleArn:        execRole.Arn,
		TaskRoleArn:             taskRole.Arn,
		ContainerDefinitions:    containerDefs,
		Tags:                    tags("task-def", appEnv),
	})
	if err != nil {
		return nil, err
	}

	// --- Cloud Map service discovery (for API Gateway VPC Link) ---

	namespace, err := servicediscovery.NewPrivateDnsNamespace(ctx, "namespace", &servicediscovery.PrivateDnsNamespaceArgs{
		Name:        pulumi.Sprintf("subject-data-%s.local", appEnv),
		Description: pulumi.Sprintf("subject-data %s service discovery", appEnv),
		Vpc:         net.VPC.ID(),
		Tags:        tags("namespace", appEnv),
	})
	if err != nil {
		return nil, err
	}

	cmService, err := servicediscovery.NewService(ctx, "cm-service", &servicediscovery.ServiceArgs{
		Name:        pulumi.String("property"),
		NamespaceId: namespace.ID(),
		DnsConfig: &servicediscovery.ServiceDnsConfigArgs{
			NamespaceId:   namespace.ID(),
			RoutingPolicy: pulumi.String("MULTIVALUE"),
			DnsRecords: servicediscovery.ServiceDnsConfigDnsRecordArray{
				&servicediscovery.ServiceDnsConfigDnsRecordArgs{
					Type: pulumi.String("A"),
					Ttl:  pulumi.Int(10),
				},
			},
		},
		HealthCheckCustomConfig: &servicediscovery.ServiceHealthCheckCustomConfigArgs{
			FailureThreshold: pulumi.Int(1),
		},
		Tags: tags("cm-service", appEnv),
	})
	if err != nil {
		return nil, err
	}

	// --- ECS service ---

	_, err = ecs.NewService(ctx, "ecs-service", &ecs.ServiceArgs{
		Name:           pulumi.Sprintf("subject-data-%s", appEnv),
		Cluster:        cluster.Arn,
		TaskDefinition: taskDef.Arn,
		DesiredCount:   pulumi.Int(1),
		LaunchType:     pulumi.String("FARGATE"),
		NetworkConfiguration: &ecs.ServiceNetworkConfigurationArgs{
			Subnets:        net.SubnetIDs,
			SecurityGroups: pulumi.StringArray{net.ECSSG.ID()},
			AssignPublicIp: pulumi.Bool(true),
		},
		ServiceRegistries: &ecs.ServiceServiceRegistriesArgs{
			RegistryArn: cmService.Arn,
		},
		Tags: tags("ecs-service", appEnv),
	})
	if err != nil {
		return nil, err
	}

	return &ECSResources{
		Cluster:         cluster,
		CloudMapService: cmService,
	}, nil
}
