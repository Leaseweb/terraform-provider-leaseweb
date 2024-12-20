# List DNS resource record sets for example.com
data "leaseweb_dns_resource_record_sets" "all" {
  domain_name = "example.com"
}
