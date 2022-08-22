# Table: steampipecloud_workspace_snapshot

A Steampipe snapshot is a point in time view of a benchmark. It can be shared across workspaces or made public.

**Important notes:**

This table supports optional quals. Queries with optional quals in the `where` clause are optimised to use Steampipe Cloud filters.

Optional quals are supported for the following columns:

- `created_at`
- `dashboard_name`
- `dashboard_title`
- `id`
- `inputs`
- `query_where` - Allows use of [query filters](https://steampipe.io/docs/cloud/reference/query-filter). For a list of supported columns for snapshots, please see [Supported APIs and Columns](https://steampipe.io/docs/cloud/reference/query-filter#supported-apis--columns). Please note that any query filter passed into the `query_where` qual will be combined with other optional quals.
- `tags`
- `visibility`

## Examples

### Basic info

```sql
select
  id,
  identity_handle,
  workspace_handle,
  state,
  visibility,
  dashboard_name,
  dashboard_title,
  schema_version,
  tags
from
  steampipecloud_workspace_snapshot;
```

### List snapshots for a specific workspace

```sql
select
  id,
  identity_handle,
  workspace_handle,
  state,
  visibility,
  dashboard_name,
  dashboard_title,
  schema_version,
  tags
from
  steampipecloud_workspace_snapshot
where
  workspace_handle = 'dev';
```

### List public snapshots for the AWS Tags Limit benchmark dashboard across all workspaces

```sql
select
  id,
  identity_handle,
  workspace_handle,
  state,
  visibility,
  dashboard_name,
  dashboard_title,
  schema_version,
  tags
from
  steampipecloud_workspace_snapshot
where
  dashboard_name = 'aws_tags.benchmark.limit'
  and visibility = 'anyone_with_link';
```

### List snapshots for the AWS Compliance CIS v1.4.0 dashboard executed in the last 7 days

```sql
select
  id,
  identity_handle,
  workspace_handle,
  state,
  visibility,
  dashboard_name,
  dashboard_title,
  schema_version,
  tags
from
  steampipecloud_workspace_snapshot
where
  dashboard_name = 'aws_compliance.benchmark.cis_v140'
  and created_at >= now() - interval '7 days';
```

### Get the raw data for a particular snapshot

```sql
select
  data
from
  steampipecloud_workspace_snapshot
where
  identity_handle = 'myuser'
  and workspace_handle = 'dev'
  and id = 'snap_cc1ini7m1tujk0r0oqvg_12fie4ah78yl5rwadj7p6j63';
```

### List snapshots for the AWS Tags Limit benchmark dashboard executed in the last 7 days using [query filter](https://steampipe.io/docs/cloud/reference/query-filter)

```sql
select
  id,
  identity_handle,
  workspace_handle,
  state,
  visibility,
  dashboard_name,
  dashboard_title,
  schema_version,
  tags
from
  steampipecloud_workspace_snapshot
where
  query_where = 'dashboard_name = ''aws_tags.benchmark.limit'' and created_at >= now() - interval ''7 days''';
```
