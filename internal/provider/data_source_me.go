package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hwakabh/terraform-provider-mockapis/internal/apiclient"
)

// https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-data-source-read#add-data-source-to-provider
type meDataSource struct {
	client *apiclient.Client
}

// https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-data-source-read#implement-data-source-data-models
type meDataSourceModel struct {
	Name        types.String `tfsdk:"name"`
	Year        types.Int64  `tfsdk:"year"`
	HomepageUrl types.String `tfsdk:"homepage"`
	ApiPath     types.String `tfsdk:"path"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &meDataSource{}
	_ datasource.DataSourceWithConfigure = &meDataSource{}
)

// NewCoffeesDataSource is a helper function to simplify the provider implementation.
func NewMeDataSource() datasource.DataSource {
	return &meDataSource{}
}

// Metadata returns the data source type name.
func (d *meDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_me"
	fmt.Println(resp.TypeName)
}

// Configure adds the provider configured client to the data source.
// https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-data-source-read#implement-data-source-client-functionality
func (d *meDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*apiclient.Client)
	if ok != true {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = client
}

// Schema defines the schema for the data source.
// note: this schema should be equal to target API schema such as REST-APIs
// in our case, we can see its schema by calling https://testapi.io/api/${username}/me on testapi.io
// https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-data-source-read#implement-data-source-schema
func (d *meDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of accounts on testapi.io",
			},
			"year": schema.Int32Attribute{
				Computed:    true,
				Description: "Born year for accounts",
			},
			"homepage": schema.StringAttribute{
				Computed:    true,
				Description: "URL of account homepage",
			},
			"path": schema.StringAttribute{
				Computed:    true,
				Description: "API Path of testapi.io",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
// https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-data-source-read#implement-data-source-data-models
func (d *meDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state meDataSourceModel

	r, err := d.client.GetResponse()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read response from /api/:username/me",
			err.Error(),
		)
		return
	}

	// Map response body to model
	meState := meDataSourceModel{
		Name:        types.StringValue(r.Name),
		Year:        types.Int64Value(int64(r.Year)),
		HomepageUrl: types.StringValue(r.HomepageUrl),
		ApiPath:     types.StringValue(r.ApiPath),
	}

	state.Name = meState.Name
	state.Year = meState.Year
	state.HomepageUrl = meState.HomepageUrl
	state.ApiPath = meState.ApiPath

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
