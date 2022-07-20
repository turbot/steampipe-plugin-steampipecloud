## v0.3.0 [2022-07-20]

_What's new?_

- New table added
  - [steampipecloud_organization_workspace_member](https://hub.steampipe.io/plugins/turbot/steampipecloud/tables/steampipecloud_organization_workspace_member) ([#19](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/19))

_Enhancements_

- Added the following new columns to `steampipecloud_workspace` table: ([#19](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/19))
  - created_by
  - updated_by
  - version_id
- Added the following new columns to `steampipecloud_organization_member` table: ([#19](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/19))
  - created_by
  - updated_by
  - org_handle
- Added Get method to `steampipecloud_organization_member` table for getting a particular org member ([#19](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/19))
- Added the following new columns to `steampipecloud_organization` table: ([#19](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/19))
  - created_by
  - updated_by
- Added Get method to `steampipecloud_organization` table for getting a particular organization ([#19](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/19))

_Dependencies_

- Upgraded [steampipe-cloud-sdk-go](https://github.com/turbot/steampipe-cloud-sdk-go/blob/main/CHANGELOG.md#012-2022-07-19) to v0.1.2 ([#19](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/19))

## v0.2.0 [2022-06-09]

_What's new?_

- New tables added
  - [steampipecloud_workspace_mod](https://hub.steampipe.io/plugins/turbot/steampipecloud/tables/steampipecloud_workspace_mod) ([#14](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/14))
  - [steampipecloud_workspace_mod_variable](https://hub.steampipe.io/plugins/turbot/steampipecloud/tables/steampipecloud_workspace_mod_variable) ([#14](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/14))

## v0.1.0 [2022-04-28]

_Enhancements_

- Added support for native Linux ARM and Mac M1 builds. ([#12](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/12))
- Recompiled plugin with [steampipe-plugin-sdk v3.1.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v310--2022-03-30) and Go version `1.18`. ([#11](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/11))

## v0.0.1 [2021-12-16]

_What's new?_

- New tables added
  - [steampipecloud_audit_log](https://hub.steampipe.io/plugins/turbot/steampipecloud/tables/steampipecloud_audit_log)
  - [steampipecloud_connection](https://hub.steampipe.io/plugins/turbot/steampipecloud/tables/steampipecloud_connection)
  - [steampipecloud_organization](https://hub.steampipe.io/plugins/turbot/steampipecloud/tables/steampipecloud_organization)
  - [steampipecloud_organization_member](https://hub.steampipe.io/plugins/turbot/steampipecloud/tables/steampipecloud_organization_member)
  - [steampipecloud_token](https://hub.steampipe.io/plugins/turbot/steampipecloud/tables/steampipecloud_token)
  - [steampipecloud_user](https://hub.steampipe.io/plugins/turbot/steampipecloud/tables/steampipecloud_user)
  - [steampipecloud_workspace](https://hub.steampipe.io/plugins/turbot/steampipecloud/tables/steampipecloud_workspace)
  - [steampipecloud_workspace_connection](https://hub.steampipe.io/plugins/turbot/steampipecloud/tables/steampipecloud_connection)
