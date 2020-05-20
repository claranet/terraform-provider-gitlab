package gitlab

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabUserImpersonationToken_basic(t *testing.T) {
	var token gitlab.ImpersonationToken
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGitlabUserImpersonationTokenDestroy,
		Steps: []resource.TestStep{
			// Create a user and impersonation token
			// gitlab_user is already tested as part of `resource_gitlab_user`
			{
				Config: testAccGitlabUserImpersonationTokenConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabUserImpersonationTokenExists("gitlab_user_impersonation_token.bar", &token),
					testAccCheckGitlabUserImpersonationTokenAttributes(&token, &testAccCheckGitlabUserImpersonationTokenExpectedAttributes{
						Name:    fmt.Sprintf("Token bar %d", rInt),
						Active:  true,
						Scopes:  []string{"api"},
						Revoked: false,
					}),
				),
			},
		},
	})
}

func TestAccGitlabUserImpersonationToken_withexpiration(t *testing.T) {
	var token gitlab.ImpersonationToken
	rInt := acctest.RandInt()
	layout := "2006-01-02"
	d, _ := time.Parse(layout, "2222-12-31")
	iso_d := gitlab.ISOTime(d)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGitlabUserImpersonationTokenDestroy,
		Steps: []resource.TestStep{
			// Create a user and impersonation token with expiration
			// gitlab_user is already tested as part of `resource_gitlab_user`
			{
				Config: testAccGitlabUserImpersonationTokenWithExpirationConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabUserImpersonationTokenExists("gitlab_user_impersonation_token.will_expire", &token),
					testAccCheckGitlabUserImpersonationTokenAttributes(&token, &testAccCheckGitlabUserImpersonationTokenExpectedAttributes{
						Name:      fmt.Sprintf("Token will_expire %d", rInt),
						Active:    true,
						Scopes:    []string{"api", "read_user"},
						Revoked:   false,
						ExpiresAt: &iso_d,
					}),
				),
			},
		},
	})
}

func testAccCheckGitlabUserImpersonationTokenExists(n string, token *gitlab.ImpersonationToken) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		tokenID := rs.Primary.ID
		if tokenID == "" {
			return fmt.Errorf("No token ID is set")
		}

		userID := rs.Primary.Attributes["user"]
		if userID == "" {
			return fmt.Errorf("No user ID is set")
		}
		conn := testAccProvider.Meta().(*gitlab.Client)

		token_id, err := strconv.Atoi(tokenID)
		user_id, err := strconv.Atoi(userID)
		gotToken, _, err := conn.Users.GetImpersonationToken(user_id, token_id, nil)
		if err != nil {
			return err
		}
		*token = *gotToken
		return nil
	}
}

func testAccGitlabUserImpersonationTokenDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*gitlab.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_user_impersonation_token" {
			continue
		}

		userID := rs.Primary.Attributes["user"]

		token_id, err := strconv.Atoi(rs.Primary.ID)
		user_id, err := strconv.Atoi(userID)

		token, resp, err := conn.Users.GetImpersonationToken(user_id, token_id, nil)
		if err == nil {
			if token != nil && fmt.Sprintf("%d", token.ID) == rs.Primary.ID && token.Revoked == false {
				return fmt.Errorf("Impersonation token still exists")
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

type testAccCheckGitlabUserImpersonationTokenExpectedAttributes struct {
	Name      string
	Active    bool
	Scopes    []string
	Revoked   bool
	ExpiresAt *gitlab.ISOTime
}

func testAccCheckGitlabUserImpersonationTokenAttributes(token *gitlab.ImpersonationToken, want *testAccCheckGitlabUserImpersonationTokenExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if token.Name != want.Name {
			return fmt.Errorf("got name %q; want %q", token.Name, want.Name)
		}

		if token.Revoked != want.Revoked {
			return fmt.Errorf("got revoked %t; want %t", token.Revoked, want.Revoked)
		}

		if token.ExpiresAt.String() != want.ExpiresAt.String() {
			return fmt.Errorf("got expires %q; want %q", token.ExpiresAt, want.ExpiresAt)
		}

		return nil
	}
}

func testAccGitlabUserImpersonationTokenConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_user" "foo" {
  name             = "foo %d"
  username         = "listest%d"
  password         = "test%dtt"
  email            = "listest%d@ssss.com"
  is_admin         = false
  projects_limit   = 0
  can_create_group = false
  is_external      = false
}

resource "gitlab_user_impersonation_token" "bar" {
    user = gitlab_user.foo.id
    name = "Token bar %d"
    scopes = ["api"]
}
  `, rInt, rInt, rInt, rInt, rInt)
}

func testAccGitlabUserImpersonationTokenWithExpirationConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_user" "foo" {
  name             = "foo %d"
  username         = "listest%d"
  password         = "test%dtt"
  email            = "listest%d@ssss.com"
  is_admin         = false
  projects_limit   = 0
  can_create_group = false
  is_external      = false
}

resource "gitlab_user_impersonation_token" "will_expire" {
    user = gitlab_user.foo.id
    name = "Token will_expire %d"
    scopes = ["api", "read_user"]
    expires_at = "2222-12-31"
}
  `, rInt, rInt, rInt, rInt, rInt)
}
