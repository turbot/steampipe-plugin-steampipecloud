# Table: steampipecloud_token




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
  status <> 'active';
```