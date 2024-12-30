# Manage a DNS record
resource "leaseweb_dns_resource_record_set" "example" {
  domain_name = "example.com"
  content = [
    "85.17.150.51",
    "85.17.150.52",
    "85.17.150.53"
  ]
  name = "example.com."
  type = "A"
  ttl  = 3600
}
