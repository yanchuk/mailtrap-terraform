package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/yanchuk/mailtrap-terraform/internal/client"
)

func TestSendingDomainResource_Metadata(t *testing.T) {
	r := &SendingDomainResource{}
	
	req := resource.MetadataRequest{
		ProviderTypeName: "mailtrap",
	}
	resp := &resource.MetadataResponse{}
	
	r.Metadata(context.Background(), req, resp)
	
	expected := "mailtrap_sending_domain"
	if resp.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, resp.TypeName)
	}
}

func TestSendingDomainResource_Schema(t *testing.T) {
	r := &SendingDomainResource{}
	
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}
	
	r.Schema(context.Background(), req, resp)
	
	if resp.Schema.Attributes == nil {
		t.Fatal("Expected schema attributes to be defined")
	}
	
	// Check required attributes
	requiredAttrs := []string{"name"}
	for _, attr := range requiredAttrs {
		if _, exists := resp.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected required attribute %s to exist", attr)
		}
	}
	
	// Check computed attributes
	computedAttrs := []string{"id", "cname", "status", "compliance_status", "dns_records", "dns_status"}
	for _, attr := range computedAttrs {
		if _, exists := resp.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected computed attribute %s to exist", attr)
		}
	}
	
	// Check optional attributes
	optionalAttrs := []string{"account_id"}
	for _, attr := range optionalAttrs {
		if _, exists := resp.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected optional attribute %s to exist", attr)
		}
	}
}

func TestSendingDomainResource_Configure(t *testing.T) {
	r := &SendingDomainResource{}
	
	// Test with valid provider data
	providerData := &ProviderData{
		Client:    &client.Client{},
		AccountID: 12345,
	}
	
	req := resource.ConfigureRequest{
		ProviderData: providerData,
	}
	resp := &resource.ConfigureResponse{}
	
	r.Configure(context.Background(), req, resp)
	
	if resp.Diagnostics.HasError() {
		t.Errorf("Expected no errors, got %v", resp.Diagnostics.Errors())
	}
	
	if r.client != providerData.Client {
		t.Error("Expected client to be set from provider data")
	}
	
	if r.accountID != providerData.AccountID {
		t.Errorf("Expected account ID %d, got %d", providerData.AccountID, r.accountID)
	}
}

func TestSendingDomainResource_Configure_InvalidProviderData(t *testing.T) {
	r := &SendingDomainResource{}
	
	req := resource.ConfigureRequest{
		ProviderData: "invalid",
	}
	resp := &resource.ConfigureResponse{}
	
	r.Configure(context.Background(), req, resp)
	
	if !resp.Diagnostics.HasError() {
		t.Error("Expected error for invalid provider data")
	}
}

func TestSendingDomainResource_ImportState(t *testing.T) {
	// Skip this test as it requires proper Terraform state setup
	t.Skip("ImportState testing requires terraform-plugin-testing framework with proper state setup")
}

func TestNewSendingDomainResource(t *testing.T) {
	r := NewSendingDomainResource()
	
	if r == nil {
		t.Fatal("Expected resource to be created")
	}
	
	_, ok := r.(*SendingDomainResource)
	if !ok {
		t.Error("Expected SendingDomainResource type")
	}
}

func TestDNSRecordModel(t *testing.T) {
	model := DNSRecordModel{
		Priority:   types.Int64Value(10),
		RecordType: types.StringValue("MX"),
		Hostname:   types.StringValue("example.com"),
		Value:      types.StringValue("mail.example.com"),
		Status:     types.StringValue("verified"),
	}
	
	if model.Priority.ValueInt64() != 10 {
		t.Errorf("Expected priority 10, got %d", model.Priority.ValueInt64())
	}
	
	if model.RecordType.ValueString() != "MX" {
		t.Errorf("Expected record type 'MX', got %s", model.RecordType.ValueString())
	}
	
	if model.Hostname.ValueString() != "example.com" {
		t.Errorf("Expected hostname 'example.com', got %s", model.Hostname.ValueString())
	}
	
	if model.Value.ValueString() != "mail.example.com" {
		t.Errorf("Expected value 'mail.example.com', got %s", model.Value.ValueString())
	}
	
	if model.Status.ValueString() != "verified" {
		t.Errorf("Expected status 'verified', got %s", model.Status.ValueString())
	}
}

func TestDNSStatusModel(t *testing.T) {
	model := DNSStatusModel{
		CNAME: types.BoolValue(true),
		MX:    types.BoolValue(false),
		TXT:   types.BoolValue(true),
	}
	
	if !model.CNAME.ValueBool() {
		t.Error("Expected CNAME to be true")
	}
	
	if model.MX.ValueBool() {
		t.Error("Expected MX to be false")
	}
	
	if !model.TXT.ValueBool() {
		t.Error("Expected TXT to be true")
	}
}

func TestSendingDomainResourceModel(t *testing.T) {
	// Create DNS record models
	dnsRecordAttrTypes := map[string]attr.Type{
		"priority":    types.Int64Type,
		"record_type": types.StringType,
		"hostname":    types.StringType,
		"value":       types.StringType,
		"status":      types.StringType,
	}
	
	cnameRecord, _ := types.ObjectValueFrom(context.Background(), dnsRecordAttrTypes, DNSRecordModel{
		Priority:   types.Int64Null(),
		RecordType: types.StringValue("CNAME"),
		Hostname:   types.StringValue("mail.example.com"),
		Value:      types.StringValue("target.mailtrap.io"),
		Status:     types.StringValue("verified"),
	})
	
	cnameList := types.ListValueMust(types.ObjectType{AttrTypes: dnsRecordAttrTypes}, []attr.Value{cnameRecord})
	mxList := types.ListValueMust(types.ObjectType{AttrTypes: dnsRecordAttrTypes}, []attr.Value{})
	txtList := types.ListValueMust(types.ObjectType{AttrTypes: dnsRecordAttrTypes}, []attr.Value{})
	
	dnsRecordsAttrTypes := map[string]attr.Type{
		"cname": types.ListType{ElemType: types.ObjectType{AttrTypes: dnsRecordAttrTypes}},
		"mx":    types.ListType{ElemType: types.ObjectType{AttrTypes: dnsRecordAttrTypes}},
		"txt":   types.ListType{ElemType: types.ObjectType{AttrTypes: dnsRecordAttrTypes}},
	}
	
	dnsRecords, _ := types.ObjectValue(dnsRecordsAttrTypes, map[string]attr.Value{
		"cname": cnameList,
		"mx":    mxList,
		"txt":   txtList,
	})
	
	dnsStatusAttrTypes := map[string]attr.Type{
		"cname": types.BoolType,
		"mx":    types.BoolType,
		"txt":   types.BoolType,
	}
	
	dnsStatus, _ := types.ObjectValueFrom(context.Background(), dnsStatusAttrTypes, DNSStatusModel{
		CNAME: types.BoolValue(true),
		MX:    types.BoolValue(false),
		TXT:   types.BoolValue(false),
	})
	
	model := SendingDomainResourceModel{
		ID:               types.Int64Value(123),
		AccountID:        types.Int64Value(456),
		Name:             types.StringValue("example.com"),
		CNAME:            types.StringValue("cname.mailtrap.io"),
		Status:           types.StringValue("verified"),
		ComplianceStatus: types.StringValue("compliant"),
		DNSRecords:       dnsRecords,
		DNSStatus:        dnsStatus,
	}
	
	if model.ID.ValueInt64() != 123 {
		t.Errorf("Expected ID 123, got %d", model.ID.ValueInt64())
	}
	
	if model.Name.ValueString() != "example.com" {
		t.Errorf("Expected name 'example.com', got %s", model.Name.ValueString())
	}
	
	if model.CNAME.ValueString() != "cname.mailtrap.io" {
		t.Errorf("Expected CNAME 'cname.mailtrap.io', got %s", model.CNAME.ValueString())
	}
	
	if model.Status.ValueString() != "verified" {
		t.Errorf("Expected status 'verified', got %s", model.Status.ValueString())
	}
	
	if model.ComplianceStatus.ValueString() != "compliant" {
		t.Errorf("Expected compliance status 'compliant', got %s", model.ComplianceStatus.ValueString())
	}
	
	if model.DNSRecords.IsNull() {
		t.Error("Expected DNS records to be set")
	}
	
	if model.DNSStatus.IsNull() {
		t.Error("Expected DNS status to be set")
	}
}