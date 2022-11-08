# Table: steampipecloud_user_email

User Emails allows a user to manage their emails.

The `steampipecloud_user_email` table returns a list of emails added by a user to their profile.

## Examples

### Basic info

```sql
select
  id,
  email,
  status,
  created_at
from
  steampipecloud_user_email;
```
