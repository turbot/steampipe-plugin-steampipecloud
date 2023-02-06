# Table: steampipecloud_workspace_process

Allows to track various activities performed on a workspace of an identity in Steampipe Cloud.

**Important notes:**

This table supports optional quals. Queries with optional quals in the `where` clause are optimised to use Steampipe Cloud filters.

Optional quals are supported for the following columns:

- `created_at`
- `id`
- `identity_handle`
- `identity_id`
- `pipeline_id`
- `query_where` - Allows use of [query filters](https://steampipe.io/docs/cloud/reference/query-filter). For a list of supported columns for workspace proceses, please see [Supported APIs and Columns](https://steampipe.io/docs/cloud/reference/query-filter#supported-apis--columns). Please note that any query filter passed into the `query_where` qual will be combined with other optional quals.
- `state`
- `type`
- `updated_at`
- `workspace_handle`
- `workspace_id`

## Examples

### Basic info

```sql
select
  id,
  identity_handle,
  workspace_handle,
  pipeline_id,
  type,
  state,
  created_at
from
  steampipecloud_workspace_process;
```

### List processes for a pipeline

```sql
select
  id,
  identity_handle,
  workspace_handle,
  pipeline_id,
  type,
  state,
  created_at
from
  steampipecloud_workspace_process
where
  pipeline_id = 'pipe_cfcgiefm1tumv1dis7lg';
```

### List running processes for a pipeline

```sql
select
  id,
  identity_handle,
  workspace_handle,
  pipeline_id,
  type,
  state,
  created_at
from
  steampipecloud_workspace_process
where
  pipeline_id = 'pipe_cfcgiefm1tumv1dis7lg'
  and state = 'running';
```

### List running processes for a pipeline using [query filter](https://steampipe.io/docs/cloud/reference/query-filter)

```sql
select
  id,
  identity_handle,
  workspace_handle,
  pipeline_id,
  type,
  state,
  created_at
from
  steampipecloud_workspace_process
where
  query_where = 'pipeline_id = ''pipe_cfcgiefm1tumv1dis7lg'' and state = ''running''';
```
