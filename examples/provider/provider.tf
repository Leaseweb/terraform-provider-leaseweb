terraform {
  required_providers {
    leaseweb = {
      version = "0.0.1"
      source  = "git.ocom.com/infra/leaseweb"
    }
  }
}

provider "leaseweb" {}

resource "leaseweb_dedicated_server" "web" {
  reference = "web"
}
