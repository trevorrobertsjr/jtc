# Use an official Python image
FROM python:3.12-slim

# Set working directory
WORKDIR /app

# Copy the requirements file to the container
COPY requirements.txt ./

# Install dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy any remaining files that are needed for the application.
COPY . .

# Run the script
CMD ["python", "pipeline.py"]