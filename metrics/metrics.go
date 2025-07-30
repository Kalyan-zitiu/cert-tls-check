package metrics

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

var (
	expiringCerts      int64
	expiringWebhookCAs int64
)

// SetExpiringCerts sets the number of expiring TLS secrets.
func SetExpiringCerts(n int) {
	atomic.StoreInt64(&expiringCerts, int64(n))
}

// SetExpiringWebhookCAs sets the number of expiring webhook CA bundles.
func SetExpiringWebhookCAs(n int) {
	atomic.StoreInt64(&expiringWebhookCAs, int64(n))
}

// StartServer starts an HTTP server exposing metrics at the given address.
func StartServer(addr string) {
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		fmt.Fprintf(w, "tls_expiring_certs_total %d\n", atomic.LoadInt64(&expiringCerts))
		fmt.Fprintf(w, "tls_expiring_webhook_ca_total %d\n", atomic.LoadInt64(&expiringWebhookCAs))
	})
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			// errors only logged
			fmt.Printf("metrics server error: %v\n", err)
		}
	}()
}
