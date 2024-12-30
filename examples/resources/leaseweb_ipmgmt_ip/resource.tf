# Manage an IP address
resource "leaseweb_ipmgmt_ip" "example" {
  ip = "192.0.2.1"
}

# Update reverse lookup for an ip address
resource "leaseweb_ipmgmt_ip" "example" {
  ip             = "192.0.2.1"
  reverse_lookup = "mydomain1.example.com"
}

