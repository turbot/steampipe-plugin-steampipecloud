# Table: steampipecloud_workspace_mod

A Steampipe mod is a portable, versioned collection of related Steampipe resources such as dashboards, benchmarks, queries, and controls. Steampipe mods and mod resources are defined in HCL, and distributed as simple text files. Modules can be found on the Steampipe Hub, and may be shared with others from any public git repository.

## Examples

### Basic information about mods across all workspaces

```sql
select
  id,
  path,
  alias,
  constraint,
  installed_version,
  state
from
  steampipecloud_workspace_mod;
```

### List mods for all workspaces of the user

```sql
select
  id,
  path,
  alias,
  constraint,
  installed_version,
  state
from
  steampipecloud_workspace_mod
where
  identity_type = 'user';
```

### List mods for all workspaces belonging to all organizations that the user is a member of

```sql
select
  id,
  path,
  alias,
  constraint,
  installed_version,
  state
from
  steampipecloud_workspace_mod
where
  identity_type = 'org';
```

### List mods for a particular workspace belonging to an organization

```sql
select 
  swm.id,
  swm.path,
  swm.alias,
  swm.constraint,
  swm.installed_version,
  swm.state
from 
  steampipecloud_workspace_mod as swm 
  inner join steampipecloud_organization as so on so.id = swm.identity_id
  inner join steampipecloud_workspace as sw on sw.id = swm.workspace_id
where
  so.handle = 'testorg'
  and sw.handle = 'dev';
```
