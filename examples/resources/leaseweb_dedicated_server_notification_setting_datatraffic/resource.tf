resource "leaseweb_dedicated_server_notification_setting_datatraffic" "alert" {
  dedicated_server_id = "1234567"
  frequency           = "WEEKLY"
  threshold           = 2
  unit                = "TB"
}
