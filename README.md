# :warning: DEPRECATED

The Steampipe Cloud plugin has been deprecated as part of our [renaming](https://turbot.com/blog/2023/07/introducing-turbot-guardrails-and-pipes) of Steampipe Cloud to Turbot Pipes. Please use the [Turbot Pipes plugin](https://hub.steampipe.io/plugins/turbot/pipes) instead.

---
![image](https://hub.steampipe.io/images/plugins/turbot/steampipecloud-social-graphic.png)

# Steampipe Cloud Plugin for Steampipe

Use SQL to query workspaces, connections and more from Steampipe Cloud.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/steampipecloud)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/steampipecloud/tables)
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-steampipecloud/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install steampipecloud
```

Run a query:

```sql
select
  user_handle,
  email,
  status
from
  steampipecloud_organization_member
where
  status = 'accepted'
```

```
> select user_handle, email, status from steampipecloud_organization_member where status = 'accepted'
+-------------+------------------+----------+
| user_handle | email            | status   |
+-------------+------------------+----------+
| mario       | mario@turbot.com | accepted |
| yoshi       | yoshi@turbot.com | accepted |
+-------------+------------------+----------+
```

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-steampipecloud.git
cd steampipe-plugin-steampipecloud
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```
make
```

Configure the plugin:

```
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/steampipecloud.spc
```

Try it!

```
steampipe query
> .inspect steampipecloud
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-steampipecloud/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Steampipe Cloud Plugin](https://github.com/turbot/steampipe-plugin-steampipecloud/labels/help%20wanted)
