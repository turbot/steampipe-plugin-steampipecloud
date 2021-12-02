# Table: steampipecloud_workspace_connection




## Examples

### List workspace connections

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

### List org workspace connections

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