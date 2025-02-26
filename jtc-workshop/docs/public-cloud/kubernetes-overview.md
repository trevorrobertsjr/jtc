---
sidebar_position: 3
---

# Kubernetes Primer

## What is Kubernetes?
An open-source container orchestration platform that automates deployment, scaling, and management of containerized applications across a cluster of nodes (i.e. servers). It provides features like self-healing, service discovery, load balancing, and automated rollouts/rollbacks.

You may sometimes see Kubernetes abbreviated as "K8s" due to the 8 letters (ubernete) between "K" and "s" in the name. Folks in the Kubernetes community will pronounce that as "Kates" or "K-eights"

![Architecture Diagram](/img/components-of-kubernetes.svg)

There are many important concepts of Kubernetes, but we will focus on just a few to get us started:

**Control Plane** - the Kubernetes components that handle user requests, maintain the cluster state info, and schedule work on the nodes.

**Data Plane** - the Kubernetes components that run the user workloads (i.e. the **nodes** or servers).

**kubelet** - a process that runs on every node in a Kubernetes cluster. It receives instructions from the control plane on which containers to run. If a container dies unexpectedly, the kubelet restarts that container for you automatically.

**kube-proxy** - a process that runs on every node in a Kubernetes cluster to control network access to workloads running on its node. By default, the kube-proxy uses iptables rules to do this, but there are alternatives like ipvs and ebpf.

**kubectl** - the CLI utility for developers to interact with the Kubernetes API. The Kubernetes community is a little divided over whether to pronounce kubectl as "kube cuttle" or "kube control"

**pod** - a grouping of one or more related containers than run on a single node. I say "one or more" because you may have an application container along with a logging container to direct the logs to your central storage service. For that logging container to get access to the logs, it needs to run on the same machine as the application, and the operations writes the pod definition file (aka **manifest**) accordingly.

**deployment** - a K8s resource to define how a pod should be run and if there should be replicas of the pod

**service** - a K8s resource to serve the application running in a pod, or in a group of pods, to other workloads in the cluster or to users. It essentially acts like a load balancer.

Let's connect to our instance and install Docker and Kubernetes on our cloud instance.



