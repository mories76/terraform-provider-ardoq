package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ardoq "github.com/mories76/ardoq-client-go/pkg"
)

func dataSourceArdoqReference() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(resourceArdoqReference().Schema)
	addRequiredFieldsToSchema(dsSchema, "id")

	return &schema.Resource{
		Description: "`arodq_reference` returns a reference",
		ReadContext: dataSourceReferenceRead,
		Schema:      dsSchema,
	}
}

func dataSourceArdoqReferences() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(resourceArdoqReference().Schema)

	return &schema.Resource{
		Description: "`arodq_references` returns all references",
		ReadContext: dataSourceReferencesRead,
		Schema: map[string]*schema.Schema{
			"references": {
				Description: "References describe relationship between components. References can have types (defined by the model) to represent different kinds of relationship i.e. Synchronized or Asynchroinzed.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dsSchema,
				},
			},
		},
	}
}

func dataSourceReferenceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(ardoq.Client)
	referenceID := d.Get("id").(string)

	reference, err := c.References().Read(ctx, referenceID)
	if err != nil {
		// return diag.FromErr(err)
		return handleNotFoundError(err, d, d.Id())
	}

	ref := flattenReference(reference)

	for key, val := range ref {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(reference.ID)
	return diags
}

func dataSourceReferencesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(ardoq.Client)

	references, err := c.References().GetAll(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("references", flattenReferences(references)); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func flattenReference(reference *ardoq.Reference) map[string]interface{} {
	return map[string]interface{}{
		"root_workspace":   reference.RootWorkspace,
		"source":           reference.Source,
		"target_workspace": reference.TargetWorkspace,
		"target":           reference.Target,
		"type":             reference.Type,
		"description":      reference.Description,
		"display_text":     reference.DisplayText,
		"id":               reference.ID,
	}
}

func flattenReferences(references *[]ardoq.Reference) []interface{} {
	var result []interface{}
	for _, reference := range *references {
		result = append(result, flattenReference(&reference))
	}

	return result
}
