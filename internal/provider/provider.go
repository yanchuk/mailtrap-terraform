package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mailtrap/terraform-provider-mailtrap/internal/client"
)

// Ensure MailtrapProvider satisfies various provider interfaces.
var (
	_ provider.Provider = &MailtrapProvider{}
)

// MailtrapProvider defines the provider implementation.
type MailtrapProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// MailtrapProviderModel describes the provider data model.
type MailtrapProviderModel struct {
	APIToken  types.String `tfsdk:"api_token"`
	AccountID types.Int64  `tfsdk:"account_id"`
}

func (p *MailtrapProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "mailtrap"
	resp.Version = p.version
}

func (p *MailtrapProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				MarkdownDescription: "API token for Mailtrap authentication. Can also be set via MAILTRAP_API_TOKEN environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"account_id": schema.Int64Attribute{
				MarkdownDescription: "Default account ID to use for resources. Can also be set via MAILTRAP_ACCOUNT_ID environment variable.",
				Optional:            true,
			},
		},
	}
}

func (p *MailtrapProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data MailtrapProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Check for API token
	apiToken := os.Getenv("MAILTRAP_API_TOKEN")
	if !data.APIToken.IsNull() {
		apiToken = data.APIToken.ValueString()
	}

	if apiToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Missing Mailtrap API Token",
			"The provider cannot create the Mailtrap API client as there is a missing or empty value for the Mailtrap API token. "+
				"Set the api_token value in the configuration or use the MAILTRAP_API_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
		return
	}

	// Check for account ID
	var accountID int64
	accountIDStr := os.Getenv("MAILTRAP_ACCOUNT_ID")
	if accountIDStr != "" {
		// Parse account ID from environment variable
		var err error
		accountID, err = parseInt64(accountIDStr)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Account ID",
				"The MAILTRAP_ACCOUNT_ID environment variable contains an invalid value: "+accountIDStr,
			)
			return
		}
	}
	if !data.AccountID.IsNull() {
		accountID = data.AccountID.ValueInt64()
	}

	// Create Mailtrap client
	client := client.NewClient(apiToken)

	// Create provider data
	providerData := &ProviderData{
		Client:    client,
		AccountID: accountID,
	}

	// Make provider data available to resources and data sources
	resp.DataSourceData = providerData
	resp.ResourceData = providerData
}

func (p *MailtrapProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProjectResource,
		NewInboxResource,
		NewSendingDomainResource,
	}
}

func (p *MailtrapProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAccountDataSource,
		NewProjectDataSource,
		NewInboxDataSource,
		NewSendingDomainDataSource,
	}
}

// New creates a new instance of the provider
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MailtrapProvider{
			version: version,
		}
	}
}

// ProviderData contains data passed to resources and data sources
type ProviderData struct {
	Client    *client.Client
	AccountID int64
}

// Helper function to parse int64
func parseInt64(s string) (int64, error) {
	var i int64
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}
