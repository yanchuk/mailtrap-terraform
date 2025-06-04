package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func TestMailtrapProvider_Metadata(t *testing.T) {
	p := &MailtrapProvider{version: "test"}
	
	req := provider.MetadataRequest{}
	resp := &provider.MetadataResponse{}
	
	p.Metadata(context.Background(), req, resp)
	
	if resp.TypeName != "mailtrap" {
		t.Errorf("Expected type name 'mailtrap', got %s", resp.TypeName)
	}
	
	if resp.Version != "test" {
		t.Errorf("Expected version 'test', got %s", resp.Version)
	}
}

func TestMailtrapProvider_Schema(t *testing.T) {
	p := &MailtrapProvider{}
	
	req := provider.SchemaRequest{}
	resp := &provider.SchemaResponse{}
	
	p.Schema(context.Background(), req, resp)
	
	if resp.Schema.Attributes == nil {
		t.Fatal("Expected schema attributes to be defined")
	}
	
	// Check api_token attribute
	apiTokenAttr, exists := resp.Schema.Attributes["api_token"]
	if !exists {
		t.Fatal("Expected api_token attribute to exist")
	}
	
	if !apiTokenAttr.IsOptional() {
		t.Error("Expected api_token to be optional")
	}
	
	if !apiTokenAttr.IsSensitive() {
		t.Error("Expected api_token to be sensitive")
	}
	
	// Check account_id attribute
	accountIDAttr, exists := resp.Schema.Attributes["account_id"]
	if !exists {
		t.Fatal("Expected account_id attribute to exist")
	}
	
	if !accountIDAttr.IsOptional() {
		t.Error("Expected account_id to be optional")
	}
}

func TestMailtrapProvider_Configure_Success(t *testing.T) {
	// Skip this test as it requires proper Terraform configuration setup
	t.Skip("Provider configuration testing requires terraform-plugin-testing framework")
}

func TestMailtrapProvider_Configure_MissingAPIToken(t *testing.T) {
	// Skip this test as it requires proper Terraform configuration setup
	t.Skip("Provider configuration testing requires terraform-plugin-testing framework")
}

func TestMailtrapProvider_Resources(t *testing.T) {
	p := &MailtrapProvider{}
	
	resources := p.Resources(context.Background())
	
	expectedCount := 3 // project, inbox, sending_domain
	if len(resources) != expectedCount {
		t.Errorf("Expected %d resources, got %d", expectedCount, len(resources))
	}
	
	// Test that resources can be instantiated
	for i, resourceFunc := range resources {
		resource := resourceFunc()
		if resource == nil {
			t.Errorf("Resource %d returned nil", i)
		}
	}
}

func TestMailtrapProvider_DataSources(t *testing.T) {
	p := &MailtrapProvider{}
	
	dataSources := p.DataSources(context.Background())
	
	expectedCount := 4 // account, project, inbox, sending_domain
	if len(dataSources) != expectedCount {
		t.Errorf("Expected %d data sources, got %d", expectedCount, len(dataSources))
	}
	
	// Test that data sources can be instantiated
	for i, dataSourceFunc := range dataSources {
		dataSource := dataSourceFunc()
		if dataSource == nil {
			t.Errorf("Data source %d returned nil", i)
		}
	}
}

func TestNew(t *testing.T) {
	version := "1.0.0"
	providerFunc := New(version)
	
	if providerFunc == nil {
		t.Fatal("Expected provider function to be returned")
	}
	
	provider := providerFunc()
	if provider == nil {
		t.Fatal("Expected provider to be returned")
	}
	
	mailtrapProvider, ok := provider.(*MailtrapProvider)
	if !ok {
		t.Fatal("Expected MailtrapProvider type")
	}
	
	if mailtrapProvider.version != version {
		t.Errorf("Expected version %s, got %s", version, mailtrapProvider.version)
	}
}

func TestParseInt64(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
		hasError bool
	}{
		{"123", 123, false},
		{"0", 0, false},
		{"-456", -456, false},
		{"abc", 0, true},
		{"", 0, true},
		{"123.45", 0, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := parseInt64(tt.input)
			
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error for input %s, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for input %s, got %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("Expected %d for input %s, got %d", tt.expected, tt.input, result)
				}
			}
		})
	}
}

// TestProviderServer tests that the provider can be served
func TestProviderServer(t *testing.T) {
	// Create a test provider server
	server := providerserver.NewProtocol6(New("test")())
	if server == nil {
		t.Fatal("Failed to create provider server")
	}
}

// Helper function to create a test provider for integration tests
func testProvider() *MailtrapProvider {
	return &MailtrapProvider{version: "test"}
}

// Helper function to create test provider data
func testProviderData() *ProviderData {
	return &ProviderData{
		Client:    nil, // Would be mocked in real tests
		AccountID: 12345,
	}
}