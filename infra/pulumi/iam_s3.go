package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// S3IamResources holds the IAM resources for S3 ingestion access.
type S3IamResources struct {
	Role *iam.Role
}

// createS3Iam creates a cross-account IAM role for the Data Ingest pipeline.
//
// - bucketOwnerAccountID: the AWS account that owns the S3 bucket (current account).
// - ingestWriterAccountID: the external AWS account whose users assume this role to write data.
//
// Flow:
//
//	Pipeline operator's Identity Center session (ingestWriterAccountID)
//	  → aws sts assume-role --role-arn <this role>
//	  → scoped temporary creds for S3 published/* writes
//	  → uploads parquet + manifest.json
//	  → calls POST /v1/ingest-jobs
func createS3Iam(ctx *pulumi.Context, appEnv string, bucketName pulumi.StringOutput, bucketOwnerAccountID, ingestWriterAccountID string) (*S3IamResources, error) {
	// Trust policy: allow both the bucket owner and the ingest writer to assume this role.
	trustPolicy := pulumi.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [{
			"Effect": "Allow",
			"Principal": {
				"AWS": [
					"arn:aws:iam::%s:root",
					"arn:aws:iam::%s:root"
				]
			},
			"Action": "sts:AssumeRole"
		}]
	}`, bucketOwnerAccountID, ingestWriterAccountID)

	role, err := iam.NewRole(ctx, fmt.Sprintf("subject-data-ingestion-role-%s", appEnv), &iam.RoleArgs{
		Name:             pulumi.Sprintf("subject-data-ingestion-role-%s", appEnv),
		AssumeRolePolicy: trustPolicy,
		Tags:             tags("role", appEnv),
	})
	if err != nil {
		return nil, err
	}

	// S3 permissions policy for the Data Ingest pipeline.
	// Layout: published/<batch>/<dataset>_<dataset_version>/{parts, manifest.json}
	// - PutObject with s3:if-none-match to prevent overwrites (create-only)
	// - GetObject for read-back / verification
	// - ListBucket scoped to published/ prefix
	// - Explicit Deny on DeleteObject
	rolePolicyDoc := bucketName.ApplyT(func(bucket string) (string, error) {
		arn := fmt.Sprintf("arn:aws:s3:::%s", bucket)
		return fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Sid": "WriteParquetAndManifest",
					"Effect": "Allow",
					"Action": "s3:PutObject",
					"Resource": "%s/published/*",
					"Condition": {
						"StringEquals": {
							"s3:if-none-match": "*"
						}
					}
				},
				{
					"Sid": "ListBucketScoped",
					"Effect": "Allow",
					"Action": "s3:ListBucket",
					"Resource": "%s",
					"Condition": {
						"StringLike": {
							"s3:prefix": "published/*"
						}
					}
				},
				{
					"Sid": "ReadOwnData",
					"Effect": "Allow",
					"Action": "s3:GetObject",
					"Resource": "%s/published/*"
				},
				{
					"Sid": "DenyDelete",
					"Effect": "Deny",
					"Action": "s3:DeleteObject",
					"Resource": "%s/*"
				}
			]
		}`, arn, arn, arn, arn), nil
	}).(pulumi.StringOutput)

	policy, err := iam.NewPolicy(ctx, fmt.Sprintf("subject-data-ingestion-policy-%s", appEnv), &iam.PolicyArgs{
		Name:   pulumi.Sprintf("subject-data-ingestion-policy-%s", appEnv),
		Policy: rolePolicyDoc,
		Tags:   tags("policy", appEnv),
	})
	if err != nil {
		return nil, err
	}

	_, err = iam.NewRolePolicyAttachment(ctx, "role-policy-attachment", &iam.RolePolicyAttachmentArgs{
		Role:      role.Name,
		PolicyArn: policy.Arn,
	})
	if err != nil {
		return nil, err
	}

	return &S3IamResources{
		Role: role,
	}, nil
}
