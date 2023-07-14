terraform {
  required_providers {
    leaseweb = {
      version = "0.3.1"
      source  = "leaseweb/leaseweb"
    }
  }
}

provider "leaseweb" {}

resource "leaseweb_dedicated_server" "web" {
  reference = "web"
}
