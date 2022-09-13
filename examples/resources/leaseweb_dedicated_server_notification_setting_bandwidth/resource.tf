resource "leaseweb_dedicated_server_notification_setting_bandwidth" "alert" {
  dedicated_server_id = "1234567"
  frequency           = "DAILY"
  threshold           = 1.5
  unit                = "Gbps"
}
