---
organization: Turbot
category: ["saas"]
icon_url: "/images/plugins/turbot/steampipecloud.svg"
brand_color: "#a42a2d"
display_name: "Steampipe Cloud"
short_name: "steampipecloud"
description: "Steampipe plugin for querying workspaces, connections and more from Steampipe Cloud."
og_description: "Query Steampipe Cloud with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/steampipecloud-social-graphic.png"
---

# Steampipe Cloud + Steampipe [DEPRECATED]

This plugin has been deprecated as part of our [renaming](https://turbot.com/blog/2023/07/introducing-turbot-guardrails-and-pipes) of Steampipe Cloud to Turbot Pipes. Please use the [Turbot Pipes plugin](https://hub.steampipe.io/plugins/turbot/pipes) instead.

---
[Steampipe Cloud](https://cloud.steampipe.io/) is a fully managed SaaS platform for hosting Steampipe instances.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

For example:

```sql
select
  user_handle,
  email,
  status
from
  steampipecloud_organization_member
where
  status = 'accepted'
```

```
> select user_handle, email, status from steampipecloud_organization_member where status = 'accepted'
+-------------+------------------+----------+
| user_handle | email            | status   |
+-------------+------------------+----------+
| mario       | mario@turbot.com | accepted |
| yoshi       | yoshi@turbot.com | accepted |
+-------------+------------------+----------+
```

## Documentation

- **[Table definitions & examples →](/plugins/turbot/steampipecloud/tables)**

## Get started

### Install

Download and install the latest Steampipe Cloud plugin:

```bash
steampipe plugin install steampipecloud
```

### Configuration

Installing the latest Steampipe Cloud plugin will create a config file (`~/.steampipe/config/steampipecloud.spc`) with a single connection named `steampipecloud`:

```hcl
connection "steampipecloud" {
  plugin = "steampipecloud"

  # Steampipe Cloud API token. If `token` is not specified, it will be loaded
  # from the `STEAMPIPE_CLOUD_TOKEN` environment variable.
  # token = "spt_thisisnotarealtoken_123"

  # Steampipe Cloud host URL. This defaults to "https://cloud.steampipe.io/".
  # You only need to set this if connecting to a remote Steampipe Cloud database
  # not hosted in "https://cloud.steampipe.io/".
  # If `host` is not specified, it will be loaded from the `STEAMPIPE_CLOUD_HOST`
  # environment variable.
  # host = "https://cloud.steampipe.io"
}
```

- `token` (required) - [API tokens](https://steampipe.io/docs/cloud/profile#api-tokens) can be used to access the Steampipe Cloud API or to connect to Steampipe Cloud workspaces from the Steampipe CLI. May alternatively be set via the `STEAMPIPE_CLOUD_TOKEN` environment variable.
- `host` (optional) The Steampipe Cloud Host URL. This defaults to `https://cloud.steampipe.io/`. You only need to set this if you are connecting to a remote Steampipe Cloud database that is NOT hosted in `https://cloud.steampipe.io/`. This can also be set via the `STEAMPIPE_CLOUD_HOST` environment variable.

## Get Involved

- Open source: https://github.com/turbot/steampipe-plugin-steampipe-cloud
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
