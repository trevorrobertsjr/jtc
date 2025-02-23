provider "google" {
  region  = "us-central1"
  project = var.gcp_project
}

provider "aws" {
  region = "us-east-2"
}

data "http" "my_ip" {
  url = "https://ifconfig.me/ip"
}

variable "gcp_project" {
  description = "The GCP project ID"
  type        = string
}

variable "aws_ami" {
  description = "The AWS AMI ID"
  type        = string
}

variable "gce_image" {
  description = "The GCE Image Name"
  type        = string
}

resource "random_string" "gcs_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "google_storage_bucket" "gcs_bucket" {
  name     = "gcs-bucket-${random_string.gcs_suffix.result}"
  location = "US"
  force_destroy = true  # Allows deletion of non-empty bucket
}

resource "aws_s3_bucket" "s3_bucket" {
  force_destroy = true  # Allows bucket deletion even if it contains objects
}

resource "tls_private_key" "gcp_ssh_key" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "tls_private_key" "aws_ssh_key" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "google_compute_instance" "gce_instance" {
  name         = "gce-instance"
  machine_type = "e2-small"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = var.gce_image
    }
  }

  network_interface {
    network = "default"
    access_config {}
  }

  metadata = {
    ssh-keys = "ubuntu:${tls_private_key.gcp_ssh_key.public_key_openssh}"
  }

  metadata_startup_script = <<-EOT
    #!/bin/bash
    set -e  # Exit on any error

    # Update APT and install dependencies
    apt-get update
    apt-get install -y ca-certificates curl gnupg

    # Ensure the keyrings directory exists
    install -m 0755 -d /etc/apt/keyrings

    # Add the Google Cloud SDK repository
    echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | \
    tee -a /etc/apt/sources.list.d/google-cloud-sdk.list

    # Import the Google Cloud public key
    # curl -fsSL https://packages.cloud.google.com/apt/doc/apt-key.gpg | tee /usr/share/keyrings/cloud.google.gpg > /dev/null

    # Download Google Cloud SDK’s GPG key and move it to the correct location
    curl -fsSL https://packages.cloud.google.com/apt/doc/apt-key.gpg -o /etc/apt/keyrings/cloud.google.gpg
    chmod a+r /etc/apt/keyrings/cloud.google.gpg

    # Update and install the Google Cloud SDK (`gcloud` CLI)
    apt-get update && apt-get install -y google-cloud-sdk

    # Download Docker's GPG key
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    chmod a+r /etc/apt/keyrings/docker.asc

    # Get Ubuntu release codename
    UBUNTU_CODENAME=$(source /etc/os-release && echo "$${UBUNTU_CODENAME:-$VERSION_CODENAME}")

    # Add the Docker APT repository correctly
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] \
    https://download.docker.com/linux/ubuntu $UBUNTU_CODENAME stable" | tee /etc/apt/sources.list.d/docker.list

    # Update package lists again and install Docker
    apt-get update
    apt-get install -y docker-ce docker-ce-cli containerd.io unzip

    # Add the default user to the Docker group
    usermod -aG docker ubuntu

    # Install kubectl
    curl -fsSLo /usr/local/bin/kubectl "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
    chmod +x /usr/local/bin/kubectl

    # Install Minikube
    curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
    chmod +x minikube
    mv minikube /usr/local/bin/

    # Verify installations
    docker --version
    kubectl version --client
    minikube version
  EOT

  service_account {
    email  = google_service_account.gce_sa.email
    scopes = ["https://www.googleapis.com/auth/devstorage.read_write"]
  }
}

resource "google_service_account" "gce_sa" {
  account_id   = "gce-instance-sa"
  display_name = "GCE Instance Service Account"
}

resource "google_project_iam_member" "gce_storage_access" {
  project = var.gcp_project
  role    = "roles/storage.objectAdmin"
  member  = "serviceAccount:${google_service_account.gce_sa.email}"
}

resource "aws_key_pair" "aws_key" {
  key_name   = "aws-ssh-key"
  public_key = tls_private_key.aws_ssh_key.public_key_openssh
}

resource "aws_instance" "ec2_instance" {
  ami           = var.aws_ami
  instance_type = "t3.small"
  key_name      = aws_key_pair.aws_key.key_name

  vpc_security_group_ids = [aws_security_group.allow_ssh.id]

  user_data = <<-EOT
    #!/bin/bash
    set -e  # Exit on any error

    # Update APT and install dependencies
    apt-get update
    apt-get install -y ca-certificates curl gnupg

    # Ensure the keyrings directory exists
    install -m 0755 -d /etc/apt/keyrings

    # Download Docker's GPG key
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    chmod a+r /etc/apt/keyrings/docker.asc

    # Get Ubuntu release codename
    UBUNTU_CODENAME=$(source /etc/os-release && echo "$${UBUNTU_CODENAME:-$VERSION_CODENAME}")

    # Add the Docker APT repository correctly
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] \
    https://download.docker.com/linux/ubuntu $UBUNTU_CODENAME stable" | tee /etc/apt/sources.list.d/docker.list

    # Update package lists again and install Docker
    apt-get update
    apt-get install -y docker-ce docker-ce-cli containerd.io unzip

    # Add the default user to the Docker group
    usermod -aG docker ubuntu

    # Install AWS CLI
    curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
    unzip awscliv2.zip
    ./aws/install

    # Install kubectl
    curl -fsSLo /usr/local/bin/kubectl "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
    chmod +x /usr/local/bin/kubectl

    # Install Minikube
    curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
    chmod +x minikube
    mv minikube /usr/local/bin/

    # Verify installations
    docker --version
    kubectl version --client
    minikube version
  EOT

  tags = {
    Name = "ec2-instance"
  }

  iam_instance_profile = aws_iam_instance_profile.ec2_profile.name
}

resource "aws_security_group" "allow_ssh" {
  name_prefix = "allow-ssh"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["${chomp(data.http.my_ip.response_body)}/32"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_iam_role" "ec2_role" {
  name = "ec2-s3-access-role"

  assume_role_policy = <<-EOT
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Action": "sts:AssumeRole",
          "Principal": {
            "Service": "ec2.amazonaws.com"
          },
          "Effect": "Allow",
          "Sid": ""
        }
      ]
    }
  EOT
}

resource "aws_iam_policy" "s3_access_policy" {
  name        = "s3-access-policy"
  description = "Allows EC2 instance to access S3 bucket"
  
  policy = <<-EOT
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": [
            "s3:PutObject",
            "s3:GetObject"
          ],
          "Resource": "${aws_s3_bucket.s3_bucket.arn}/*"
        },
        {
          "Effect": "Allow",
          "Action": [
            "s3:ListBucket"
          ],
          "Resource": "${aws_s3_bucket.s3_bucket.arn}"
        }
      ]
    }
  EOT
}

resource "aws_iam_role_policy_attachment" "s3_attach" {
  role       = aws_iam_role.ec2_role.name
  policy_arn = aws_iam_policy.s3_access_policy.arn
}

resource "aws_iam_instance_profile" "ec2_profile" {
  name = "ec2-profile"
  role = aws_iam_role.ec2_role.name
}

resource "local_file" "gcp_private_key" {
  content  = tls_private_key.gcp_ssh_key.private_key_pem
  filename = "${path.module}/gcp-ssh-key.pem"
}

resource "local_file" "aws_private_key" {
  content  = tls_private_key.aws_ssh_key.private_key_pem
  filename = "${path.module}/aws-ssh-key.pem"
}

resource "google_service_account" "gcs_sa" {
  account_id   = "gcs-storage-sa"
  display_name = "GCS Storage Service Account"
}

resource "google_project_iam_member" "gcs_admin_role" {
  project = var.gcp_project
  role    = "roles/storage.objectAdmin"
  member  = "serviceAccount:${google_service_account.gcs_sa.email}"
}

resource "google_service_account_key" "gcs_sa_key" {
  service_account_id = google_service_account.gcs_sa.name
  private_key_type   = "TYPE_GOOGLE_CREDENTIALS_FILE"
}

resource "local_file" "gcs_service_account_json" {
  content  = google_service_account_key.gcs_sa_key.private_key
  filename = "${path.module}/gcs-key.json"
}

output "gcs_bucket_name" {
  value = google_storage_bucket.gcs_bucket.name
}

output "s3_bucket_name" {
  value = aws_s3_bucket.s3_bucket.id
}

output "gce_instance_public_ip" {
  value = google_compute_instance.gce_instance.network_interface[0].access_config[0].nat_ip
}

output "ec2_instance_public_ip" {
  value = aws_instance.ec2_instance.public_ip
}

output "gcs_service_account_key_file" {
  value     = local_file.gcs_service_account_json.filename
  sensitive = false
}