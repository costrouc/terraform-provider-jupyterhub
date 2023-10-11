// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/costrouc/go-jupyterhub-api/api"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure JupyterHubProvider satisfies various provider interfaces.
var _ provider.Provider = &JupyterHubProvider{}

// JupyterHubProvider defines the provider implementation.
type JupyterHubProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// JupyterHubProviderModel describes the provider data model.
type JupyterHubProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Protocol types.String `tfsdk:"protocol"`
	Prefix   types.String `tfsdk:"prefix"`
	Token    types.String `tfsdk:"token"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (p *JupyterHubProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "jupyterhub"
	resp.Version = p.version
}

func (p *JupyterHubProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with JupyterHub.",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "Hostname for JupyterHub API. Default is 'localhost:8000'. May also be provided via JUPYTERHUB_HOST environment variable.",
				Optional:    true,
			},
			"protocol": schema.StringAttribute{
				Description: "Protocol for JupyterHub API. Default is 'http'. May also be provided via JUPYTERHUB_PROTOCOL environment variable.",
				Optional:    true,
			},
			"prefix": schema.StringAttribute{
				Description: "Prefix for JupyterHub API. Default is '/'. May also be provided via JUPYTERHUB_PREFIX environment variable.",
				Optional:    true,
			},
			"token": schema.StringAttribute{
				Description: "API Token for JupyterHub API. Optional if username and password are set. May also be provided via JUPYTERHUB_TOKEN environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"username": schema.StringAttribute{
				Description: "API Token for JupyterHub API. Optional may also be provided via JUPYTERHUB_USERNAME environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"password": schema.StringAttribute{
				Description: "API Token for JupyterHub API. Optional may also be provided via JUPYTERHUB_PASSWORD environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *JupyterHubProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data JupyterHubProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if data.Protocol.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("protocol"),
			"Unknown JupyterHub API Protocol",
			"The provider cannot create the JupyterHub API client as there is an unknown configuration value for the JupyterHub API protocol. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the JUPYTERHUB_PROTOCOL environment variable.",
		)
	}

	if data.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown JupyterHub API Host",
			"The provider cannot create the JupyterHub API client as there is an unknown configuration value for the JupyterHub API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the JUPYTERHUB_HOST environment variable.",
		)
	}

	if data.Prefix.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("prefix"),
			"Unknown JupyterHub API Prefix",
			"The provider cannot create the JupyterHub API client as there is an unknown configuration value for the JupyterHub API prefix. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the JUPYTERHUB_PREFIX environment variable.",
		)
	}

	if data.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown JupyterHub API Token",
			"The provider cannot create the JupyterHub API client as there is an unknown configuration value for the JupyterHub API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the JUPYTERHUB_TOKEN environment variable.",
		)
	}

	if data.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown JupyterHub API Username",
			"The provider cannot create the JupyterHub API client as there is an unknown configuration value for the JupyterHub API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the JUPYTERHUB_USERNAME environment variable.",
		)
	}

	if data.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown JupyterHub API Password",
			"The provider cannot create the JupyterHub API client as there is an unknown configuration value for the JupyterHub API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the JUPYTERHUB_PASSWORD environment variable.",
		)
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	protocol := os.Getenv("HASHICUPS_PROTOCOL")
	host := os.Getenv("HASHICUPS_HOST")
	prefix := os.Getenv("HASHICUPS_PREFIX")
	token := os.Getenv("HASHICUPS_TOKEN")
	username := os.Getenv("HASHICUPS_USERNAME")
	password := os.Getenv("HASHICUPS_PASSWORD")

	if !data.Protocol.IsNull() {
		protocol = data.Protocol.ValueString()
	} else if protocol == "" {
		protocol = "http"
	}

	if !data.Host.IsNull() {
		host = data.Host.ValueString()
	} else if host == "" {
		host = "localhost:8000"
	}

	if !data.Prefix.IsNull() {
		prefix = data.Prefix.ValueString()
	} else if prefix == "" {
		prefix = "/"
	}

	if !data.Token.IsNull() {
		token = data.Token.ValueString()
	}

	if !data.Username.IsNull() {
		username = data.Username.ValueString()
	}

	if !data.Password.IsNull() {
		password = data.Password.ValueString()
	}

	if token == "" && username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing JupyterHub API Token and Username",
			"The provider cannot create the JupyterHub API client as there is a missing or empty value for both JupyterHub API username and token (one is needed). "+
				"Set the username value in the configuration or use the JUPYTERHUB_USERNAME environment variable. "+
				"Set the token value in the configuration or use the JUPYTERHUB_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if token == "" && password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing JupyterHub API Token and Password",
			"The provider cannot create the JupyterHub API client as there is a missing or empty value for both JupyterHub API password and token (one is needed). "+
				"Set the password value in the configuration or use the JUPYTERHUB_PASSWORD environment variable. "+
				"Set the token value in the configuration or use the JUPYTERHUB_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	ctx = tflog.SetField(ctx, "jupyterhub_uri", fmt.Sprintf("%s://%s/%s", protocol, host, prefix))
	tflog.Debug(ctx, "Creating JupyterHub client")

	// Create a new HashiCups client using the configuration values
	client, err := api.CreateClient(&api.ClientConfig{Protocol: protocol, Host: host, Prefix: prefix, Token: token, Username: username, Password: password})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Jupyterhub API Client",
			"An unexpected error occurred when creating the JuypterHub API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"HashiCups Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured HashiCups client", map[string]any{"success": true})
}

func (p *JupyterHubProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *JupyterHubProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewUserDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &JupyterHubProvider{
			version: version,
		}
	}
}
