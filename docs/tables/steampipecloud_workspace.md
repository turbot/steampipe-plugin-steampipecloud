# Table: steampipecloud_workspace




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
  identity_id like 'u%';
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
  identity_id like 'o%';
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