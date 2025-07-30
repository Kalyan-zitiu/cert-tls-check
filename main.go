package main

import (
	"flag"
	"log"
	"time"
	"time-tls-checker/cert"
	"time-tls-checker/client"
	"time-tls-checker/metrics"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "", "Path to kubeconfig for local debugging")
	alertThreshold := flag.Int("alert-threshold", 30, "Days before expiry to warn")
	metricsAddr := flag.String("metrics-addr", ":8080", "Address to serve Prometheus metrics")
	flag.Parse()

	// 初始化 Kubernetes 客户端
	clientset, err := client.InitK8SClient(*kubeconfig)
	if err != nil {
		log.Fatalf("Failed to init Kubernetes client: %v", err)
	}

	// start metrics server
	metrics.StartServer(*metricsAddr)

	// 定时检查逻辑
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	// 立即执行一次
	expiring := cert.CheckAllNamespaces(clientset, *alertThreshold)
	metrics.SetExpiringCerts(expiring)
	webhookExpiring := cert.CheckMutatingWebhookCABundles(clientset, *alertThreshold)
	metrics.SetExpiringWebhookCAs(webhookExpiring)

	// 循环检查
	for {
		select {
		case <-ticker.C:
			expiring := cert.CheckAllNamespaces(clientset, *alertThreshold)
			metrics.SetExpiringCerts(expiring)
			webhookExpiring := cert.CheckMutatingWebhookCABundles(clientset, *alertThreshold)
			metrics.SetExpiringWebhookCAs(webhookExpiring)
		}
	}
}
