package provider

import (
	"context"
	"fmt"

	"github.com/costrouc/go-jupyterhub-api/api"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &userDataSource{}
	_ datasource.DataSourceWithConfigure = &userDataSource{}
)

// NewCoffeesDataSource is a helper function to simplify the provider implementation.
func NewUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

// userDataSource is the data source implementation.
type userDataSource struct {
	client *api.ClientConfig
}

// userDataSourceModel maps the data source schema data.
type userDataSourceModel struct {
	Name   types.String   `tfsdk:"name"`
	Admin  types.Bool     `tfsdk:"admin"`
	Roles  []types.String `tfsdk:"roles"`
	Groups []types.String `tfsdk:"groups"`
}

// Configure adds the provider configured client to the data source.
func (d *userDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.ClientConfig)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Metadata returns the data source type name.
func (d *userDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "jupyterhub_user"
}

// Schema defines the schema for the data source.
func (d *userDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Jupyterhub User.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "JupyterHub username.",
				Required:    true,
			},
			"admin": schema.BoolAttribute{
				Description: "User is administrator.",
				Computed:    true,
			},
			"roles": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "Roles assigned to user",
				Computed:    true,
			},
			"groups": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "Groups assigned to user",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state userDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	username := state.Name.String()
	tflog.Info(ctx, fmt.Sprintf("expected username %s", username[1:len(username)-1]))
	user, err := d.client.GetUser(username[1 : len(username)-1])
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read JupyterHub User",
			err.Error(),
		)
		return
	}

	state.Admin = types.BoolValue(user.Admin)

	state.Roles = []types.String{}
	for _, role := range user.Roles {
		state.Roles = append(state.Roles, types.StringValue(role))
	}

	state.Groups = []types.String{}
	for _, group := range user.Groups {
		state.Roles = append(state.Groups, types.StringValue(group))
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
