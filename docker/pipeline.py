import os
import pandas as pd
import pyarrow.parquet as pq
from urllib.parse import urlparse
import boto3
from google.cloud import storage
from io import BytesIO
import logging

# Configure logging
logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")
logger = logging.getLogger(__name__)

# Get the input path from environment variable
SOURCE_PATH = os.getenv("SOURCE_PATH")

if not SOURCE_PATH:
    logger.error("Environment variable SOURCE_PATH is required.")
    raise ValueError("Environment variable SOURCE_PATH is required.")

def is_s3_url(url):
    return url.startswith("s3://")

def is_gcs_url(url):
    return url.startswith("gs://")

def read_from_s3(bucket_name, object_key):
    """Reads a CSV file from S3 into a Pandas DataFrame."""
    logger.info(f"Fetching data from S3: s3://{bucket_name}/{object_key}")
    try:
        s3_client = boto3.client("s3")
        response = s3_client.get_object(Bucket=bucket_name, Key=object_key)
        df = pd.read_csv(response["Body"])
        logger.info(f"Successfully read {len(df)} rows from S3.")
        return df
    except Exception as e:
        logger.error(f"Failed to read from S3: {e}")
        raise

def read_from_gcs(bucket_name, object_key):
    """Reads a CSV file from GCS into a Pandas DataFrame."""
    logger.info(f"Fetching data from GCS: gs://{bucket_name}/{object_key}")
    try:
        client = storage.Client()
        bucket = client.bucket(bucket_name)
        blob = bucket.blob(object_key)
        data = blob.download_as_bytes()
        df = pd.read_csv(BytesIO(data))
        logger.info(f"Successfully read {len(df)} rows from GCS.")
        return df
    except Exception as e:
        logger.error(f"Failed to read from GCS: {e}")
        raise

def write_to_s3(bucket_name, output_key, df):
    """Writes a DataFrame to S3 in Parquet format."""
    logger.info(f"Writing data to S3: s3://{bucket_name}/{output_key}")
    try:
        buffer = BytesIO()
        df.to_parquet(buffer, engine="pyarrow", index=False)
        buffer.seek(0)
        s3_client = boto3.client("s3")
        s3_client.put_object(Bucket=bucket_name, Key=output_key, Body=buffer.getvalue())
        logger.info(f"Successfully wrote data to S3.")
    except Exception as e:
        logger.error(f"Failed to write to S3: {e}")
        raise

def write_to_gcs(bucket_name, output_key, df):
    """Writes a DataFrame to GCS in Parquet format."""
    logger.info(f"Writing data to GCS: gs://{bucket_name}/{output_key}")
    try:
        client = storage.Client()
        bucket = client.bucket(bucket_name)
        blob = bucket.blob(output_key)
        buffer = BytesIO()
        df.to_parquet(buffer, engine="pyarrow", index=False)
        buffer.seek(0)
        blob.upload_from_file(buffer, content_type="application/octet-stream")
        logger.info(f"Successfully wrote data to GCS.")
    except Exception as e:
        logger.error(f"Failed to write to GCS: {e}")
        raise

def process_data():
    """Determines the source, reads the CSV, and writes as a Parquet file."""
    parsed_url = urlparse(SOURCE_PATH)
    bucket_name = parsed_url.netloc
    object_key = parsed_url.path.lstrip("/")

    # Define output key with 'cleaned/' prefix
    output_key = f"cleaned/{object_key.rsplit('/', 1)[-1].replace('.csv', '.parquet')}"

    try:
        if is_s3_url(SOURCE_PATH):
            logger.info(f"Detected S3 source: {SOURCE_PATH}")
            df = read_from_s3(bucket_name, object_key)
            write_to_s3(bucket_name, output_key, df)
            logger.info(f"Pipeline completed successfully. Parquet file saved at s3://{bucket_name}/{output_key}")

        elif is_gcs_url(SOURCE_PATH):
            logger.info(f"Detected GCS source: {SOURCE_PATH}")
            df = read_from_gcs(bucket_name, object_key)
            write_to_gcs(bucket_name, output_key, df)
            logger.info(f"Pipeline completed successfully. Parquet file saved at gs://{bucket_name}/{output_key}")

        else:
            logger.error("Invalid source path. Must start with 's3://' or 'gs://'.")
            raise ValueError("Invalid source path. Must start with 's3://' or 'gs://'.")

    except Exception as e:
        logger.error(f"Pipeline execution failed: {e}")
        raise

if __name__ == "__main__":
    process_data()
