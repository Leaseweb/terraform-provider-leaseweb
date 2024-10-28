# Manage example public cloud credential
resource "leaseweb_public_cloud_credential" "example" {
  instance_id = "12345"
  username    = "root"
  type        = "OPERATING_SYSTEM"
  password    = "mys3cr3tp@ssw0rd"
}
