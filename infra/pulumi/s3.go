package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// S3BucketResources holds the S3 bucket resources.
type S3BucketResources struct {
	Bucket *s3.BucketV2
}

// createS3Bucket creates an S3 bucket for ingestion with:
// - versioning enabled
// - SSE-S3 default encryption (AES-256, free, no key management)
// - lifecycle rule (expire noncurrent versions after 90 days)
// - public access block (fully private)
func createS3Bucket(ctx *pulumi.Context, appEnv string) (*S3BucketResources, error) {
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

	// SSE-S3 encryption (AES-256, Amazon-managed keys — free, no key policy needed)
	_, err = s3.NewBucketServerSideEncryptionConfigurationV2(ctx, "bucket-encryption", &s3.BucketServerSideEncryptionConfigurationV2Args{
		Bucket: bucket.Bucket,
		Rules: s3.BucketServerSideEncryptionConfigurationV2RuleArray{
			&s3.BucketServerSideEncryptionConfigurationV2RuleArgs{
				ApplyServerSideEncryptionByDefault: &s3.BucketServerSideEncryptionConfigurationV2RuleApplyServerSideEncryptionByDefaultArgs{
					SseAlgorithm: pulumi.String("AES256"),
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

	return &S3BucketResources{
		Bucket: bucket,
	}, nil
}
