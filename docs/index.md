---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "jupyterhub Provider"
subcategory: ""
description: |-
  Interact with JupyterHub.
---

# jupyterhub Provider

Interact with JupyterHub.

## Example Usage

```terraform
provider "jupyterhub" {
  protocol = "http"
  host     = "localhost:8000"
  prefix   = "/"
  token    = "abcdefghijklmnopqrstuvxyz"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `host` (String) Hostname for JupyterHub API. Default is 'localhost:8000'. May also be provided via JUPYTERHUB_HOST environment variable.
- `password` (String, Sensitive) API Token for JupyterHub API. Optional may also be provided via JUPYTERHUB_PASSWORD environment variable.
- `prefix` (String) Prefix for JupyterHub API. Default is '/'. May also be provided via JUPYTERHUB_PREFIX environment variable.
- `protocol` (String) Protocol for JupyterHub API. Default is 'http'. May also be provided via JUPYTERHUB_PROTOCOL environment variable.
- `token` (String, Sensitive) API Token for JupyterHub API. Optional if username and password are set. May also be provided via JUPYTERHUB_TOKEN environment variable.
- `username` (String, Sensitive) API Token for JupyterHub API. Optional may also be provided via JUPYTERHUB_USERNAME environment variable.
