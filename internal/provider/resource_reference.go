package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ardoq "github.com/mories76/ardoq-client-go/pkg"
)

func resourceArdoqReference() *schema.Resource {
	return &schema.Resource{
		Description:   "`ardoq_reference` resource lets you create a reference",
		CreateContext: resourceArdoqReferenceCreate,
		ReadContext:   resourceArdoqReferenceRead,
		UpdateContext: resourceArdoqReferenceUpdate,
		DeleteContext: resourceArdoqReferenceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Description: "Text field describing the reference",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"display_text": {
				Description: "Short label describing the reference, is visible in some visualizations",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"id": {
				Description: "The unique ID of the reference",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"root_workspace": {
				Description: "Id of the source component's workspace",
				Type:        schema.TypeString,
				Required:    true,
			},
			"source": {
				Description: "Id of the source component",
				Type:        schema.TypeString,
				Required:    true,
			},
			"target": {
				Description: "Id of the target component",
				Type:        schema.TypeString,
				Required:    true,
			},
			"target_workspace": {
				Description: "Id of the target component's workspace",
				Type:        schema.TypeString,
				Required:    true,
			},
			"type": {
				Description: "Type (as defined by the model) i.e. Synchronous, Implicit etc.",
				Type:        schema.TypeInt,
				Required:    true,
				// Computed: true,
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

func resourceArdoqReferenceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(ardoq.Client)

	// var diags diag.Diagnostics

	req := ardoq.ReferenceRequest{
		RootWorkspace:   d.Get("root_workspace").(string),
		Source:          d.Get("source").(string),
		TargetWorkspace: d.Get("target_workspace").(string),
		Target:          d.Get("target").(string),
		Type:            d.Get("type").(int),
	}

	if v, ok := d.GetOk("description"); ok {
		req.Description = v.(string)
	}

	if v, ok := d.GetOk("display_text"); ok {
		req.DisplayText = v.(string)
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

	reference, err := c.References().Create(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(reference.ID)

	return resourceArdoqReferenceRead(ctx, d, m)
}

func resourceArdoqReferenceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(ardoq.Client)

	reference, err := c.References().Read(ctx, d.Id())
	if err != nil {
		return handleNotFoundError(err, d, d.Id())
	}

	flatRefence := flattenReference(reference)

	for key, val := range flatRefence {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceArdoqReferenceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	// var diags diag.Diagnostics

	c := m.(ardoq.Client)
	id := d.Id()

	// create new request for update
	req := ardoq.ReferenceRequest{}

	// update field if changes are detected
	if d.HasChange("description") {
		req.Description = d.Get("description").(string)
	}

	if d.HasChange("display_text") {
		req.DisplayText = d.Get("display_text").(string)
	}

	if d.HasChange("type") {
		req.Type = d.Get("type").(int)
	}

	if d.HasChange("source") {
		req.Source = d.Get("source").(string)
	}

	if d.HasChange("target") {
		req.Target = d.Get("target").(string)
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

	_, err := c.References().Update(ctx, id, req)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceArdoqReferenceRead(ctx, d, m)
}

func resourceArdoqReferenceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(ardoq.Client)
	id := d.Id()

	err := c.References().Delete(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diag.Diagnostics{}
}
