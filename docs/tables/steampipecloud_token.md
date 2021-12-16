# Table: steampipecloud_token

API tokens can be used to access the Steampipe Cloud API or to connect to Steampipe Cloud workspaces from the Steampipe CLI.

## Examples

### Basic info

```sql
select
  id,
  user_id,
  status,
  last4
from
  steampipecloud_token;
```

### List inactive tokens

```sql
select
  id,
  user_id,
  status,
  last4
from
  steampipecloud_token
where
  status = 'inactive';
```

### List tokens older than 90 days

```sql
select
  id,
  user_id,
  status,
  created_at,
  last4
from
  steampipecloud_token
where
  created_at <= (current_date - interval '90' day);
```
