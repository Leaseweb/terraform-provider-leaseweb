terraform {
  required_providers {
    leaseweb = {
      source = "registry.terraform.io/LeaseWeb/leaseweb"
    }
  }
}

provider "leaseweb" {
  host  = "127.0.0.1:4010"
  token = "tralala"
}

data "leaseweb_instances" "example" {}
