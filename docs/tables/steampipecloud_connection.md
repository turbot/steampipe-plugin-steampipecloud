# Table: steampipecloud_connection

A Steampipe Connection represents a set of tables for a single data source. Each connection is represented as a distinct Postgres schema.

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
  identity_type = 'user';
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
  identity_type = 'org';
```
