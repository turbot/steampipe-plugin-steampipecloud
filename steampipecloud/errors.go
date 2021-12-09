package steampipecloud

import (
	"log"
	"strings"
)

func shouldRetryError(err error) bool {
	if strings.Contains(err.Error(), "429") {
		log.Printf("[WARN] Received Rate Limit Error")
		return true
	}
	return false
}
