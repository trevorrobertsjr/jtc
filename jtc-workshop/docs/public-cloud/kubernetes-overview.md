---
sidebar_position: 2
---

# Kubernetes Primer

## What is Kubernetes?
An open-source container orchestration platform that automates deployment, scaling, and management of containerized applications across a cluster of nodes (i.e. servers). It provides features like self-healing, service discovery, load balancing, and automated rollouts/rollbacks.

You may sometimes see Kubernetes abbreviated as "K8s" due to the 8 letters (ubernete) between "K" and "s" in the name. Folks in the Kubernetes community will pronounce that as "Kates" or "K-eights"

![Architecture Diagram](/img/components-of-kubernetes.svg)

There are many important concepts of Kubernetes, but we will focus on just a few to get us started:
**Control Plane** - the portion of Kubernetes components that are responsible for accepting user requests and making sure the worker nodes execute them.
**Data Plane** - the nodes that are part of the Kubernetes cluster that run the containers that the user requests to be run.
**kubelet** - a process that runs on every node in a Kubernetes cluster. It is responsible for executing the containers and making sure they stay running. The kubelet regularly polls the control plane to determine which containers should be running. If a container dies unexpectedly, the kubelet restarts that container at the next polling interval.
**kube-proxy** - a process that runs on every node in a Kubernetes cluster to control network access to workloads running on its node. By default, the kube-proxy uses iptables rules to do this, but there are alternatives like ipvs and ebpf.
**pod** - a grouping of one or more related containers than run on a single node. I say "one or more" because you may have an application container along with a logging container that need to run together for the logs to be properly propagated. The operations team may choose to run both of these containers together in a Kubernetes pod definition.
**kubectl** - the CLI utility for developers to interact with the Kubernetes API. The Kubernetes community is a little divided over whether to pronounce kubectl as "kube cuttle" or "kube control"

Let's install Docker and Kubernetes on our cloud instance:

```bash
# Update APT and install dependencies
sudo apt-get update
sudo apt-get install -y ca-certificates curl gnupg

# Ensure the keyrings directory exists
sudo install -m 0755 -d /etc/apt/keyrings

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

# Install kubectl
sudo curl -fsSLo /usr/local/bin/kubectl "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo chmod 555 /usr/local/bin/kubectl

# Install Minikube
curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
chmod +x minikube
sudo mv minikube /usr/local/bin/

# Verify installations
docker --version
kubectl version --client
minikube version
```

We have all the tools installed to run workloads on Kubernetes. Let's get a simple test workload going.



