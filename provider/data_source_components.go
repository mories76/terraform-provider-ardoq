package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ardoq "github.com/mories76/ardoq-client-go/client"
)

func dataSourceArdoqComponent() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceArdoqComponentRead,
		Schema: map[string]*schema.Schema{
			"workspace_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"fields": { // this is the place for the custom fields
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceArdoqComponents() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceArdoqComponentsRead,
		Schema: map[string]*schema.Schema{
			"workspace_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"fields": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"components": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// "id": &schema.Schema{
						// 	Type:     schema.TypeString,
						// 	Computed: true,
						// },
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"fields": { // this is the place for the custom/unmapped fields
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"type_id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceArdoqComponentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(ardoq.Client)
	component_name := d.Get("name").(string)
	workspace_id := d.Get("workspace_id").(string)

	components, err := c.Components().Search(ctx, &ardoq.ComponentSearchQuery{Name: component_name, Workspace: workspace_id})

	if err != nil {
		return diag.FromErr(err)
	}

	cmp := flattenComponent(&(*components)[0])

	// loop through map "cmp", and update the schema for each key value pair
	for key, val := range cmp {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId((*components)[0].ID)
	return diags
}

func dataSourceArdoqComponentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(ardoq.Client)
	workspace_id := d.Get("workspace_id").(string)

	var qry = &ardoq.ComponentSearchQuery{
		Workspace: workspace_id,
	}
	components, err := c.Components().Search(ctx, qry)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("components", flattenComponents(components)); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func flattenComponent(component *ardoq.Component) map[string]interface{} {
	return map[string]interface{}{
		"name":        component.Name,
		"description": component.Description,
		"fields":      component.Fields,
		// "fields":  component.GetFields(), //TODO figure something out, that if there are no additional fields. the object "Fields: """ doesn't get added
		"parent":  component.Parent,
		"type_id": component.TypeID,
	}
}

func flattenComponents(components *[]ardoq.Component) []interface{} {
	var result []interface{}
	for _, component := range *components {
		result = append(result, flattenComponent(&component))
	}

	return result
}
