# List all Public Cloud target groups
data "leaseweb_public_target_groups" "all" {}

# Get Public Cloud target groups filtered by id
data "leaseweb_public_target_groups" "example" {
  id = "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"
}

# Get Public Cloud target groups filtered by name
data "leaseweb_public_target_groups" "example2" {
  name = "Foo bar"
}

# Get Public Cloud target groups filtered by protocol
data "leaseweb_public_target_groups" "example3" {
  protocol = "HTTP"
}

# Get Public Cloud target groups filtered by port
data "leaseweb_public_target_groups" "example4" {
  port = 80
}

# Get Public Cloud target groups filtered by region
data "leaseweb_public_target_groups" "example5" {
  region = "eu-west-3"
}
