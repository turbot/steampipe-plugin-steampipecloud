connection "steampipecloud" {
  plugin = "steampipecloud"

  # Steampipe Cloud API token. If `token` is not specified, it will be loaded
  # from the `STEAMPIPE_CLOUD_TOKEN` environment variable.
  # token = "spt_thisisnotarealtoken_123"

  # Steampipe Cloud host URL. This defaults to "https://cloud.steampipe.io/".
  # You only need to set this if connecting to a remote Steampipe Cloud database
  # not hosted in "https://cloud.steampipe.io/".
  # If `host` is not specified, it will be loaded from the `STEAMPIPE_CLOUD_HOST`
  # environment variable.
  # host = "https://cloud.steampipe.io"
}