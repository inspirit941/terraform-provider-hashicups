package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
    _ provider.Provider = &hashicupsProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
    return func() provider.Provider {
        return &hashicupsProvider{
            version: version,
        }
    }
}

// hashicupsProvider is the provider implementation.
type hashicupsProvider struct {
    // version is set to the provider version on release, "dev" when the
    // provider is built and ran locally, and "test" when running acceptance
    // testing.
    version string
}

// hashicupsProviderModel maps provider schema data to a Go type.
// schema에 정의한 필드를 provider 서버에서 Marshal / unmarshal하기 위한 struct
type hashicupsProviderModel struct {
    Host     types.String `tfsdk:"host"`
    Username types.String `tfsdk:"username"`
    Password types.String `tfsdk:"password"`
}


// Metadata returns the provider type name.
func (p *hashicupsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	// provider에 정의한 모든 resource / datasource에 포함되는 이름이 typeName.
	// AWS의 경우 aws_ 형태. 이 예시의 경우 hashicups_ 형태로 리소스가 정의된다.
    resp.TypeName = "hashicups"
    resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *hashicupsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	// attributes available on your provider.
	// 예시의 경우 username, password를 입력받아서 authentication을 수행하도록 한다.
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "host": schema.StringAttribute{
                Optional: true,
            },
            "username": schema.StringAttribute{
                Optional: true,
            },
            "password": schema.StringAttribute{
                Optional:  true,
                Sensitive: true,
            },
        },
    }
}

// Configure prepares a HashiCups API client for data sources and resources.
// -> provider가 running하기 위해 필요한 로직을 넣는 함수. 보통 authentication이 들어감.
func (p *hashicupsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

// DataSources defines the data sources implemented in the provider.
func (p *hashicupsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
    return nil
}

// Resources defines the resources implemented in the provider.
func (p *hashicupsProvider) Resources(_ context.Context) []func() resource.Resource {
    return nil
}
