# Manage example Public Cloud load balancer
resource "leaseweb_public_cloud_load_balancer" "example" {
  contract = {
    billing_frequency = 1
    term              = 0
    type              = "HOURLY"
  }
  reference = "my webserver"
  region    = "eu-west-3"
  type      = "lsw.m3.large"
}
