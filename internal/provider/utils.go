package provider

import (
	"log"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ardoq "github.com/mories76/ardoq-client-go/pkg"
)

// Check Error Code
func isAPIErrorWithCode(err error, errCode int) bool {
	gerr, ok := errwrap.GetType(err, &ardoq.Error{}).(*ardoq.Error)
	return ok && gerr != nil && gerr.Code == errCode
}

func handleNotFoundError(err error, d *schema.ResourceData, resource string) diag.Diagnostics {
	if isAPIErrorWithCode(err, 404) {
		log.Printf("[WARN] Removing %s because it's gone", resource)
		// The resource doesn't exist anymore
		d.SetId("")

		return nil
	}

	return diag.Errorf("Error when reading or editing %s: %s", resource, err.Error())
}
