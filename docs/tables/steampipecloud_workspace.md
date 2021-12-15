# Table: steampipecloud_workspace

A Workspace provides a bounded context for managing, operating, and securing Steampipe resources. A workspace comprises a single Steampipe database instance as well as a directory of mod resources such as queries, benchmarks, and controls. Workspaces allow you to separate your Steampipe instances for security, operational, or organizational purposes.

## Examples

### List workspaces

```sql
select
  id,
  workspace_state,
  handle,
  identity_handle
from
  steampipecloud_workspace;
```

### List user workspaces

```sql
select
  id,
  workspace_state,
  handle,
  identity_handle
from
  steampipecloud_workspace
where
  identity_type = 'user';
```

### List org workspaces

```sql
select
  id,
  workspace_state,
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
  workspace_state,
  handle,
  identity_handle
from
  steampipecloud_workspace
where
  workspace_state <> 'running';
```