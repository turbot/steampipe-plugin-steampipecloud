# Table: steampipecloud_workspace_connection

Workspace connections are the associations between workspaces and connections.

## Examples

### Basic info

```sql
select
  id,
  connection_id,
  workspace_id,
  identity_id
from
  steampipecloud_workspace_connection;
```

### List user workspace connections

```sql
select
  id,
  connection_id,
  workspace_id,
  identity_id
from
  steampipecloud_workspace_connection
where
  identity_id like 'u%';
```

### List organization workspace connections

```sql
select
  id,
  connection_id,
  workspace_id,
  identity_id
from
  steampipecloud_workspace_connection
where
  identity_id like 'o%';
```
