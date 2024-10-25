# Manage example Public Cloud Instance
resource "leaseweb_public_cloud_loadbalancer" "example" {
  contract = {
    billing_frequency = 1
    term              = 0
    type              = "HOURLY"
  }
  reference = "my webserver"
  region    = "eu-west-3"
  type      = "lsw.m3.large"
}
