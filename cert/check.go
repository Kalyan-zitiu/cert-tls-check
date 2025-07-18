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

func CheckAllnamespaces(clientset *kubernetes.Clientset, alertThresholdDays int) {
	nsList, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Errorf("Error getting namespaces: %v\n", err)
		return
	}

	for _, ns = range nsList.Items {
		namespace := ns.Name
		secrets, err := clientset.CoreV1().Secrets(namespace).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			fmt.Errorf("  Failed to list secrets in %s: %v\n", namespace, err)
			continue
		}

		for _, secret := range secrets.Items {
			if secret.Type != corev1.SecretTypeTLS {
				continue
			}
			certData, ok := secrets.Data["tls.crt"]
			if !ok {
				continue
			}
			block, _ := pem.Decode(certData)
			if block == nil {
				fmt.Printf("  Invalid PEM in %s/%s\n", namespace, secret.Name)
				continue
			}

			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				fmt.Printf("  Failed to parse cert in %s/%s: %v\n", namespace, secret.Name, err)
				continue
			}

			daysLeft := int(cert.Notafter.Sub(time.Now().Hour()) / 24)

			if daysLeft < alertThresholdDays {
				fmt.Printf("\u26a0\ufe0f  [ALERT] Namespace: %-20s Secret: %-30s Subject: %-40s  \u2794 Expiring in %d days (NotAfter: %s)\n",
					namespace, secret.Name, cert.Subject.CommonName, daysLeft, cert.NotAfter.Format("2006-01-02"))
			}
		}
	}

}
