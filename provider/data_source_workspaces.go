package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ardoq "github.com/mories76/ardoq-client-go/client"
)

var workspaceSchema = map[string]*schema.Schema{
	"id": &schema.Schema{
		Description: "The unique ID of the workspace",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"name": &schema.Schema{
		Description: "Name of workspace",
		Type:        schema.TypeString,
		Required:    true,
	},
	"component_model": &schema.Schema{
		Description: "Id of the model the workspace is based on",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"component_template": &schema.Schema{
		Description: "Id of the template the workspace is based on",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"description": &schema.Schema{
		Description: "Text field describing the workspace",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"fields": {
		Description: "All custom fields from the model end up here",
		Type:        schema.TypeMap,
		Computed:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
}

func dataSourceArdoqWorkspace() *schema.Resource {
	return &schema.Resource{
		Description: "`arodq_workspace` data source returns a workspace",
		ReadContext: dataSourceWorkspaceRead,
		Schema:      workspaceSchema,
	}
}

func dataSourceArdoqWorkspaces() *schema.Resource {
	return &schema.Resource{
		Description: "`arodq_workspaces` data source returns all workspaces",
		ReadContext: dataSourceWorkspacesRead,
		Schema: map[string]*schema.Schema{
			"workspaces": &schema.Schema{
				Description: "Ardoq groups documentation into workspaces. A workspace contains all the resources that Ardoq needs to render the textual and visual documentation.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: workspaceSchema,
				},
			},
		},
	}
}

func dataSourceWorkspaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(ardoq.Client)
	workspaceName := d.Get("name").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	workspace, err := c.Workspaces().Search(ctx, &ardoq.WorkspaceSearchQuery{Name: workspaceName})

	if err != nil {
		return diag.FromErr(err)
	}

	flatWorkspace := flattenWorkspace(workspace)

	for key, val := range flatWorkspace {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(workspace.ID)
	return diags
}

func dataSourceWorkspacesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(ardoq.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	workspaces, err := c.Workspaces().List(ctx, &ardoq.WorkspaceSearchQuery{})
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("workspaces", flattenWorkspaces(workspaces)); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func flattenWorkspace(workspace *ardoq.Workspace) map[string]interface{} {
	return map[string]interface{}{
		"id":                 workspace.ID,
		"name":               workspace.Name,
		"component_model":    workspace.ComponentModel,
		"component_template": workspace.ComponentTemplate,
		"description":        workspace.Description,
		"fields":             workspace.Fields,
	}
}

func flattenWorkspaces(workspaces *[]ardoq.Workspace) []interface{} {
	var result []interface{}
	for _, workspace := range *workspaces {
		result = append(result, flattenWorkspace(&workspace))
	}

	return result
}
