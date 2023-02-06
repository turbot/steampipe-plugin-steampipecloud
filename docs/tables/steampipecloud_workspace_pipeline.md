# Table: steampipecloud_workspace_pipeline

Pipelines allow users to run different kinds of activities in Steampipe Cloud on a schedule.

**Important notes:**

This table supports optional quals. Queries with optional quals in the `where` clause are optimised to use Steampipe Cloud filters.

Optional quals are supported for the following columns:

- `created_at`
- `id`
- `identity_handle`
- `identity_id`
- `pipeline`
- `query_where` - Allows use of [query filters](https://steampipe.io/docs/cloud/reference/query-filter). Please note that any query filter passed into the `query_where` qual will be combined with other optional quals.
- `title`
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
  title,
  frequency,
  pipeline,
  args,
  tags,
  last_process_id
from
  steampipecloud_workspace_pipeline;
```

### List pipelines for a specific workspace

```sql
select
  id,
  identity_handle,
  workspace_handle,
  title,
  frequency,
  pipeline,
  args,
  tags,
  last_process_id
from
  steampipecloud_workspace_pipeline
where
  workspace_handle = 'dev';
```

### List pipelines of frequency type `interval` for a specific workspace

```sql
select
  id,
  identity_handle,
  workspace_handle,
  title,
  frequency,
  pipeline,
  args,
  tags,
  last_process_id
from
  steampipecloud_workspace_pipeline
where
  workspace_handle = 'dev'
  and frequency->>'type' = 'interval';
```

### List pipelines for the `AWS Compliance CIS v1.4.0` dashboard created in the last 7 days

```sql
select
  id,
  identity_handle,
  workspace_handle,
  title,
  frequency,
  pipeline,
  args,
  tags,
  last_process_id
from
  steampipecloud_workspace_pipeline
where
  args->>'resource' = 'aws_compliance.benchmark.cis_v140'
  and created_at >= now() - interval '7 days';
```

### List pipelines for the `AWS Compliance CIS v1.4.0` dashboard created in the last 7 days using [query filter](https://steampipe.io/docs/cloud/reference/query-filter)

```sql
select
  id,
  identity_handle,
  workspace_handle,
  title,
  frequency,
  pipeline,
  args,
  tags,
  last_process_id
from
  steampipecloud_workspace_pipeline
where
  query_where = 'title = ''Scheduled snapshot: CIS v1.4.0'' and created_at >= now() - interval ''7 days''';
```
