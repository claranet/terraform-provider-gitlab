<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Terraform Provider for Gitlab
=============================

- [Documentation](https://www.terraform.io/docs/providers/gitlab/index.html)
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)
- Build status:
  - ![Unit Tests](https://github.com/terraform-providers/terraform-provider-gitlab/workflows/Unit%20Tests/badge.svg?branch=master)
  - ![Acceptance Tests](https://github.com/terraform-providers/terraform-provider-gitlab/workflows/Acceptance%20Tests/badge.svg?branch=master)
  - ![Website Build](https://github.com/terraform-providers/terraform-provider-gitlab/workflows/Website%20Build/badge.svg?branch=master)

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) >= 1.14 (to build the provider plugin)

Building The Provider
---------------------

This is a fork of the official Terraform Provider for Gitlab, with additional changes from Claranet.

Most notably, this fork:

- includes not yet released changes from the official provider.
- adds the `gitlab_group_members` resource.
- has its `master` branch on par with the upstream `master`, with no additional changes.
- has its default branch set to `claranet`, which is where custom development happens.

Unofficial tags made by Claranet are suffixed `-claranet`, and eventually an increment.
They are based on the most recent upstream tag name on which the `claranet` branch was last rebased.

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-gitlab`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:terraform-providers/terraform-provider-gitlab
$ cd terraform-provider-gitlab
$ git remote rename origin upstream
$ git remote add origin git@github.com:claranet/terraform-provider-gitlab
```

Note: this is done in that order so that master tracks the upstream one.

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-gitlab
$ make build
```

Rebasing The Provider
---------------------

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-gitlab
$ git checkout master
$ git pull
$ git push origin master
$ git checkout claranet
$ git rebase upstream/master
```


Using the provider
----------------------
## Fill in for each provider

# Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.14+ is *required*).

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-gitlab
...
```

### Running tests

The Terraform Provider only has acceptance tests, these can run against a gitlab instance where you have a token with administrator permissions (likely not gitlab.com).
There is excellent documentation on [how to run gitlab from docker at gitlab.com](https://docs.gitlab.com/omnibus/docker/)

In order to run the full suite of acceptance tests, export the environment variables: 

- `GITLAB_TOKEN` //token for account with admin priviliges
- `GITLAB_BASE_URL` //URL with api part e.g. `http://localhost:8929/api/v4/`
- `GITLAB_USER_ID` //ID of the user running tests. Needed for `resource_gitlab_group_members` tests

and run `make testacc`.

```sh
$ make testacc
```

### Gitlab Community Edition and Gitlab Entreprise Edition

This module supports both Gitlab CE and Gitlab EE. We run tests on Gitlab EE,
but can't run them on pull requests from forks.

Features that only work on one flavour can use the following helpers as
SkipFunc: `isRunningInEE` and `isRunningInCE`. You can see an exemple of this
for [gitlab_project_push_rules](gitlab/resource_gitlab_project_push_rules_test.go)
tests.
