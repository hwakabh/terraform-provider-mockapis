package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type mockapisProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// note: this requires import of github.com/hashicorp/terraform-plugin-framework/types
// https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider-configure#implement-provider-data-model
type mockapisProviderModel struct {
	Username types.String `tfsdk:"username"`
	Token    types.String `tfsdk:"token"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &mockapisProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &mockapisProvider{
			version: version,
		}
	}
}

// Metadata returns the provider type name.
func (p *mockapisProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "mockapis"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
// note:
// - this is required for your terraform-provider as API client to interact with target APIs
// - the values of this Schema would be provided in `provider` stanza in .tf files
// https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider-configure#implement-provider-schema
func (p *mockapisProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				Required:    true,
				Description: "Your account name in testapi.io",
			},
			"token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "testapi.io tokens for private endpoint",
			},
		},
	}
}

// Configure prepares a HashiCups API client for data sources and resources.
func (p *mockapisProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

// Resources defines the resources implemented in the provider.
func (p *mockapisProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

// DataSources defines the data sources implemented in the provider.
// note: this is required even if you would not implement data source for your provider
// https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider#implement-initial-provider-type
func (p *mockapisProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}
