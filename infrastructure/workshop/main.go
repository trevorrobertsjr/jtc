package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/cloudfront"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/route53"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create S3 bucket
		s3Bucket, err := s3.NewBucket(ctx, "websiteBucket", &s3.BucketArgs{
			Bucket: pulumi.String("jtc-wanfooru-com"),
			Website: &s3.BucketWebsiteArgs{
				IndexDocument: pulumi.String("index.html"),
			},
		})
		if err != nil {
			return err
		}

		// Create CloudFront Origin Access Control (OAC)
		oac, err := cloudfront.NewOriginAccessControl(ctx, "oac", &cloudfront.OriginAccessControlArgs{
			Name:                          pulumi.String("JTC-OAC"),
			Description:                   pulumi.String("Access control for S3 website"),
			OriginAccessControlOriginType: pulumi.String("s3"),
			SigningBehavior:               pulumi.String("always"),
			SigningProtocol:               pulumi.String("sigv4"),
		})
		if err != nil {
			return err
		}

		// Create CloudFront Distribution
		cf, err := cloudfront.NewDistribution(ctx, "websiteDistribution", &cloudfront.DistributionArgs{
			Enabled:           pulumi.Bool(true),
			DefaultRootObject: pulumi.String("index.html"),
			Origins: cloudfront.DistributionOriginArray{
				&cloudfront.DistributionOriginArgs{
					DomainName:            s3Bucket.BucketRegionalDomainName,
					OriginId:              pulumi.String("S3Origin"),
					OriginAccessControlId: oac.ID(),
					S3OriginConfig: &cloudfront.DistributionOriginS3OriginConfigArgs{
						OriginAccessIdentity: pulumi.String(""),
					},
				},
			},
			DefaultCacheBehavior: &cloudfront.DistributionDefaultCacheBehaviorArgs{
				TargetOriginId:       pulumi.String("S3Origin"),
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

		// Create Lambda function
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

		// Attach AWSLambdaBasicExecutionRole policy
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

		// Create Route 53 Record to point to CloudFront
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

		ctx.Export("CloudFrontDomain", cf.DomainName)
		ctx.Export("S3BucketName", s3Bucket.Bucket)
		ctx.Export("LambdaFunction", lambdaFunc.Arn)
		ctx.Export("CloudFrontDistributionID", cf.ID())
		return nil
	})
}
