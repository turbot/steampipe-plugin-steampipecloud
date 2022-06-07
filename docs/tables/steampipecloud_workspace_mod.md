# Table: steampipecloud_workspace_mod

A Steampipe mod is a portable, versioned collection of related Steampipe resources such as dashboards, benchmarks, queries, and controls. Steampipe mods and mod resources are defined in HCL, and distributed as simple text files. Modules can be found on the Steampipe Hub, and may be shared with others from any public git repository.

## Examples

### Basic info

```sql
select
  id,
  path,
  alias,
  mod_constraint,
  installed_version,
  state
from
  steampipecloud_workspace_mod;
```

### List user workspaces

```sql
select
  id,
  path,
  alias,
  mod_constraint,
  installed_version,
  state
from
  steampipecloud_workspace_mod
where
  identity_type = 'user';
```

### List organization workspaces

```sql
select
  id,
  path,
  alias,
  mod_constraint,
  installed_version,
  state
from
  steampipecloud_workspace_mod
where
  identity_type = 'org';
```

### List workspaces mods which are not installed

```sql
select
  id,
  path,
  alias,
  mod_constraint,
  installed_version,
  state
from
  steampipecloud_workspace_mod
where
  state <> 'installed';
```
