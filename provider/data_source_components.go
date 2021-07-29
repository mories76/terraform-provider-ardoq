package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ardoq "github.com/mories76/ardoq-client-go/client"
)

func dataSourceArdoqComponent() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(resourceArdoqComponent().Schema)
	addRequiredFieldsToSchema(dsSchema, "root_workspace", "name")

	return &schema.Resource{
		Description: "`ardoq_component` data source can be used to retrieve information for a component by name and workspace.",
		ReadContext: dataSourceArdoqComponentRead,
		Schema:      dsSchema,
	}
}

func dataSourceArdoqComponents() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(resourceArdoqComponent().Schema)

	return &schema.Resource{
		Description: "`ardoq_components` data source can be used to retrieve all components from a specific workspace.",
		ReadContext: dataSourceArdoqComponentsRead,
		Schema: map[string]*schema.Schema{
			"root_workspace": &schema.Schema{
				Description: "Id of the workspace where to retrieve components from",
				Type:        schema.TypeString,
				Required:    true,
			},
			"components": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: dsSchema,
				},
			},
		},
	}
}

func dataSourceArdoqComponentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(ardoq.Client)
	component_name := d.Get("name").(string)
	root_workspace := d.Get("root_workspace").(string)

	components, err := c.Components().Search(ctx, &ardoq.ComponentSearchQuery{Name: component_name, Workspace: root_workspace})
	if err != nil {
		return diag.FromErr(err)
	}
	// TODO: check other datasource/resources for error handling
	if len(*components) != 1 { // check if components result is 1, if 0 then no result was found, if more then 1 was found, the query was not specific enough
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("%d components found, ardoq_component should return 1", len(*components)),
		})
		return diags
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
	root_workspace := d.Get("root_workspace").(string)

	var qry = &ardoq.ComponentSearchQuery{
		Workspace: root_workspace,
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
	// d.Set("root_workspace", component.RootWorkspace)
	// d.Set("name", component.Name)
	// d.Set("description", component.Description)
	// d.Set("fields", component.Fields)
	// d.Set("parent", component.Parent)
	// d.Set("type_id", component.TypeID)

	return map[string]interface{}{
		"root_workspace": component.RootWorkspace,
		"name":           component.Name,
		"description":    component.Description,
		"fields":         component.Fields,
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
