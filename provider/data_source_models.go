package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ardoq "github.com/mories76/ardoq-client-go/client"
)

func dataSourceArdoqModel() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceModelRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"component_types": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"fields": {
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
			"reference_types": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceArdoqModels() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceModelsRead,
		Schema: map[string]*schema.Schema{
			// "name": &schema.Schema{
			// 	Type:     schema.TypeString,
			// 	Optional: true,
			// },
			"models": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"fields": {
							Type:     schema.TypeMap,
							Computed: true,
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

func dataSourceModelRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(ardoq.Client)
	model_id := d.Get("id").(string)

	model, err := c.Models().Read(ctx, model_id)
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
		// "fields":          model.GetFields(),
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
