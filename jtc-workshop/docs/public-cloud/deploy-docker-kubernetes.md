---
sidebar_position: 3
---

# Kubernetes Primer

## Installing Kubernetes on Our Instance

Make sure you are SSH'ed in to your Amazon EC2 instance

Let's prepare our instance to use GPG keys to verify the integrity of the packages we will install:

```bash
# Update APT and install dependencies
sudo apt-get update
sudo apt-get install -y ca-certificates curl gnupg

# Ensure the keyrings directory exists
sudo install -m 0755 -d /etc/apt/keyrings
```

Let's install Docker
```bash
# Download Docker's GPG key
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo tee /etc/apt/keyrings/docker.asc > /dev/null
sudo chmod a+r /etc/apt/keyrings/docker.asc

# Get Ubuntu release codename
UBUNTU_CODENAME=$(source /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}")

# Add the Docker APT repository correctly
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] \
https://download.docker.com/linux/ubuntu $UBUNTU_CODENAME stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Update package lists again and install Docker
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io unzip vim

# Add the default user to the Docker group
sudo usermod -aG docker ubuntu
newgrp docker
```
Let's install Kubernetes
```bash
# Install kubectl - the Kubernetes CLI
sudo curl -fsSLo /usr/local/bin/kubectl "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo chmod 555 /usr/local/bin/kubectl

# Install Minikube
curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
chmod +x minikube
sudo mv minikube /usr/local/bin/

minikube start
```

As minikube starts, you should see output similar to the following:
```bash
ubuntu@ip-172-31-28-91:~$ minikube start
😄  minikube v1.35.0 on Ubuntu 24.04
✨  Automatically selected the docker driver. Other choices: none, ssh

🧯  The requested memory allocation of 1910MiB does not leave room for system overhead (total system memory: 1910MiB). You may face stability issues.
💡  Suggestion: Start minikube with less memory allocated: 'minikube start --memory=1910mb'

📌  Using Docker driver with root privileges
👍  Starting "minikube" primary control-plane node in "minikube" cluster
🚜  Pulling base image v0.0.46 ...
💾  Downloading Kubernetes v1.32.0 preload ...
    > preloaded-images-k8s-v18-v1...:  333.57 MiB / 333.57 MiB  100.00% 150.52 
    > gcr.io/k8s-minikube/kicbase...:  464.06 MiB / 500.31 MiB  92.76% 45.02 Mi
🔥  Creating docker container (CPUs=2, Memory=1910MB) ...
🧯  Docker is nearly out of disk space, which may cause deployments to fail! (90% of capacity). You can pass '--force' to skip this check.
💡  Suggestion: 

    Try one or more of the following to free up space on the device:
    
    1. Run "docker system prune" to remove unused Docker data (optionally with "-a")
    2. Increase the storage allocated to Docker for Desktop by clicking on:
    Docker icon > Preferences > Resources > Disk Image Size
    3. Run "minikube ssh -- docker system prune" if using the Docker container runtime
🍿  Related issue: https://github.com/kubernetes/minikube/issues/9024

🐳  Preparing Kubernetes v1.32.0 on Docker 27.4.1 ...
    ▪ Generating certificates and keys ...
    ▪ Booting up control plane ...
    ▪ Configuring RBAC rules ...
🔗  Configuring bridge CNI (Container Networking Interface) ...
🔎  Verifying Kubernetes components...
    ▪ Using image gcr.io/k8s-minikube/storage-provisioner:v5
🌟  Enabled addons: storage-provisioner, default-storageclass
🏄  Done! kubectl is now configured to use "minikube" cluster and "default" namespace by default
```

Let's make sure everything is installed and running properly
```bash
# Verify installations
minikube version
docker --version
kubectl version --client
kubectl get nodes
```

Your output should look similar to this
```bash
minikube version: v1.35.0
commit: dd5d320e41b5451cdf3c01891bc4e13d189586ed-dirty
Docker version 28.0.0, build f9ced58
Client Version: v1.32.2
Kustomize Version: v5.5.0
NAME       STATUS   ROLES           AGE   VERSION
minikube   Ready    control-plane   41s   v1.32.0
```

We have all the tools installed to run workloads on Kubernetes. Let's get a simple test workload going.



