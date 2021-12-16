# Table: steampipecloud_connection

Connections represent a set of tables for a single data source. Each connection is represented as a distinct Postgres schema.

## Examples

### Basic info

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

### List organization workspaces

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
