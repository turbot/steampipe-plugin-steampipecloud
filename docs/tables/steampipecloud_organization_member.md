# Table: steampipecloud_organization_member

Organization members can collaborate and share workspaces and connections.

## Examples

### Basic info

```sql
select
  id,
  org_id,
  user_handle,
  status
from
  steampipecloud_organization_member;
```

### List invited members

```sql
select
  id,
  org_id,
  user_handle,
  status
from
  steampipecloud_organization_member
where
  status = 'invited';
```

### List members with owner role

```sql
select
  id,
  org_id,
  user_handle,
  status
from
  steampipecloud_organization_member
where
  role = 'owner';
```
