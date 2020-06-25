package gitlab

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	tfawsresource "github.com/terraform-providers/terraform-provider-gitlab/gitlab/internal/helpers"
)

func TestAccGitlabGroupMembers_basic(t *testing.T) {
	resourceName := "gitlab_group_members.test-group-members"
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{PreCheck: func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabGroupMembersConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "members.#", "2"),
					tfawsresource.TestCheckTypeSetElemNestedAttrs(resourceName, "members.*", map[string]string{
						"access_level": "owner",
						"expires_at":   "",
					}),
					tfawsresource.TestCheckTypeSetElemNestedAttrs(resourceName, "members.*", map[string]string{
						"access_level": "developer",
						"expires_at":   "",
					}),
				),
			},
			{
				Config: testAccGitlabGroupMembersUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "members.#", "2"),
					tfawsresource.TestCheckTypeSetElemNestedAttrs(resourceName, "members.*", map[string]string{
						"access_level": "owner",
						"expires_at":   "",
					}),
					tfawsresource.TestCheckTypeSetElemNestedAttrs(resourceName, "members.*", map[string]string{
						"access_level": "guest",
						"expires_at":   "2099-01-01",
					}),
				),
			},
			{
				Config: testAccGitlabGroupMembersUpdateConfig2(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "members.#", "2"),
					tfawsresource.TestCheckTypeSetElemNestedAttrs(resourceName, "members.*", map[string]string{
						"access_level": "owner",
						"expires_at":   "",
					}),
					tfawsresource.TestCheckTypeSetElemNestedAttrs(resourceName, "members.*", map[string]string{
						"access_level": "maintainer",
						"expires_at":   "",
					}),
				),
			},
			{
				Config: testAccGitlabGroupMembersConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "members.#", "2"),
					tfawsresource.TestCheckTypeSetElemNestedAttrs(resourceName, "members.*", map[string]string{
						"access_level": "owner",
						"expires_at":   "",
					}),
					tfawsresource.TestCheckTypeSetElemNestedAttrs(resourceName, "members.*", map[string]string{
						"access_level": "developer",
						"expires_at":   "",
					}),
				),
			},
		},
	})
}

func testAccGitlabGroupMembersConfig(rInt int) string {
	return fmt.Sprintf(`
data "gitlab_users" "all" {
  sort     = "asc"
  search   = ""
  order_by = "id"
}

resource "gitlab_user" "test-user" {
  name     = "foo%d"
  username = "listest%d"
  password = "test%dtt"
  email    = "listest%d@ssss.com"
}

resource "gitlab_group_members" "test-group-members" {
  group_id = "${gitlab_group.test-group.id}"

  members {
    id           = data.gitlab_users.all.users[0].id
    access_level = "owner"
  }

  members {
    id           = gitlab_user.test-user.id
    access_level = "developer"
  }
}

resource "gitlab_group" "test-group" {
  name             = "bar-name-%d"
  path             = "bar-path-%d"
  description      = "Terraform acceptance tests - group members"
  visibility_level = "public"
}
`, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccGitlabGroupMembersUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
data "gitlab_users" "all" {
  sort     = "asc"
  search   = ""
  order_by = "id"
}

resource "gitlab_user" "test-user" {
  name     = "foo%d"
  username = "listest%d"
  password = "test%dtt"
  email    = "listest%d@ssss.com"
}

resource "gitlab_group_members" "test-group-members" {
  group_id = "${gitlab_group.test-group.id}"

  members {
    id           = data.gitlab_users.all.users[0].id
    access_level = "owner"
  }

  members {
    id           = gitlab_user.test-user.id
    access_level = "guest"
    expires_at   = "2099-01-01"
  }
}

resource "gitlab_group" "test-group" {
  name             = "bar-name-%d"
  path             = "bar-path-%d"
  description      = "Terraform acceptance tests - group members"
  visibility_level = "public"
}
`, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccGitlabGroupMembersUpdateConfig2(rInt int) string {
	return fmt.Sprintf(`
data "gitlab_users" "all" {
  sort     = "asc"
  search   = ""
  order_by = "id"
}

resource "gitlab_user" "test-user" {
  name     = "foo%d"
  username = "listest%d"
  password = "test%dtt"
  email    = "listest%d@ssss.com"
}

resource "gitlab_group_members" "test-group-members" {
  group_id = "${gitlab_group.test-group.id}"

  members {
    id           = data.gitlab_users.all.users[0].id
    access_level = "owner"
  }

  members {
    id           = gitlab_user.test-user.id
    access_level = "maintainer"
  }
}

resource "gitlab_group" "test-group" {
  name             = "bar-name-%d"
  path             = "bar-path-%d"
  description      = "Terraform acceptance tests - group members"
  visibility_level = "public"
}
`, rInt, rInt, rInt, rInt, rInt, rInt)
}
