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

# Steampipe Cloud + Steampipe

[Steampipe Cloud](https://cloud.steampipe.io/) is a fully managed SaaS platform for hosting Steampipe instances.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

For example:

```sql
select
  user_handle,
  email,
  status
from
  steampipecloud_org_member
where
  status = 'pending'
```

```
> select user_handle, email, status from steampipecloud_org_member
+-------------+------------------+----------+
| user_handle | email            | status   |
+-------------+------------------+----------+
| mario       | mario@turbot.com | pending  |
| yoshi       | yoshi@turbot.com | pending  |
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

  # Token for your Steampipe Cloud user
  token = "YOUR_STEAMPIPECLOUD_ACCESS_TOKEN"
}
```

- `token` (required) - [API tokens](https://steampipe.io/docs/cloud/profile#api-tokens) can be used to access the Steampipe Cloud API or to connect to Steampipe Cloud workspaces from the Steampipe CLI. May alternatively be set via the `STEAMPIPE_CLOUD_TOKEN` environment variable.

## Get Involved

- Open source: https://github.com/turbot/steampipe-plugin-steampipe-cloud
- Community: [Slack Channel](https://steampipe.io/community/join)
