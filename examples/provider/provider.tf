terraform {
  required_providers {
    leaseweb = {
      version = "0.2.0"
      source  = "leaseweb/leaseweb"
    }
  }
}

provider "leaseweb" {}

resource "leaseweb_dedicated_server" "web" {
  reference = "web"
}
