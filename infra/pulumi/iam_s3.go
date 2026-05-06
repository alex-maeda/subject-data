package main

import (
	"fmt"
	"strings"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// S3IamResources holds the IAM resources for S3 ingestion access.
type S3IamResources struct {
	Role *iam.Role
}

// createS3Iam creates a cross-account IAM role for writers to the S3 ingestion
// bucket.
//
//   - bucketOwnerAccountID: the AWS account that owns the bucket (current
//     account); included in the trust policy so admins in this account can
//     exercise the role.
//   - ingestWriterAccountIDs: external AWS accounts whose principals can assume
//     this role to PutObject under the bucket. Multiple writers are supported;
//     the trust policy lists each as a Principal.AWS entry.
//
// Flow:
//
//	Writer's Identity Center session (ingestWriterAccountIDs[i])
//	  → aws sts assume-role --role-arn <this role>
//	  → scoped temporary creds for S3 published/* writes
//	  → uploads parquet + manifest.json
//	  → calls POST /v1/ingest-jobs
func createS3Iam(ctx *pulumi.Context, appEnv string, bucketName pulumi.StringOutput, bucketOwnerAccountID string, ingestWriterAccountIDs []string) (*S3IamResources, error) {
	// Build the Principal.AWS array: bucket owner + each writer account root.
	principalARNs := make([]string, 0, len(ingestWriterAccountIDs)+1)
	principalARNs = append(principalARNs, fmt.Sprintf(`"arn:aws:iam::%s:root"`, bucketOwnerAccountID))
	for _, id := range ingestWriterAccountIDs {
		principalARNs = append(principalARNs, fmt.Sprintf(`"arn:aws:iam::%s:root"`, id))
	}
	trustPolicy := pulumi.String(fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [{
			"Effect": "Allow",
			"Principal": {
				"AWS": [%s]
			},
			"Action": "sts:AssumeRole"
		}]
	}`, strings.Join(principalARNs, ", ")))

	role, err := iam.NewRole(ctx, fmt.Sprintf("subject-data-ingestion-role-%s", appEnv), &iam.RoleArgs{
		Name:             pulumi.Sprintf("subject-data-ingestion-role-%s", appEnv),
		AssumeRolePolicy: trustPolicy,
		Tags:             tags("role", appEnv),
	})
	if err != nil {
		return nil, err
	}

	// S3 permissions policy for ingest writers.
	// Layout: published/<batch>/<dataset>_<dataset_version>/{parts, manifest.json}
	// - PutObject with s3:if-none-match to prevent overwrites (create-only)
	// - ListBucket scoped to published/ prefix
	// - Explicit Deny on DeleteObject
	// Note: writers do not need GetObject — job status is exposed via the SDS
	// REST API (GET /v1/ingest-jobs/{id}), not via S3 objects.
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
					"Sid": "DenyDelete",
					"Effect": "Deny",
					"Action": "s3:DeleteObject",
					"Resource": "%s/*"
				}
			]
		}`, arn, arn, arn), nil
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
