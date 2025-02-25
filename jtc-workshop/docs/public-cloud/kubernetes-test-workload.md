---
sidebar_position: 4
---

# Kubernetes Workloads

Let's verify that 

Deploy an nginx web server and a busybox pod for network troubleshooting

```bash
kubectl get pods
```
```bash
NAME                              READY   STATUS    RESTARTS   AGE
nginx-deployment-96b9d695-ssw4m   1/1     Running   0          7s
```

Copy your pod name to another window. We will be using it for the next few commands...

```bash
kubectl describe pod your-pod-name-here
```

```bash
Name:             nginx-deployment-96b9d695-ssw4m
Namespace:        default
Priority:         0
Service Account:  default
Node:             minikube/192.168.49.2
Start Time:       Tue, 25 Feb 2025 13:25:04 +0000
Labels:           app=nginx
                  pod-template-hash=96b9d695
Annotations:      <none>
Status:           Running
IP:               10.244.0.5
IPs:
  IP:           10.244.0.5
Controlled By:  ReplicaSet/nginx-deployment-96b9d695
Containers:
  nginx:
    Container ID:   docker://9be8fd9f9105daedfa4fba6871431a57003d3860d5df23c6dd987d8900893cc1
    Image:          nginx:latest
    Image ID:       docker-pullable://nginx@sha256:9d6b58feebd2dbd3c56ab5853333d627cc6e281011cfd6050fa4bcf2072c9496
    Port:           80/TCP
    Host Port:      0/TCP
    State:          Running
      Started:      Tue, 25 Feb 2025 13:25:06 +0000
    Ready:          True
    Restart Count:  0
    Environment:    <none>
    Mounts:
      /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-wpwsx (ro)
Conditions:
  Type                        Status
  PodReadyToStartContainers   True 
  Initialized                 True 
  Ready                       True 
  ContainersReady             True 
  PodScheduled                True 
Volumes:
  kube-api-access-wpwsx:
    Type:                    Projected (a volume that contains injected data from multiple sources)
    TokenExpirationSeconds:  3607
    ConfigMapName:           kube-root-ca.crt
    ConfigMapOptional:       <nil>
    DownwardAPI:             true
QoS Class:                   BestEffort
Node-Selectors:              <none>
Tolerations:                 node.kubernetes.io/not-ready:NoExecute op=Exists for 300s
                             node.kubernetes.io/unreachable:NoExecute op=Exists for 300s
Events:
  Type    Reason     Age   From               Message
  ----    ------     ----  ----               -------
  Normal  Scheduled  44s   default-scheduler  Successfully assigned default/nginx-deployment-96b9d695-ssw4m to minikube
  Normal  Pulling    43s   kubelet            Pulling image "nginx:latest"
  Normal  Pulled     43s   kubelet            Successfully pulled image "nginx:latest" in 225ms (225ms including waiting). Image size: 191998640 bytes.
  Normal  Created    43s   kubelet            Created container: nginx
  Normal  Started    42s   kubelet            Started container nginx
```
```bash
kubectl logs your-pod-name
```

```bash
/docker-entrypoint.sh: /docker-entrypoint.d/ is not empty, will attempt to perform configuration
/docker-entrypoint.sh: Looking for shell scripts in /docker-entrypoint.d/
/docker-entrypoint.sh: Launching /docker-entrypoint.d/10-listen-on-ipv6-by-default.sh
10-listen-on-ipv6-by-default.sh: info: Getting the checksum of /etc/nginx/conf.d/default.conf
10-listen-on-ipv6-by-default.sh: info: Enabled listen on IPv6 in /etc/nginx/conf.d/default.conf
/docker-entrypoint.sh: Sourcing /docker-entrypoint.d/15-local-resolvers.envsh
/docker-entrypoint.sh: Launching /docker-entrypoint.d/20-envsubst-on-templates.sh
/docker-entrypoint.sh: Launching /docker-entrypoint.d/30-tune-worker-processes.sh
/docker-entrypoint.sh: Configuration complete; ready for start up
2025/02/25 13:25:06 [notice] 1#1: using the "epoll" event method
2025/02/25 13:25:06 [notice] 1#1: nginx/1.27.4
2025/02/25 13:25:06 [notice] 1#1: built by gcc 12.2.0 (Debian 12.2.0-14) 
2025/02/25 13:25:06 [notice] 1#1: OS: Linux 6.8.0-1021-aws
2025/02/25 13:25:06 [notice] 1#1: getrlimit(RLIMIT_NOFILE): 1048576:1048576
2025/02/25 13:25:06 [notice] 1#1: start worker processes
2025/02/25 13:25:06 [notice] 1#1: start worker process 29
2025/02/25 13:25:06 [notice] 1#1: start worker process 30
```

Let's try accessing our web server
```bash
kubectl port-forward svc/nginx-service 8080:80 &
```
```bash
Forwarding from 127.0.0.1:8080 -> 80
Forwarding from [::1]:8080 -> 80
```

```bash
curl localhost:8080
```

```bash
Handling connection for 8080
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
html { color-scheme: light dark; }
body { width: 35em; margin: 0 auto;
font-family: Tahoma, Verdana, Arial, sans-serif; }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>
```

Now that we verified we can deploy a workload on Kubernetes, let's deploy our pipeline as a Kubernetes pod.
