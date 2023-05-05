package steampipecloud

import (
	"context"
	"log"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func shouldIgnoreErrors(notFoundErrors []string) plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		for _, pattern := range notFoundErrors {
			// handle not found error
			if strings.Contains(err.Error(), pattern) {
				return true
			}
		}
		return false
	}
}

func shouldRetryError(err error) bool {
	if strings.Contains(err.Error(), "429") {
		log.Printf("[WARN] Received Rate Limit Error")
		return true
	}
	return false
}
