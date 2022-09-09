# Access all the servers from AMS-01
data "leaseweb_dedicated_servers" "ams_01_servers" {
  site = "AMS-01"
}
