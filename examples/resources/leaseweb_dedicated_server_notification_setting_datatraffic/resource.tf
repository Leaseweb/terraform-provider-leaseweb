# Manage example Dedicated server data traffic notification
resource "leaseweb_dedicated_server_notification_setting_datatraffic" "example" {
  dedicated_server_id = "12345"
  frequency           = "WEEKLY"
  threshold           = "12"
  unit                = "GB"
}
