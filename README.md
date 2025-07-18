# cert-tls-check

A small utility that scans Kubernetes TLS secrets and webhook configurations for
certificates and reports how many days remain until they expire. By default it
runs inside the cluster but you can provide a kubeconfig for local debugging.

## Usage

```bash
go run main.go --kubeconfig=/path/to/kubeconfig
```

The tool prints `INFO` lines for certificates that are still valid, `ALERT` when
a certificate is within the threshold (30 days by default), and `EXPIRED` for
any that are already expired.
