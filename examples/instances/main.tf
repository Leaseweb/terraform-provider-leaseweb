terraform {
  required_providers {
    leaseweb = {
      source = "registry.terraform.io/LeaseWeb/leaseweb"
    }
  }
}

provider "leaseweb" {
  host   = "localhost:8080"
  scheme = "http"
  token  = "tralala"
}

data "leaseweb_instances" "public_cloud_instances" {}

output "public_cloud" { value = data.leaseweb_instances.public_cloud_instances }
