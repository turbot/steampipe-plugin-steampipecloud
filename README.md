# SteampipeCloud Plugin for Steampipe

Use SQL to query infrastructure including servers, networks, identity and more from SteampipeCloud.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/steampipecloud)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/steampipecloud/tables)
- Community: [Slack Channel](https://join.slack.com/t/steampipe/shared_invite/zt-oij778tv-lYyRTWOTMQYBVAbtPSWs3g)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-steampipecloud/issues)

## Quick Start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install steampipecloud
```

Run a query:

```sql
select
  id,
  workspace_state,
  identity_handle
from
  steampipecloud_workspace;
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

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-digitalocean/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [SteampipeCloud Plugin](https://github.com/turbot/steampipe-plugin-steampipecloud/labels/help%20wanted)