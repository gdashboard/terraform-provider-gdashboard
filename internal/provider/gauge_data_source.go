package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/iRevive/terraform-provider-gdashboard/internal/provider/grafana"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &GaugeDataSource{}

func NewGaugeDataSource() datasource.DataSource {
	return &GaugeDataSource{}
}

// GaugeDataSource defines the data source implementation.
type GaugeDataSource struct {
	Defaults GaugeDefaults
}

type GaugeDefaults struct {
	Field FieldDefaults
	Graph GaugeGraphDefault
}

type GaugeGraphDefault struct {
	Orientation          string
	ShowThresholdLabels  bool
	ShowThresholdMarkers bool
	TextSize             TextSizeDefaults
	ReduceOptions        ReduceOptionDefaults
}

// GaugeDataSourceModel describes the data source data model.
type GaugeDataSourceModel struct {
	Id              types.String           `tfsdk:"id"`
	Json            types.String           `tfsdk:"json"`
	Title           types.String           `tfsdk:"title"`
	Description     types.String           `tfsdk:"description"`
	Queries         []Query                `tfsdk:"queries"`
	Field           []FieldOptions         `tfsdk:"field"`
	Graph           []GaugeOptions         `tfsdk:"graph"`
	Overrides       []FieldOverrideOptions `tfsdk:"overrides"`
	Transformations []Transformations      `tfsdk:"transform"`
}

type GaugeOptions struct {
	Orientation          types.String      `tfsdk:"orientation"`
	ShowThresholdLabels  types.Bool        `tfsdk:"show_threshold_labels"`
	ShowThresholdMarkers types.Bool        `tfsdk:"show_threshold_markers"`
	TextSize             []TextSizeOptions `tfsdk:"text_size"`
	ReduceOptions        []ReduceOptions   `tfsdk:"options"`
}

func (d *GaugeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gauge"
}

func gaugeGraphBlock() schema.Block {
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
				"show_threshold_labels": schema.BoolAttribute{
					Optional:    true,
					Description: "Whether to render the threshold values around the gauge bar or not.",
				},
				"show_threshold_markers": schema.BoolAttribute{
					Optional:    true,
					Description: "Whether to render the thresholds as an outer bar or not.",
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func (d *GaugeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description:         "Gauge panel data source.",
		MarkdownDescription: "Gauge panel data source. See Grafana [documentation](https://grafana.com/docs/grafana/latest/panels-visualizations/visualizations/gauge/). for more details",

		Blocks: map[string]schema.Block{
			"queries":   queryBlock(),
			"field":     fieldBlock(),
			"graph":     gaugeGraphBlock(),
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

func (d *GaugeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.Defaults = defaults.Gauge
}

func (d *GaugeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GaugeDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	targets := createTargets(data.Queries)
	fieldConfig := createFieldConfig(d.Defaults.Field, data.Field)
	transformations := createTransformations(data.Transformations)

	options := grafana.Options{
		Orientation:          d.Defaults.Graph.Orientation,
		ShowThresholdLabels:  &d.Defaults.Graph.ShowThresholdLabels,
		ShowThresholdMarkers: &d.Defaults.Graph.ShowThresholdMarkers,
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

		if !graph.ShowThresholdLabels.IsNull() {
			options.ShowThresholdLabels = graph.ShowThresholdLabels.ValueBoolPointer()
		}

		if !graph.ShowThresholdMarkers.IsNull() {
			options.ShowThresholdMarkers = graph.ShowThresholdMarkers.ValueBoolPointer()
		}

		updateTextSize(&options.TextSize, graph.TextSize)
		updateReduceOptions(&options.ReduceOptions, graph.ReduceOptions)
	}

	panel := &grafana.Panel{
		CommonPanel: grafana.CommonPanel{
			OfType:          grafana.GaugeType,
			Title:           data.Title.ValueString(),
			Type:            "gauge",
			Span:            12,
			IsNew:           true,
			Transformations: transformations,
		},
		GaugePanel: &grafana.GaugePanel{
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
