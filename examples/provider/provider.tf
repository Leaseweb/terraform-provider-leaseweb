terraform {
  required_providers {
    leaseweb = {
      version = "~> 1.2.0"
      source  = "leaseweb/leaseweb"
    }
  }
}

provider "leaseweb" {
  token = "527070ca-8449-4f06-b609-ec6797bd8222"
}
