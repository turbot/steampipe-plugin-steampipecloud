## v1.0.0 [2022-08-12]

_What's new?_

- New tables added
  - [steampipecloud_workspace_snapshot](https://hub.steampipe.io/plugins/turbot/steampipecloud/tables/steampipecloud_workspace_snapshot) ([#22](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/22))
  - [steampipecloud_workspace_db_log](https://hub.steampipe.io/plugins/turbot/steampipecloud/tables/steampipecloud_workspace_db_log) ([#23](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/23))

_Enhancements_

- Added `process_id` column to `steampipecloud_audit_log` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Added `created_by_id`, `created_by`, `updated_by_id`, `updated_by` columns to `steampipecloud_connection` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Added `created_by_id`, `updated_by_id` columns to `steampipecloud_organization` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Modified columns `created_by`, `updated_by` to store additional information about the user who created or updated the organization in `steampipecloud_organization` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Added `user`, `scope`, `created_by_id`, `updated_by_id` columns to `steampipecloud_organization_member` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Modified columns `created_by`, `updated_by` to store additional information about the user who created or updated the organization member in `steampipecloud_organization_member` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Removed column `email` from `steampipecloud_organization_member` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Added `user`, `created_by_id`, `updated_by_id` columns to `steampipecloud_organization_workspace_member` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Modified columns `created_by`, `updated_by` to store additional information about the user who created or updated the organization workspace member in `steampipecloud_organization_workspace_member` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Removed column `email` from `steampipecloud_organization_workspace_member` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Removed column `email` from `steampipecloud_user` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Added `host`, `api_version`, `cli_version`, `created_by_id`, `updated_by_id` columns to `steampipecloud_workspace` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Modified column `workspace_state` renamed to `state` in `steampipecloud_workspace` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Modified columns `created_by`, `updated_by` to store additional information about the user who created or updated the organization workspace member in `steampipecloud_workspace` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Added `identity_handle`, `identity_type`, `workspace_handle`, `connection_handle`, `connection`, `created_by_id`, `created_by`, `updated_by_id`, `updated_by` columns to `steampipecloud_workspace_connection` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Added `identity_handle`, `workspace_handle`, `created_by_id`, `created_by`, `updated_by_id`, `updated_by`, `version_id` columns to `steampipecloud_workspace_mod` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Modified column `mod_constraint` renamed to `constraint` in `steampipecloud_workspace_mod` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Added `version_id` column to `steampipecloud_workspace_mod_variable` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Modified column `created_by_handle` renamed to `created_by_id` to store identifier of the person who created the setting in `steampipecloud_workspace_mod_variable` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))
- Modified column `updated_by_handle` renamed to `updated_by_id` to store identifier of the person who last updated the setting in `steampipecloud_workspace_mod_variable` table. ([#24](https://github.com/turbot/steampipe-plugin-steampipecloud/pull/24))

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
