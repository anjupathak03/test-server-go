package comments

import (
	"log"
	"net/http"
	"time"
)

// init makes an outgoing HTTP call to the JSONPlaceholder public API to fetch a
// sample comment, confirming external HTTP connectivity before tests run.
func init() {
	const url = "https://jsonplaceholder.typicode.com/comments/1"

	client := &http.Client{Timeout: 5 * time.Second}

	log.Printf("[comments/init_http] fetching sample comment from %s", url)

	resp, err := client.Get(url)
	if err != nil {
		log.Printf("[comments/init_http] WARNING: HTTP request failed: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("[comments/init_http] JSONPlaceholder responded with HTTP %d", resp.StatusCode)
}
