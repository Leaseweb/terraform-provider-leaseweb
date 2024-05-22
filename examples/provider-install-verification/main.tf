terraform {
  required_providers {
    leaseweb = {
      source = "registry.terraform.io/LeaseWeb/leaseweb"
    }
  }
}

provider "leaseweb" {}

data "leaseweb_instances" "example" {}
