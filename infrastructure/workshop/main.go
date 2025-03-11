package main

import (
	"encoding/json"
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/cloudfront"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/route53"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Retrieve AWS Account ID
		callerIdentity, err := aws.GetCallerIdentity(ctx, nil)
		if err != nil {
			return err
		}

		awsAccountID := callerIdentity.AccountId
		// Load configuration variables
		conf := config.New(ctx, "")
		bucketName := conf.Require("bucketName")         // Pulumi config for S3 bucket name
		hostedZoneId := conf.Require("hostedZoneId")     // Pulumi config for Route 53 Hosted Zone ID
		acmCertificate := conf.Require("acmCertificate") // Pulumi config for ACM Certificate

		// Define constants
		githubRepo := "trevorrobertsjr/jtc"
		awsRegion := "us-east-1"
		lambdaFunctionName := "invalidateCacheLambda"
		siteName := "jtc.wanfooru.com"

		// Create S3 bucket for website hosting
		s3Bucket, err := s3.NewBucket(ctx, "websiteBucket", &s3.BucketArgs{
			Bucket: pulumi.String(bucketName),
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
			}`, bucketName)),
		})
		if err != nil {
			return err
		}

		// CloudFront Distribution for S3 website
		cf, err := cloudfront.NewDistribution(ctx, "websiteDistribution", &cloudfront.DistributionArgs{
			Enabled:           pulumi.Bool(true),
			DefaultRootObject: pulumi.String("index.html"),
			Origins: cloudfront.DistributionOriginArray{
				&cloudfront.DistributionOriginArgs{
					DomainName: s3Bucket.WebsiteEndpoint,
					OriginId:   pulumi.String("S3WebsiteOrigin"),
					CustomOriginConfig: &cloudfront.DistributionOriginCustomOriginConfigArgs{
						OriginProtocolPolicy: pulumi.String("http-only"),
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
				pulumi.String(siteName),
			},
			ViewerCertificate: &cloudfront.DistributionViewerCertificateArgs{
				AcmCertificateArn: pulumi.String(acmCertificate),
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

		// GitHub OIDC IAM Role for GitHub Actions
		trustPolicy, err := json.Marshal(map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []map[string]interface{}{
				{
					"Effect": "Allow",
					"Principal": map[string]string{
						"Federated": "arn:aws:iam::" + awsAccountID + ":oidc-provider/token.actions.githubusercontent.com",
					},
					"Action": "sts:AssumeRoleWithWebIdentity",
					"Condition": map[string]interface{}{
						"StringEquals": map[string]string{
							"token.actions.githubusercontent.com:aud": "sts.amazonaws.com",
							"token.actions.githubusercontent.com:sub": "repo:" + githubRepo + ":ref:refs/heads/main",
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}

		iamRole, err := iam.NewRole(ctx, "githubActionsRole", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(trustPolicy),
		})
		if err != nil {
			return err
		}

		// Attach IAM Policy using NewRolePolicy and ApplyT()
		_, err = iam.NewRolePolicy(ctx, "githubActionsPolicy", &iam.RolePolicyArgs{
			Role: iamRole.Name,
			Policy: cf.ID().ApplyT(func(cfID pulumi.ID) (string, error) {
				return fmt.Sprintf(`{
					"Version": "2012-10-17",
					"Statement": [
						{
							"Effect": "Allow",
							"Action": ["s3:PutObject", "s3:ListBucket", "s3:DeleteObject", "s3:GetObject"],
							"Resource": ["arn:aws:s3:::%s", "arn:aws:s3:::%s/*"]
						},
						{
							"Effect": "Allow",
							"Action": ["lambda:InvokeFunction"],
							"Resource": "arn:aws:lambda:%s:*:function:%s"
						},
						{
							"Effect": "Allow",
							"Action": ["cloudfront:CreateInvalidation"],
							"Resource": "arn:aws:cloudfront::*:distribution/%s"
						},
						{
							"Effect": "Allow",
							"Action": "sts:AssumeRoleWithWebIdentity",
							"Resource": "*"
						}
					]
				}`, bucketName, bucketName, awsRegion, lambdaFunctionName, string(cfID)), nil
			}).(pulumi.StringOutput),
		})

		if err != nil {
			return err
		}

		// Route 53 Record
		_, err = route53.NewRecord(ctx, "websiteRecord", &route53.RecordArgs{
			ZoneId: pulumi.String(hostedZoneId), // Using Pulumi config for Hosted Zone ID
			Name:   pulumi.String(siteName),
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

		// Export values
		ctx.Export("S3WebsiteURL", pulumi.Sprintf("http://%s.s3-website-%s.amazonaws.com", bucketName, awsRegion))
		ctx.Export("CloudFrontDomain", cf.DomainName)
		ctx.Export("S3BucketName", s3Bucket.Bucket)
		ctx.Export("CloudFrontDistributionID", cf.ID())
		ctx.Export("GitHubActionsRoleARN", iamRole.Arn)

		return nil
	})
}
