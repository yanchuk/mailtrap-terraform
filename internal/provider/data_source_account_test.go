package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mailtrap/terraform-provider-mailtrap/internal/client"
)

func TestAccountDataSource_Metadata(t *testing.T) {
	d := &AccountDataSource{}
	
	req := datasource.MetadataRequest{
		ProviderTypeName: "mailtrap",
	}
	resp := &datasource.MetadataResponse{}
	
	d.Metadata(context.Background(), req, resp)
	
	expected := "mailtrap_account"
	if resp.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, resp.TypeName)
	}
}

func TestAccountDataSource_Schema(t *testing.T) {
	d := &AccountDataSource{}
	
	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}
	
	d.Schema(context.Background(), req, resp)
	
	if resp.Schema.Attributes == nil {
		t.Fatal("Expected schema attributes to be defined")
	}
	
	// Check required attributes
	requiredAttrs := []string{"id"}
	for _, attr := range requiredAttrs {
		if _, exists := resp.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected required attribute %s to exist", attr)
		}
	}
	
	// Check computed attributes
	computedAttrs := []string{"name"}
	for _, attr := range computedAttrs {
		if _, exists := resp.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected computed attribute %s to exist", attr)
		}
	}
}

func TestAccountDataSource_Configure(t *testing.T) {
	d := &AccountDataSource{}
	
	// Test with valid provider data
	providerData := &ProviderData{
		Client:    &client.Client{},
		AccountID: 12345,
	}
	
	req := datasource.ConfigureRequest{
		ProviderData: providerData,
	}
	resp := &datasource.ConfigureResponse{}
	
	d.Configure(context.Background(), req, resp)
	
	if resp.Diagnostics.HasError() {
		t.Errorf("Expected no errors, got %v", resp.Diagnostics.Errors())
	}
	
	if d.client != providerData.Client {
		t.Error("Expected client to be set from provider data")
	}
}

func TestAccountDataSource_Configure_InvalidProviderData(t *testing.T) {
	d := &AccountDataSource{}
	
	req := datasource.ConfigureRequest{
		ProviderData: "invalid",
	}
	resp := &datasource.ConfigureResponse{}
	
	d.Configure(context.Background(), req, resp)
	
	if !resp.Diagnostics.HasError() {
		t.Error("Expected error for invalid provider data")
	}
}

func TestAccountDataSource_Configure_NilProviderData(t *testing.T) {
	d := &AccountDataSource{}
	
	req := datasource.ConfigureRequest{
		ProviderData: nil,
	}
	resp := &datasource.ConfigureResponse{}
	
	d.Configure(context.Background(), req, resp)
	
	if resp.Diagnostics.HasError() {
		t.Error("Expected no error for nil provider data")
	}
}

func TestNewAccountDataSource(t *testing.T) {
	d := NewAccountDataSource()
	
	if d == nil {
		t.Fatal("Expected data source to be created")
	}
	
	_, ok := d.(*AccountDataSource)
	if !ok {
		t.Error("Expected AccountDataSource type")
	}
}

func TestAccountDataSourceModel(t *testing.T) {
	model := AccountDataSourceModel{
		ID:   types.Int64Value(123),
		Name: types.StringValue("Test Account"),
	}
	
	if model.ID.ValueInt64() != 123 {
		t.Errorf("Expected ID 123, got %d", model.ID.ValueInt64())
	}
	
	if model.Name.ValueString() != "Test Account" {
		t.Errorf("Expected name 'Test Account', got %s", model.Name.ValueString())
	}
}