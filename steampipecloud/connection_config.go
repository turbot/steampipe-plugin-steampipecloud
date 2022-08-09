package steampipecloud

import (
	"context"
	"errors"
	"fmt"
	"net/url"
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

	if host == "" {
		return nil, errors.New("'host' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}

	if !strings.Contains(host, "cloud.steampipe.io") {
		parsedURL, parseErr := url.Parse(host)
		if parseErr != nil {
			return nil, errors.New(fmt.Sprintf(`invalid host: %v`, parseErr))
		}
		if parsedURL.Host == "" {
			return nil, errors.New(`missing protocol or host`)
		}

		// Parse and frame the Primary Servers
		var primaryServers []openapiclient.ServerConfiguration
		for _, server := range configuration.Servers {
			serverURL, parseErr := url.Parse(server.URL)
			if parseErr != nil {
				return nil, errors.New(fmt.Sprintf(`invalid host: %v`, parseErr))
			}
			primaryServers = append(primaryServers, openapiclient.ServerConfiguration{URL: fmt.Sprintf("%s://%s%s", serverURL.Scheme, parsedURL.Host, serverURL.Path), Description: "Local API"})
		}
		configuration.Servers = primaryServers

		// Parse and frame the Operation Servers
		operationServers := make(map[string]openapiclient.ServerConfigurations)
		for service, servers := range configuration.OperationServers {
			var serviceServers []openapiclient.ServerConfiguration
			for _, server := range servers {
				serverURL, parseErr := url.Parse(server.URL)
				if parseErr != nil {
					return nil, errors.New(fmt.Sprintf(`invalid host: %v`, parseErr))
				}
				serviceServers = append(serviceServers, openapiclient.ServerConfiguration{URL: fmt.Sprintf("%s://%s%s", serverURL.Scheme, parsedURL.Host, serverURL.Path), Description: "Local API"})
			}
			operationServers[service] = serviceServers
		}
		configuration.OperationServers = operationServers
	}

	apiClient := openapiclient.NewAPIClient(configuration)

	return apiClient, nil
}
