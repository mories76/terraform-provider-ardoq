package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ardoq "github.com/mories76/ardoq-client-go/pkg"
)

func resourceArdoqComponent() *schema.Resource {
	return &schema.Resource{
		Description:   "`ardoq_component` resource lets you create a component",
		CreateContext: resourceArdoqComponentCreate,
		ReadContext:   resourceArdoqComponentRead,
		UpdateContext: resourceArdoqComponentUpdate,
		DeleteContext: resourceArdoqComponentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of the component",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Text field describing the component",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"id": {
				Description: "The unique ID of the component",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"parent": {
				Description: "Id of the component's parent",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"type_id": {
				Description: "Id of the component's type",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"root_workspace": {
				Description: "Id of the workspace the component belongs to",
				Type:        schema.TypeString,
				Required:    true,
			},
			"fields": {
				Description: "All custom fields from the model end up here",
				Type:        schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
		},
	}
}

func resourceArdoqComponentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	// var diags diag.Diagnostics

	c := m.(ardoq.Client)

	// get all required fields
	req := ardoq.ComponentRequest{
		Name:          d.Get("name").(string),
		RootWorkspace: d.Get("root_workspace").(string),
	}

	// to get optional fields, first check if thay are set with GetOK, if set the set te request value
	// if not set, the req value will be nill, en therefore left out of the json body ",omitempty"
	if v, ok := d.GetOk("description"); ok {
		req.Description = v.(string)
	}

	if v, ok := d.GetOk("parent"); ok {
		req.Parent = v.(string)
	}

	if v, ok := d.GetOk("type_id"); ok {
		req.TypeID = v.(string)
	}

	// check if custom fields are specified by checking len of the schema field "fields"
	// if so, loop map and add each field to the request
	if len(d.Get("fields").(map[string]interface{})) > 0 {
		//create temporary map
		fields := make(map[string]interface{})

		// foreach key and value append to temporary fields maps
		for k, v := range d.Get("fields").(map[string]interface{}) {
			fields[k] = v.(string)
		}
		// assign temporary fields map to req
		req.Fields = fields
	}

	component, err := c.Components().Create(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(component.ID)

	return resourceArdoqComponentRead(ctx, d, m)
}

func resourceArdoqComponentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(ardoq.Client)

	component, err := c.Components().Read(ctx, d.Id())
	if err != nil {
		// return diag.FromErr(err)
		return handleNotFoundError(err, d, d.Id())
	}

	cmp := flattenComponent(component)

	for key, val := range cmp {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceArdoqComponentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	// var diags diag.Diagnostics

	c := m.(ardoq.Client)
	id := d.Id()

	// create new request for update
	req := ardoq.ComponentRequest{}

	// update field if changes are detected
	if d.HasChange("description") {
		req.Description = d.Get("description").(string)
	}

	if d.HasChange("name") {
		req.Name = d.Get("name").(string)
	}

	if d.HasChange("parent") {
		req.Parent = d.Get("parent").(string)
	}

	if d.HasChange("type_id") {
		req.TypeID = d.Get("type_id").(string)
	}

	// check if custom fields are specified by checking len of the schema field "fields"
	// if so, loop map and add each field to the request
	if len(d.Get("fields").(map[string]interface{})) > 0 {
		//create temporary map
		fields := make(map[string]interface{})

		// foreach key and value append to temporary fields maps
		for k, v := range d.Get("fields").(map[string]interface{}) {
			fields[k] = v.(string)
		}
		// assign temporary fields map to req
		req.Fields = fields
	}

	_, err := c.Components().Update(ctx, id, req)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceArdoqComponentRead(ctx, d, m)
}

func resourceArdoqComponentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(ardoq.Client)
	id := d.Id()

	err := c.Components().Delete(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diag.Diagnostics{}
}
