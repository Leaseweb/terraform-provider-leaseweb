# Credential for dedicated server
data "leaseweb_dedicated_server_credential" "all" {
  dedicated_server_id = "12345"
  type                = "OPERATING_SYSTEM"
  username            = "root"
}
