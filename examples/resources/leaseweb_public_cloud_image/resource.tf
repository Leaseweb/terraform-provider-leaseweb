# Manage example Public Cloud image
resource "leaseweb_public_cloud_image" "example" {
  instance_id = "396a3299-1795-464b-aa10-e1f179db1926"
  name        = "Custom image"
}
