package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ardoq "github.com/mories76/ardoq-client-go/pkg"
)

var modelSchema = map[string]*schema.Schema{
	"id": {
		Description: "The unique ID of the model",
		Type:        schema.TypeString,
		Required:    true,
	},
	"component_types": {
		Description: "An array of component types and their id's",
		Type:        schema.TypeMap,
		Computed:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	"description": {
		Description: "Text field describing the model,",
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
	"name": {
		Description: "Name of the model",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"reference_types": {
		Description: "An array of reference types and their id's",
		Type:        schema.TypeMap,
		Computed:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
}

func dataSourceArdoqModel() *schema.Resource {
	return &schema.Resource{
		Description: "`ardoq_model` returns a model",
		ReadContext: dataSourceModelRead,
		Schema:      modelSchema,
	}
}

func dataSourceArdoqModels() *schema.Resource {
	return &schema.Resource{
		Description: "`ardoq_models` returns all models",
		ReadContext: dataSourceModelsRead,
		Schema: map[string]*schema.Schema{
			"models": {
				// Description: "TODO", //TODOC
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: modelSchema,
				},
			},
		},
	}
}

func dataSourceModelRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(ardoq.Client)
	modelID := d.Get("id").(string)

	model, err := c.Models().Read(ctx, modelID)
	if err != nil {
		return diag.FromErr(err)
	}

	flatModel := flattenModel(model)

	for key, val := range flatModel {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(model.ID)
	return diags
}

func dataSourceModelsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(ardoq.Client)

	models, err := c.Models().GetAll(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("models", flattenModels(models)); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func flattenModel(model *ardoq.Model) map[string]interface{} {
	return map[string]interface{}{
		"component_types": model.GetComponentTypeID(),
		"description":     model.Description,
		"fields":          model.Fields,
		"name":            model.Name,
		"reference_types": model.GetReferenceTypes(),
	}
}

func flattenModels(models *[]ardoq.Model) []interface{} {
	var result []interface{}
	for _, model := range *models {
		result = append(result, flattenModel(&model))
	}

	return result
}
