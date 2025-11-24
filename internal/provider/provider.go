package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hwakabh/terraform-provider-mockapis/internal/apiclient"
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
// note:
// - this is a kind of validations for provider configurations
// - with this Configure() methods, we can implement logics to fetch configurations from envars
// https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider-configure#implement-client-configuration-functionality
func (p *mockapisProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config mockapisProviderModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// validations
	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Username for TestAPI.io",
			"The provider cannot create testapi.io client as there is an unknown configuration value for testapi.io username.\n"+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the TESTAPI_USERNAME environment variable.",
		)
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Token for TestAPI.io private endpoint",
			"The provider cannot create testapi.io client as there is an unknown configuration value for testapi.io tokens.\n"+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the TESTAPI_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// fallback to envars
	username := os.Getenv("TESTAPI_USERNAME")
	token := os.Getenv("TESTAPI_TOKEN")
	if config.Username.IsNull() != true {
		username = config.Username.ValueString()
	}
	if config.Token.IsNull() != true {
		token = config.Token.ValueString()
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Username for TestAPI.io is missing",
			"The provider cannot create testapi.io client as there is an unknown configuration value for testapi.io username.\n"+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the TESTAPI_USERNAME environment variable.\n"+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Token for TestAPI.io private endpoint is missing",
			"The provider cannot create testapi.io client as there is an unknown configuration value for testapi.io tokens.\n"+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the TESTAPI_TOKEN environment variable.\n"+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// instantiate TestAPI client for interact with your APIs
	client, err := apiclient.NewClient(username, token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create TestAPI Client",
			"An unexpected error occurred when creating the TestAPI client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"TestAPI Client Error: "+err.Error(),
		)
		return
	}

	// Make the HashiCups client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// Resources defines the resources implemented in the provider.
func (p *mockapisProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

// DataSources defines the data sources implemented in the provider.
// note: this is required even if you would not implement data source for your provider
// https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider#implement-initial-provider-type
func (p *mockapisProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewMeDataSource,
	}
}
