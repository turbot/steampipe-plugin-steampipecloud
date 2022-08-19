# Table: steampipecloud_workspace

Workspaces provide a bounded context for managing, operating, and securing Steampipe resources. A workspace comprises a single Steampipe database instance as well as a directory of mod resources such as queries, benchmarks, and controls. Workspaces allow you to separate your Steampipe instances for security, operational, or organizational purposes.

## Examples

### Basic info

```sql
select
  id,
  state,
  handle,
  identity_handle
from
  steampipecloud_workspace;
```

### List user workspaces

```sql
select
  id,
  state,
  handle,
  identity_handle
from
  steampipecloud_workspace
where
  identity_type = 'user';
```

### List organization workspaces

```sql
select
  id,
  state,
  handle,
  identity_handle
from
  steampipecloud_workspace
where
  identity_type = 'org';
```

### List workspaces which are not running

```sql
select
  id,
  state,
  handle,
  identity_handle
from
  steampipecloud_workspace
where
  state <> 'running';
```
