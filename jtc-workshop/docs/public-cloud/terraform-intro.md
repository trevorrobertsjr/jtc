---
sidebar_position: 6
---

# Automate Cloud Infrastructure Deployment with Terraform

Even though I simplified the experience of creating an environment by providing screenshots and the Linux commands in advance, there are many screens to point and click through and many commands to run. This translates to numerous opportunities to introduce human error to the equation and potentially cause a vulnerability from a misconfiguration.

How do we avoid those complications? Enter Infrastructure as Code (IaC). The most popular tool on the market is [HashiCorp Terraform](https://www.terraform.io/)

Terraform uses a domain-specific language (DSL) called the HashiCorp Configuration Language (HCL).

Here is a very simple example of how to get an EC2 instance up and running:

```hcl
# main.tf
provider "aws" {
  region = "us-east-2"
}

# EC2 Instance
resource "aws_instance" "hello_world" {
  ami           = "ami-0c55b159cbfafe1f0" # Amazon Linux 2 AMI (Change as needed)
  instance_type = "t2.micro"
  user_data = <<-EOF
              #!/bin/bash
              yum install -y httpd
              echo "Hello, World!" > /var/www/html/index.html
              systemctl start httpd
              systemctl enable httpd
              EOF

  tags = {
    Name = "HelloWorldInstance"
  }
}

```

If you want to try out this Terraform script, make a new subdirectory on your computer and save this content to a filename `main.tf`

In that same directory, run this command to prepare your computer to run the script

```bash
terraform init
```

```bash title="Sample Output"
Initializing the backend...

Initializing provider plugins...
- Finding latest version of hashicorp/aws...
- Installing hashicorp/aws v5.88.0...
- Installed hashicorp/aws v5.88.0 (signed by HashiCorp)

Terraform has created a lock file .terraform.lock.hcl to record the provider
selections it made above. Include this file in your version control repository
so that Terraform can guarantee to make the same selections by default when
you run "terraform init" in the future.

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
```

Next, run this command to apply the Terraform manifest.
```bash
terraform apply
```

```bash title="Sample Output"
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following
symbols:
  + create

Terraform will perform the following actions:

  # aws_instance.hello_world will be created
  + resource "aws_instance" "hello_world" {
      + ami                                  = "ami-0c55b159cbfafe1f0"
      + arn                                  = (known after apply)
      + associate_public_ip_address          = (known after apply)
      + availability_zone                    = (known after apply)
      + cpu_core_count                       = (known after apply)
      + cpu_threads_per_core                 = (known after apply)
      + disable_api_stop                     = (known after apply)
      + disable_api_termination              = (known after apply)
      + ebs_optimized                        = (known after apply)
      + enable_primary_ipv6                  = (known after apply)
      + get_password_data                    = false
      + host_id                              = (known after apply)
      + host_resource_group_arn              = (known after apply)
      + iam_instance_profile                 = (known after apply)
      + id                                   = (known after apply)
      + instance_initiated_shutdown_behavior = (known after apply)
      + instance_lifecycle                   = (known after apply)
      + instance_state                       = (known after apply)
      + instance_type                        = "t2.micro"
      + ipv6_address_count                   = (known after apply)
      + ipv6_addresses                       = (known after apply)
      + key_name                             = (known after apply)
      + monitoring                           = (known after apply)
      + outpost_arn                          = (known after apply)
      + password_data                        = (known after apply)
      + placement_group                      = (known after apply)
      + placement_partition_number           = (known after apply)
      + primary_network_interface_id         = (known after apply)
      + private_dns                          = (known after apply)
      + private_ip                           = (known after apply)
      + public_dns                           = (known after apply)
      + public_ip                            = (known after apply)
      + secondary_private_ips                = (known after apply)
      + security_groups                      = (known after apply)
      + source_dest_check                    = true
      + spot_instance_request_id             = (known after apply)
      + subnet_id                            = (known after apply)
      + tags                                 = {
          + "Name" = "HelloWorldInstance"
        }
      + tags_all                             = {
          + "Name" = "HelloWorldInstance"
        }
      + tenancy                              = (known after apply)
      + user_data                            = "4c434995ac652805d4de830dbcacaf3155ef7be5"
      + user_data_base64                     = (known after apply)
      + user_data_replace_on_change          = false
      + vpc_security_group_ids               = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: 
```
Enter `yes` to proceed

There you have it, an automated deployment of a Terraform server

Before we move on, let's destroy what we created.
```bash
terraform destroy
```

Make sure to type `yes` at the prompt!

Let's see what a Terraform file looks like to deploy our infrastructure on AWS and on GCP!