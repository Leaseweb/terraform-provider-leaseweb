# List all Dedicated servers
data "leaseweb_dedicated_servers" "all" {
  reference               = "test-reference"
  ip                      = "127.0.0.4"
  mac_address             = "aa:bb:cc:dd:ee:ff"
  site                    = "ams-01"
  private_rack_id         = "rack id"
  private_network_capable = "true"
  private_network_enabled = "true"
}
