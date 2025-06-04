package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mailtrap/terraform-provider-mailtrap/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SendingDomainDataSource{}

func NewSendingDomainDataSource() datasource.DataSource {
	return &SendingDomainDataSource{}
}

// SendingDomainDataSource defines the data source implementation.
type SendingDomainDataSource struct {
	client    *client.Client
	accountID int64
}

// SendingDomainDataSourceModel describes the data source data model.
type SendingDomainDataSourceModel struct {
	ID               types.Int64  `tfsdk:"id"`
	AccountID        types.Int64  `tfsdk:"account_id"`
	Name             types.String `tfsdk:"name"`
	CNAME            types.String `tfsdk:"cname"`
	Status           types.String `tfsdk:"status"`
	ComplianceStatus types.String `tfsdk:"compliance_status"`
	DNSRecords       types.Object `tfsdk:"dns_records"`
	DNSStatus        types.Object `tfsdk:"dns_status"`
}


func (d *SendingDomainDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sending_domain"
}

func (d *SendingDomainDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Sending domain data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Sending domain identifier",
				Required:            true,
			},
			"account_id": schema.Int64Attribute{
				MarkdownDescription: "Account ID for the sending domain",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Domain name",
				Computed:            true,
			},
			"cname": schema.StringAttribute{
				MarkdownDescription: "CNAME value for domain verification",
				Computed:            true,
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
				Attributes: map[string]schema.Attribute{
					"cname": schema.ListNestedAttribute{
						MarkdownDescription: "CNAME records",
						Computed:            true,
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
				Attributes: map[string]schema.Attribute{
					"cname": schema.BoolAttribute{
						MarkdownDescription: "CNAME verification status",
						Computed:            true,
					},
					"mx": schema.BoolAttribute{
						MarkdownDescription: "MX verification status",
						Computed:            true,
					},
					"txt": schema.BoolAttribute{
						MarkdownDescription: "TXT verification status",
						Computed:            true,
					},
				},
			},
		},
	}
}

func (d *SendingDomainDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SendingDomainDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SendingDomainDataSourceModel

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

	// Get sending domain
	endpoint := fmt.Sprintf("/api/accounts/%d/sending_domains/%d", accountID, data.ID.ValueInt64())
	
	var domain client.SendingDomain
	err := d.client.Get(endpoint, &domain)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read sending domain, got error: %s", err))
		return
	}

	// Update model with response data
	data.AccountID = types.Int64Value(accountID)
	data.Name = types.StringValue(domain.Name)
	data.CNAME = types.StringValue(domain.CNAME)
	data.Status = types.StringValue(domain.Status)
	data.ComplianceStatus = types.StringValue(domain.ComplianceStatus)

	// Convert DNS records
	dnsRecordsValue, diags := d.convertDNSRecordsToTerraform(ctx, &domain.DNSRecords)
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

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Helper function to convert DNS records to Terraform types
func (d *SendingDomainDataSource) convertDNSRecordsToTerraform(ctx context.Context, records *client.DNSRecords) (types.Object, diag.Diagnostics) {
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
	cnameList, cnameListDiags := types.ListValue(types.ObjectType{AttrTypes: dnsRecordAttrTypes}, cnameRecords)
	diags.Append(cnameListDiags...)
	
	mxList, mxListDiags := types.ListValue(types.ObjectType{AttrTypes: dnsRecordAttrTypes}, mxRecords)
	diags.Append(mxListDiags...)
	
	txtList, txtListDiags := types.ListValue(types.ObjectType{AttrTypes: dnsRecordAttrTypes}, txtRecords)
	diags.Append(txtListDiags...)

	// Create the DNS records object
	dnsRecordsAttrTypes := map[string]attr.Type{
		"cname": types.ListType{ElemType: types.ObjectType{AttrTypes: dnsRecordAttrTypes}},
		"mx":    types.ListType{ElemType: types.ObjectType{AttrTypes: dnsRecordAttrTypes}},
		"txt":   types.ListType{ElemType: types.ObjectType{AttrTypes: dnsRecordAttrTypes}},
	}

	dnsRecordsObj, dnsRecordsDiags := types.ObjectValue(dnsRecordsAttrTypes, map[string]attr.Value{
		"cname": cnameList,
		"mx":    mxList,
		"txt":   txtList,
	})
	diags.Append(dnsRecordsDiags...)

	return dnsRecordsObj, diags
}
