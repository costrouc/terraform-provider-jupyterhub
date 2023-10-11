terraform {
  required_providers {
    jupyterhub = {
      source = "registry.terraform.io/costrouc/jupyterhub"
    }
  }
}

provider "jupyterhub" {
  protocol = "http"
  host = "localhost:8000"
  prefix = "/"
  token = " b3a2d0844af3413f972dceb46d7a33f7 "
}
