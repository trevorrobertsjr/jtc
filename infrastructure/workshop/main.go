package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/cloudfront"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/route53"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create S3 bucket with website hosting enabled
		s3Bucket, err := s3.NewBucket(ctx, "websiteBucket", &s3.BucketArgs{
			Bucket: pulumi.String("jtc-wanfooru-com"),
			Website: &s3.BucketWebsiteArgs{
				IndexDocument: pulumi.String("index.html"),
				ErrorDocument: pulumi.String("error.html"),
			},
		})
		if err != nil {
			return err
		}

		// Enable public access for S3 website endpoint
		_, err = s3.NewBucketPublicAccessBlock(ctx, "publicAccessBlock", &s3.BucketPublicAccessBlockArgs{
			Bucket:                s3Bucket.ID(),
			BlockPublicAcls:       pulumi.Bool(false),
			BlockPublicPolicy:     pulumi.Bool(false),
			IgnorePublicAcls:      pulumi.Bool(false),
			RestrictPublicBuckets: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}

		// S3 Bucket Policy to allow public access for the S3 website
		_, err = s3.NewBucketPolicy(ctx, "bucketPolicy", &s3.BucketPolicyArgs{
			Bucket: s3Bucket.ID(),
			Policy: pulumi.String(fmt.Sprintf(`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Effect": "Allow",
						"Principal": "*",
						"Action": "s3:GetObject",
						"Resource": "arn:aws:s3:::%s/*"
					}
				]
			}`, "jtc-wanfooru-com")),
		})
		if err != nil {
			return err
		}

		// CloudFront Distribution pointing to S3 website endpoint
		cf, err := cloudfront.NewDistribution(ctx, "websiteDistribution", &cloudfront.DistributionArgs{
			Enabled:           pulumi.Bool(true),
			DefaultRootObject: pulumi.String("index.html"),
			Origins: cloudfront.DistributionOriginArray{
				&cloudfront.DistributionOriginArgs{
					DomainName: s3Bucket.WebsiteEndpoint,
					OriginId:   pulumi.String("S3WebsiteOrigin"),
					CustomOriginConfig: &cloudfront.DistributionOriginCustomOriginConfigArgs{
						OriginProtocolPolicy: pulumi.String("http-only"), // S3 website only supports HTTP
						HttpPort:             pulumi.Int(80),
						HttpsPort:            pulumi.Int(443),
						OriginSslProtocols: pulumi.StringArray{
							pulumi.String("TLSv1.2"),
						},
					},
				},
			},
			DefaultCacheBehavior: &cloudfront.DistributionDefaultCacheBehaviorArgs{
				TargetOriginId:       pulumi.String("S3WebsiteOrigin"),
				ViewerProtocolPolicy: pulumi.String("redirect-to-https"),
				AllowedMethods: pulumi.StringArray{
					pulumi.String("GET"),
					pulumi.String("HEAD"),
				},
				CachedMethods: pulumi.StringArray{
					pulumi.String("GET"),
					pulumi.String("HEAD"),
				},
				ForwardedValues: &cloudfront.DistributionDefaultCacheBehaviorForwardedValuesArgs{
					QueryString: pulumi.Bool(false),
					Cookies: &cloudfront.DistributionDefaultCacheBehaviorForwardedValuesCookiesArgs{
						Forward: pulumi.String("none"),
					},
				},
			},
			Aliases: pulumi.StringArray{
				pulumi.String("jtc.wanfooru.com"),
			},
			ViewerCertificate: &cloudfront.DistributionViewerCertificateArgs{
				AcmCertificateArn: pulumi.String("arn:aws:acm:us-east-1:318168271290:certificate/e6a5a02b-9537-4111-8753-9ca709d7b480"),
				SslSupportMethod:  pulumi.String("sni-only"),
			},
			Restrictions: &cloudfront.DistributionRestrictionsArgs{
				GeoRestriction: &cloudfront.DistributionRestrictionsGeoRestrictionArgs{
					RestrictionType: pulumi.String("none"),
				},
			},
		})
		if err != nil {
			return err
		}

		// Preserve the existing Lambda function
		lambdaRole, err := iam.NewRole(ctx, "lambdaRole", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(`{
				"Version": "2012-10-17",
				"Statement": [{
					"Action": "sts:AssumeRole",
					"Principal": { "Service": "lambda.amazonaws.com" },
					"Effect": "Allow"
				}]
			}`),
		})
		if err != nil {
			return err
		}

		_, err = iam.NewRolePolicyAttachment(ctx, "lambdaBasicExecutionAttachment", &iam.RolePolicyAttachmentArgs{
			Role:      lambdaRole.Name,
			PolicyArn: pulumi.String("arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"),
		})
		if err != nil {
			return err
		}

		lambdaFunc, err := lambda.NewFunction(ctx, "invalidateCacheLambda", &lambda.FunctionArgs{
			Runtime: pulumi.String("python3.12"),
			Handler: pulumi.String("index.handler"),
			Role:    lambdaRole.Arn,
			Code: pulumi.NewAssetArchive(map[string]interface{}{
				"folder": pulumi.NewFileArchive("./lambda"),
			}),
			Environment: &lambda.FunctionEnvironmentArgs{
				Variables: pulumi.StringMap{
					"DISTRIBUTION_ID": cf.ID(),
				},
			},
		})
		if err != nil {
			return err
		}

		// Keep Route 53 pointing to CloudFront
		_, err = route53.NewRecord(ctx, "websiteRecord", &route53.RecordArgs{
			ZoneId: pulumi.String("Z01322833I96GV0X0FMD7"),
			Name:   pulumi.String("jtc.wanfooru.com"),
			Type:   pulumi.String("A"),
			Aliases: route53.RecordAliasArray{
				&route53.RecordAliasArgs{
					Name:                 cf.DomainName,
					ZoneId:               cf.HostedZoneId,
					EvaluateTargetHealth: pulumi.Bool(false),
				},
			},
		})
		if err != nil {
			return err
		}

		// Export all values (preserving everything)
		ctx.Export("S3WebsiteURL", pulumi.Sprintf("http://%s.s3-website-us-east-1.amazonaws.com", "jtc-wanfooru-com"))
		ctx.Export("CloudFrontDomain", cf.DomainName)
		ctx.Export("S3BucketName", s3Bucket.Bucket)
		ctx.Export("LambdaFunction", lambdaFunc.Arn)
		ctx.Export("CloudFrontDistributionID", cf.ID())

		return nil
	})
}
