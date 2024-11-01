# Manage example Dedicated server bandwidth notification setting
resource "leaseweb_dedicated_server_notification_setting_bandwidth" "example" {
  dedicated_server_id = "12345"
  frequency           = "WEEKLY"
  threshold           = "1"
  unit                = "Gbps"
}
