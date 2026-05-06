// Package main is the Pulumi entry point for subject-data infrastructure.
//
// Resources: VPC (public subnets), RDS PostgreSQL, ECS Fargate cluster + service,
// API Gateway HTTP API with IAM auth, Tailscale sidecar for backend connectivity.
//
// This stack is self-contained — no StackReferences to other services.
// subject-data lives in its own AWS account (per-service-per-env topology).
package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
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

		// Accounts allowed to assume the S3 ingestion bucket writer role.
		// Distinct from callerAccounts (those are API Gateway invokers).
		var ingestWriterAccountIDs []string
		cfg.RequireObject("ingestWriterAccounts", &ingestWriterAccountIDs)

		// Bucket owner is the current account (where this stack is deployed)
		currentIdentity, err := aws.GetCallerIdentity(ctx, &aws.GetCallerIdentityArgs{})
		if err != nil {
			return err
		}
		bucketOwnerAccountID := currentIdentity.AccountId

		net, err := createNetworking(ctx, appEnv, vpcCidr)
		if err != nil {
			return err
		}

		// 2. S3 bucket for ingestion
		s3Bucket, err := createIngestionS3Bucket(ctx, appEnv)
		if err != nil {
			return err
		}

		// 3. IAM role for cross-account S3 access (Data Ingest pipeline)
		s3Iam, err := createS3Iam(ctx, appEnv, s3Bucket.Bucket.Bucket, bucketOwnerAccountID, ingestWriterAccountIDs)
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

		ecsRes, err := createECS(ctx, appEnv, net, ecrRes.RepositoryUrl, rdsRes, s3Bucket)
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

		// Ingestion bucket and IAM
		ctx.Export("ingestBucketName", s3Bucket.Bucket.Bucket)
		ctx.Export("ingestRoleArn", s3Iam.Role.Arn)

		return nil
	})
}
