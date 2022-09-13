# Access all control panels available with Ubuntu 22.04
data "leaseweb_dedicated_server_control_panels" "ubuntu_cps" {
  operating_system_id = "UBUNTU_22_04_64BIT"
}
