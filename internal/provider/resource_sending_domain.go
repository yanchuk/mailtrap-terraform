package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mailtrap/terraform-provider-mailtrap/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &SendingDomainResource{}
	_ resource.ResourceWithImportState = &SendingDomainResource{}
)

func NewSendingDomainResource() resource.Resource {
	return &SendingDomainResource{}
}

// SendingDomainResource defines the resource implementation.
type SendingDomainResource struct {
	client    *client.Client
	accountID int64
}

// SendingDomainResourceModel describes the resource data model.
type SendingDomainResourceModel struct {
	ID               types.Int64  `tfsdk:"id"`
	AccountID        types.Int64  `tfsdk:"account_id"`
	Name             types.String `tfsdk:"name"`
	CNAME            types.String `tfsdk:"cname"`
	Status           types.String `tfsdk:"status"`
	ComplianceStatus types.String `tfsdk:"compliance_status"`
	DNSRecords       types.Object `tfsdk:"dns_records"`
	DNSStatus        types.Object `tfsdk:"dns_status"`
}

// DNSRecordsModel describes the DNS records structure
type DNSRecordsModel struct {
	CNAME types.List `tfsdk:"cname"`
	MX    types.List `tfsdk:"mx"`
	TXT   types.List `tfsdk:"txt"`
}

// DNSRecordModel describes a single DNS record
type DNSRecordModel struct {
	Priority   types.Int64  `tfsdk:"priority"`
	RecordType types.String `tfsdk:"record_type"`
	Hostname   types.String `tfsdk:"hostname"`
	Value      types.String `tfsdk:"value"`
	Status     types.String `tfsdk:"status"`
}

// DNSStatusModel describes the DNS status
type DNSStatusModel struct {
	CNAME types.Bool `tfsdk:"cname"`
	MX    types.Bool `tfsdk:"mx"`
	TXT   types.Bool `tfsdk:"txt"`
}

func (r *SendingDomainResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sending_domain"
}

func (r *SendingDomainResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Sending domain resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Sending domain identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"account_id": schema.Int64Attribute{
				MarkdownDescription: "Account ID for the sending domain",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Domain name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cname": schema.StringAttribute{
				MarkdownDescription: "CNAME value for domain verification",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Domain status",
				Computed:            true,
			},
			"compliance_status": schema.StringAttribute{
				MarkdownDescription: "Compliance status",
				Computed:            true,
			},
			"dns_records": schema.SingleNestedAttribute{
				MarkdownDescription: "DNS records for domain verification",
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"cname": schema.ListNestedAttribute{
						MarkdownDescription: "CNAME records",
						Computed:            true,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"priority": schema.Int64Attribute{
									MarkdownDescription: "Record priority",
									Computed:            true,
								},
								"record_type": schema.StringAttribute{
									MarkdownDescription: "Record type",
									Computed:            true,
								},
								"hostname": schema.StringAttribute{
									MarkdownDescription: "Hostname",
									Computed:            true,
								},
								"value": schema.StringAttribute{
									MarkdownDescription: "Record value",
									Computed:            true,
								},
								"status": schema.StringAttribute{
									MarkdownDescription: "Record status",
									Computed:            true,
								},
							},
						},
					},
					"mx": schema.ListNestedAttribute{
						MarkdownDescription: "MX records",
						Computed:            true,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"priority": schema.Int64Attribute{
									MarkdownDescription: "Record priority",
									Computed:            true,
								},
								"record_type": schema.StringAttribute{
									MarkdownDescription: "Record type",
									Computed:            true,
								},
								"hostname": schema.StringAttribute{
									MarkdownDescription: "Hostname",
									Computed:            true,
								},
								"value": schema.StringAttribute{
									MarkdownDescription: "Record value",
									Computed:            true,
								},
								"status": schema.StringAttribute{
									MarkdownDescription: "Record status",
									Computed:            true,
								},
							},
						},
					},
					"txt": schema.ListNestedAttribute{
						MarkdownDescription: "TXT records",
						Computed:            true,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"priority": schema.Int64Attribute{
									MarkdownDescription: "Record priority",
									Computed:            true,
								},
								"record_type": schema.StringAttribute{
									MarkdownDescription: "Record type",
									Computed:            true,
								},
								"hostname": schema.StringAttribute{
									MarkdownDescription: "Hostname",
									Computed:            true,
								},
								"value": schema.StringAttribute{
									MarkdownDescription: "Record value",
									Computed:            true,
								},
								"status": schema.StringAttribute{
									MarkdownDescription: "Record status",
									Computed:            true,
								},
							},
						},
					},
				},
			},
			"dns_status": schema.SingleNestedAttribute{
				MarkdownDescription: "DNS verification status",
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"cname": schema.BoolAttribute{
						MarkdownDescription: "CNAME verification status",
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"mx": schema.BoolAttribute{
						MarkdownDescription: "MX verification status",
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"txt": schema.BoolAttribute{
						MarkdownDescription: "TXT verification status",
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
		},
	}
}

func (r *SendingDomainResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SendingDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SendingDomainResourceModel

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
	createReq := client.SendingDomainRequest{
		SendingDomain: struct {
			DomainName string `json:"domain_name"`
		}{
			DomainName: data.Name.ValueString(),
		},
	}

	endpoint := fmt.Sprintf("/api/accounts/%d/sending_domains", accountID)
	
	var domain client.SendingDomain
	err := r.client.Post(endpoint, createReq, &domain)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create sending domain, got error: %s", err))
		return
	}

	// Update model with response data
	data.ID = types.Int64Value(int64(domain.ID))
	data.AccountID = types.Int64Value(accountID)
	data.Name = types.StringValue(domain.Name)
	data.CNAME = types.StringValue(domain.CNAME)
	data.Status = types.StringValue(domain.Status)
	data.ComplianceStatus = types.StringValue(domain.ComplianceStatus)

	// Convert DNS records
	dnsRecordsValue, diags := r.convertDNSRecordsToTerraform(ctx, &domain.DNSRecords)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.DNSRecords = dnsRecordsValue

	// Convert DNS status
	dnsStatusObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
		"cname": types.BoolType,
		"mx":    types.BoolType,
		"txt":   types.BoolType,
	}, &DNSStatusModel{
		CNAME: types.BoolValue(domain.DNSStatus.CNAME),
		MX:    types.BoolValue(domain.DNSStatus.MX),
		TXT:   types.BoolValue(domain.DNSStatus.TXT),
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.DNSStatus = dnsStatusObj

	tflog.Trace(ctx, "created a sending domain resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SendingDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SendingDomainResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get current domain state
	endpoint := fmt.Sprintf("/api/accounts/%d/sending_domains/%d", data.AccountID.ValueInt64(), data.ID.ValueInt64())
	
	var domain client.SendingDomain
	err := r.client.Get(endpoint, &domain)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read sending domain, got error: %s", err))
		return
	}

	// Update model with response data
	data.Name = types.StringValue(domain.Name)
	data.CNAME = types.StringValue(domain.CNAME)
	data.Status = types.StringValue(domain.Status)
	data.ComplianceStatus = types.StringValue(domain.ComplianceStatus)

	// Convert DNS records
	dnsRecordsValue, diags := r.convertDNSRecordsToTerraform(ctx, &domain.DNSRecords)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.DNSRecords = dnsRecordsValue

	// Convert DNS status
	dnsStatusObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
		"cname": types.BoolType,
		"mx":    types.BoolType,
		"txt":   types.BoolType,
	}, &DNSStatusModel{
		CNAME: types.BoolValue(domain.DNSStatus.CNAME),
		MX:    types.BoolValue(domain.DNSStatus.MX),
		TXT:   types.BoolValue(domain.DNSStatus.TXT),
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.DNSStatus = dnsStatusObj

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SendingDomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Sending domains don't support updates through the API
	// Just read the current state
	r.Read(ctx, resource.ReadRequest{State: req.State}, (*resource.ReadResponse)(resp))
}

func (r *SendingDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SendingDomainResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// The API doesn't provide a delete endpoint for sending domains
	// We'll just remove it from state
	tflog.Warn(ctx, "Sending domains cannot be deleted via API. The domain will be removed from Terraform state but will remain in Mailtrap.")
}

func (r *SendingDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expected format: account_id/domain_id
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Incorrect Import ID",
			"Import ID must be in the format: account_id/domain_id",
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

	domainID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Domain ID",
			fmt.Sprintf("Could not parse domain ID: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account_id"), accountID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), domainID)...)
}

// Helper function to convert DNS records to Terraform types
func (r *SendingDomainResource) convertDNSRecordsToTerraform(ctx context.Context, records *client.DNSRecords) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Define the attribute types for a DNS record
	dnsRecordAttrTypes := map[string]attr.Type{
		"priority":    types.Int64Type,
		"record_type": types.StringType,
		"hostname":    types.StringType,
		"value":       types.StringType,
		"status":      types.StringType,
	}

	// Convert CNAME records
	cnameRecords := make([]attr.Value, len(records.CNAME))
	for i, record := range records.CNAME {
		var priority types.Int64
		if record.Priority != nil {
			priority = types.Int64Value(int64(*record.Priority))
		} else {
			priority = types.Int64Null()
		}

		recordObj, d := types.ObjectValueFrom(ctx, dnsRecordAttrTypes, DNSRecordModel{
			Priority:   priority,
			RecordType: types.StringValue(record.RecordType),
			Hostname:   types.StringValue(record.Hostname),
			Value:      types.StringValue(record.Value),
			Status:     types.StringValue(record.Status),
		})
		diags.Append(d...)
		cnameRecords[i] = recordObj
	}

	// Convert MX records
	mxRecords := make([]attr.Value, len(records.MX))
	for i, record := range records.MX {
		var priority types.Int64
		if record.Priority != nil {
			priority = types.Int64Value(int64(*record.Priority))
		} else {
			priority = types.Int64Null()
		}

		recordObj, d := types.ObjectValueFrom(ctx, dnsRecordAttrTypes, DNSRecordModel{
			Priority:   priority,
			RecordType: types.StringValue(record.RecordType),
			Hostname:   types.StringValue(record.Hostname),
			Value:      types.StringValue(record.Value),
			Status:     types.StringValue(record.Status),
		})
		diags.Append(d...)
		mxRecords[i] = recordObj
	}

	// Convert TXT records
	txtRecords := make([]attr.Value, len(records.TXT))
	for i, record := range records.TXT {
		var priority types.Int64
		if record.Priority != nil {
			priority = types.Int64Value(int64(*record.Priority))
		} else {
			priority = types.Int64Null()
		}

		recordObj, d := types.ObjectValueFrom(ctx, dnsRecordAttrTypes, DNSRecordModel{
			Priority:   priority,
			RecordType: types.StringValue(record.RecordType),
			Hostname:   types.StringValue(record.Hostname),
			Value:      types.StringValue(record.Value),
			Status:     types.StringValue(record.Status),
		})
		diags.Append(d...)
		txtRecords[i] = recordObj
	}

	// Create lists
	cnameList, d := types.ListValue(types.ObjectType{AttrTypes: dnsRecordAttrTypes}, cnameRecords)
	diags.Append(d...)
	
	mxList, d := types.ListValue(types.ObjectType{AttrTypes: dnsRecordAttrTypes}, mxRecords)
	diags.Append(d...)
	
	txtList, d := types.ListValue(types.ObjectType{AttrTypes: dnsRecordAttrTypes}, txtRecords)
	diags.Append(d...)

	// Create the DNS records object
	dnsRecordsAttrTypes := map[string]attr.Type{
		"cname": types.ListType{ElemType: types.ObjectType{AttrTypes: dnsRecordAttrTypes}},
		"mx":    types.ListType{ElemType: types.ObjectType{AttrTypes: dnsRecordAttrTypes}},
		"txt":   types.ListType{ElemType: types.ObjectType{AttrTypes: dnsRecordAttrTypes}},
	}

	dnsRecordsObj, d := types.ObjectValue(dnsRecordsAttrTypes, map[string]attr.Value{
		"cname": cnameList,
		"mx":    mxList,
		"txt":   txtList,
	})
	diags.Append(d...)

	return dnsRecordsObj, diags
}
