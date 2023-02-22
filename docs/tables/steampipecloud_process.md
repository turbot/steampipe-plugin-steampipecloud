# Table: steampipecloud_process

Allows to track various activities performed on an identity in Steampipe Cloud.

## Examples

### Basic info

```sql
select
  id,
  identity_handle,
  identity_type,
  pipeline_id,
  type,
  state,
  created_at
from
  steampipecloud_process;
```

### List processes that are being run by an identity pipeline

```sql
select
  id,
  identity_handle,
  identity_type,
  pipeline_id,
  type,
  state,
  created_at
from
  steampipecloud_process
where
  pipeline_id is not null;
```

### List user processes

```sql
select
  id,
  identity_handle,
  identity_type,
  pipeline_id,
  type,
  state,
  created_at
from
  steampipecloud_process
where
  identity_type = 'user';
```

### List running processes

```sql
select
  id,
  identity_handle,
  identity_type,
  pipeline_id,
  type,
  state,
  created_at
from
  steampipecloud_process
where
  state = 'running';
```
