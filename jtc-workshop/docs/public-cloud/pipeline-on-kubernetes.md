---
sidebar_position: 5
---

# Data Pipeline on Kubernetes

First, make sure we have the image name that we used when we created our pipline container image to Docker Hub. If you followed my naming convention, it should be in the format `dockerhub-username/pipeline`.

Remember the concerns we ran into while running the container ourselves with Docker:
1. We needed to provide the AWS credentials to container.
2. We needed to provide the path to the file via an environment variable.

## Kubernetes Secrets
Kubernetes has a resource type called **Secret** that allows you to store base64-encoded strings for secure use by your pods as a mounted volume or as environment variables.

:::danger

Base64 encoding is not a recommended strategy for protecting secrets. The K8s built support for integrating proven secure vendor tools like HashiCorp Vault, AWS Secrets Manager, etc. with a capability call the [External Secrets Operator](https://external-secrets.io/latest/). For production environments, use these types of solutions instead. 

:::

On your commandline, use the base64 command to encode your your AWS credentials. Windows users please look up the Powershell equivalent to do this if you are not using Windows Subsystem for Linux.

If these are my credentials (they certainly are not):
```bash
AWS Access Key ID:     AKIAEXAMPLE123456789
AWS Secret Access Key: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
```

Run these commands to encode them:
```bash
echo -n "AKIAEXAMPLE123456789" | base64
echo -n "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" | base64
```

Sample output:
```
QUtJQUVYQU1QTEUxMjM0NTY3ODk=
d0phbHJYVXRuRkVNSS9LN01ERU5HL2JQeFJmaUNZRVhBTVBMRUtFWQ==
```

Create a Kubernetes **secret** resource to contain our AWS credentials
```yaml title="aws-secret.yaml"
apiVersion: v1
kind: Secret
metadata:
  name: aws-secret
  namespace: default
type: Opaque
data:
  AWS_ACCESS_KEY_ID: QUtJQUVYQU1QTEUxMjM0NTY3ODk=
  AWS_SECRET_ACCESS_KEY: d0phbHJYVXRuRkVNSS9LN01ERU5HL2JQeFJmaUNZRVhBTVBMRUtFWQ==
```

Now, let's use this secret in our pipeline pod manifest along with specifying an environment variable:

```yaml title="aws-pipeline.yaml"
apiVersion: v1
kind: Pod
metadata:
  name: csv-to-parquet
spec:
  containers:
  - name: csv-to-parquet
    image: trevorrobertsjr/pipeline:latest
    env:
    - name: SOURCE_PATH
      value: "s3://your-bucket-name/nyc_reviews.csv"
    - name: AWS_ACCESS_KEY_ID
      valueFrom:
        secretKeyRef:
          name: aws-secret
          key: AWS_ACCESS_KEY_ID
    - name: AWS_SECRET_ACCESS_KEY
      valueFrom:
        secretKeyRef:
          name: aws-secret
          key: AWS_SECRET_ACCESS_KEY
```

```bash
kubectl apply -f aws-secret.yaml
kubectl apply -f aws-pipeline.yaml
```

Check the logs of your pod to see if you have any issues
```bash
kubectl get pods
kubectl logs your-pod-name
```

You may be thinking, "that was an awful lot of setup just to get a workload running," and you would be right!

In industry, it is not a recommended practice to do this much manual configuration of a server that will be deployed into production. Let's see how Infrastructure as Code (IaC) can simplify our lives by deploying cloud resources for us in a repeatable manner.



