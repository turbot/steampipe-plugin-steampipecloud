package steampipecloud

import (
	"context"
	"net/http"
	"time"

	"github.com/sethvargo/go-retry"
	openapi "github.com/turbot/steampipe-cloud-sdk-go"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func getUserIdentity(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	cacheKey := "GetUserIdentity"

	// if found in cache, return the result
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(openapi.TypesUser), nil
	}

	// get the service connection for the service
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("GetUserIdentity", "connection_error", err)
		return nil, err
	}

	var user openapi.TypesUser
	var httpResp *http.Response
	b, err := retry.NewFibonacci(100 * time.Millisecond)
	err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
		user, httpResp, err = svc.UsersApi.GetActor(ctx).Execute()
		// 429 too many request
		if httpResp.StatusCode == 429 {
			return retry.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		plugin.Logger(ctx).Error("GetUserIdentity", "error", err)
		// 404 Not Found
		if httpResp.StatusCode == 404 {
			return nil, nil
		}
		return nil, err
	}

	// save to extension cache
	d.ConnectionManager.Cache.Set(cacheKey, user)
	return user, nil
}
