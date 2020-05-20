---
layout: "gitlab"
page_title: "GitLab: gitlab_user_impersonation_token"
sidebar_current: "docs-gitlab-resource-user_impersonation_token"
description: |-
  Manage impersonation tokens for a user.
---

# gitlab\_user_impersonation_token

This resource allows you to manage impersonation tokens of an existing user.

## Example Usage

```hcl
resource "gitlab_user" "example" {
  name     = "Example Foo"
  username = "example"
  password = "superPassword"
  email    = "gitlab@user.create"
}

resource "gitlab_user_impersonation_token" "my-new-token" {
    user   = gitlab_user.example.id
    name   = "Token bar %d"
    scopes = ["api"]
}
```

## Argument Reference

The following arguments are supported:

* `user` - (Required) The id of user
* `name` - (Required) The name of the token
* `scopes` - (Optional) Array, scopes of the token, can be any `api`, `user_read` or both.
* `expires_at` - (Optional) Expiration date, format is `YYYY-MM-DD`

## Attributes Reference

The resource exports the following attributes:

* `id` The unique id given by the Gitlab server.
* `active` (Boolean) Is the token active or expired
* `revoked` (Boolean) Has the token been revoked
* `created_ad` Time of token creation
