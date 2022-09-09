resource "leaseweb_dedicated_server_credential" "control_panel" {
  dedicated_server_id = "1234567"
  type                = "CONTROL_PANEL"
  username            = "AzureDiamond"
  password            = "hunter2"
}
