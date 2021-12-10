package steampipecloud

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/schema"
	openapiclient "github.com/turbot/steampipecloud-sdk-go"
)

type steampipecloudConfig struct {
	Token *string `cty:"token"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"token": {
		Type: schema.TypeString,
	},
}

func ConfigInstance() interface{} {
	return &steampipecloudConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) steampipecloudConfig {
	if connection == nil || connection.Config == nil {
		return steampipecloudConfig{}
	}
	config, _ := connection.Config.(steampipecloudConfig)
	return config
}

func connect(_ context.Context, d *plugin.QueryData) (*openapiclient.APIClient, error) {
	configuration := openapiclient.NewConfiguration()
	steampipecloudConfig := GetConfig(d.Connection)

	if steampipecloudConfig.Token != nil {
		configuration.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", *steampipecloudConfig.Token))
	} else if os.Getenv("STEAMPIPE_CLOUD_TOKEN") != "" {
		configuration.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("STEAMPIPE_CLOUD_TOKEN")))
	} else {
		return nil, errors.New("'token' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}

	apiClient := openapiclient.NewAPIClient(configuration)

	return apiClient, nil
}
