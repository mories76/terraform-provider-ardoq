package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceReference_basic(t *testing.T) {
	// t.Skip("resource not yet implemented, remove this once you add your own code")

	t.Parallel()

	RootWorkspace := os.Getenv("ARDOQ_WORKSPACE")

	if RootWorkspace == "" {
		t.Skip("ARDOQ_WORKSPACE needs to be set to run this test")
	}

	componentName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceReference_basic(RootWorkspace, componentName),
			},
			{
				ResourceName:      "ardoq_reference.my-reference",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceReference_basic(RootWorkspace, componentName string) string {
	return fmt.Sprintf(`
resource "ardoq_component" "my-component-1" {
	root_workspace = "%[1]s"
	name = "%[2]s"
	description = "TestAcc"
}
resource "ardoq_component" "my-component-2" {
	root_workspace = "%[1]s"
	name = "%[2]s"
	description = "TestAcc"
}
resource "ardoq_reference" "my-reference" {
	source           = ardoq_component.my-component-1.id
	root_workspace   = "%[1]s"
	target           = ardoq_component.my-component-2.id
	target_workspace = "%[1]s"
	type = 2
	description  = "TestAcc"
	display_text = "TestAcc"
}
`, RootWorkspace, componentName)
}
