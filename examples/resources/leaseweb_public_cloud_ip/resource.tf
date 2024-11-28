# Manage example Public Cloud ip
resource "leaseweb_public_cloud_ip" "example" {
  instance_id    = "695ddd91-051f-4dd6-9120-938a927a47d0"
  ip             = "10.0.0.1"
  reverse_lookup = "example.com"
}
