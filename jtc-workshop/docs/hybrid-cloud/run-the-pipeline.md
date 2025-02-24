---
sidebar_position: 4
---

# Run our Data Pipeline



<details><summary>What do we have to do before we run our python script?</summary><p>1. Verify our cloud credentials are set</p><p>2. Install the required python packages (possibly with a virtualenv)</p><p>3. Set the required environment variable `SOURCE_PATH`; optionally with a .env file</p></details>

Once we have all of our prerequisites complete, let's run the python script!

```bash
python pipeline.py
```

If the script executes successfully, you should see the following output:

```bash
❯ python pipeline.py
2025-02-24 09:12:30,965 - INFO - Detected S3 source: s3://unique-bucket-name-trevorjtc/nyc_reviews.csv
2025-02-24 09:12:30,965 - INFO - Fetching data from S3: s3://unique-bucket-name-trevorjtc/nyc_reviews.csv
2025-02-24 09:12:30,972 - INFO - Found credentials in shared credentials file: ~/.aws/credentials
2025-02-24 09:12:31,530 - INFO - Successfully read 1000 rows from S3.
2025-02-24 09:12:31,532 - INFO - Dropping the country column.
2025-02-24 09:12:31,535 - INFO - Writing Parquet data to S3: s3://unique-bucket-name-trevorjtc/cleaned/nyc_reviews.parquet
2025-02-24 09:12:32,330 - INFO - Successfully wrote Parquet file to S3.
2025-02-24 09:12:32,332 - INFO - Writing JSON data to S3: s3://unique-bucket-name-trevorjtc/cleaned/nyc_reviews.json
2025-02-24 09:12:33,576 - INFO - Successfully wrote JSON file to S3.
2025-02-24 09:12:33,580 - INFO - Pipeline completed successfully. Files saved to S3.
```

When I look at my bucket to verify, I should see my new files there:
```bash
❯ aws s3 ls s3://unique-bucket-name-trevorjtc/cleaned/
2025-02-24 09:12:33     441276 nyc_reviews.json
2025-02-24 09:12:32     157812 nyc_reviews.parquet
```

We have our working pipeline, and we want to share it with our co-workers. Let's take a look at a popular way to do this: Containers