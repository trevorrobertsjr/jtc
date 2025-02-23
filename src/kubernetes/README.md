```bash
minikube start --driver=docker
kubectl get nodes
```

## AWS
open your ~/.aws/credentiasl file
```bash
echo -n "your-access-key" | base64
echo -n "your-secret-access-key" | base64
```

edit `aws-secret.yaml` and put the base64-encoded value with its corresponding key

on your cloud instance, run
```bash
kubectl apply -f aws-secret.yaml
```

## Google Cloud
