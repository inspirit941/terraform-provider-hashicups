package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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


/// Marshal / Unmarshal go struct with terraform data struct, by the name of the attribute.
// coffeesDataSourceModel maps the data source schema data. -> list of coffees
type coffeesDataSourceModel struct {
	Coffees []coffeesModel `tfsdk:"coffees"`
}

// coffeesModel maps coffees schema data. -> coffee element
type coffeesModel struct {
	ID          types.Int64               `tfsdk:"id"`
	Name        types.String              `tfsdk:"name"`
	Teaser      types.String              `tfsdk:"teaser"`
	Description types.String              `tfsdk:"description"`
	Price       types.Float64             `tfsdk:"price"`
	Image       types.String              `tfsdk:"image"`
	Ingredients []coffeesIngredientsModel `tfsdk:"ingredients"`
}

// coffeesIngredientsModel maps coffee ingredients data
type coffeesIngredientsModel struct {
	ID types.Int64 `tfsdk:"id"`
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

// terraform의 plan / apply / destroy 메소드로 lifecycle을 정상적으로 동작하려면 CRUD를 구현해야 함.
// 이건 data source이므로, read만 있으면 된다. 
// (그냥 resource라면 CRUD를 전부 구현해야 한다.)
// Read refreshes the Terraform state with the latest data.
func (d *coffeesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
  var state coffeesDataSourceModel

	// client API로 필요한 리소스 가져온다.
  coffees, err := d.client.GetCoffees()
  if err != nil {
    resp.Diagnostics.AddError(
      "Unable to Read HashiCups Coffees",
      err.Error(),
    )
    return
  }

  // Map response body to model
	// client API 리소스를 terraform struct에 매핑한다 (conversion)
  for _, coffee := range coffees {
    coffeeState := coffeesModel{
      ID:          types.Int64Value(int64(coffee.ID)),
      Name:        types.StringValue(coffee.Name),
      Teaser:      types.StringValue(coffee.Teaser),
      Description: types.StringValue(coffee.Description),
      Price:       types.Float64Value(coffee.Price),
      Image:       types.StringValue(coffee.Image),
    }

    for _, ingredient := range coffee.Ingredient {
      coffeeState.Ingredients = append(coffeeState.Ingredients, coffeesIngredientsModel{
        ID: types.Int64Value(int64(ingredient.ID)),
      })
    }

    state.Coffees = append(state.Coffees, coffeeState)
  }

  // Set state. 매핑한 정보를 terraform state에 저장한다.
  diags := resp.State.Set(ctx, &state)
  resp.Diagnostics.Append(diags...)
  if resp.Diagnostics.HasError() {
    return
  }
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
