# Example install operating system on dedicated server
resource "leaseweb_dedicated_server_installation" "example" {
  dedicated_server_id = "12345"
  operating_system_id = "UBUNTU_22_04_64BIT"
}

# Example install operating system on dedicated server with post_install_script
resource "leaseweb_dedicated_server_installation" "example" {
  dedicated_server_id = "12345"
  operating_system_id = "UBUNTU_22_04_64BIT"
  post_install_script = <<-EOS
      #!/bin/sh
      apt install nginx -y -qq
  EOS
}