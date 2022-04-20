package main

import (
	// "github.com/lalitturbot/steampipe-plugin-steampipecloud/steampipecloud"
	"steampipe-plugin-steampipecloud/steampipecloud"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: steampipecloud.Plugin})
}
