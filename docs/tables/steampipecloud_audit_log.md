# Table: steampipecloud_audit_log

Audit logs record series of events performed on the identity.

Note: You must specify an organization or user ID, or an organization or user handle, in the where or join clause using the `identity_id` or `identity_handle` columns respectively.

## Examples

### List audit logs for a user handle

```sql
select
  id,
  action_type,
  jsonb_pretty(data) as data
from
  steampipecloud_audit_log
where
  identity_handle = 'myuser';
```

### List audit logs for a user ID

```sql
select
  id,
  action_type,
  jsonb_pretty(data) as data
from
  steampipecloud_audit_log
where
  identity_id = 'u_c6fdjke232example';
```

### List audit logs for an organization handle

```sql
select
  id,
  action_type,
  jsonb_pretty(data) as data
from
  steampipecloud_audit_log
where
  identity_handle = 'myorg';
```

### List audit logs for an organization ID

```sql
select
  id,
  action_type,
  jsonb_pretty(data) as data
from
  steampipecloud_audit_log
where
  identity_id = 'o_c6qjjsaa6guexample';
```
