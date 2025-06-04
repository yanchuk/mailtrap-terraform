package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mailtrap/terraform-provider-mailtrap/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ProjectResource{}
	_ resource.ResourceWithImportState = &ProjectResource{}
)

func NewProjectResource() resource.Resource {
	return &ProjectResource{}
}

// ProjectResource defines the resource implementation.
type ProjectResource struct {
	client    *client.Client
	accountID int64
}

// ProjectResourceModel describes the resource data model.
type ProjectResourceModel struct {
	ID         types.Int64  `tfsdk:"id"`
	AccountID  types.Int64  `tfsdk:"account_id"`
	Name       types.String `tfsdk:"name"`
	ShareLinks types.Object `tfsdk:"share_links"`
}

// ShareLinksModel describes the share links data model
type ShareLinksModel struct {
	Admin  types.String `tfsdk:"admin"`
	Viewer types.String `tfsdk:"viewer"`
}

func (r *ProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *ProjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Project resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Project identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"account_id": schema.Int64Attribute{
				MarkdownDescription: "Account ID for the project",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Project name (min 2 characters, max 100 characters)",
				Required:            true,
			},
			"share_links": schema.SingleNestedAttribute{
				MarkdownDescription: "Share links for the project",
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
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

func (r *ProjectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(*ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = providerData.Client
	r.accountID = providerData.AccountID
}

func (r *ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProjectResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Determine account ID
	accountID := r.accountID
	if !data.AccountID.IsNull() && !data.AccountID.IsUnknown() {
		accountID = data.AccountID.ValueInt64()
	}

	if accountID == 0 {
		resp.Diagnostics.AddError(
			"Missing Account ID",
			"Account ID must be provided either in the resource configuration or provider configuration",
		)
		return
	}

	// Create API request
	createReq := client.ProjectRequest{
		Project: struct {
			Name string `json:"name"`
		}{
			Name: data.Name.ValueString(),
		},
	}

	endpoint := fmt.Sprintf("/api/accounts/%d/projects", accountID)
	
	var project client.Project
	err := r.client.Post(endpoint, createReq, &project)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create project, got error: %s", err))
		return
	}

	// Update model with response data
	data.ID = types.Int64Value(int64(project.ID))
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

	tflog.Trace(ctx, "created a project resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProjectResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get current project state
	endpoint := fmt.Sprintf("/api/accounts/%d/projects/%d", data.AccountID.ValueInt64(), data.ID.ValueInt64())
	
	var project client.Project
	err := r.client.Get(endpoint, &project)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read project, got error: %s", err))
		return
	}

	// Update model with response data
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

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ProjectResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API request
	updateReq := client.ProjectRequest{
		Project: struct {
			Name string `json:"name"`
		}{
			Name: data.Name.ValueString(),
		},
	}

	endpoint := fmt.Sprintf("/api/accounts/%d/projects/%d", data.AccountID.ValueInt64(), data.ID.ValueInt64())
	
	var project client.Project
	err := r.client.Patch(endpoint, updateReq, &project)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update project, got error: %s", err))
		return
	}

	// Update model with response data
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

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProjectResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := fmt.Sprintf("/api/accounts/%d/projects/%d", data.AccountID.ValueInt64(), data.ID.ValueInt64())
	
	err := r.client.Delete(endpoint, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete project, got error: %s", err))
		return
	}
}

func (r *ProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expected format: account_id/project_id
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Incorrect Import ID",
			"Import ID must be in the format: account_id/project_id",
		)
		return
	}

	accountID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Account ID",
			fmt.Sprintf("Could not parse account ID: %s", err),
		)
		return
	}

	projectID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Project ID",
			fmt.Sprintf("Could not parse project ID: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account_id"), accountID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), projectID)...)
}
