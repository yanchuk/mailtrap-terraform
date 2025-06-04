package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/yanchuk/mailtrap-terraform/internal/client"
)

func TestProjectResource_Metadata(t *testing.T) {
	r := &ProjectResource{}
	
	req := resource.MetadataRequest{
		ProviderTypeName: "mailtrap",
	}
	resp := &resource.MetadataResponse{}
	
	r.Metadata(context.Background(), req, resp)
	
	expected := "mailtrap_project"
	if resp.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, resp.TypeName)
	}
}

func TestProjectResource_Schema(t *testing.T) {
	r := &ProjectResource{}
	
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
	computedAttrs := []string{"id", "share_links"}
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

func TestProjectResource_Configure(t *testing.T) {
	r := &ProjectResource{}
	
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

func TestProjectResource_Configure_InvalidProviderData(t *testing.T) {
	r := &ProjectResource{}
	
	req := resource.ConfigureRequest{
		ProviderData: "invalid",
	}
	resp := &resource.ConfigureResponse{}
	
	r.Configure(context.Background(), req, resp)
	
	if !resp.Diagnostics.HasError() {
		t.Error("Expected error for invalid provider data")
	}
}

func TestProjectResource_Configure_NilProviderData(t *testing.T) {
	r := &ProjectResource{}
	
	req := resource.ConfigureRequest{
		ProviderData: nil,
	}
	resp := &resource.ConfigureResponse{}
	
	r.Configure(context.Background(), req, resp)
	
	if resp.Diagnostics.HasError() {
		t.Error("Expected no error for nil provider data")
	}
}

func TestProjectResource_ImportState(t *testing.T) {
	// Skip this test as it requires proper Terraform state setup
	t.Skip("ImportState testing requires terraform-plugin-testing framework with proper state setup")
}

func TestNewProjectResource(t *testing.T) {
	r := NewProjectResource()
	
	if r == nil {
		t.Fatal("Expected resource to be created")
	}
	
	_, ok := r.(*ProjectResource)
	if !ok {
		t.Error("Expected ProjectResource type")
	}
}

func TestShareLinksModel(t *testing.T) {
	model := ShareLinksModel{
		Admin:  types.StringValue("https://admin.link"),
		Viewer: types.StringValue("https://viewer.link"),
	}
	
	if model.Admin.ValueString() != "https://admin.link" {
		t.Errorf("Expected admin link 'https://admin.link', got %s", model.Admin.ValueString())
	}
	
	if model.Viewer.ValueString() != "https://viewer.link" {
		t.Errorf("Expected viewer link 'https://viewer.link', got %s", model.Viewer.ValueString())
	}
}

func TestProjectResourceModel(t *testing.T) {
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
	
	model := ProjectResourceModel{
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