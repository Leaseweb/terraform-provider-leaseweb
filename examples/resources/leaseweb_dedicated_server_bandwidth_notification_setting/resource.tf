# Manage example Dedicated server bandwidth notification setting
resource "leaseweb_dedicated_server_bandwidth_notification_setting" "example" {
  dedicated_server_id = "12345"
  frequency           = "WEEKLY"
  threshold           = "1"
  unit                = "Gbps"
}
