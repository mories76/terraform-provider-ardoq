package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ardoq "github.com/mories76/ardoq-client-go/client"
)

func dataSourceArdoqReference() *schema.Resource {
	return &schema.Resource{
		Description: "`arodq_reference` returns a reference",
		ReadContext: dataSourceReferenceRead,
		Schema: map[string]*schema.Schema{
			"description": {
				Description: "Text field describing the reference",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"display_text": {
				Description: "Short label describing the reference, is visible in some visualizations",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"id": &schema.Schema{
				Description: "The unique ID of the reference",
				Type:        schema.TypeString,
				Required:    true,
			},
			"root_workspace": {
				Description: "Id of the source component's workspace",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"source": {
				Description: "Id of the source component",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"target": {
				Description: "Id of the target component",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"target_workspace": {
				Description: "Id of the target component's workspace",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"type": {
				Description: "Type (as defined by the model) i.e. Synchronous, Implicit etc.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func dataSourceArdoqReferences() *schema.Resource {
	return &schema.Resource{
		Description: "`arodq_references` returns all references",
		ReadContext: dataSourceReferencesRead,
		Schema: map[string]*schema.Schema{
			"references": &schema.Schema{
				Description: "References describe relationship between components. References can have types (defined by the model) to represent different kinds of relationship i.e. Synchronized or Asynchroinzed.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Description: "The unique ID of the reference",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"source": {
							Description: "Id of the source component",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"root_workspace": {
							Description: "Id of the source component's workspace",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"target": {
							Description: "Id of the target component",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"target_workspace": {
							Description: "Id of the target component's workspace",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"type": {
							Description: "Type (as defined by the model) i.e. Synchronous, Implicit etc.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"description": {
							Description: "Text field describing the reference",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"display_text": {
							Description: "Short label describing the reference, is visible in some visual",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceReferenceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(ardoq.Client)
	reference_id := d.Get("id").(string)

	reference, err := c.References().Read(ctx, reference_id)

	if err != nil {
		return diag.FromErr(err)
	}

	svc := map[string]interface{}{
		"root_workspace":   reference.RootWorkspace,
		"source":           reference.Source,
		"target_workspace": reference.TargetWorkspace,
		"target":           reference.Target,
		"type":             reference.Type,
		"description":      reference.Description,
		"display_text":     reference.DisplayText,
	}

	for key, val := range svc {
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
		"root_workspace_id":   reference.RootWorkspace,
		"source":              reference.Source,
		"target_workspace_id": reference.TargetWorkspace,
		"target":              reference.Target,
		"type":                reference.Type,
		"description":         reference.Description,
		"display_text":        reference.DisplayText,
	}
}

func flattenReferences(references *[]ardoq.Reference) []interface{} {
	var result []interface{}
	for _, reference := range *references {
		result = append(result, flattenReference(&reference))
	}

	return result
}
