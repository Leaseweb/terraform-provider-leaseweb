# Attach an ISO to a Public Cloud instance
resource "leaseweb_public_cloud_instance_iso" "example" {
  instance_id = "695ddd91-051f-4dd6-9120-938a927a47d0"
  desired_id  = "ACRONIS_BOOT_MEDIA"
}
