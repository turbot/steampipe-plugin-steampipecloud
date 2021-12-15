package steampipecloud

import (
	"context"

	openapi "github.com/turbot/steampipe-cloud-sdk-go"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func getUserIdentity(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cacheKey := "GetUserIdentity"

	// get the service connection for the service
	svc, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("GetUserIdentity", "connection_error", err)
		return nil, err
	}

	// if found in cache, return the result
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(openapi.User), nil
	}

	var resp openapi.User

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		resp, _, err = svc.Actors.Get(ctx).Execute()
		return resp, err
	}

	response, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	user := response.(openapi.User)

	if err != nil {
		plugin.Logger(ctx).Error("GetUserIdentity", "error", err)
		return nil, err
	}

	// save to extension cache
	d.ConnectionManager.Cache.Set(cacheKey, user)

	return user, nil
}
