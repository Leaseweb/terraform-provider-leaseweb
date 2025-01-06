# Manage null route history
resource "leaseweb_ipmgmt_null_route" "nr" {
  ip = "127.0.0.1"
}

# Update comment for a null route history
resource "leaseweb_ipmgmt_null_route" "nr" {
  id      = "123456"
  comment = "this is comment"
}

