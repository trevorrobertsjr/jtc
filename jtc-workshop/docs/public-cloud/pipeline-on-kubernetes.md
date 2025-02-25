---
sidebar_position: 5
---

# Data Pipeline on Kubernetes

We will need the image name that we used when pushing our pipleine container image to Docker Hub. If you followed my naming convention, it should be in the format `dockerhub-username/pipeline`.

Deploy an nginx web server and a busybox pod for network troubleshooting

kubectl logs

kubectl get pods

kubectl describe pod

We did it! We have our data pipeline running in the cloud as a Kubernetes pod.

You may be thinking, "that was an awful lot of setup just to get a workload running," and you would be right. In industry, it is not a recommended practice to do this much manual configuration of a server. Let's see how Infrastructure as Code (IaC) can simplify our lives.



