---
organization: Turbot
category: ["software development"]
icon_url: "/images/plugins/turbot/steampipecloud.svg"
brand_color: "#008bcf"
display_name: "SteampipeCloud"
name: "steampipecloud"
description: "Steampipe plugin for querying steampipecloud workspaces, connections and other resources."
---

# SteampipeCloud

Query your SteampipeCloud infrastructure including workspaces, connections, members and more.

## Installation

Download and install the latest SteampipeCloud plugin:

```bash
$ steampipe plugin install steampipecloud
Installing plugin steampipecloud...
$
```

## Connection Configuration

Connection configurations are defined using HCL in one or more Steampipe config files. Steampipe will load ALL configuration files from `~/.steampipe/config` that have a `.spc` extension. A config file may contain multiple connections.

Installing the latest steampipecloud plugin will create a connection file (`~/.steampipe/config/steampipecloud.spc`) with a single connection named `steampipecloud`. You must modify this connection to include your Token for SteampipeCloud account.

```hcl
connection "steampipecloud" {
  plugin  = "steampipecloud"
  token   = "17ImlCYdfZ3WJIrGk96gCpJn1fi1pLwVdrb23kj4"
}
```

### Configuration Arguments

The SteampipeCloud plugin allows you set static credentials with the `token` argument. You can use them to authenticate to the API by including it in a bearer-type Authorization header with your request. 

To use the plugin, you'll first need to `create token` in your SteampipeCloud console.

If the `token` argument is not specified for a connection, the plugin will look for the `STEAMPIPE_CLOUD_TOKEN` environment variable.
