import json
import boto3
import os

cloudfront = boto3.client("cloudfront")

def handler(event, context):
    distribution_id = os.environ["DISTRIBUTION_ID"]

    try:
        invalidation = cloudfront.create_invalidation(
            DistributionId=distribution_id,
            InvalidationBatch={
                "Paths": {"Quantity": 1, "Items": ["/*"]},
                "CallerReference": str(context.aws_request_id),
            },
        )

        return {
            "statusCode": 200,
            "body": json.dumps("Invalidation created: " + invalidation["Invalidation"]["Id"]),
        }

    except Exception as e:
        return {
            "statusCode": 500,
            "body": json.dumps("Error: " + str(e)),
        }
