# Configuration-based authentication
terraform {
  required_providers {
    leaseweb = {
      version = "1.6.0-alpha"
      source  = "leaseweb/leaseweb"
    }
  }
}

provider "leaseweb" {
  token = "super-secret-token-value"
}
