# cert-tls-check

A small utility that scans Kubernetes secrets and webhook configurations for TLS
certificates that are nearing expiration. It runs inside a cluster and prints
warnings for any certificates expiring within the configured threshold (30 days
by default).

When running outside of a cluster you can provide a kubeconfig file for local
debugging:

```bash
go run main.go --kubeconfig=/path/to/config
```
