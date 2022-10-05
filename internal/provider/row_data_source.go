package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/iRevive/terraform-provider-gdashboard/internal/provider/grafana"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &RowDataSource{}

func NewRowDataSource() datasource.DataSource {
	return &RowDataSource{}
}

// RowDataSource defines the data source implementation.
type RowDataSource struct {
}

// RowDataSourceModel describes the data source data model.
type RowDataSourceModel struct {
	Id    types.String `tfsdk:"id"`
	Json  types.String `tfsdk:"json"`
	Title types.String `tfsdk:"title"`
	Graph []RowOptions `tfsdk:"graph"`
}

type RowOptions struct {
	Collapsed types.Bool `tfsdk:"collapsed"`
}

func (d *RowDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_row"
}

func rowGraphBlock() tfsdk.Block {
	return tfsdk.Block{
		NestingMode: tfsdk.BlockNestingModeList,
		MinItems:    0,
		MaxItems:    1,
		Attributes: map[string]tfsdk.Attribute{
			"collapsed": {
				Type:        types.BoolType,
				Optional:    true,
				Description: "Whether to render row collapsed or not",
			},
		},
	}
}

func (d *RowDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Row panel data source",

		Blocks: map[string]tfsdk.Block{
			"graph": rowGraphBlock(),
		},

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"json": {
				Type:     types.StringType,
				Computed: true,
			},
			"title": {
				Type:        types.StringType,
				Required:    true,
				Description: "The title of the panel",
			},
		},
	}, nil
}

func (d *RowDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
}

func (d *RowDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RowDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rowPanel := grafana.RowPanel{}

	for _, graph := range data.Graph {
		if !graph.Collapsed.Null {
			rowPanel.Collapsed = graph.Collapsed.Value
		}
	}

	panel := &grafana.Panel{
		CommonPanel: grafana.CommonPanel{
			OfType: grafana.RowType,
			Title:  data.Title.Value,
			Type:   "row",
			Span:   12,
			IsNew:  true,
		},
		RowPanel: &rowPanel,
	}

	jsonData, err := json.MarshalIndent(panel, "", "  ")
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Could not marshal json: %s", err))
		return
	}

	data.Json = types.String{Value: string(jsonData)}
	data.Id = types.String{Value: strconv.Itoa(hashcode(jsonData))}

	//resp.Diagnostics.AddError("Client Error", fmt.Sprintf("%s", string(jsonData)))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
