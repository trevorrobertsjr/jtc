# Use official Python image
FROM python:3.12-slim

# Set working directory
WORKDIR /app

# Copy files to container
COPY pipeline.py requirements.txt ./

# Install dependencies
RUN pip install --no-cache-dir -r requirements.txt

# # Set environment variable (example, replace when running)
# ENV SOURCE_PATH="s3://my-bucket/path/to/data.csv"

# Run the script
CMD ["python", "pipeline.py"]
