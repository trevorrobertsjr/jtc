# jtc-workshop-infrastructure

```bash
pulumi config set bucketName YOUR-BUCKET-NAME
pulumi config set hostedZoneId YOURID1234567890
pulumi config set acmCertificate arn:aws:acm:us-east-1:1234567890:certificate/e6a5a02b-use-your-cert-1234567890
```

## Constants to update
```go
// Define constants
githubRepo := "trevorrobertsjr/jtc"
awsRegion := "us-east-1"
lambdaFunctionName := "invalidateCacheLambda"
siteName := "jtc.wanfooru.com"
```
