# Manage example Public Cloud Instance
resource "leaseweb_public_cloud_instance" "example" {
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
