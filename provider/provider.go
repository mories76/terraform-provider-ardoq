package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ardoq "github.com/mories76/ardoq-client-go/client"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"apikey": &schema.Schema{
				Description: "API key. Can be specified with the `ARDOQ_APIKEY` " +
					"environment variable.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("ARDOQ_APIKEY", nil),
			},
			"baseuri": &schema.Schema{
				Description: "Base URI for the Ardoq API. For example https://mycompany.ardoq.com/api/ Can be specified with the `ARDOQ_BASEURI` " +
					"environment variable.",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARDOQ_BASEURI", nil),
			},
			"org": &schema.Schema{
				Description: "You can specify an organization for your API requests. Can be specified with the `ARDOQ_ORG` " +
					"environment variable.",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARDOQ_ORG", nil),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ardoq_component":  dataSourceArdoqComponent(),
			"ardoq_components": dataSourceArdoqComponents(),
			"ardoq_model":      dataSourceArdoqModel(),
			"ardoq_models":     dataSourceArdoqModels(),
			"ardoq_reference":  dataSourceArdoqReference(),
			"ardoq_references": dataSourceArdoqReferences(),
			"ardoq_workspace":  dataSourceArdoqWorkspace(),
			"ardoq_workspaces": dataSourceArdoqWorkspaces(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"ardoq_component": resourceArdoqComponent(),
			"ardoq_reference": resourceArdoqReference(),
		},
		ConfigureContextFunc: configure,
	}
}

func configure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apikey := d.Get("apikey").(string)
	baseuri := d.Get("baseuri").(string)
	org := d.Get("org").(string)

	var diags diag.Diagnostics

	if (apikey != "") && (baseuri != "") && (org != "") {
		c, err := ardoq.NewRestClient(baseuri, apikey, org)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return c, diags
	}
	// TODO : add diag error
	return nil, diags
}
