# Table: steampipecloud_workspace_mod_variable

Variables are module level objects that allow you to pass values to your module at runtime. When running Steampipe, you can pass values on the command line or from a .spvars file, and you will be prompted for any variables that have no values.

## Examples

### List basic information for all variables for a mod in a workspace

```sql
select
  id,
  name,
  description,
  value_default,
  value_setting,
  value,
  type
from
  steampipecloud_workspace_mod_variable
where
  workspace_id = 'w_cafeina2ip835d2eoacg' 
and 
  mod_alias = 'aws_thrifty';
```

### List all variables which have an explicit setting in a workspace mod

```sql
select
  id,
  name,
  description,
  value_default,
  value_setting,
  value,
  type
from
  steampipecloud_workspace_mod_variable
where
  workspace_id = 'w_cafeina2ip835d2eoacg'
and
  mod_alias = 'aws_thrifty' 
and 
  value_setting is not null;
```

### List details about a particular variable in a workspace mod

```sql
select
  id,
  name,
  description,
  value_default,
  value_setting,
  value,
  type
from
  steampipecloud_workspace_mod_variable
where
  workspace_id = 'w_cafeina2ip835d2eoacg'
and
  mod_alias = 'aws_tags' 
and 
  name = 'mandatory_tags';
```
