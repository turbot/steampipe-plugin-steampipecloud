# Table: steampipecloud_token

A steampipe cloud token which can be used to sign programmatic request and perform action on steampipe cloud.


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