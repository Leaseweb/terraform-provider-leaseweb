terraform {
  required_providers {
    leaseweb = {
      version = "0.1.0"
      source  = "leaseweb/leaseweb"
    }
  }
}

provider "leaseweb" {}

resource "leaseweb_dedicated_server" "web" {
  reference = "web"
}
