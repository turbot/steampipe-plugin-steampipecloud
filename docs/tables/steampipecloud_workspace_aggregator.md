# Table: steampipecloud_workspace_aggregator

Aggregators allow users to aggregate or search for data across multiple connections in a workspace.

## Examples

### Basic info

```sql
select
  id,
  handle,
  identity_handle,
  workspace_handle,
  type,
  plugin,
  connections
from
  steampipecloud_workspace_aggregator;
```

### List aggregators for a specific workspace

```sql
select
  id,
  handle,
  identity_handle,
  workspace_handle,
  type,
  plugin,
  connections
from
  steampipecloud_workspace_aggregator
where
  workspace_handle = 'dev';
```

### List aggregators of plugin type `aws` for a specific workspace

```sql
select
  id,
  handle,
  identity_handle,
  workspace_handle,
  type,
  plugin,
  connections
from
  steampipecloud_workspace_aggregator
where
  workspace_handle = 'dev'
  and plugin = 'aws';
```

### List aggregators created in the last 7 days for a specific workspace

```sql
select
  id,
  handle,
  identity_handle,
  workspace_handle,
  type,
  plugin,
  connections
from
  steampipecloud_workspace_aggregator
where
  workspace_handle = 'dev'
  and created_at >= now() - interval '7 days';
```
