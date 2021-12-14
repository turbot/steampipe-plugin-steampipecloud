# Table: steampipecloud_organization

Steampipe Cloud users can create their own connections and workspaces, but they are not shared with other users. Steampipe Organizations, on the other hand, include multiple users and are intended for organizations to collaborate and share workspaces and connections.

## Examples

### Basic info

```sql
select
  id,
  org_id,
  org_handle,
  status
from
  steampipecloud_organization;
```

### List organizations with owner role

```sql
select
  id,
  org_id,
  org_handle,
  status
from
  steampipecloud_organization
where
  role = 'owner';
```
