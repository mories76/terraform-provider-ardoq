package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceComponent_basic(t *testing.T) {
	t.Parallel()

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
				Config: testAccResourceComponent_basic(RootWorkspace, componentName),
			},
			{
				ResourceName:      "ardoq_component.my-component",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceComponent_full(t *testing.T) {
	t.Parallel()

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
				Config: testAccResourceComponent_full(RootWorkspace, componentName),
			},
			{
				ResourceName:      "ardoq_component.my-component",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourceComponent_fullUpdate(RootWorkspace, componentName),
			},
			{
				ResourceName:      "ardoq_component.my-component",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceComponent_basic(RootWorkspace, componentName string) string {
	return fmt.Sprintf(`
resource "ardoq_component" "my-component" {
  root_workspace = "%s"
  name = "%s"
  description = "TestAcc"
}
`, RootWorkspace, componentName)
}

func testAccResourceComponent_full(RootWorkspace, componentName string) string {
	return fmt.Sprintf(`
resource "ardoq_component" "my-parent" {
	root_workspace = "%[1]s"
	name = "%[2]s-parent"
	description = "TestAcc"
	}
resource "ardoq_component" "my-component" {
  root_workspace = "%[1]s"
  name = "%[2]s"
  description = "TestAcc"
  parent = ardoq_component.my-parent.id
  fields = {
	  "bewaartermijn" = "heel lang"
  }
}
`, RootWorkspace, componentName)
}

func testAccResourceComponent_fullUpdate(RootWorkspace, componentName string) string {
	return fmt.Sprintf(`
	resource "ardoq_component" "my-parent" {
		root_workspace = "%[1]s"
		name = "%[2]s-parent"
		description = "TestAcc"
		}
	resource "ardoq_component" "my-component" {
	  root_workspace = "%[1]s"
	  name = "%[2]s"
	  description = "TestAcc updated"
	  fields = {
		"bewaartermijn" = "heel lang updated"
	}
	}
`, RootWorkspace, componentName)
}
