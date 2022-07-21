# Table: steampipecloud_organization_workspace_member

Organization Workspace members can collaborate and share connections and dashboards.

## Examples

### Basic info

```sql
select
  id,
  org_handle,
  workspace_handle,
  user_handle,
  status,
  role
from
  steampipecloud_organization_workspace_member;
```

### List invited members

```sql
select
  id,
  org_handle,
  workspace_handle,
  user_handle,
  status,
  role
from
  steampipecloud_organization_workspace_member
where
  status = 'invited';
```

### List owners of an organization workspace

```sql
select
  id,
  org_handle,
  workspace_handle,
  user_handle,
  status,
  role 
from
  steampipecloud_organization_workspace_member 
where
  org_handle = 'testorg' 
  and workspace_handle = 'dev' 
  and role = 'owner';
```

### Get details of a particular member in an organization workspace

```sql
select
  id,
  org_handle,
  workspace_handle,
  user_handle,
  status,
  role 
from
  steampipecloud_organization_workspace_member 
where
  org_handle = 'testorg' 
  and workspace_handle = 'dev' 
  and user_handle = 'myuser';
```
