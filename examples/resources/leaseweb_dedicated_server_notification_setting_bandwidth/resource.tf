# Manage example Dedicated Server Notification Setting Bandwidth
resource "leaseweb_dedicated_server_notification_setting_bandwidth" "example" {
  server_id = "1234567"
  frequency = "DAILY"
  threshold = "1"
  unit      = "Gbps"
}
