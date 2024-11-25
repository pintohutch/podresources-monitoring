# podresources-monitoring

A basic go program to poll the local [podresources API](https://kubernetes.io/blog/2023/08/23/kubelet-podresources-api-ga/)
and log device info for any assigned containers by the kubelet.

## Deploy
Simply apply the DaemonSet manifest to ensure an instance is running on every node.
```bash
kubectl apply -f daemonset.yaml
```

Then view the logs to debug assigned devices by container.
```bash
kubectl logs -f ds/podresources-monitoring
```