package anypoint

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/tech-nico/terraform-provider-anypoint/anypoint/sdk"
	"testing"
)

func TestAccBusinessGroup_create(t *testing.T) {
	var providers []*schema.Provider

	bgName := fmt.Sprintf("test-bg-create-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	parentPath := "RootOrg/Sub Org 2/Sub Org 2.1"
	bg := sdk.BusinessGroup{}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(&providers),
		CheckDestroy:      testAccCheckWithProviders(testAccCheckBusinessGroupDestroyWithProvider, &providers),
		//IDRefreshName:     "anypoint_business_group.test",
		Steps: []resource.TestStep{
			{
				Config: testAccBusinessGroupConfig_basic(bgName, parentPath),
				Check: resource.ComposeTestCheckFunc(
					testBGExists("ap_bg.test", &bg)),
			},
		},
	})

}
func testBGExists(resourceName string, bg *sdk.BusinessGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ress := s.RootModule().Resources
		for k, v := range ress {
			fmt.Println("k: ", k, " - v: ", v)
		}
		rs, ok := ress[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Business Group ID has been set")
		}

		auth := testAccProvider.Meta().(*Config).AnypointClient.AccessManagement

		bg, err := auth.GetBusinessGroupByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if bg.ID != rs.Primary.ID {
			return fmt.Errorf("Business Group not found")
		}

		return nil
	}
}

func testAccCheckBusinessGroupDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	conn := provider.Meta().(*Config).AnypointClient.AccessManagement

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ap_bg" {
			continue
		}

		// Try to find the resource
		bg, err := conn.GetBusinessGroupByID(rs.Primary.ID)

		//We assume that if err != nil it could be because the bG doesn't exist.
		//Unfortunately the API seems to be returning 401 if the BG does not exist :o|

		if err == nil && bg.ID != "" {
			return fmt.Errorf("Found business group with ID %s (name: %s)", rs.Primary.ID, bg.Name)
		}
	}

	return nil
}

func testAccBusinessGroupConfig_basic(bgName, parentPath string) string {

	return fmt.Sprintf(`
		provider "ap" {
			http_debug_log = true
		}

		resource "ap_bg" "test" {
  			name = "%s"
			parent_path = "%s"
		}
	`, bgName, parentPath)
}
