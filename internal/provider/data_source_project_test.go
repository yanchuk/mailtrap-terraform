package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/yanchuk/mailtrap-terraform/internal/client"
)

func TestProjectDataSource_Metadata(t *testing.T) {
	d := &ProjectDataSource{}
	
	req := datasource.MetadataRequest{
		ProviderTypeName: "mailtrap",
	}
	resp := &datasource.MetadataResponse{}
	
	d.Metadata(context.Background(), req, resp)
	
	expected := "mailtrap_project"
	if resp.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, resp.TypeName)
	}
}

func TestProjectDataSource_Schema(t *testing.T) {
	d := &ProjectDataSource{}
	
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
	computedAttrs := []string{"name", "share_links"}
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

func TestProjectDataSource_Configure(t *testing.T) {
	d := &ProjectDataSource{}
	
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
	
	if d.accountID != providerData.AccountID {
		t.Errorf("Expected account ID %d, got %d", providerData.AccountID, d.accountID)
	}
}

func TestProjectDataSource_Configure_InvalidProviderData(t *testing.T) {
	d := &ProjectDataSource{}
	
	req := datasource.ConfigureRequest{
		ProviderData: "invalid",
	}
	resp := &datasource.ConfigureResponse{}
	
	d.Configure(context.Background(), req, resp)
	
	if !resp.Diagnostics.HasError() {
		t.Error("Expected error for invalid provider data")
	}
}

func TestProjectDataSource_Configure_NilProviderData(t *testing.T) {
	d := &ProjectDataSource{}
	
	req := datasource.ConfigureRequest{
		ProviderData: nil,
	}
	resp := &datasource.ConfigureResponse{}
	
	d.Configure(context.Background(), req, resp)
	
	if resp.Diagnostics.HasError() {
		t.Error("Expected no error for nil provider data")
	}
}

func TestNewProjectDataSource(t *testing.T) {
	d := NewProjectDataSource()
	
	if d == nil {
		t.Fatal("Expected data source to be created")
	}
	
	_, ok := d.(*ProjectDataSource)
	if !ok {
		t.Error("Expected ProjectDataSource type")
	}
}

func TestProjectDataSourceModel(t *testing.T) {
	// Test creating share links object
	shareLinksAttrTypes := map[string]attr.Type{
		"admin":  types.StringType,
		"viewer": types.StringType,
	}
	
	shareLinksObj, diags := types.ObjectValueFrom(context.Background(), shareLinksAttrTypes, &ShareLinksModel{
		Admin:  types.StringValue("https://admin.link"),
		Viewer: types.StringValue("https://viewer.link"),
	})
	
	if diags.HasError() {
		t.Fatalf("Expected no errors creating share links object, got %v", diags.Errors())
	}
	
	model := ProjectDataSourceModel{
		ID:         types.Int64Value(123),
		AccountID:  types.Int64Value(456),
		Name:       types.StringValue("Test Project"),
		ShareLinks: shareLinksObj,
	}
	
	if model.ID.ValueInt64() != 123 {
		t.Errorf("Expected ID 123, got %d", model.ID.ValueInt64())
	}
	
	if model.AccountID.ValueInt64() != 456 {
		t.Errorf("Expected AccountID 456, got %d", model.AccountID.ValueInt64())
	}
	
	if model.Name.ValueString() != "Test Project" {
		t.Errorf("Expected name 'Test Project', got %s", model.Name.ValueString())
	}
	
	if model.ShareLinks.IsNull() {
		t.Error("Expected share links to be set")
	}
}