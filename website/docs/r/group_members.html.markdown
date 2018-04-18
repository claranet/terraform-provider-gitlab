---
layout: "gitlab"
page_title: "GitLab: gitlab_group_members"
sidebar_current: "docs-gitlab-resource-group_members"
description: |-
  Manage members of a group.
---

# gitlab\_group_members

This resource allows you to manage users of an existing group.

## Example Usage

```hcl
resource "gitlab_group_members" "my-group-members" {
  group_id       = "${gitlab_group.mygroup.id}"
  group_owner_id = 1
  access_level   = "developer"

  members {
    id           = 1
    access_level = "owner"
  }

  members {
    id = "29"
  }

  members {
    id           = "40"
    access_level = "guest"
    expires_at   = "2019-05-01"
  }
}
```

List syntax is **NO LONGER WORKING** since 0.12 Terraform release. Ex:
```hcl
members = [
  {
    id = 1
    access_level = "owner"
  },
  {
    id = "29"
  }
]
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required) The id of the group.

* `group_owner_id` - (Required) The id of the group owner. Necessary since the group owner is aumatically part of the group.

* `access_level` - (Optional) Default access level applied to all members. Acceptable values are: guest, reporter, developer, maintainer, master (deprecated) and owner.

* `expires_at` - (Optional) Default expiration date applied to all group members. Format: `YYYY-MM-DD`.

* `members` - (Required) List of members that are part of the group.

  * `id` - (Required) The id of the user.

  * `access_level` - (Optional) Access level applied to a single member. Acceptable values are: guest, reporter, developer, maintainer, master (deprecated) and owner.

  * `expires_at` - (Optional) Expiration date applied to a single member. Format: `YYYY-MM-DD`.

## Import

GitLab group members can be imported using the group's id.

```
$ terraform import gitlab_group_members.test 100
```
