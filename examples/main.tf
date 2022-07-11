terraform {
  required_providers {
    leaseweb = {
      version = "0.0.1"
      source  = "git.ocom.com/infra/leaseweb"
    }
  }
}

provider "leaseweb" {}

resource "leaseweb_dedicatedserver" "web01" {
  # reference = "web01"
}

resource "leaseweb_dedicatedserver" "web02" {
  # reverse_lookup = "web02.example.com"
}

resource "leaseweb_dedicatedserver" "db01" {
  # powered_on = true
}

resource "leaseweb_dedicatedserver" "db02" {
  # dhcp_lease = "https://boot.netboot.xyz"
}
