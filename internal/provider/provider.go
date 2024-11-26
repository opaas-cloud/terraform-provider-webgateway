package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-webgateway/tools"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &webGatewayProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &webGatewayProvider{
			version: version,
		}
	}
}

// webgatewayProvider is the provider implementation.
type webGatewayProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type webGatewayProviderModel struct {
	HOST  types.String `tfsdk:"host"`
	TOKEN types.String `tfsdk:"token"`
}

// Metadata returns the provider type name.
func (p *webGatewayProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "webgateway"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *webGatewayProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Required: true,
			},
			"token": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (p *webGatewayProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config webGatewayProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	webgateway := tools.WebGateway{
		Host:  config.HOST.ValueString(),
		Token: config.TOKEN.ValueString(),
	}

	resp.DataSourceData = &webgateway
	resp.ResourceData = &webgateway
}

// DataSources defines the data sources implemented in the provider.
func (p *webGatewayProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Resources defines the resources implemented in the provider.
func (p *webGatewayProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDomainResource,
	}
}
