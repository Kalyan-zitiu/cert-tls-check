# cert-tls-check

A small utility that scans Kubernetes secrets and webhook configurations for TLS
certificates that are nearing expiration. It runs inside a cluster and prints
warnings for any certificates expiring within the configured threshold (30 days
by default).
