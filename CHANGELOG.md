## v0.10.0 [2023-05-09]

_What's new?_

- New tables added
  - [steampipecloud_workspace_aggregator](https://hub.steampipe.io/plugins/turbot/steampipecloud/tables/steampipecloud_workspace_aggregator) ([#35](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/35))

## v0.9.0 [2023-04-10]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.3.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v530-2023-03-16) which includes fixes for query cache pending item mechanism and aggregator connections not working for dynamic tables. ([#33](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/33))

## v0.8.0 [2023-02-22]

_What's new?_

- New tables added
  - [steampipecloud_process](https://hub.steampipe.io/plugins/turbot/steampipecloud/tables/steampipecloud_process) ([#31](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/31))
  - [steampipecloud_workspace_pipeline](https://hub.steampipe.io/plugins/turbot/steampipecloud/tables/steampipecloud_workspace_pipeline) ([#31](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/31))
  - [steampipecloud_workspace_process](https://hub.steampipe.io/plugins/turbot/steampipecloud/tables/steampipecloud_workspace_process) ([#31](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/31))

## v0.7.0 [2022-12-27]

_Enhancements_

- Added column `expires_at` to `steampipecloud_workspace_snapshot` table. ([#29](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/29))

## v0.6.0 [2022-11-24]

_What's new?_

- New tables added
  - [steampipecloud_user_email](https://hub.steampipe.io/plugins/turbot/steampipecloud/tables/steampipecloud_user_email) ([#27](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/27))
  - [steampipecloud_user_preferences](https://hub.steampipe.io/plugins/turbot/steampipecloud/tables/steampipecloud_user_preferences) ([#27](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/27))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.8](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v418-2022-09-08) which increases the default open file limit. ([#28](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/28))

## v0.5.0 [2022-09-27]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.7](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v417-2022-09-08) which includes several caching and memory management improvements. ([#25](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/25))
- Recompiled plugin with Go version `1.19`. ([#25](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/25))

## v0.4.0 [2022-08-22]

_Breaking changes_

- Removed column `email` from `steampipecloud_organization_member` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Removed column `email` from `steampipecloud_organization_workspace_member` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Removed column `email` from `steampipecloud_user` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Renamed column `workspace_state` to `state` in `steampipecloud_workspace` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Renamed column `mod_constraint` to `constraint` in `steampipecloud_workspace_mod` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Renamed column `created_by_handle` to `created_by_id` to store the identifier of the person who created the setting for the variable in `steampipecloud_workspace_mod_variable` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Renamed column `updated_by_handle` to `updated_by_id` to store the identifier of the person who last updated the setting for the variable in `steampipecloud_workspace_mod_variable` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))

_What's new?_

- New tables added
  - [steampipecloud_workspace_db_log](https://hub.steampipe.io/plugins/turbot/steampipecloud/tables/steampipecloud_workspace_db_log) ([#23](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/23))
  - [steampipecloud_workspace_snapshot](https://hub.steampipe.io/plugins/turbot/steampipecloud/tables/steampipecloud_workspace_snapshot) ([#22](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/22))

_Enhancements_

- Added `process_id` column to `steampipecloud_audit_log` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Added `created_by_id`, `created_by`, `updated_by_id`, `updated_by` columns to `steampipecloud_connection` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Added `created_by_id`, `updated_by_id` columns to `steampipecloud_organization` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Updated columns `created_by`, `updated_by` to store additional information about the user who created or updated the organization in `steampipecloud_organization` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Added `created_by_id`, `scope`, `updated_by_id`, `user` columns to `steampipecloud_organization_member` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Updated columns `created_by`, `updated_by` to store additional information about the user who created or updated the organization member in `steampipecloud_organization_member` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Added `created_by_id`, `updated_by_id`, `user` columns to `steampipecloud_organization_workspace_member` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Updated columns `created_by`, `updated_by` to store additional information about the user who created or updated the organization workspace member in `steampipecloud_organization_workspace_member` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Added `api_version`, `cli_version`, `created_by_id`, `host`, `updated_by_id` columns to `steampipecloud_workspace` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Updated columns `created_by`, `updated_by` to store additional information about the user who created or updated the organization workspace member in `steampipecloud_workspace` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Added `connection_handle`, `connection`, `created_by_id`, `created_by`, `identity_handle`, `identity_type`, `updated_by_id`, `updated_by`, `workspace_handle` columns to `steampipecloud_workspace_connection` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Added `created_by_id`, `created_by`, `identity_handle`, `updated_by_id`, `updated_by`, `version_id`, `workspace_handle` columns to `steampipecloud_workspace_mod` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Added `version_id` column to `steampipecloud_workspace_mod_variable` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))

_Dependencies_

- Recompiled plugin with [steampipe-cloud-sdk-go v0.1.3](https://github.com/turbot/steampipe-cloud-sdk-go/blob/main/CHANGELOG.md#013-2022-08-12). ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))

## v0.3.0 [2022-07-20]

_What's new?_

- New tables added
  - [steampipecloud_organization_workspace_member](https://hub.steampipe.io/plugins/turbot/steampipecloud/tables/steampipecloud_organization_workspace_member) ([#19](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/19))

_Enhancements_

- Added `created_by` and `updated_by` columns to the `steampipecloud_organization` table. ([#19](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/19))
- Added `created_by`, `org_handle`, and `updated_by` columns to the `steampipecloud_organization_member` table. ([#19](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/19))
- Added `created_by`, `updated_by`, and `version_id` columns to the `steampipecloud_workspace` table. ([#19](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/19))
- Added `GetConfig` to `steampipecloud_organization` and `steampipecloud_organization_member` tables. ([#19](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/19))

_Dependencies_

- Recompiled plugin with [steampipe-cloud-sdk-go v0.1.2](https://github.com/turbot/steampipe-cloud-sdk-go/blob/main/CHANGELOG.md#012-2022-07-19). ([#19](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/19))

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
