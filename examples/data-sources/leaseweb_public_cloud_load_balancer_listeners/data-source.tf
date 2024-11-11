# List Public Cloud load balancer listeners that belong to the load balancer
data "leaseweb_public_cloud_load_balancer_listeners" "example" {
  load_balancer_id = "695ddd91-051f-4dd6-9120-938a927a47d0"
}
