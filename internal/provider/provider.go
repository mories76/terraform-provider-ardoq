package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ardoq "github.com/mories76/ardoq-client-go/pkg"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description
		if s.Default != nil {
			desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
		}
		return strings.TrimSpace(desc)
	}
}

// Provider -
func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
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
					Optional:    true,
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
			// ConfigureContextFunc: configure,
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		var diags diag.Diagnostics
		var apikey, baseuri, org string

		// Get apikey
		if v, ok := d.GetOk("apikey"); ok {
			apikey = v.(string)
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "apikey is required",
			})

			return nil, diags
		}

		// Get baseuri
		if v, ok := d.GetOk("baseuri"); ok {
			baseuri = v.(string)
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "baseuri is required",
			})

			return nil, diags
		}

		// Get org
		if v, ok := d.GetOk("org"); ok {
			org = v.(string)
		}

		// create new client
		c, err := ardoq.NewRestClient(baseuri, apikey, org, version)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return c, diags
	}
}
