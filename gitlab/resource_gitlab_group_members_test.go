package gitlab

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccGitlabGroupMembers_basic(t *testing.T) {
	resourceName := "gitlab_group_members.test-group-members"
	rInt := acctest.RandInt()

	// TODO: add setup step to create first user
	resource.Test(t, resource.TestCase{PreCheck: func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabGroupMembersConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "members.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "members.2031542183.access_level", "developer"),
				),
			},
			{
				Config: testAccGitlabGroupMembersUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "members.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "members.2922300817.access_level", "guest"),
					resource.TestCheckResourceAttr(resourceName, "members.2922300817.expires_at", "2099-01-01"),
				),
			},
			{
				Config: testAccGitlabGroupMembersConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "members.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "members.3534297817.access_level", "developer"),
					resource.TestCheckResourceAttr(resourceName, "members.3534297817.expires_at", ""),
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

resource "gitlab_group_members" "test-group-members" {
  group_id       = "${gitlab_group.test-group.id}"
  group_owner_id = data.gitlab_users.all.users[0].id
  access_level   = "developer"

  members {
    id           = data.gitlab_users.all.users[0].id
    access_level = "owner"
  }

  members {
    // Use the second user which should be the "Ghost User", created to hold all issues authored by
    // users that have since been deleted. This user cannot be removed.
    id = data.gitlab_users.all.users[1].id
  }
}

resource "gitlab_group" "test-group" {
  name             = "bar-name-%d"
  path             = "bar-path-%d"
  description      = "Terraform acceptance tests - group members"
  visibility_level = "public"
}
`, rInt, rInt)
}

func testAccGitlabGroupMembersUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
data "gitlab_users" "all" {
  sort     = "asc"
  search   = ""
  order_by = "id"
}

resource "gitlab_group_members" "test-group-members" {
  group_id       = "${gitlab_group.test-group.id}"
  group_owner_id = data.gitlab_users.all.users[0].id
  access_level   = "guest"

  members {
    id           = data.gitlab_users.all.users[0].id
    access_level = "owner"
  }

  members {
    // Use the second user which should be the "Ghost User", created to hold all issues authored by
    // users that have since been deleted. This user cannot be removed.
    id         = data.gitlab_users.all.users[1].id
    expires_at = "2099-01-01"
  }
}

resource "gitlab_group" "test-group" {
  name             = "bar-name-%d"
  path             = "bar-path-%d"
  description      = "Terraform acceptance tests - group members"
  visibility_level = "public"
}
`, rInt, rInt)
}
