package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/kms"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// S3BucketResources holds the S3 bucket and its KMS key.
type S3BucketResources struct {
	Bucket *s3.BucketV2
	KmsKey *kms.Key
}

// createS3Bucket creates an S3 bucket for ingestion with:
// - versioning enabled
// - KMS encryption (customer-managed key)
// - lifecycle rule (expire noncurrent versions after 90 days)
// - public access block (fully private)
func createS3Bucket(ctx *pulumi.Context, appEnv string) (*S3BucketResources, error) {
	// KMS key for bucket encryption
	kmsKey, err := kms.NewKey(ctx, "ingestion-bucket-key", &kms.KeyArgs{
		Description: pulumi.Sprintf("KMS key for subject-data-ingestion-%s bucket", appEnv),
		Tags:        tags("kms", appEnv),
	})
	if err != nil {
		return nil, err
	}

	// S3 bucket (name includes environment suffix)
	bucket, err := s3.NewBucketV2(ctx, "bucket", &s3.BucketV2Args{
		Bucket: pulumi.Sprintf("subject-data-ingestion-%s", appEnv),
		Tags:   tags("bucket", appEnv),
	})
	if err != nil {
		return nil, err
	}

	// Enable versioning
	_, err = s3.NewBucketVersioningV2(ctx, "bucket-versioning", &s3.BucketVersioningV2Args{
		Bucket: bucket.Bucket,
		VersioningConfiguration: &s3.BucketVersioningV2VersioningConfigurationArgs{
			Status: pulumi.String("Enabled"),
		},
	})
	if err != nil {
		return nil, err
	}

	// KMS encryption configuration
	_, err = s3.NewBucketServerSideEncryptionConfigurationV2(ctx, "bucket-encryption", &s3.BucketServerSideEncryptionConfigurationV2Args{
		Bucket: bucket.Bucket,
		Rules: s3.BucketServerSideEncryptionConfigurationV2RuleArray{
			&s3.BucketServerSideEncryptionConfigurationV2RuleArgs{
				ApplyServerSideEncryptionByDefault: &s3.BucketServerSideEncryptionConfigurationV2RuleApplyServerSideEncryptionByDefaultArgs{
					SseAlgorithm:   pulumi.String("aws:kms"),
					KmsMasterKeyId: kmsKey.Arn,
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// Lifecycle: expire noncurrent versions after 90 days
	_, err = s3.NewBucketLifecycleConfigurationV2(ctx, "bucket-lifecycle", &s3.BucketLifecycleConfigurationV2Args{
		Bucket: bucket.Bucket,
		Rules: s3.BucketLifecycleConfigurationV2RuleArray{
			&s3.BucketLifecycleConfigurationV2RuleArgs{
				Id:     pulumi.String("expire-old-versions"),
				Status: pulumi.String("Enabled"),
				NoncurrentVersionExpiration: &s3.BucketLifecycleConfigurationV2RuleNoncurrentVersionExpirationArgs{
					NoncurrentDays: pulumi.Int(90),
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// Block all public access
	_, err = s3.NewBucketPublicAccessBlock(ctx, "bucket-public-access-block", &s3.BucketPublicAccessBlockArgs{
		Bucket:                bucket.ID(),
		BlockPublicAcls:       pulumi.Bool(true),
		BlockPublicPolicy:     pulumi.Bool(true),
		IgnorePublicAcls:      pulumi.Bool(true),
		RestrictPublicBuckets: pulumi.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	// Export values for other stacks (e.g., IAM)
	ctx.Export("ingestBucketName", bucket.Bucket)
	ctx.Export("ingestKmsKeyArn", kmsKey.Arn)

	return &S3BucketResources{
		Bucket: bucket,
		KmsKey: kmsKey,
	}, nil
}