package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/secretsmanager"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// S3IamResources holds the IAM resources for S3 access.
type S3IamResources struct {
	Role      *iam.Role
	User      *iam.User
	AccessKey *iam.AccessKey
	Secret    *secretsmanager.Secret
}

// createS3Iam creates:
// - IAM role that can be assumed to access S3 (with conditional write policy)
// - IAM user with only sts:AssumeRole (temporary solution for Hetzner)
// - Access key for the user
// - Secrets Manager entry for the secret key
func createS3Iam(ctx *pulumi.Context, appEnv string, bucketName pulumi.StringOutput, kmsKeyArn pulumi.StringOutput) (*S3IamResources, error) {
	// ---------- IAM ROLE ----------
	// Trust policy: allow the AWS account root to assume the role
	// (stable principal; the user's permission policy gates who can actually call AssumeRole)
	trustPolicy := pulumi.String(fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [{
			"Effect": "Allow",
			"Principal": { "AWS": "arn:aws:iam::174581551884:root" },
			"Action": "sts:AssumeRole"
		}]
	}`))

	role, err := iam.NewRole(ctx, fmt.Sprintf("subject-data-ingestion-role-%s", appEnv), &iam.RoleArgs{
		Name:             pulumi.Sprintf("subject-data-ingestion-role-%s", appEnv),
		AssumeRolePolicy: trustPolicy,
		Tags:             tags("role", appEnv),
	})
	if err != nil {
		return nil, err
	}

	// S3 permissions policy (inline JSON)
	// Uses s3:if-none-match condition to prevent overwrites (create-only)
	rolePolicyDoc := bucketName.ApplyT(func(bucket string) (string, error) {
		arn := fmt.Sprintf("arn:aws:s3:::%s", bucket)
		return fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Sid": "EnforceIfNoneMatchForWrites",
					"Effect": "Allow",
					"Action": "s3:PutObject",
					"Resource": "%s/imports/*",
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
							"s3:prefix": "imports/*"
						}
					}
				},
				{
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

	// ---------- IAM USER (temporary for Hetzner) ----------
	// This is a pragmatic solution because the pipeline runs outside AWS.
	// Post-MVP we intend to replace it with OIDC federation.
	user, err := iam.NewUser(ctx, fmt.Sprintf("subject-data-ingestion-user-%s", appEnv), &iam.UserArgs{
		Name: pulumi.Sprintf("subject-data-ingestion-user-%s", appEnv),
		Tags: tags("user", appEnv),
	})
	if err != nil {
		return nil, err
	}

	// User policy: only allow sts:AssumeRole on the above role
	userAssumePolicy := role.Arn.ApplyT(func(roleArn string) (string, error) {
		return fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Action": "sts:AssumeRole",
				"Resource": "%s"
			}]
		}`, roleArn), nil
	}).(pulumi.StringOutput)

	_, err = iam.NewUserPolicy(ctx, "user-assume-role-policy", &iam.UserPolicyArgs{
		User:   user.Name,
		Policy: userAssumePolicy,
	})
	if err != nil {
		return nil, err
	}

	// Access key for the user
	accessKey, err := iam.NewAccessKey(ctx, fmt.Sprintf("subject-data-ingestion-accesskey-%s", appEnv), &iam.AccessKeyArgs{
		User: user.Name,
	})
	if err != nil {
		return nil, err
	}

	// Store the secret access key in AWS Secrets Manager
	secret, err := secretsmanager.NewSecret(ctx, fmt.Sprintf("ingestion-access-key-secret-%s", appEnv), &secretsmanager.SecretArgs{
		Name: pulumi.Sprintf("subject-data-ingestion-key-%s", appEnv),
		Tags: tags("secret", appEnv),
	})
	if err != nil {
		return nil, err
	}
	_, err = secretsmanager.NewSecretVersion(ctx, fmt.Sprintf("ingestion-access-key-version-%s", appEnv), &secretsmanager.SecretVersionArgs{
		SecretId:     secret.ID(),
		SecretString: accessKey.Secret,
	})
	if err != nil {
		return nil, err
	}

	// Export IAM resources for other stacks
	ctx.Export("ingestRoleArn", role.Arn)
	ctx.Export("ingestUserName", user.Name)
	ctx.Export("ingestUserAccessKeyId", accessKey.ID())
	ctx.Export("ingestUserSecretArn", secret.Arn)

	return &S3IamResources{
		Role:      role,
		User:      user,
		AccessKey: accessKey,
		Secret:    secret,
	}, nil
}