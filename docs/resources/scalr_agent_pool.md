---
layout: "scalr"
page_title: "Scalr: scalr_agent_pool"
sidebar_current: "docs-resource-scalr-agent-pool"
description: |-
  Manages agent pools.
---

# scalr_agent_pool Resource

Manage the state of agent pools in Scalr. Create, update and destroy.

## Example Usage

Basic usage:

```hcl
resource "scalr_agent_pool" "default" {
  name       = "default-pool"
  account_id = "acc-xxxxxxxx"
}
```

## Argument Reference

* `name` - (Required) Name of the agent pool.
* `account_id` - (Required) ID of the account.
* `environment_id` - (Optional) ID of the environment.

## Attribute Reference

All arguments plus:

* `id` - The ID of the agent pool.

## Import

To import agent pool use agent pool ID as the import ID. For example:
```shell
terraform import scalr_agent_pool.default apool-xxxxxxxxx
```