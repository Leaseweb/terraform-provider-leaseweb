# Manage example Public Cloud non TCP target group
resource "leaseweb_public_cloud_target_group" "example" {
  name     = "test"
  protocol = "HTTP"
  port     = 80
  region   = "eu-west-3"
  health_check = {
    port     = 80
    protocol = "HTTP"
    method   = "GET"
    host     = "example.com"
    uri      = "/endpoint"
  }
}

# Manage example Public Cloud target group with a health check
resource "leaseweb_public_cloud_target_group" "example" {
  name     = "test"
  protocol = "HTTP"
  port     = 80
  region   = "eu-west-3"
  health_check = {
    port     = 80
    protocol = "TCP"
    host     = "example.com"
  }
}

# Manage example Public Cloud without a health check
resource "leaseweb_public_cloud_target_group" "example" {
  name     = "test"
  protocol = "HTTP"
  port     = 80
  region   = "eu-west-3"
}
