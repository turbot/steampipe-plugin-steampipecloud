# Table: steampipecloud_workspace_snapshot

A Steampipe snapshot is a point in time view of a benchmark. It can be shared across workspaces or made public.

**Important notes:**

This table supports optional quals. Queries with optional quals in `where` clause are optimised to use steampipecloud filters.

Optional quals are supported for the following columns:

- `query_where` - Allows use of explicit query where clause.
- `identity_id`
- `workspace_id`
- `state`
- `visibility`
- `dashboard_name`
- `schema_version`
- `start_time`
- `end_time`

## Examples

### Basic information about snapshots across all workspaces

```sql
select
  id,
  state,
  visibility,
  dashboard_name,
  schema_version,
  start_time,
  end_time
from
  steampipecloud_workspace_snapshot;
```

### List snapshots of a workspace with handle

```sql
select
  id,
  state,
  visibility,
  dashboard_name,
  schema_version,
  start_time,
  end_time
from
  steampipecloud_workspace_snapshot
where
  workspace_handle = 'dev';
```

### List public snapshots for a dashboard across all workspaces

```sql
select
  id,
  state,
  visibility,
  dashboard_name,
  schema_version,
  start_time,
  end_time
from
  steampipecloud_workspace_snapshot
where
  dashboard_name = 'aws_tags.benchmark.limit'
  and visibility = 'anyone_with_link';
```

### List snapshots for a dashboard executed in the last 7 days

```sql
select
  id,
  state,
  visibility,
  dashboard_name,
  schema_version,
  start_time,
  end_time
from
  steampipecloud_workspace_snapshot
where
  dashboard_name = 'aws_tags.benchmark.limit'
  and start_time >= now() - interval '7 days'
  and end_time is not null;
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
  and id = 'snap_cbocshgsmpv0mhlfm12g_7bi1qkns085n333lclfdi26o';
```

### List snapshots for a dashboard executed in the last 7 days using an explicit query_where

```sql
select
  id,
  state,
  visibility,
  dashboard_name,
  schema_version,
  start_time,
  end_time
from
  steampipecloud_workspace_snapshot
where
  query_where = 'dashboard_name = ''aws_tags.benchmark.limit'' and start_time >= now() - interval ''7 days''';
```