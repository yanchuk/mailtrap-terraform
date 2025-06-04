package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mailtrap/terraform-provider-mailtrap/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &InboxResource{}
	_ resource.ResourceWithImportState = &InboxResource{}
)

func NewInboxResource() resource.Resource {
	return &InboxResource{}
}

// InboxResource defines the resource implementation.
type InboxResource struct {
	client    *client.Client
	accountID int64
}

// InboxResourceModel describes the resource data model.
type InboxResourceModel struct {
	ID                      types.Int64  `tfsdk:"id"`
	AccountID               types.Int64  `tfsdk:"account_id"`
	ProjectID               types.Int64  `tfsdk:"project_id"`
	Name                    types.String `tfsdk:"name"`
	Username                types.String `tfsdk:"username"`
	Password                types.String `tfsdk:"password"`
	EmailUsername           types.String `tfsdk:"email_username"`
	EmailUsernameEnabled    types.Bool   `tfsdk:"email_username_enabled"`
	Domain                  types.String `tfsdk:"domain"`
	EmailDomain             types.String `tfsdk:"email_domain"`
	POP3Domain              types.String `tfsdk:"pop3_domain"`
	SMTPPorts               types.List   `tfsdk:"smtp_ports"`
	POP3Ports               types.List   `tfsdk:"pop3_ports"`
	Status                  types.String `tfsdk:"status"`
	MaxSize                 types.Int64  `tfsdk:"max_size"`
	SentMessagesCount       types.Int64  `tfsdk:"sent_messages_count"`
	ForwardedMessagesCount  types.Int64  `tfsdk:"forwarded_messages_count"`
	ForwardFromEmailAddress types.String `tfsdk:"forward_from_email_address"`
}

func (r *InboxResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inbox"
}

func (r *InboxResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Inbox resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Inbox identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"account_id": schema.Int64Attribute{
				MarkdownDescription: "Account ID for the inbox",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.Int64Attribute{
				MarkdownDescription: "Project ID for the inbox",
				Required:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Inbox name",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "SMTP username for the inbox",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "SMTP password for the inbox",
				Computed:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email_username": schema.StringAttribute{
				MarkdownDescription: "Email username part (before @) for the inbox email address",
				Optional:            true,
				Computed:            true,
			},
			"email_username_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether email username is enabled",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "Domain for SMTP",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email_domain": schema.StringAttribute{
				MarkdownDescription: "Email domain",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"pop3_domain": schema.StringAttribute{
				MarkdownDescription: "POP3 domain",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"smtp_ports": schema.ListAttribute{
				MarkdownDescription: "Available SMTP ports",
				Computed:            true,
				ElementType:         types.Int64Type,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"pop3_ports": schema.ListAttribute{
				MarkdownDescription: "Available POP3 ports",
				Computed:            true,
				ElementType:         types.Int64Type,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Inbox status",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"max_size": schema.Int64Attribute{
				MarkdownDescription: "Maximum inbox size",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"sent_messages_count": schema.Int64Attribute{
				MarkdownDescription: "Number of sent messages",
				Computed:            true,
			},
			"forwarded_messages_count": schema.Int64Attribute{
				MarkdownDescription: "Number of forwarded messages",
				Computed:            true,
			},
			"forward_from_email_address": schema.StringAttribute{
				MarkdownDescription: "Email address used for forwarding",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *InboxResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InboxResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data InboxResourceModel

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
	createReq := client.InboxRequest{
		Inbox: struct {
			Name          string `json:"name"`
			EmailUsername string `json:"email_username,omitempty"`
		}{
			Name: data.Name.ValueString(),
		},
	}

	// Add email username if provided
	if !data.EmailUsername.IsNull() && !data.EmailUsername.IsUnknown() {
		createReq.Inbox.EmailUsername = data.EmailUsername.ValueString()
	}

	endpoint := fmt.Sprintf("/api/accounts/%d/projects/%d/inboxes", accountID, data.ProjectID.ValueInt64())
	
	var inbox client.Inbox
	err := r.client.Post(endpoint, createReq, &inbox)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create inbox, got error: %s", err))
		return
	}

	// Update model with response data
	r.updateModelFromInbox(ctx, &data, &inbox, accountID)

	tflog.Trace(ctx, "created an inbox resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InboxResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data InboxResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get current inbox state
	endpoint := fmt.Sprintf("/api/accounts/%d/inboxes/%d", data.AccountID.ValueInt64(), data.ID.ValueInt64())
	
	var inbox client.Inbox
	err := r.client.Get(endpoint, &inbox)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read inbox, got error: %s", err))
		return
	}

	// Update model with response data
	r.updateModelFromInbox(ctx, &data, &inbox, data.AccountID.ValueInt64())

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InboxResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data InboxResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API request
	updateReq := client.InboxRequest{
		Inbox: struct {
			Name          string `json:"name"`
			EmailUsername string `json:"email_username,omitempty"`
		}{
			Name: data.Name.ValueString(),
		},
	}

	// Add email username if provided
	if !data.EmailUsername.IsNull() && !data.EmailUsername.IsUnknown() {
		updateReq.Inbox.EmailUsername = data.EmailUsername.ValueString()
	}

	endpoint := fmt.Sprintf("/api/accounts/%d/inboxes/%d", data.AccountID.ValueInt64(), data.ID.ValueInt64())
	
	var inbox client.Inbox
	err := r.client.Patch(endpoint, updateReq, &inbox)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update inbox, got error: %s", err))
		return
	}

	// Update model with response data
	r.updateModelFromInbox(ctx, &data, &inbox, data.AccountID.ValueInt64())

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InboxResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InboxResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := fmt.Sprintf("/api/accounts/%d/inboxes/%d", data.AccountID.ValueInt64(), data.ID.ValueInt64())
	
	err := r.client.Delete(endpoint, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete inbox, got error: %s", err))
		return
	}
}

func (r *InboxResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expected format: account_id/inbox_id
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Incorrect Import ID",
			"Import ID must be in the format: account_id/inbox_id",
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

	inboxID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Inbox ID",
			fmt.Sprintf("Could not parse inbox ID: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account_id"), accountID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), inboxID)...)
}

// Helper function to update model from inbox data
func (r *InboxResource) updateModelFromInbox(ctx context.Context, data *InboxResourceModel, inbox *client.Inbox, accountID int64) {
	data.ID = types.Int64Value(int64(inbox.ID))
	data.AccountID = types.Int64Value(accountID)
	data.ProjectID = types.Int64Value(int64(inbox.ProjectID))
	data.Name = types.StringValue(inbox.Name)
	data.Username = types.StringValue(inbox.Username)
	data.Password = types.StringValue(inbox.Password)
	data.EmailUsername = types.StringValue(inbox.EmailUsername)
	data.EmailUsernameEnabled = types.BoolValue(inbox.EmailUsernameEnabled)
	data.Domain = types.StringValue(inbox.Domain)
	data.EmailDomain = types.StringValue(inbox.EmailDomain)
	data.POP3Domain = types.StringValue(inbox.POP3Domain)
	data.Status = types.StringValue(inbox.Status)
	data.MaxSize = types.Int64Value(int64(inbox.MaxSize))
	data.SentMessagesCount = types.Int64Value(int64(inbox.SentMessagesCount))
	data.ForwardedMessagesCount = types.Int64Value(int64(inbox.ForwardedMessagesCount))
	data.ForwardFromEmailAddress = types.StringValue(inbox.ForwardFromEmailAddress)

	// Convert SMTP ports
	smtpPorts := make([]types.Int64, len(inbox.SMTPPorts))
	for i, port := range inbox.SMTPPorts {
		smtpPorts[i] = types.Int64Value(int64(port))
	}
	data.SMTPPorts, _ = types.ListValueFrom(ctx, types.Int64Type, smtpPorts)

	// Convert POP3 ports
	pop3Ports := make([]types.Int64, len(inbox.POP3Ports))
	for i, port := range inbox.POP3Ports {
		pop3Ports[i] = types.Int64Value(int64(port))
	}
	data.POP3Ports, _ = types.ListValueFrom(ctx, types.Int64Type, pop3Ports)
}
