package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mailtrap/terraform-provider-mailtrap/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ProjectDataSource{}

func NewProjectDataSource() datasource.DataSource {
	return &ProjectDataSource{}
}

// ProjectDataSource defines the data source implementation.
type ProjectDataSource struct {
	client    *client.Client
	accountID int64
}

// ProjectDataSourceModel describes the data source data model.
type ProjectDataSourceModel struct {
	ID         types.Int64  `tfsdk:"id"`
	AccountID  types.Int64  `tfsdk:"account_id"`
	Name       types.String `tfsdk:"name"`
	ShareLinks types.Object `tfsdk:"share_links"`
}

func (d *ProjectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (d *ProjectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Project data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Project identifier",
				Required:            true,
			},
			"account_id": schema.Int64Attribute{
				MarkdownDescription: "Account ID for the project",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Project name",
				Computed:            true,
			},
			"share_links": schema.SingleNestedAttribute{
				MarkdownDescription: "Share links for the project",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"admin": schema.StringAttribute{
						MarkdownDescription: "Admin share link",
						Computed:            true,
					},
					"viewer": schema.StringAttribute{
						MarkdownDescription: "Viewer share link",
						Computed:            true,
					},
				},
			},
		},
	}
}

func (d *ProjectDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.accountID = providerData.AccountID
}

func (d *ProjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProjectDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Determine account ID
	accountID := d.accountID
	if !data.AccountID.IsNull() && !data.AccountID.IsUnknown() {
		accountID = data.AccountID.ValueInt64()
	}

	if accountID == 0 {
		resp.Diagnostics.AddError(
			"Missing Account ID",
			"Account ID must be provided either in the data source configuration or provider configuration",
		)
		return
	}

	// Get project
	endpoint := fmt.Sprintf("/api/accounts/%d/projects/%d", accountID, data.ID.ValueInt64())
	
	var project client.Project
	err := d.client.Get(endpoint, &project)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read project, got error: %s", err))
		return
	}

	// Update model with response data
	data.AccountID = types.Int64Value(accountID)
	data.Name = types.StringValue(project.Name)

	// Convert share links
	shareLinksObj, diags := types.ObjectValueFrom(ctx, map[string]types.Type{
		"admin":  types.StringType,
		"viewer": types.StringType,
	}, &ShareLinksModel{
		Admin:  types.StringValue(project.ShareLinks.Admin),
		Viewer: types.StringValue(project.ShareLinks.Viewer),
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.ShareLinks = shareLinksObj

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
