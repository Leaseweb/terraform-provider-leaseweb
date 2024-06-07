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

resource "leaseweb_instance" "public_cloud_instance" {
  region    = "eu-west-3"
  type      = "lsw.m3.large"
  reference = "my webserver"
  operating_system = {
    id = "UBUNTU_22_04_64BIT"
  }
  root_disk_storage_type = "CENTRAL"
  contract = {
    billing_frequency = 1
    term              = 0
    type              = "HOURLY"
  }
}

output "public_cloud" {
  value = leaseweb_instance.public_cloud_instance
}
