package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceComponent(t *testing.T) {

	RootWorkspace := os.Getenv("ARDOQ_WORKSPACE")

	if RootWorkspace == "" {
		t.Skip("ARDOQ_WORKSPACE needs to be set to run this test")
	}

	componentName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceComponent(RootWorkspace, componentName),
				Check:  resource.TestCheckResourceAttr("data.ardoq_component.my-component", "name", componentName),
			},
		},
	})
}

func testAccResourceComponent(RootWorkspace, componentName string) string {
	return fmt.Sprintf(`
resource "ardoq_component" "my-component" {
  root_workspace = "%s"
  name = "%s"
  description = "TestAcc"
}
data "ardoq_component" "my-component" {
  root_workspace = ardoq_component.my-component.root_workspace
  name = ardoq_component.my-component.name
}
`, RootWorkspace, componentName)
}
