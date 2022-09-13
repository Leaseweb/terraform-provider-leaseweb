resource "leaseweb_dedicated_server" "web01" {
  reference                       = "web01"
  reverse_lookup                  = "web01.example.com"
  dhcp_lease                      = "https://boot.netboot.xyz"
  powered_on                      = true
  public_network_interface_opened = true
  public_ip_null_routed           = false
}
