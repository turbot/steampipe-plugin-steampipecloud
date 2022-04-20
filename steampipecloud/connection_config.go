package steampipecloud

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	openapiclient "github.com/turbot/steampipe-cloud-sdk-go"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/schema"
)

type steampipecloudConfig struct {
	Token *string `cty:"token"`
	Host  *string `cty:"host"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"token": {
		Type: schema.TypeString,
	},
	"host": {
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
	steampipecloudConfig := GetConfig(d.Connection)

	token := os.Getenv("STEAMPIPE_CLOUD_TOKEN")
	if steampipecloudConfig.Token != nil {
		token = *steampipecloudConfig.Token
	}
	if token == "" {
		return nil, errors.New("'token' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}

	configuration := openapiclient.NewConfiguration()
	configuration.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", token))

	host := os.Getenv("STEAMPIPE_CLOUD_HOST")
	if steampipecloudConfig.Host != nil {
		host = *steampipecloudConfig.Host
	}

	if host != "" && !strings.Contains(host, "cloud.steampipe.io") {
		configuration.Servers = []openapiclient.ServerConfiguration{
			{
				URL:         fmt.Sprintf("%s/api/v0", host),
				Description: "Local API",
			},
		}
	}

	apiClient := openapiclient.NewAPIClient(configuration)

	return apiClient, nil
}
