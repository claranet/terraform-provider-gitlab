package gitlab

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccGitlabGroupMembers_basic(t *testing.T) {
	resourceName := "gitlab_group_members.test-group-members"
	userResourceName := "gitlab_user.test-user"
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{PreCheck: func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabGroupMembersConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "members.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "members.2667931517.access_level", "owner"),
					resource.TestCheckResourceAttr(resourceName, "members.2667931517.expires_at", ""),
					testCheckResourceAttrKeyedTypeSet(resourceName, userResourceName, "members", "id", map[string]string{
						"access_level": "developer",
						"expires_at":   "",
					}),
				),
			},
			{
				Config: testAccGitlabGroupMembersUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "members.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "members.2667931517.access_level", "owner"),
					resource.TestCheckResourceAttr(resourceName, "members.2667931517.expires_at", ""),
					testCheckResourceAttrKeyedTypeSet(resourceName, userResourceName, "members", "id", map[string]string{
						"access_level": "guest",
						"expires_at":   "2099-01-01",
					}),
				),
			},
			{
				Config: testAccGitlabGroupMembersUpdateConfig2(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "members.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "members.2667931517.access_level", "owner"),
					resource.TestCheckResourceAttr(resourceName, "members.2667931517.expires_at", ""),
					testCheckResourceAttrKeyedTypeSet(resourceName, userResourceName, "members", "id", map[string]string{
						"access_level": "maintainer",
						"expires_at":   "",
					}),
				),
			},
			{
				Config: testAccGitlabGroupMembersConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "members.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "members.2667931517.access_level", "owner"),
					resource.TestCheckResourceAttr(resourceName, "members.2667931517.expires_at", ""),
					testCheckResourceAttrKeyedTypeSet(resourceName, userResourceName, "members", "id", map[string]string{
						"access_level": "developer",
						"expires_at":   "",
					}),
				),
			},
		},
	})
}

// This custom testCheckResourceAttrKeyedTypeSet function may be reused in other use cases
// where the hash of the set element cannot be determined in advance but the element in the
// set can be identified by a key attribute.
// Here, the user ID will not be stable because users are created/destroyed in undetermined
// orders across test runs, so the hash of the member element will quite always change.
// The new function also to dynamically get the test user's ID
func testCheckResourceAttrKeyedTypeSet(resourceName, keyResourceName, setName, keyAttribute string, values map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ms := s.RootModule()
		keyAttributeValue := ms.Resources[keyResourceName].Primary.Attributes[keyAttribute]
		log.Printf("[DEBUG] Dynamic key attribute value: %#v", keyAttributeValue)

		rs, ok := ms.Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s in %s", resourceName, ms.Path)
		}

		is := rs.Primary
		if is == nil {
			return fmt.Errorf("No primary instance: %s in %s", resourceName, ms.Path)
		}

		for k, v := range is.Attributes {
			addr := strings.Split(k, ".")
			if len(addr) == 3 && addr[0] == setName && addr[2] == keyAttribute {
				prefix := addr[0] + "." + addr[1]
				if v == keyAttributeValue {
					log.Printf("[DEBUG] Found matching %s: %#v, v: %#v", setName, k, v)
					for k2, v2 := range values {
						if is.Attributes[prefix+"."+k2] != v2 {
							return fmt.Errorf("State content does not match for %s resource set %s with key attribute %s in %s", resourceName, setName, keyAttributeValue, ms.Path)
						}
					}
					log.Printf("[DEBUG] Found matching %s attributes for key %s with value %s", setName, k, v)
					return nil
				}
			}
		}

		return fmt.Errorf("State content does not match for %s in %s", resourceName, ms.Path)
	}
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
