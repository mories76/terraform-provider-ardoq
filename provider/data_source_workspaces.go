package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ardoq "github.com/mories76/ardoq-client-go/client"
)

func dataSourceArdoqWorkspace() *schema.Resource {
	return &schema.Resource{
		Description: "TODO", //TODOC
		ReadContext: dataSourceWorkspaceRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Description: "TODO", //TODOC
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": &schema.Schema{
				Description: "TODO", //TODOC
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": &schema.Schema{
				Description: "TODO", //TODOC
				Type:        schema.TypeString,
				Computed:    true,
			},
			"component_model": &schema.Schema{
				Description: "TODO", //TODOC
				Type:        schema.TypeString,
				Computed:    true,
			},
			"fields": {
				Description: "TODO", //TODOC
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceArdoqWorkspaces() *schema.Resource {
	return &schema.Resource{
		Description: "TODO", //TODOC
		ReadContext: dataSourceWorkspacesRead,
		Schema: map[string]*schema.Schema{
			"workspaces": &schema.Schema{
				Description: "TODO", //TODOC
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Description: "TODO", //TODOC
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": &schema.Schema{
							Description: "TODO", //TODOC
							Type:        schema.TypeString,
							Computed:    true,
						},
						"component_model": &schema.Schema{
							Description: "TODO", //TODOC
							Type:        schema.TypeString,
							Computed:    true,
						},
						"description": &schema.Schema{
							Description: "TODO", //TODOC
							Type:        schema.TypeString,
							Computed:    true,
						},
						"fields": {
							Description: "TODO", //TODOC
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
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
		"id":              workspace.ID,
		"name":            workspace.Name,
		"description":     workspace.Description,
		"component_model": workspace.ComponentModel,
	}
}

func flattenWorkspaces(workspaces *[]ardoq.Workspace) []interface{} {
	var result []interface{}
	for _, workspace := range *workspaces {
		result = append(result, flattenWorkspace(&workspace))
	}

	return result
}
