package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

func (d *RowDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_row"
}

func rowGraphBlock() schema.Block {
	return schema.ListNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"collapsed": schema.BoolAttribute{
					Optional:    true,
					Description: "Whether to render row collapsed or not",
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func (d *RowDataSource) GetSchema(_ context.Context) (schema.Schema, diag.Diagnostics) {
	return schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Row panel data source.",

		Blocks: map[string]schema.Block{
			"graph": rowGraphBlock(),
		},

		Attributes: map[string]schema.Attribute{
			"id":    idAttribute(),
			"json":  jsonAttribute(),
			"title": titleAttribute(),
		},
	}, nil
}

func (d *RowDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
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
		if !graph.Collapsed.IsNull() {
			rowPanel.Collapsed = graph.Collapsed.ValueBool()
		}
	}

	panel := &grafana.Panel{
		CommonPanel: grafana.CommonPanel{
			OfType: grafana.RowType,
			Title:  data.Title.ValueString(),
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

	data.Json = types.StringValue(string(jsonData))
	data.Id = types.StringValue(strconv.Itoa(hashcode(jsonData)))

	//resp.Diagnostics.AddError("Client Error", fmt.Sprintf("%s", string(jsonData)))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
