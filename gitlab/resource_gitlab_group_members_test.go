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

	resource.Test(t, resource.TestCase{PreCheck: func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabGroupMembersSetupUser(rInt),
			},
			{
				Config:  "provider gitlab {}\n",
				Destroy: true,
			},
			{
				Config: testAccGitlabGroupMembersConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "members.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "members.2667931517.access_level", "owner"),
					resource.TestCheckResourceAttr(resourceName, "members.2667931517.expires_at", ""),
					resource.TestCheckResourceAttr(resourceName, "members.2031542183.access_level", "developer"),
					resource.TestCheckResourceAttr(resourceName, "members.2031542183.expires_at", ""),
				),
			},
			{
				Config: testAccGitlabGroupMembersUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "members.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "members.2667931517.access_level", "owner"),
					resource.TestCheckResourceAttr(resourceName, "members.2667931517.expires_at", ""),
					resource.TestCheckResourceAttr(resourceName, "members.2922300817.access_level", "guest"),
					resource.TestCheckResourceAttr(resourceName, "members.2922300817.expires_at", "2099-01-01"),
				),
			},
			{
				Config: testAccGitlabGroupMembersUpdateConfig2(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "members.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "members.2667931517.access_level", "owner"),
					resource.TestCheckResourceAttr(resourceName, "members.2667931517.expires_at", ""),
					resource.TestCheckResourceAttr(resourceName, "members.3940079224.access_level", "maintainer"),
					resource.TestCheckResourceAttr(resourceName, "members.3940079224.expires_at", ""),
				),
			},
			{
				Config: testAccGitlabGroupMembersConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "members.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "members.2667931517.access_level", "owner"),
					resource.TestCheckResourceAttr(resourceName, "members.2667931517.expires_at", ""),
					resource.TestCheckResourceAttr(resourceName, "members.2031542183.access_level", "developer"),
					resource.TestCheckResourceAttr(resourceName, "members.2031542183.expires_at", ""),
				),
			},
		},
	})
}

func testAccGitlabGroupMembersSetupUser(rInt int) string {
	return fmt.Sprintf(`
// Create a random user to initialize the Gitlab "Ghost User"
resource "gitlab_user" "test-user" {
  name     = "foo%d"
  username = "listest%d"
  password = "test%dtt"
  email    = "listest%d@ssss.com"
}
`, rInt, rInt, rInt, rInt)
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

  members {
    id           = data.gitlab_users.all.users[0].id
    access_level = "owner"
  }

  members {
    // Use the second user which should be the "Ghost User" with a stable id 3 which 
    // is important for hashes used in tests.
    // Note: this user is created to hold all issues authored by users that have
    // since been deleted. This user cannot be removed.
    id = data.gitlab_users.all.users[1].id
    access_level = "developer"
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

  members {
    id           = data.gitlab_users.all.users[0].id
    access_level = "owner"
  }

  members {
    // Use the second user which should be the "Ghost User" with a stable id 3 which 
    // is important for hashes used in tests.
    // Note: this user is created to hold all issues authored by users that have
    // since been deleted. This user cannot be removed.
    id         = data.gitlab_users.all.users[1].id
    access_level = "guest"
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

func testAccGitlabGroupMembersUpdateConfig2(rInt int) string {
	return fmt.Sprintf(`
data "gitlab_users" "all" {
  sort     = "asc"
  search   = ""
  order_by = "id"
}

resource "gitlab_group_members" "test-group-members" {
  group_id       = "${gitlab_group.test-group.id}"

  members {
    id           = data.gitlab_users.all.users[0].id
    access_level = "owner"
  }

  members {
    // Use the second user which should be the "Ghost User" with a stable id 3 which 
    // is important for hashes used in tests.
    // Note: this user is created to hold all issues authored by users that have
    // since been deleted. This user cannot be removed.
    id         = data.gitlab_users.all.users[1].id
    access_level = "maintainer"
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
