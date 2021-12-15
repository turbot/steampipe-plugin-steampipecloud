# Table: steampipecloud_organization_member

 The member of an organization who can collaborate and share workspaces and connections.

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
  status <> 'accepted';
```