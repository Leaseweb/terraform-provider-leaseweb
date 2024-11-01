# Manage example Dedicated server credential
resource "leaseweb_dedicated_server_credential" "example" {
  dedicated_server_id = "12345"
  username            = "root"
  type                = "OPERATING_SYSTEM"
  password            = "mys3cr3tp@ssw0rd"
}
