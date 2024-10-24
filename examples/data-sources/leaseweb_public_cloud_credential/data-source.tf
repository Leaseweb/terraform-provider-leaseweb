# Credential for public cloud
data "leaseweb_public_cloud_credential" "all" {
  instance_id = "12345"
  type        = "OPERATING_SYSTEM"
  username    = "root"
}
