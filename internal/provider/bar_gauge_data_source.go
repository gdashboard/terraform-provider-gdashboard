package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/iRevive/terraform-provider-gdashboard/internal/provider/grafana"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &BarGaugeDataSource{}

func NewBarGaugeDataSource() datasource.DataSource {
	return &BarGaugeDataSource{}
}

// BarGaugeDataSource defines the data source implementation.
type BarGaugeDataSource struct {
	Defaults BarGaugeDefaults
}

type BarGaugeDefaults struct {
	Field FieldDefaults
	Graph BarGaugeGraphDefault
}

type BarGaugeGraphDefault struct {
	Orientation   string
	DisplayMode   string
	TextAlignment string
	TextSize      TextSizeDefaults
	ReduceOptions ReduceOptionDefaults
}

// BarGaugeDataSourceModel describes the data source data model.
type BarGaugeDataSourceModel struct {
	Id              types.String           `tfsdk:"id"`
	Json            types.String           `tfsdk:"json"`
	Title           types.String           `tfsdk:"title"`
	Description     types.String           `tfsdk:"description"`
	Queries         []Query                `tfsdk:"queries"`
	Field           []FieldOptions         `tfsdk:"field"`
	Graph           []BarGaugeOptions      `tfsdk:"graph"`
	Overrides       []FieldOverrideOptions `tfsdk:"overrides"`
	Transformations []Transformations      `tfsdk:"transform"`
}

type BarGaugeOptions struct {
	Orientation   types.String      `tfsdk:"orientation"`
	DisplayMode   types.String      `tfsdk:"display_mode"`
	TextAlignment types.String      `tfsdk:"text_alignment"`
	TextSize      []TextSizeOptions `tfsdk:"text_size"`
	ReduceOptions []ReduceOptions   `tfsdk:"options"`
}

func (d *BarGaugeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bar_gauge"
}

func barGaugeGraphBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: "The visualization options.",
		NestedObject: schema.NestedBlockObject{
			Blocks: map[string]schema.Block{
				"options":   reduceOptionsBlock(),
				"text_size": textSizeBlock(),
			},
			Attributes: map[string]schema.Attribute{
				"orientation": schema.StringAttribute{
					Optional:            true,
					Description:         "The layout orientation. The choices are: auto, horizontal, vertical.",
					MarkdownDescription: "The layout orientation. The choices are: `auto`, `horizontal`, `vertical`.",
					Validators: []validator.String{
						stringvalidator.OneOf("auto", "horizontal", "vertical"),
					},
				},
				"display_mode": schema.StringAttribute{
					Optional:            true,
					Description:         "The display mode. The choices are: gradient, lcd, basic.",
					MarkdownDescription: "The display mode. The choices are: `gradient`, `lcd`, `basic`.",
					Validators: []validator.String{
						stringvalidator.OneOf("gradient", "lcd", "basic"),
					},
				},
				"text_alignment": schema.StringAttribute{
					Optional:            true,
					Description:         "The text alignment. The choices are: auto, center.",
					MarkdownDescription: "The text alignment. The choices are: `auto`, `center`.",
					Validators: []validator.String{
						stringvalidator.OneOf("auto", "center"),
					},
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func (d *BarGaugeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description:         "Bar gauge panel data source.",
		MarkdownDescription: "Bar gauge panel data source. See Grafana [documentation](https://grafana.com/docs/grafana/latest/panels-visualizations/visualizations/bar-gauge/) for more details.",

		Blocks: map[string]schema.Block{
			"queries":   queryBlock(),
			"field":     fieldBlock(),
			"graph":     barGaugeGraphBlock(),
			"overrides": fieldOverrideBlock(),
			"transform": transformationsBlock(),
		},

		Attributes: map[string]schema.Attribute{
			"id":          idAttribute(),
			"json":        jsonAttribute(),
			"title":       titleAttribute(),
			"description": descriptionAttribute(),
		},
	}
}

func (d *BarGaugeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	defaults, ok := req.ProviderData.(Defaults)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected Defaults, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
	}

	d.Defaults = defaults.BarGauge
}

func (d *BarGaugeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data BarGaugeDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	targets := createTargets(data.Queries)
	fieldConfig := createFieldConfig(d.Defaults.Field, data.Field)
	transformations := createTransformations(data.Transformations)

	options := grafana.Options{
		Orientation: d.Defaults.Graph.Orientation,
		DisplayMode: d.Defaults.Graph.DisplayMode,
		JustifyMode: d.Defaults.Graph.TextAlignment,
		ReduceOptions: grafana.ReduceOptions{
			Values: d.Defaults.Graph.ReduceOptions.Values,
			Fields: d.Defaults.Graph.ReduceOptions.Fields,
			Limit:  d.Defaults.Graph.ReduceOptions.Limit,
			Calcs:  []string{d.Defaults.Graph.ReduceOptions.Calculation},
		},
		TextSize: grafana.TextSize{
			TitleSize: d.Defaults.Graph.TextSize.Title,
			ValueSize: d.Defaults.Graph.TextSize.Value,
		},
	}

	for _, graph := range data.Graph {
		if !graph.Orientation.IsNull() {
			options.Orientation = graph.Orientation.ValueString()
		}

		if !graph.DisplayMode.IsNull() {
			options.DisplayMode = graph.DisplayMode.ValueString()
		}

		if !graph.TextAlignment.IsNull() {
			options.JustifyMode = graph.TextAlignment.ValueString()
		}

		updateTextSize(&options.TextSize, graph.TextSize)
		updateReduceOptions(&options.ReduceOptions, graph.ReduceOptions)
	}

	panel := &grafana.Panel{
		CommonPanel: grafana.CommonPanel{
			OfType:          grafana.BarGaugeType,
			Title:           data.Title.ValueString(),
			Type:            "bargauge",
			Span:            12,
			IsNew:           true,
			Transformations: transformations,
		},
		BarGaugePanel: &grafana.BarGaugePanel{
			Targets: targets,
			Options: options,
			FieldConfig: grafana.FieldConfig{
				Defaults:  fieldConfig,
				Overrides: createOverrides(data.Overrides),
			},
		},
	}

	if !data.Description.IsNull() {
		description := data.Description.ValueString()
		panel.CommonPanel.Description = &description
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
