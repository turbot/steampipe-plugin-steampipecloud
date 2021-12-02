# Table: steampipecloud_user

A steampipe cloud user is an entity that interacts with the application.


## Examples

### Basic info

```sql
select
  id,
  display_name,
  status,
  handle
from
  steampipecloud_user;
```