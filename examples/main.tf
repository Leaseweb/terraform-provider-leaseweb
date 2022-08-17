terraform {
  required_providers {
    leaseweb = {
      version = "0.0.1"
      source  = "git.ocom.com/infra/leaseweb"
    }
  }
}

provider "leaseweb" {}

data "leaseweb_dedicated_server_operating_systems" "all_os" {
}


locals {
  latest_ubuntu_os_id = reverse(sort([
    for id in data.leaseweb_dedicated_server_operating_systems.all_os.ids : id
    if length(regexall("^UBUNTU_.*", id)) > 0
  ]))[0]
  supported_controlpanel_id = reverse(sort([
    for id in data.leaseweb_dedicated_server_control_panels.all_cp.ids : id
    if length(regexall("^PLESK_.*", id)) > 0
  ]))[0]
  hostname = "web01.example.org"
}

data "leaseweb_dedicated_server_control_panels" "all_cp" {
  # Providing an Operating System Id is optional. If provided, Control Panels supported by the provided Operating System will be listed.
  operating_system_id = local.latest_ubuntu_os_id
}

resource "leaseweb_dedicated_server" "my-test" {
  # reference = "web01"
  # reverse_lookup = "web02.example.com"
  # dhcp_lease = "https://boot.netboot.xyz"
  # powered_on = true
  # main_ip_nulled = false
}

resource "leaseweb_dedicated_server_credential" "os" {
  dedicated_server_id = leaseweb_dedicated_server.my-test.id
  type                = "OPERATING_SYSTEM"
  username            = "root"
  password            = "Abcdef.123456"
}

resource "leaseweb_dedicated_server_installation" "my-ubuntu" {
  dedicated_server_id = leaseweb_dedicated_server.my-test.id
  operating_system_id = local.latest_ubuntu_os_id
  control_panel_id    = local.supported_controlpanel_id
  password            = leaseweb_dedicated_server_credential.os.password

  hostname = local.hostname
  timezone = "Europe/Amsterdam"
  ssh_keys = [
    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCyA/oHo5JPiWPbXFjnHs06kTGVP5dcx6gfBnB8Fg6NJg5sbnHIywd1kXY0XbS4hDKpvbnEIBBs9kX0ps8Hra0GniJs/FdI+9T+/15VkAdzgSuI+oEi2M2oydwwVwR+YancG7NKpmHa3dtCifRpC0EPvHJUfe4r+aQ+7FSXSXp4PbAhd0zoE8fXS/sUoPflMcMyffxdRCCg8krGk0757FHAGttQRvWST1lv9w42CaInkRgV3ncTZy/buoZJ2YnaonpzoExFaDJU7HDi49yUN3S/PptdF0Ce7f6fCKd826wQBcz9ilmHOiXOYb3RHIXaEdJEuz99EWO09S7aV5dSOhbh4VHZTQESLCvcJXif9aeFY80Nz924k1HiGEtNow96CNwlIm1cWmNFdIK+y/DJVJOoZYZGyT0L8Hp/ggVK9aTn5BAi+4HR4kAZsEMP/6/C65aXvIo3f/L7CkcW0kuQmlisjY8Ak3jsofhKGuLguB7kx2v3BrX1udO4M7p4YdYpjN8= user1@example.org",
    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDPyQ01ld53VQyp/ie9md6XNvT+Ix1AAyFykNDXg6z6aMDe/LLzNhR9cDsxzIEaypMrS283R3H8nv7c/TKkVME3ACwYjlB6sq9aWOSEHIHQ451DC0vLtTA3ZYu3dP7E7ygKGgqSHjBaItm9ettWLhU4Ffa8z55vIIIAP6qFdiq7z4FOIGrWEBG9DO9ulr0foWm5tUvqQJh38zH3FU96RvQu4q9+K99XYEJA/ib6tfuLCwpBBBkou1T9hX4ChHqeNcg5EwuKitcW7xm7OFjv0RqDgzsbmyIWeB3ZmRcTwglz+ZMUgnxql1xxhQyJnmjY8dgauLF8Q5zXBox3ZMHYiD0EOwpi/oEoVOiQsk95hYqrdssqNEDksW/UwA70yvQoaBheDDxmXsQCf5aZz1kU4DMVjSEiR0k1xX3i5tugfClr5oaqFQwGsqBeyKOwdbUp0CKF/Bs8F5SnY0G2/j5TMir4y9vWscg7AMmbc50OKiQxvWg10+Qpnmw9ewZXjxmlzSE= user2@example.org",
  ]
  post_install_script = <<-EOS
    #!/bin/sh
    apt install nginx -y -qq
    echo "${local.hostname} on ${leaseweb_dedicated_server.my-test.main_ip}" > /var/www/html/index.html
  EOS

  device = "SATA_SAS"

  partition {
    mountpoint = "/boot"
    size       = 1024
    filesystem = "ext2"
    bootable   = true
    primary    = true
  }

  partition {
    size       = 4096
    filesystem = "swap"
  }

  partition {
    mountpoint = "/tmp"
    size       = 4096
    filesystem = "ext4"
  }

  # order matters: this partition needs to be at the end because of the * size
  partition {
    mountpoint = "/"
    size       = "*"
    filesystem = "ext4"
    primary    = true
  }

  raid {
    type            = "HW"
    level           = 1
    number_of_disks = 2
  }

  timeouts {
    create = "30m"
  }
}

resource "leaseweb_dedicated_server_notification_setting_bandwidth" "alert" {
  dedicated_server_id = leaseweb_dedicated_server.my-test.id
  frequency           = "DAILY"
  threshold           = 1.5
  unit                = "Gbps"
}

resource "leaseweb_dedicated_server_notification_setting_datatraffic" "alert" {
  dedicated_server_id = leaseweb_dedicated_server.my-test.id
  frequency           = "WEEKLY"
  threshold           = 2
  unit                = "TB"
}

resource "leaseweb_dedicated_server_credential" "firewall" {
  dedicated_server_id = leaseweb_dedicated_server.my-test.id
  type                = "FIREWALL"
  username            = "admin"
  password            = "abcdef"

  # Installation will delete all credentials, so this resource needs to be created afterwards
  depends_on = [
    leaseweb_dedicated_server_installation.my-ubuntu
  ]
}

output "latest_ubuntu_os_name" {
  value = data.leaseweb_dedicated_server_operating_systems.all_os.names[local.latest_ubuntu_os_id]
}

output "control_panels_list" {
  value = data.leaseweb_dedicated_server_control_panels.all_cp.ids
}
