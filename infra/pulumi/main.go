// Package main is the Pulumi entry point for subject-data infrastructure.
//
// Resources: VPC (public subnets), RDS PostgreSQL, ECS Fargate cluster + service,
// API Gateway HTTP API with IAM auth, Tailscale sidecar for backend connectivity.
//
// This stack is self-contained — no StackReferences to other services.
// subject-data lives in its own AWS account (per-service-per-env topology).
package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")
		appEnv := cfg.Require("appEnv")
		vpcCidr := cfg.Require("vpcCidr")

		var callerAccounts map[string]string
		cfg.RequireObject("callerAccounts", &callerAccounts)

		net, err := createNetworking(ctx, appEnv, vpcCidr)
		if err != nil {
			return err
		}

		ecrRes, err := lookupECR(ctx)
		if err != nil {
			return err
		}

		rdsRes, err := createRDS(ctx, appEnv, net)
		if err != nil {
			return err
		}

		ecsRes, err := createECS(ctx, appEnv, net, ecrRes.RepositoryUrl, rdsRes)
		if err != nil {
			return err
		}

		apigw, err := createAPIGateway(ctx, appEnv, callerAccounts, net, ecsRes.CloudMapService)
		if err != nil {
			return err
		}

		// --- Outputs ---

		ctx.Export("apiEndpoint", apigw.Stage.InvokeUrl)
		ctx.Export("apiId", apigw.API.ID())
		invokeRoleArns := pulumi.StringMap{}
		for name, role := range apigw.InvokeRoles {
			invokeRoleArns[name] = role.Arn
		}
		ctx.Export("invokeRoleArns", invokeRoleArns)
		ctx.Export("ecrRepoUrl", pulumi.String(ecrRes.RepositoryUrl))
		ctx.Export("clusterName", ecsRes.Cluster.Name)
		ctx.Export("vpcId", net.VPC.ID())
		ctx.Export("dbEndpoint", rdsRes.Instance.Endpoint)

		return nil
	})
}
