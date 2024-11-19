# Manage non HTTPS example Public Cloud load balancer listener
resource "leaseweb_public_cloud_load_balancer_listener" "example" {
  protocol = "HTTP"
  port     = 80
  default_rule = {
    target_group_id = "b05917e1-96a4-442a-900c-c41f273d95c9"
  }
  load_balancer_id = "695ddd91-051f-4dd6-9120-938a927a47d0"
}

# Manage HTTPS example Public Cloud load balancer listener
resource "leaseweb_public_cloud_load_balancer_listener" "example2" {
  protocol = "HTTPS"
  port     = 443
  certificate = {
    certificate = "-----BEGIN CERTIFICATE-----MIIBhDCB7gIBADBFMQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCtWdKNbZxvkXKAADjJMJ7VTJz6uFoMD403C+gMIF8hwqIsHggzCao6iXrW9sZoyZtUBVBiiq5RumHbbpwvOdMmXrShEB4sTJkWRMDy7yD4D91WCU1fc10E/zBJMwssAvmHZo5kGW1Pj1N9ktb+O/TMsEc6yd5suvdQj6aaJbQlTQIDAQABoAAwDQYJKoZIhvcNAQELBQADgYEAWOQ2CJLRo8MQgJgvhdoSIkHITnrbjB5hS3f/dx0lIcnyI6Q9nOyuQHXkCgkdBaV8lz7l+IbqcGc3CaIRP2ZIVFvo2252n630tOOSsqoqJS1tYIoIKsohi3T3d8T1i/s0BWbTJi8Xgd186wyUn/jHwXROKx2rq6yYsAO6fISDKw8=-----END CERTIFICATE-----"
    private_key = "-----BEGIN EC PRIVATE KEY-----MIHcAgEBBEIBVlC0IObonfQZIQ81l/WILKfWT5Fv96eNnYmQZ7uleu73igfiVESVuPfNlbW9oNEK1XcXli4YNZMxWMkKuzC3w8CgBwYFK4EEACOhgYkDgYYABAHvOqz9d2xeSpm1FNdw0NR5j/q6PMd6whZFsTPNYNj0/PsTpsHk78ZB4MYnJUXwHJjpj+gnKkLNc02f4w/vSF8VXADX4l40XU/w82TAOCftQwoxO5o0jZcwEUIYzl02Zd7uNxhjtKJQnYFi9x8WI8L8zQ6GZB/fJKYwoHaUr0I1h/5LzQ==-----END EC PRIVATE KEY-----"
  }
  default_rule = {
    target_group_id = "b05917e1-96a4-442a-900c-c41f273d95c9"
  }
  load_balancer_id = "695ddd91-051f-4dd6-9120-938a927a47d0"
}
