---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "leaseweb_public_cloud_instance Resource - leaseweb"
subcategory: ""
description: |-
  
---

# leaseweb_public_cloud_instance (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `contract` (Attributes) (see [below for nested schema](#nestedatt--contract))
- `operating_system` (Attributes) (see [below for nested schema](#nestedatt--operating_system))
- `region` (String) Region to launch the instance into
- `root_disk_storage_type` (String) The root disk's storage type
- `type` (String) Instance type

### Optional

- `market_app_id` (String) Market App ID that must be installed into the instance
- `reference` (String) The identifying name set to the instance
- `root_disk_size` (Number) The root disk's size in GB. Must be at least 5 GB for Linux and FreeBSD instances and 50 GB for Windows instances
- `ssh_key` (String) Public SSH key to be installed into the instance. Must be used only on Linux/FreeBSD instances

### Read-Only

- `customer_id` (String) The customer ID who owns the instance
- `equipment_id` (String) Equipment's UUID
- `has_private_network` (Boolean)
- `has_public_ipv4` (Boolean)
- `id` (String) The instance unique identifier
- `ips` (Attributes List) (see [below for nested schema](#nestedatt--ips))
- `iso` (Attributes) (see [below for nested schema](#nestedatt--iso))
- `private_network` (Attributes) (see [below for nested schema](#nestedatt--private_network))
- `product_type` (String) The product type
- `resource` (Attributes) i available for the load balancer (see [below for nested schema](#nestedatt--resource))
- `sales_org_id` (String)
- `started_at` (String) Date and time when the instance was started for the first time, right after launching it
- `state` (String) The instance's current state

<a id="nestedatt--contract"></a>
### Nested Schema for `contract`

Required:

- `billing_frequency` (Number) The billing frequency (in months) of the instance.
- `term` (Number) Contract term (in months). Used only when contract type is MONTHLY
- `type` (String) Select HOURLY for billing based on hourly usage, else MONTHLY for billing per month usage

Read-Only:

- `created_at` (String) Date when the contract was created
- `ends_at` (String)
- `renewals_at` (String) Date when the contract will be automatically renewed
- `state` (String)


<a id="nestedatt--operating_system"></a>
### Nested Schema for `operating_system`

Required:

- `id` (String) Operating System ID

Read-Only:

- `architecture` (String)
- `family` (String)
- `flavour` (String)
- `market_apps` (List of String)
- `name` (String)
- `storage_types` (List of String) The supported storage types for the instance type
- `version` (String)


<a id="nestedatt--ips"></a>
### Nested Schema for `ips`

Read-Only:

- `ddos` (Attributes) (see [below for nested schema](#nestedatt--ips--ddos))
- `ip` (String)
- `main_ip` (Boolean)
- `network_type` (String)
- `null_routed` (Boolean)
- `prefix_length` (String)
- `reverse_lookup` (String)
- `version` (Number)

<a id="nestedatt--ips--ddos"></a>
### Nested Schema for `ips.ddos`

Read-Only:

- `detection_profile` (String)
- `protection_type` (String)



<a id="nestedatt--iso"></a>
### Nested Schema for `iso`

Read-Only:

- `id` (String)
- `name` (String)


<a id="nestedatt--private_network"></a>
### Nested Schema for `private_network`

Read-Only:

- `id` (String)
- `status` (String)
- `subnet` (String)


<a id="nestedatt--resource"></a>
### Nested Schema for `resource`

Read-Only:

- `cpu` (Attributes) Number of cores (see [below for nested schema](#nestedatt--resource--cpu))
- `memory` (Attributes) Total memory in GiB (see [below for nested schema](#nestedatt--resource--memory))
- `private_network_speed` (Attributes) Private network speed in Gbps (see [below for nested schema](#nestedatt--resource--private_network_speed))
- `public_network_speed` (Attributes) Public network speed in Gbps (see [below for nested schema](#nestedatt--resource--public_network_speed))

<a id="nestedatt--resource--cpu"></a>
### Nested Schema for `resource.cpu`

Read-Only:

- `unit` (String)
- `value` (Number)


<a id="nestedatt--resource--memory"></a>
### Nested Schema for `resource.memory`

Read-Only:

- `unit` (String)
- `value` (Number)


<a id="nestedatt--resource--private_network_speed"></a>
### Nested Schema for `resource.private_network_speed`

Read-Only:

- `unit` (String)
- `value` (Number)


<a id="nestedatt--resource--public_network_speed"></a>
### Nested Schema for `resource.public_network_speed`

Read-Only:

- `unit` (String)
- `value` (Number)