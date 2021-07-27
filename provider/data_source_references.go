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
		ReadContext: dataSourceReferenceRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"source": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"root_workspace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_workspace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_text": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArdoqReferences() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReferencesRead,
		Schema: map[string]*schema.Schema{
			// "name": &schema.Schema{
			// 	Type:     schema.TypeString,
			// 	Optional: true,
			// },
			"references": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// "id": &schema.Schema{
						// 	Type:     schema.TypeString,
						// 	Computed: true,
						// },
						"source": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"root_workspace_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_workspace_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_text": {
							Type:     schema.TypeString,
							Computed: true,
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
		"root_workspace_id":   reference.RootWorkspace,
		"source":              reference.Source,
		"target_workspace_id": reference.TargetWorkspace,
		"target":              reference.Target,
		"type":                reference.Type,
		"description":         reference.Description,
		"display_text":        reference.DisplayText,
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
		// "id":                  reference.ID,
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
