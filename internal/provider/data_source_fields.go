package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ardoq "github.com/mories76/ardoq-client-go/pkg"
)

var fieldSchema = map[string]*schema.Schema{
	"component_type": {
		Description: "An array of component types and their id's",
		Type:        schema.TypeList,
		Computed:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	"created": {
		Description: "Created",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"created_by": {
		Description: "Created by",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"created_by_email": {
		Description: "Created by email",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"created_by_name": {
		Description: "Created by name",
		Type:        schema.TypeString,
		Computed:    true,
	},
	// DateTimeFields
	"default_value": {
		Description: "Default value",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"description": {
		Description: "Text field describing the field",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"global": {
		Description: "Global",
		Type:        schema.TypeBool,
		Computed:    true,
	},
	"global_ref": {
		Description: "Global ref",
		Type:        schema.TypeBool,
		Computed:    true,
	},
	"id": {
		Description: "The unique ID of the field",
		Type:        schema.TypeString,
		Required:    true,
	},
	"label": {
		Description: "Label",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"last_modified_by": {
		Description: "Last modified by",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"last_modified_by_email": {
		Description: "Last modified by email",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"last_modified_by_name": {
		Description: "Last modified by name",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"last_updated": {
		Description: "Last Updated",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"name": {
		Description: "Name of the field",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"order": {
		Description: "Order",
		Type:        schema.TypeFloat,
		Computed:    true,
	},
	// Origin
	"reference_type": {
		Description: "An array of reference types",
		Type:        schema.TypeList,
		Computed:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	"type": {
		Description: "Type of the field",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"version": {
		Description: "Version of the field",
		Type:        schema.TypeInt,
		Computed:    true,
	},
	"fields": {
		Description: "All custom fields from the field end up here",
		Type:        schema.TypeMap,
		Computed:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
}

func dataSourceArdoqField() *schema.Resource {
	return &schema.Resource{
		Description: "`ardoq_field` returns a field",
		ReadContext: dataSourceFieldRead,
		Schema:      fieldSchema,
	}
}

func dataSourceArdoqFields() *schema.Resource {
	return &schema.Resource{
		Description: "`ardoq_fields` returns all fields",
		ReadContext: dataSourceFieldsRead,
		Schema: map[string]*schema.Schema{
			"fields": {
				// Description: "TODO", //TODOC
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: fieldSchema,
				},
			},
		},
	}
}

func dataSourceFieldRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(ardoq.Client)
	fieldID := d.Get("id").(string)

	field, err := c.Fields().Read(ctx, fieldID)
	if err != nil {
		return diag.FromErr(err)
	}

	flatField := flattenField(field)

	for key, val := range flatField {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(field.ID)
	return diags
}

func dataSourceFieldsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(ardoq.Client)

	fields, err := c.Fields().GetAll(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("fields", flattenFields(fields)); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func flattenField(field *ardoq.Field) map[string]interface{} {
	return map[string]interface{}{
		// DateTimeFields
		// Origin
		"component_type":         field.ComponentType,
		"created":                field.Created,
		"created_by":             field.CreatedBy,
		"created_by_email":       field.CreatedByEmail,
		"created_by_name":        field.CreatedByName,
		"default_value":          field.DefaultValue,
		"description":            field.Description,
		"global":                 field.Global,
		"global_ref":             field.GlobalRef,
		"id":                     field.ID,
		"label":                  field.Label,
		"last_modified_by":       field.LastModifiedBy,
		"last_modified_by_email": field.LastModifiedByEmail,
		"last_modified_by_name":  field.LastModifiedByName,
		"last_updated":           field.LastUpdated,
		"name":                   field.Name,
		"order":                  field.Order,
		"reference_type":         field.ReferenceType,
		"type":                   field.Type,
		"version":                field.Version,
		"fields":                 field.Fields,
	}
}

func flattenFields(fields *[]ardoq.Field) []interface{} {
	var result []interface{}
	for _, field := range *fields {
		result = append(result, flattenField(&field))
	}

	return result
}
