---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "secureworkload_port Resource - terraform-provider-secureworkload"
subcategory: ""
description: |-
  Resource for creating a new service port on Secure Workload
  Example
  An example is shown below:
  hcl
  resource "secureworkload_port" "port1" {
       policy_id = secureworkload_policies.policy1.id
      start_port = 80 
      end_port = 80 
      proto = 6 
  }
  
  Note: If creating multiple rules during a single terraform apply, remember to use depends_on to chain the rules so that terraform creates it in the same order that you intended.
---

# secureworkload_port (Resource)

Resource for creating a new service port on Secure Workload

## Example
An example is shown below: 
```hcl
resource "secureworkload_port" "port1" {
	 policy_id = secureworkload_policies.policy1.id
    start_port = 80 
    end_port = 80 
    proto = 6 
}
```
**Note:** If creating multiple rules during a single `terraform apply`, remember to use `depends_on` to chain the rules so that terraform creates it in the same order that you intended.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `end_port` (Number) End port of the range.
- `policy_id` (String) ID of the needed policy.
- `start_port` (Number) Start port of the range.

### Optional

- `description` (String) (optional) Short string about this proto and port
- `proto` (Number) Protocol Integer value (NULL means all protocols)
- `version` (String) Indicates the version of the workspace the cluster will be added to.

### Read-Only

- `id` (String) The ID of this resource.


