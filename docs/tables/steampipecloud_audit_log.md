# Table: steampipecloud_audit_log

Audit logs record series of events performed on the identity.

Note: You must specify an organization or user ID, or an organization or user handle, in the where or join clause using the `identity_id` or `identity_handle` columns respectively.

## Examples

### List user audit logs

```sql
select
  l.id as id,
  l.identity_id as identity_id,
  l.action_type as action_type,
  l.identity_handle as identity_handle
from
  steampipecloud_audit_log l,
  steampipecloud_user u
where
  l.identity_handle = u.handle;
```

### List organization workspaces

```sql
select
  l.id as id,
  l.identity_id as identity_id,
  l.action_type as action_type,
  l.identity_handle as identity_handle
from
  steampipecloud_audit_log l,
  steampipecloud_org o
where
  l.identity_handle = o.org_handle;
```
