package s3_manager

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func CreateS3Manager(ctx *pulumi.Context) error {
	var (
		s3ListInfo []AwsS3infoInterface
	)
	cfg := config.New(ctx, "")
	if err := cfg.TryObject("pulumiAWSS3", &s3ListInfo); err != nil {
		ctx.Log.Warn("PulumiAWSS3 Not Defined", nil)
		return err
	}

	for _, s3Info := range s3ListInfo {
		bucketInfo, err := s3.NewBucket(ctx, s3Info.Name, &s3.BucketArgs{
			Bucket:       pulumi.String(s3Info.Name),
			ForceDestroy: pulumi.Bool(true),
		})
		if err != nil {
			ctx.Log.Warn("S3 creation failed ", nil)
			return err
		}

		if s3Info.Notification == "true" {
			_, err = s3.NewBucketNotification(ctx, s3Info.Name+"Notification", &s3.BucketNotificationArgs{
				Bucket:      bucketInfo.ID(),
				Eventbridge: pulumi.Bool(true),
			})
			if err != nil {
				ctx.Log.Error("S3 Notification creation failed ", nil)
				return err
			}
		}

		if s3Info.PublicAccess == "true" {
			_, err = s3.NewBucketPublicAccessBlock(ctx, s3Info.Name+"PublicAccessBlock", &s3.BucketPublicAccessBlockArgs{
				Bucket:                bucketInfo.ID(),
				BlockPublicAcls:       pulumi.Bool(true),
				BlockPublicPolicy:     pulumi.Bool(true),
				IgnorePublicAcls:      pulumi.Bool(true),
				RestrictPublicBuckets: pulumi.Bool(true),
			}, pulumi.DependsOn([]pulumi.Resource{
				bucketInfo,
			}))
			if err != nil {
				ctx.Log.Error("S3 policy creation failed ", nil)
				return err
			}
			piblicS3Policy := map[string]interface{}{
				"Version": "2012-10-17",
				"Statement": map[string]interface{}{
					"Effect":    "Allow",
					"Principal": "*",
					"Action":    "s3:GetObject",
					"Resource":  pulumi.Sprintf("arn:aws:s3:::%s/*", bucketInfo.ID()),
				},
			}
			s3.NewBucketAclV2(ctx, s3Info.Name+"Acls", &s3.BucketAclV2Args{
				Bucket: bucketInfo.ID(),
				Acl:    pulumi.String(s3.CannedAclPublicRead),
			}, pulumi.DependsOn([]pulumi.Resource{
				bucketInfo,
			}))
			s3.NewBucketPolicy(ctx, s3Info.Name+"Policy", &s3.BucketPolicyArgs{
				Bucket: bucketInfo.ID(),
				Policy: pulumi.JSONMarshal(piblicS3Policy),
			}, pulumi.DependsOn([]pulumi.Resource{
				bucketInfo,
			}))
		}
	}

	return nil
}
