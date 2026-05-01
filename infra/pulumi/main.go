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

		// 2. S3 bucket (and KMS key)
		s3Bucket, err := createS3Bucket(ctx, appEnv)
		if err != nil {
			return err
		}

		// 3. IAM for S3 (role, user, policies, access key)
		s3Iam, err := createS3Iam(ctx, appEnv, s3Bucket.Bucket.Bucket, s3Bucket.KmsKey.Arn)
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

		// New outputs for ingestion bucket and IAM
		ctx.Export("ingestBucketName", s3Bucket.Bucket.Bucket)
		ctx.Export("ingestKmsKeyArn", s3Bucket.KmsKey.Arn)
		ctx.Export("ingestRoleArn", s3Iam.Role.Arn)
		ctx.Export("ingestUserName", s3Iam.User.Name)
		ctx.Export("ingestUserAccessKeyId", s3Iam.AccessKey.ID())
		ctx.Export("ingestUserSecretArn", s3Iam.Secret.Arn)

		return nil
	})
}
