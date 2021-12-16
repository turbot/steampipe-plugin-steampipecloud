# Table: steampipecloud_workspace_connection

Workspace connections are the associations between workspaces and connections.

## Examples

### Basic info

```sql
select
  id,
  connection_id,
  workspace_id,
  identity_id,
  jsonb_pretty(connection) as connection
from
  steampipecloud_workspace_connection;
```

### List workspace connections using AWS plugin

```sql
select
  id,
  connection_id,
  workspace_id,
  identity_id,
  jsonb_pretty(connection) as connection
from
  steampipecloud_workspace_connection
where
  connection ->> 'plugin' = 'aws';
```

### List user workspace connections

```sql
select
  id,
  connection_id,
  workspace_id,
  identity_id,
  jsonb_pretty(connection) as connection
from
  steampipecloud_workspace_connection
where
  identity_id like 'u_%';
```

### List organization workspace connections

```sql
select
  id,
  connection_id,
  workspace_id,
  identity_id,
  jsonb_pretty(connection) as connection
from
  steampipecloud_workspace_connection
where
  identity_id like 'o_%';
```
