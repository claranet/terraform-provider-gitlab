package gitlab

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccGitlabGroupMembers_basic(t *testing.T) {
	resourceName := "gitlab_group_members.test-group-members"
	ownerID := os.Getenv("GITLAB_USER_ID")
	rInt := acctest.RandInt()

	skipIfEnvNotSet(t, "GITLAB_USER_ID")

	resource.Test(t, resource.TestCase{PreCheck: func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabGroupMembersConfig(rInt, ownerID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "members.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "members.1.access_level", "developer"),
				),
			},
			{
				Config: testAccGitlabGroupMembersUpdateConfig(rInt, ownerID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "members.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "members.1.access_level", "guest"),
					resource.TestCheckResourceAttr(resourceName, "members.1.expires_at", "2099-01-01"),
				),
			},
			{
				Config: testAccGitlabGroupMembersConfig(rInt, ownerID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "members.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "members.1.access_level", "developer"),
					resource.TestCheckResourceAttr(resourceName, "members.1.expires_at", ""),
				),
			},
		},
	})
}

func testAccGitlabGroupMembersConfig(rInt int, ownerID string) string {
	return fmt.Sprintf(`
resource "gitlab_group_members" "test-group-members" {
	group_id       = "${gitlab_group.test-group.id}"
	group_owner_id = %s
	access_level   = "developer"

  members {
    id           = %s
    access_level = "owner"
	}

  members {
    id           = "${gitlab_user.test-user.id}"
	}

}

resource "gitlab_group" "test-group" {
  name             = "bar-name-%d"
  path             = "bar-path-%d"
  description      = "Terraform acceptance tests - group members"
  visibility_level = "public"
}

resource "gitlab_user" "test-user" {
  name     = "foo%d"
  username = "listest%d"
  password = "test%dtt"
  email    = "listest%d@ssss.com"
}
`, ownerID, ownerID, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccGitlabGroupMembersUpdateConfig(rInt int, ownerID string) string {
	return fmt.Sprintf(`
resource "gitlab_group_members" "test-group-members" {
	group_id       = "${gitlab_group.test-group.id}"
	group_owner_id = %s
	access_level   = "guest"

  members {
    id           = %s
    access_level = "owner"
	}

	members {
    id = "${gitlab_user.test-user.id}"
    expires_at   = "2099-01-01"
  }

}

resource "gitlab_group" "test-group" {
  name             = "bar-name-%d"
  path             = "bar-path-%d"
  description      = "Terraform acceptance tests - group members"
  visibility_level = "public"
}

resource "gitlab_user" "test-user" {
  name     = "foo%d"
  username = "listest%d"
  password = "test%dtt"
  email    = "listest%d@ssss.com"
}
`, ownerID, ownerID, rInt, rInt, rInt, rInt, rInt, rInt)
}

func skipIfEnvNotSet(t *testing.T, envs ...string) {
	for _, k := range envs {
		if os.Getenv(k) == "" {
			t.Skipf("Environment variable %s is not set", k)
		}
	}
}
