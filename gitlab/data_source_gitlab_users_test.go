package gitlab

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceGitlabUsers_basic(t *testing.T) {
	rInt := acctest.RandInt()
	rInt2 := acctest.RandInt()
	user2 := fmt.Sprintf("user%d@test.test", rInt2)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGitlabUsersConfig(rInt, rInt2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gitlab_user.foo", "name", "footest1"),
					resource.TestCheckResourceAttr("gitlab_user.foo2", "name", "footest2"),
				),
			},
			{
				Config: testAccDataSourceGitlabUsersConfigSort(rInt, rInt2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_users.foo", "users.#", "2"),
					resource.TestCheckResourceAttr("data.gitlab_users.foo", "users.0.email", user2),
					resource.TestCheckResourceAttr("data.gitlab_users.foo", "users.0.projects_limit", "2"),
				),
			},
			{
				Config: testAccDataSourceGitlabUsersConfigSearch(rInt, rInt2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_users.foo", "users.#", "1"),
					// resource.TestCheckResourceAttr("data.gitlab_users.foo", "users.0.email", user2),
				),
			},
			{
				Config: testAccDataSourceGitlabLotsOfUsers(),
			},
			{
				Config: testAccDataSourceGitlabLotsOfUsersSearch(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_users.foo", "users.#", "99"),
				),
			},
			{
				Config: testAccDataSourceGitlabExternalUsers(rInt),
			},
			{
				Config: testAccDataSourceGitlabExternalUsersSearch(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_users.only_external", "users.#", "10"),
				),
			},
		},
	})
}

func testAccDataSourceGitlabUsersConfig(rInt int, rInt2 int) string {
	return fmt.Sprintf(`
resource "gitlab_user" "foo" {
  name             = "footest1"
  username         = "listest%d"
  password         = "test%dtt"
  email            = "user%d@test.test"
  projects_limit   = 3
}

resource "gitlab_user" "foo2" {
  name             = "footest2"
  username         = "listest%d"
  password         = "test%dtt"
  email            = "user%d@test.test"
  projects_limit   = 2
}
	`, rInt, rInt, rInt, rInt2, rInt2, rInt2)
}

func testAccDataSourceGitlabUsersConfigSort(rInt int, rInt2 int) string {
	return fmt.Sprintf(`
resource "gitlab_user" "foo" {
  name             = "footest1"
  username         = "listest%d"
  password         = "test%dtt"
  email            = "user%d@test.test"
  projects_limit   = 3
}

resource "gitlab_user" "foo2" {
  name             = "footest2"
  username         = "listest%d"
  password         = "test%dtt"
  email            = "user%d@test.test"
  projects_limit   = 2
}

data "gitlab_users" "foo" {
  sort = "desc"
  search = "footest"
  order_by = "name"
}
	`, rInt, rInt, rInt, rInt2, rInt2, rInt2)
}

func testAccDataSourceGitlabUsersConfigSearch(rInt int, rInt2 int) string {
	return fmt.Sprintf(`
resource "gitlab_user" "foo" {
  name             = "footest1"
  username         = "listest%d"
  password         = "test%dtt"
  email            = "user%d@test.test"
  projects_limit   = 3
}

resource "gitlab_user" "foo2" {
  name             = "footest2"
  username         = "listest%d"
  password         = "test%dtt"
  email            = "user%d@test.test"
  projects_limit   = 2
}

data "gitlab_users" "foo" {
  search = "user%d@test.test"
}
	`, rInt, rInt, rInt, rInt2, rInt2, rInt2, rInt2)
}

func testAccDataSourceGitlabLotsOfUsers() string {
	return fmt.Sprintf(`
resource "gitlab_user" "foo" {
  name             = format("lots user%%02d", count.index+1)
  username         = format("user%%02d", count.index+1)
  email            = format("user%%02d@example.com", count.index+1)
  password         = "8characters"
  count            = 99
}
`)
}

func testAccDataSourceGitlabLotsOfUsersSearch() string {
	return fmt.Sprintf(`%v
data "gitlab_users" "foo" {
	search = "lots"
}
	`, testAccDataSourceGitlabLotsOfUsers())
}

func testAccDataSourceGitlabExternalUsers(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_user" "external" {
  name             = format("ext user%%02d%d", count.index+1)
  username         = format("ext%%02d%d", count.index+1)
  email            = format("ext%%02d%d@example.com", count.index+1)
  password         = "8characters"
  is_external      = true
  count            = 10
}
resource "gitlab_user" "internal" {
  name             = format("int user%%02d%d", count.index+1)
  username         = format("int%%02d%d", count.index+1)
  email            = format("int%%02d%d@example.com", count.index+1)
  password         = "8characters"
  count            = 10
}`, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccDataSourceGitlabExternalUsersSearch(rInt int) string {
	return fmt.Sprintf(`%v
data "gitlab_users" "only_external" {
	external = true
}
	`, testAccDataSourceGitlabExternalUsers(rInt))
}
