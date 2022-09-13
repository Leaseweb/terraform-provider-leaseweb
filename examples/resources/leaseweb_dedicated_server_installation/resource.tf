resource "leaseweb_dedicated_server_installation" "frontend" {
  dedicated_server_id = "1234567"
  operating_system_id = "UBUNTU_22_04_64BIT"
  control_panel_id    = "PLESK_DEDSER_WEB_ADMIN"

  hostname = "www.example.com"
  timezone = "UTC"
  password = "hunter2"

  ssh_keys = [
    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCyA/oHo5JPiWPbXFjnHs06kTGVP5dcx6gfBnB8Fg6NJg5sbnHIywd1kXY0XbS4hDKpvbnEIBBs9kX0ps8Hra0GniJs/FdI+9T+/15VkAdzgSuI+oEi2M2oydwwVwR+YancG7NKpmHa3dtCifRpC0EPvHJUfe4r+aQ+7FSXSXp4PbAhd0zoE8fXS/sUoPflMcMyffxdRCCg8krGk0757FHAGttQRvWST1lv9w42CaInkRgV3ncTZy/buoZJ2YnaonpzoExFaDJU7HDi49yUN3S/PptdF0Ce7f6fCKd826wQBcz9ilmHOiXOYb3RHIXaEdJEuz99EWO09S7aV5dSOhbh4VHZTQESLCvcJXif9aeFY80Nz924k1HiGEtNow96CNwlIm1cWmNFdIK+y/DJVJOoZYZGyT0L8Hp/ggVK9aTn5BAi+4HR4kAZsEMP/6/C65aXvIo3f/L7CkcW0kuQmlisjY8Ak3jsofhKGuLguB7kx2v3BrX1udO4M7p4YdYpjN8= user1@example.org",
    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDPyQ01ld53VQyp/ie9md6XNvT+Ix1AAyFykNDXg6z6aMDe/LLzNhR9cDsxzIEaypMrS283R3H8nv7c/TKkVME3ACwYjlB6sq9aWOSEHIHQ451DC0vLtTA3ZYu3dP7E7ygKGgqSHjBaItm9ettWLhU4Ffa8z55vIIIAP6qFdiq7z4FOIGrWEBG9DO9ulr0foWm5tUvqQJh38zH3FU96RvQu4q9+K99XYEJA/ib6tfuLCwpBBBkou1T9hX4ChHqeNcg5EwuKitcW7xm7OFjv0RqDgzsbmyIWeB3ZmRcTwglz+ZMUgnxql1xxhQyJnmjY8dgauLF8Q5zXBox3ZMHYiD0EOwpi/oEoVOiQsk95hYqrdssqNEDksW/UwA70yvQoaBheDDxmXsQCf5aZz1kU4DMVjSEiR0k1xX3i5tugfClr5oaqFQwGsqBeyKOwdbUp0CKF/Bs8F5SnY0G2/j5TMir4y9vWscg7AMmbc50OKiQxvWg10+Qpnmw9ewZXjxmlzSE= user2@example.org",
  ]

  post_install_script = <<-EOS
    #!/bin/sh
    apt install nginx -y -qq
  EOS

  callback_url = "https://www.example.com/callback"

  raid {
    type            = "HW"
    level           = 1
    number_of_disks = 2
  }

  device = "SATA_SAS"

  partition {
    mountpoint = "/boot"
    size       = 1024
    filesystem = "ext2"
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
  }

  timeouts {
    create = "30m"
  }
}
