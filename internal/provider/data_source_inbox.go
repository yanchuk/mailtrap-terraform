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
var _ datasource.DataSource = &InboxDataSource{}

func NewInboxDataSource() datasource.DataSource {
	return &InboxDataSource{}
}

// InboxDataSource defines the data source implementation.
type InboxDataSource struct {
	client    *client.Client
	accountID int64
}

// InboxDataSourceModel describes the data source data model.
type InboxDataSourceModel struct {
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

func (d *InboxDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inbox"
}

func (d *InboxDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Inbox data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Inbox identifier",
				Required:            true,
			},
			"account_id": schema.Int64Attribute{
				MarkdownDescription: "Account ID for the inbox",
				Optional:            true,
				Computed:            true,
			},
			"project_id": schema.Int64Attribute{
				MarkdownDescription: "Project ID for the inbox",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Inbox name",
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "SMTP username for the inbox",
				Computed:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "SMTP password for the inbox",
				Computed:            true,
				Sensitive:           true,
			},
			"email_username": schema.StringAttribute{
				MarkdownDescription: "Email username part (before @) for the inbox email address",
				Computed:            true,
			},
			"email_username_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether email username is enabled",
				Computed:            true,
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "Domain for SMTP",
				Computed:            true,
			},
			"email_domain": schema.StringAttribute{
				MarkdownDescription: "Email domain",
				Computed:            true,
			},
			"pop3_domain": schema.StringAttribute{
				MarkdownDescription: "POP3 domain",
				Computed:            true,
			},
			"smtp_ports": schema.ListAttribute{
				MarkdownDescription: "Available SMTP ports",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"pop3_ports": schema.ListAttribute{
				MarkdownDescription: "Available POP3 ports",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Inbox status",
				Computed:            true,
			},
			"max_size": schema.Int64Attribute{
				MarkdownDescription: "Maximum inbox size",
				Computed:            true,
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
			},
		},
	}
}

func (d *InboxDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *InboxDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data InboxDataSourceModel

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

	// Get inbox
	endpoint := fmt.Sprintf("/api/accounts/%d/inboxes/%d", accountID, data.ID.ValueInt64())
	
	var inbox client.Inbox
	err := d.client.Get(endpoint, &inbox)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read inbox, got error: %s", err))
		return
	}

	// Update model with response data
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

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
