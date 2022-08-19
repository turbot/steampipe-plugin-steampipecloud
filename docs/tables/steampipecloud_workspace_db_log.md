# Table: steampipecloud_workspace_db_log

Database logs records the underlying queries executed when a user executes a query.

## Examples

### List db logs for an actor by handle

```sql
select
  id,
  workspace_id,
  workspace_handle,
  duration,
  query,
  log_timestamp
from
  steampipecloud_workspace_db_log
where
  actor_handle = 'siddharthaturbot';
```

### List db logs for an actor by handle in a particular workspace

```sql
select
  id,
  workspace_id,
  workspace_handle,
  duration,
  query,
  log_timestamp
from
  steampipecloud_workspace_db_log
where
  actor_handle = 'siddharthaturbot'
  and workspace_handle = 'dev';
```

### List queries that took more than 30 seconds to execute

```sql
select
  id,
  workspace_id,
  workspace_handle,
  duration,
  query,
  log_timestamp
from
  steampipecloud_workspace_db_log
where
  duration > 30000;
```

### List all queries that ran in my workspace in the last hour

```sql
select
  id,
  workspace_id,
  workspace_handle,
  duration,
  query,
  log_timestamp
from
  steampipecloud_workspace_db_log
where
  workspace_handle = 'dev'
  and log_timestamp > now() - interval '1 hr';
```
