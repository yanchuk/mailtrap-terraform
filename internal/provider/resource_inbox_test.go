package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/yanchuk/mailtrap-terraform/internal/client"
)

func TestInboxResource_Metadata(t *testing.T) {
	r := &InboxResource{}
	
	req := resource.MetadataRequest{
		ProviderTypeName: "mailtrap",
	}
	resp := &resource.MetadataResponse{}
	
	r.Metadata(context.Background(), req, resp)
	
	expected := "mailtrap_inbox"
	if resp.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, resp.TypeName)
	}
}

func TestInboxResource_Schema(t *testing.T) {
	r := &InboxResource{}
	
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}
	
	r.Schema(context.Background(), req, resp)
	
	if resp.Schema.Attributes == nil {
		t.Fatal("Expected schema attributes to be defined")
	}
	
	// Check required attributes
	requiredAttrs := []string{"project_id", "name"}
	for _, attr := range requiredAttrs {
		if _, exists := resp.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected required attribute %s to exist", attr)
		}
	}
	
	// Check computed attributes
	computedAttrs := []string{"id", "username", "password", "domain", "email_domain", "pop3_domain"}
	for _, attr := range computedAttrs {
		if _, exists := resp.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected computed attribute %s to exist", attr)
		}
	}
	
	// Check optional attributes
	optionalAttrs := []string{"account_id", "email_username"}
	for _, attr := range optionalAttrs {
		if _, exists := resp.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected optional attribute %s to exist", attr)
		}
	}
}

func TestInboxResource_Configure(t *testing.T) {
	r := &InboxResource{}
	
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

func TestInboxResource_Configure_InvalidProviderData(t *testing.T) {
	r := &InboxResource{}
	
	req := resource.ConfigureRequest{
		ProviderData: "invalid",
	}
	resp := &resource.ConfigureResponse{}
	
	r.Configure(context.Background(), req, resp)
	
	if !resp.Diagnostics.HasError() {
		t.Error("Expected error for invalid provider data")
	}
}

func TestInboxResource_ImportState(t *testing.T) {
	// Skip this test as it requires proper Terraform state setup
	t.Skip("ImportState testing requires terraform-plugin-testing framework with proper state setup")
}

func TestNewInboxResource(t *testing.T) {
	r := NewInboxResource()
	
	if r == nil {
		t.Fatal("Expected resource to be created")
	}
	
	_, ok := r.(*InboxResource)
	if !ok {
		t.Error("Expected InboxResource type")
	}
}

func TestInboxResourceModel(t *testing.T) {
	model := InboxResourceModel{
		ID:                        types.Int64Value(123),
		AccountID:                 types.Int64Value(456),
		ProjectID:                 types.Int64Value(789),
		Name:                      types.StringValue("Test Inbox"),
		EmailUsername:             types.StringValue("test"),
		Username:                  types.StringValue("smtp_user"),
		Password:                  types.StringValue("smtp_pass"),
		EmailUsernameEnabled:      types.BoolValue(true),
		Domain:                    types.StringValue("smtp.mailtrap.io"),
		EmailDomain:               types.StringValue("sandbox.mailtrap.io"),
		POP3Domain:                types.StringValue("pop3.mailtrap.io"),
		SMTPPorts:                 types.ListValueMust(types.Int64Type, []attr.Value{types.Int64Value(587), types.Int64Value(2525)}),
		POP3Ports:                 types.ListValueMust(types.Int64Type, []attr.Value{types.Int64Value(1100), types.Int64Value(9950)}),
		Status:                    types.StringValue("active"),
		MaxSize:                   types.Int64Value(1000),
		SentMessagesCount:         types.Int64Value(5),
		ForwardedMessagesCount:    types.Int64Value(2),
		ForwardFromEmailAddress:   types.StringValue("test@example.com"),
	}
	
	if model.ID.ValueInt64() != 123 {
		t.Errorf("Expected ID 123, got %d", model.ID.ValueInt64())
	}
	
	if model.Name.ValueString() != "Test Inbox" {
		t.Errorf("Expected name 'Test Inbox', got %s", model.Name.ValueString())
	}
	
	if model.ProjectID.ValueInt64() != 789 {
		t.Errorf("Expected ProjectID 789, got %d", model.ProjectID.ValueInt64())
	}
	
	if !model.EmailUsernameEnabled.ValueBool() {
		t.Error("Expected EmailUsernameEnabled to be true")
	}
	
	if model.Domain.ValueString() != "smtp.mailtrap.io" {
		t.Errorf("Expected domain 'smtp.mailtrap.io', got %s", model.Domain.ValueString())
	}
}