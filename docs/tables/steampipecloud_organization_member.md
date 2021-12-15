# Table: steampipecloud_organization_member

Organization members can collaborate and share workspaces and connections.

## Examples

### Basic info

```sql
select
  id,
  org_id,
  user_id,
  status
from
  steampipecloud_organization_member;
```

### List pending members

```sql
select
  id,
  org_id,
  user_id,
  status
from
  steampipecloud_organization_member
where
  status = 'pending';
```
