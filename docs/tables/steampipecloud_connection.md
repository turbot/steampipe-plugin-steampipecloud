# Table: steampipecloud_connection




## Examples

### List connections

```sql
select
  id,
  plugin,
  handle,
  identity_handle
from
  steampipecloud_connection;
```

### List user connections

```sql
select
  id,
  plugin,
  handle,
  identity_handle
from
  steampipecloud_connection
where
  identity_id like 'u%';
```

### List org workspaces

```sql
select
  id,
  plugin,
  handle,
  identity_handle
from
  steampipecloud_connection
where
  identity_id like 'o%';
```