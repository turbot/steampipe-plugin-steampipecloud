# Table: steampipecloud_org




## Examples

### Basic info

```sql
select
  id,
  org_id,
  org_handle,
  status
from
  steampipecloud_org;
```

### List organization with owner role

```sql
select
  id,
  org_id,
  org_handle,
  status
from
  steampipecloud_org
where
  role = 'owner';
```