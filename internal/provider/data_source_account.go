package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/yanchuk/mailtrap-terraform/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AccountDataSource{}

func NewAccountDataSource() datasource.DataSource {
	return &AccountDataSource{}
}

// AccountDataSource defines the data source implementation.
type AccountDataSource struct {
	client *client.Client
}

// AccountDataSourceModel describes the data source data model.
type AccountDataSourceModel struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (d *AccountDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account"
}

func (d *AccountDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Account data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Account identifier",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Account name",
				Computed:            true,
			},
		},
	}
}

func (d *AccountDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(*ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = providerData.Client
}

func (d *AccountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AccountDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get all accounts
	var accounts []client.Account
	err := d.client.Get("/api/accounts", &accounts)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read accounts, got error: %s", err))
		return
	}

	// Find the account with matching ID
	accountID := data.ID.ValueInt64()
	var foundAccount *client.Account
	for _, account := range accounts {
		if int64(account.ID) == accountID {
			foundAccount = &account
			break
		}
	}

	if foundAccount == nil {
		resp.Diagnostics.AddError(
			"Account Not Found",
			fmt.Sprintf("Account with ID %d not found", accountID),
		)
		return
	}

	// Update model with account data
	data.Name = types.StringValue(foundAccount.Name)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
