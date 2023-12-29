package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &coffeesDataSource{}
	// configure method 구현. external api 호출하거나 state가 필요한 경우 configure method 포함한 인터페이스를 구현하면 된다.
	_ datasource.DataSourceWithConfigure = &coffeesDataSource{}
)

func NewCoffeesDataSource() datasource.DataSource {
	return &coffeesDataSource{}
}

// coffee list를 넘겨주는 data source.
// -> api call 수행하기 위한 client를 struct에 포함한다.
type coffeesDataSource struct {
	client *hashicups.Client
}

func (d *coffeesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	// returns "hashicups_coffees"
	resp.TypeName = req.ProviderTypeName + "_coffees"
}

// Schema defines the schema for the data source.
func (d *coffeesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	// data source <-> terraform-understanding data structure.
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// 리스트 object
			"coffees": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					// 리스트 object 각각의 attributes
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"teaser": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"price": schema.Float64Attribute{
							Computed: true,
						},
						"image": schema.StringAttribute{
							Computed: true,
						},
						"ingredients": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.Int64Attribute{
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *coffeesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
}

// Configure adds the provider configured client to the data source.
func (d *coffeesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*hashicups.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}
