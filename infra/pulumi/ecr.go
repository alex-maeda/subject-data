package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ecr"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// lookupECR looks up the ECR repository created during bootstrap.
// ECR is a bootstrap resource (like the S3 state bucket) because the
// CI/CD pipeline pushes images before pulumi up runs.
func lookupECR(ctx *pulumi.Context) (*ecr.LookupRepositoryResult, error) {
	return ecr.LookupRepository(ctx, &ecr.LookupRepositoryArgs{
		Name: "subject-data",
	})
}
