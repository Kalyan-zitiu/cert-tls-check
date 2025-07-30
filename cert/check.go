package cert

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

func CheckAllNamespaces(clientset *kubernetes.Clientset, alertThresholdDays int) int {
	expiring := 0
	nsList, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error getting namespaces: %v\n", err)
		return 0
	}

	for _, ns := range nsList.Items {
		namespace := ns.Name

		secrets, err := clientset.CoreV1().Secrets(namespace).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			fmt.Printf("  Failed to list secrets in %s: %v\n", namespace, err)
			continue
		}

		for _, secret := range secrets.Items {
			if secret.Type != corev1.SecretTypeTLS {
				continue
			}
			certData, ok := secret.Data["tls.crt"]
			if !ok {
				continue
			}

			for {
				var block *pem.Block
				block, certData = pem.Decode(certData)
				if block == nil {
					break
				}

				cert, err := x509.ParseCertificate(block.Bytes)
				if err != nil {
					fmt.Printf("  Failed to parse cert in %s/%s: %v\n", namespace, secret.Name, err)
					continue
				}

				daysLeft := int(cert.NotAfter.Sub(time.Now()).Hours() / 24)
				if daysLeft < alertThresholdDays {
					expiring++
					fmt.Printf("\u26A0\uFE0F  [List] Namespace: %-20s Secret: %-30s Subject: %-40s  \u2794 Expiring in %d days (NotAfter: %s)\n",
						namespace, secret.Name, cert.Subject.CommonName, daysLeft, cert.NotAfter.Format("2006-01-02"))
				}
				if daysLeft < alertThresholdDays {
					fmt.Printf("\u26A0\uFE0F  [ALERT] Namespace: %-20s Secret: %-30s Subject: %-40s  \u2794 Expiring in %d days (NotAfter: %s)\n",
						namespace, secret.Name, cert.Subject.CommonName, daysLeft, cert.NotAfter.Format("2006-01-02"))
				}
			}
		}
	}

	return expiring
}

// CheckMutatingWebhookCABundles checks the CA bundles of all MutatingWebhookConfigurations
// and warns if a certificate is expiring within alertThresholdDays.
func CheckMutatingWebhookCABundles(clientset *kubernetes.Clientset, alertThresholdDays int) int {
	expiring := 0
	configs, err := clientset.AdmissionregistrationV1().MutatingWebhookConfigurations().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing MutatingWebhookConfigurations: %v\n", err)
		return 0
	}

	for _, cfg := range configs.Items {
		for _, hook := range cfg.Webhooks {
			certData := hook.ClientConfig.CABundle
			if len(certData) == 0 {
				continue
			}

			remaining := certData
			for {
				var block *pem.Block
				block, remaining = pem.Decode(remaining)
				if block == nil {
					break
				}
				cert, err := x509.ParseCertificate(block.Bytes)
				if err != nil {
					fmt.Printf("  Failed to parse CA bundle in %s/%s: %v\n", cfg.Name, hook.Name, err)
					break
				}

				daysLeft := int(cert.NotAfter.Sub(time.Now()).Hours() / 24)
				if daysLeft >= alertThresholdDays {
					fmt.Printf("\u26A0\uFE0F  [List] Webhook: %-30s Hook: %-20s  \u2794 CA expiring in %d days (NotAfter: %s)\n",
						cfg.Name, hook.Name, daysLeft, cert.NotAfter.Format("2006-01-02"))
				}
				if daysLeft < alertThresholdDays {
					expiring++
					fmt.Printf("\u26A0\uFE0F  [ALERT] Webhook: %-30s Hook: %-20s  \u2794 CA expiring in %d days (NotAfter: %s)\n",
						cfg.Name, hook.Name, daysLeft, cert.NotAfter.Format("2006-01-02"))
				}
			}
		}
	}

	return expiring
}
