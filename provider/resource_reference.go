package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ardoq "github.com/mories76/ardoq-client-go/client"
)

func resourceArdoqReference() *schema.Resource {
	return &schema.Resource{
		Description:   "Arodq references...",
		CreateContext: resourceArdoqReferenceCreate,
		ReadContext:   resourceArdoqReferenceRead,
		UpdateContext: resourceArdoqReferenceUpdate,
		DeleteContext: resourceArdoqReferenceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"root_workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeInt,
				Required: true,
				// Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_text": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceArdoqReferenceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(ardoq.Client)

	// var diags diag.Diagnostics

	req := ardoq.ReferenceRequest{
		RootWorkspace:   d.Get("root_workspace_id").(string),
		Source:          d.Get("source").(string),
		TargetWorkspace: d.Get("target_workspace_id").(string),
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

		// foreach key and value append to temporay fields maps
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
		// return diag.FromErr(err)
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

	// check if custom fields are specified by checking len of the schema field "fields"
	// if so, loop map and add each field to the request
	if len(d.Get("fields").(map[string]interface{})) > 0 {
		//create temporary map
		fields := make(map[string]interface{})

		// foreach key and value append to temporay fields maps
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
