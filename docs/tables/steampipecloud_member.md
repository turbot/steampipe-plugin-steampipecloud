# Table: steampipecloud_member




## Examples

### Basic info

```sql
select
  id,
  org_id,
  user_id,
  status
from
  steampipecloud_member;
```

### List pending members

```sql
select
  id,
  org_id,
  user_id,
  status
from
  steampipecloud_member
where
  status <> 'accepted';
```