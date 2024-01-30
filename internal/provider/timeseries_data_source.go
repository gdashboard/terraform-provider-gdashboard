package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"strconv"

	"github.com/gdashboard/terraform-provider-gdashboard/internal/provider/grafana"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &TimeseriesDataSource{}

func NewTimeseriesDataSource() datasource.DataSource {
	return &TimeseriesDataSource{}
}

// TimeseriesDataSource defines the data source implementation.
type TimeseriesDataSource struct {
	CompactJson bool
	Defaults    TimeseriesDefaults
}

type TimeseriesDefaults struct {
	Legend  TimeseriesLegendDefault
	Tooltip TimeseriesTooltipDefaults
	Field   FieldDefaults
	Axis    AxisDefaults
	Graph   TimeseriesGraphDefault
}

type TimeseriesGraphDefault struct {
	DrawStyle         string
	LineInterpolation string
	LineWidth         int
	FillOpacity       int
	GradientMode      string
	LineStyle         string
	SpanNulls         bool
	ShowPoints        string
	PointSize         int
	StackSeries       string
}

type TimeseriesTooltipDefaults struct {
	Mode string
}

type TimeseriesLegendDefault struct {
	Calculations []string
	DisplayMode  string
	Placement    string
}

// TimeseriesDataSourceModel describes the data source data model.
type TimeseriesDataSourceModel struct {
	Id              types.String               `tfsdk:"id"`
	Json            types.String               `tfsdk:"json"`
	CompactJson     types.Bool                 `tfsdk:"compact_json"`
	Title           types.String               `tfsdk:"title"`
	Description     types.String               `tfsdk:"description"`
	Queries         []Query                    `tfsdk:"queries"`
	Legend          []TimeseriesLegendOptions  `tfsdk:"legend"`
	Tooltip         []TimeseriesTooltipOptions `tfsdk:"tooltip"`
	Field           []FieldOptions             `tfsdk:"field"`
	Axis            []AxisOptions              `tfsdk:"axis"`
	Graph           []TimeseriesGraphOptions   `tfsdk:"graph"`
	Overrides       []FieldOverrideOptions     `tfsdk:"overrides"`
	Transformations []Transformations          `tfsdk:"transform"`
}

type TimeseriesLegendOptions struct {
	Calculations []types.String `tfsdk:"calculations"`
	DisplayMode  types.String   `tfsdk:"display_mode"`
	Placement    types.String   `tfsdk:"placement"`
}

type TimeseriesTooltipOptions struct {
	Mode types.String `tfsdk:"mode"`
}

type TimeseriesGraphOptions struct {
	DrawStyle         types.String `tfsdk:"draw_style"`
	LineInterpolation types.String `tfsdk:"line_interpolation"`
	LineWidth         types.Int64  `tfsdk:"line_width"`
	FillOpacity       types.Int64  `tfsdk:"fill_opacity"`
	GradientMode      types.String `tfsdk:"gradient_mode"`
	LineStyle         types.String `tfsdk:"line_style"`
	SpanNulls         types.Bool   `tfsdk:"span_nulls"`
	ShowPoints        types.String `tfsdk:"show_points"`
	PointSize         types.Int64  `tfsdk:"point_size"`
	StackSeries       types.String `tfsdk:"stack_series"`
}

func (d *TimeseriesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_timeseries"
}

func timeseriesTooltipBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: "The tooltip visualization options.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"mode": schema.StringAttribute{
					Required:            true,
					Description:         "Choose the how to display the tooltip. The choices are: multi, single, hidden.",
					MarkdownDescription: "Choose the how to display the tooltip. The choices are: `multi`, `single`, `hidden`.",
					Validators: []validator.String{
						stringvalidator.OneOf("multi", "single", "hidden"),
					},
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func timeseriesGraphBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: "The visualization options.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"draw_style": schema.StringAttribute{
					Optional:            true,
					Description:         "Choose the visualization style. The choices are: line, bars, points.",
					MarkdownDescription: "Choose the visualization style. The choices are: `line`, `bars`, `points`.",
					Validators: []validator.String{
						stringvalidator.OneOf("line", "bars", "points"),
					},
				},
				"line_interpolation": schema.StringAttribute{
					Optional:            true,
					Description:         "Choose how to interpolation the line. The choices are: linear, smooth, stepBefore, stepAfter.",
					MarkdownDescription: "Choose how to interpolation the line. The choices are: `linear`, `smooth`, `stepBefore`, `stepAfter`.",
					Validators: []validator.String{
						stringvalidator.OneOf("linear", "smooth", "stepBefore", "stepAfter"),
					},
				},
				"line_width": schema.Int64Attribute{
					Optional:            true,
					Description:         "The width of the line. Must be between 0 and 10 (inclusive).",
					MarkdownDescription: "The width of the line. Must be between `0` and `10` (inclusive).",
					Validators: []validator.Int64{
						int64validator.Between(0, 10),
					},
				},
				"fill_opacity": schema.Int64Attribute{
					Optional:            true,
					Description:         "The opacity of the filled areas. Must be between 0 and 100 (inclusive).",
					MarkdownDescription: "The opacity of the filled areas. Must be between `0` and `100` (inclusive).",
					Validators: []validator.Int64{
						int64validator.Between(0, 100),
					},
				},
				"gradient_mode": schema.StringAttribute{
					Optional:            true,
					Description:         "The gradient mode. The choices are: none, opacity, hue, scheme.",
					MarkdownDescription: "The gradient mode. The choices are: `none`, `opacity`, `hue`, `scheme`.",
					Validators: []validator.String{
						stringvalidator.OneOf("none", "opacity", "hue", "scheme"),
					},
				},
				"line_style": schema.StringAttribute{
					Optional:            true,
					Description:         "The style of the line. The choices are: solid, dash, dots.",
					MarkdownDescription: "The style of the line. The choices are: `solid`, `dash`, `dots`.",
					Validators: []validator.String{
						stringvalidator.OneOf("solid", "dash", "dots"),
					},
				},
				"span_nulls": schema.BoolAttribute{
					Optional:    true,
					Description: "Whether to ignore or replace null values with zeroes or not.",
				},
				"show_points": schema.StringAttribute{
					Optional:            true,
					Description:         "Choose how to display data points. The choices are: auto, never, always.",
					MarkdownDescription: "Choose how to display data points. The choices are: `auto`, `never`, `always`.",
					Validators: []validator.String{
						stringvalidator.OneOf("auto", "never", "always"),
					},
				},
				"point_size": schema.Int64Attribute{
					Optional:            true,
					Description:         "The size of the data point. Must be between 1 and 40 (inclusive).",
					MarkdownDescription: "The size of the data point. Must be between `1` and `40` (inclusive).",
					Validators: []validator.Int64{
						int64validator.Between(1, 40),
					},
				},
				"stack_series": schema.StringAttribute{
					Optional:            true,
					Description:         "Choose how to stack the series. The choices are: none, normal, percent.",
					MarkdownDescription: "Choose how to stack the series. The choices are: `none`, `normal`, `percent`.",
					Validators: []validator.String{
						stringvalidator.OneOf("none", "normal", "percent"),
					},
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func timeseriesLegendBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: "Legend options.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"calculations": schema.ListAttribute{
					ElementType: types.StringType,
					Optional:    true,
					Description: "Choose which of the standard calculations to show in the legend: min, max, mean, etc.",
				},
				"display_mode": schema.StringAttribute{
					Optional:            true,
					Description:         "Choose how to display the legend. The choices are: list, table, hidden.",
					MarkdownDescription: "Choose how to display the legend. The choices are: `list`, `table`, `hidden`.",
					Validators: []validator.String{
						stringvalidator.OneOf("list", "table", "hidden"),
					},
				},
				"placement": schema.StringAttribute{
					Optional:            true,
					Description:         "Choose where to display the legend. The choice are: bottom, right.",
					MarkdownDescription: "Choose where to display the legend. The choice are: `bottom`, `right`.",
					Validators: []validator.String{
						stringvalidator.OneOf("bottom", "right"),
					},
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func (d *TimeseriesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description:         "Time series panel data source.",
		MarkdownDescription: "Time series panel data source. See Grafana [documentation](https://grafana.com/docs/grafana/latest/panels-visualizations/visualizations/time-series/).",

		Blocks: map[string]schema.Block{
			"queries":   queryBlock(),
			"legend":    timeseriesLegendBlock(),
			"tooltip":   timeseriesTooltipBlock(),
			"field":     fieldBlock(),
			"axis":      axisBlock(),
			"graph":     timeseriesGraphBlock(),
			"overrides": fieldOverrideBlock(),
			"transform": transformationsBlock(),
		},

		Attributes: map[string]schema.Attribute{
			"id":           idAttribute(),
			"json":         jsonAttribute(),
			"compact_json": compactJsonAttribute(),
			"title":        titleAttribute(),
			"description":  descriptionAttribute(),
		},
	}
}

func (d *TimeseriesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.CompactJson = defaults.CompactJson
	d.Defaults = defaults.Timeseries
}

func (d *TimeseriesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TimeseriesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	targets := createTargets(data.Queries)
	transformations := createTransformations(data.Transformations)

	legendOptions := grafana.TimeseriesLegendOptions{
		Calcs:       d.Defaults.Legend.Calculations,
		DisplayMode: d.Defaults.Legend.DisplayMode,
		Placement:   d.Defaults.Legend.Placement,
	}

	tooltipOptions := grafana.TimeseriesTooltipOptions{
		Mode: d.Defaults.Tooltip.Mode,
	}

	for _, legend := range data.Legend {
		if len(legend.Calculations) > 0 {
			calculations := make([]string, len(legend.Calculations))
			for i, calc := range legend.Calculations {
				calculations[i] = calc.ValueString()
			}

			legendOptions.Calcs = calculations
		}

		if !legend.DisplayMode.IsNull() {
			legendOptions.DisplayMode = legend.DisplayMode.ValueString()
		}

		if !legend.Placement.IsNull() {
			legendOptions.Placement = legend.Placement.ValueString()
		}
	}

	for _, tooltip := range data.Tooltip {
		tooltipOptions.Mode = tooltip.Mode.ValueString()
	}

	fieldConfig := createFieldConfig(d.Defaults.Field, data.Field)

	fieldConfig.Custom = grafana.FieldConfigCustom{
		// graph
		DrawStyle:         d.Defaults.Graph.DrawStyle,
		LineInterpolation: d.Defaults.Graph.LineInterpolation,
		LineWidth:         d.Defaults.Graph.LineWidth,
		FillOpacity:       d.Defaults.Graph.FillOpacity,
		GradientMode:      d.Defaults.Graph.GradientMode,
		SpanNulls:         d.Defaults.Graph.SpanNulls,
		ShowPoints:        d.Defaults.Graph.ShowPoints,
		PointSize:         d.Defaults.Graph.PointSize,
		// axis
		AxisLabel:     d.Defaults.Axis.Label,
		AxisPlacement: d.Defaults.Axis.Placement,
		AxisSoftMin:   d.Defaults.Axis.SoftMin,
		AxisSoftMax:   d.Defaults.Axis.SoftMax,
	}

	fieldConfig.Custom.LineStyle.Fill = d.Defaults.Graph.LineStyle
	fieldConfig.Custom.Stacking.Mode = d.Defaults.Graph.StackSeries
	fieldConfig.Custom.ScaleDistribution.Type = d.Defaults.Axis.Scale.Type
	fieldConfig.Custom.ScaleDistribution.Log = d.Defaults.Axis.Scale.Log

	for _, axis := range data.Axis {
		if !axis.Label.IsNull() {
			fieldConfig.Custom.AxisLabel = axis.Label.ValueString()
		}

		if !axis.Placement.IsNull() {
			fieldConfig.Custom.AxisPlacement = axis.Placement.ValueString()
		}

		if !axis.SoftMin.IsNull() {
			fieldConfig.Custom.AxisSoftMin = axis.SoftMin.ValueInt64Pointer()
		}

		if !axis.SoftMax.IsNull() {
			fieldConfig.Custom.AxisSoftMax = axis.SoftMax.ValueInt64Pointer()
		}

		for _, scale := range axis.Scale {
			if !scale.Type.IsNull() {
				fieldConfig.Custom.ScaleDistribution.Type = scale.Type.ValueString()
			}

			if !scale.Log.IsNull() {
				fieldConfig.Custom.ScaleDistribution.Log = int(scale.Log.ValueInt64())
			}
		}
	}

	for _, graph := range data.Graph {
		if !graph.DrawStyle.IsNull() {
			fieldConfig.Custom.DrawStyle = graph.DrawStyle.ValueString()
		}

		if !graph.LineInterpolation.IsNull() {
			fieldConfig.Custom.LineInterpolation = graph.LineInterpolation.ValueString()
		}

		if !graph.LineWidth.IsNull() {
			fieldConfig.Custom.LineWidth = int(graph.LineWidth.ValueInt64())
		}

		if !graph.FillOpacity.IsNull() {
			fieldConfig.Custom.FillOpacity = int(graph.FillOpacity.ValueInt64())
		}

		if !graph.GradientMode.IsNull() {
			fieldConfig.Custom.GradientMode = graph.GradientMode.ValueString()
		}

		if !graph.LineStyle.IsNull() {
			fieldConfig.Custom.LineStyle.Fill = graph.LineStyle.ValueString()
		}

		if !graph.SpanNulls.IsNull() {
			fieldConfig.Custom.SpanNulls = graph.SpanNulls.ValueBool()
		}

		if !graph.ShowPoints.IsNull() {
			fieldConfig.Custom.ShowPoints = graph.ShowPoints.ValueString()
		}

		if !graph.PointSize.IsNull() {
			fieldConfig.Custom.PointSize = int(graph.PointSize.ValueInt64())
		}

		if !graph.StackSeries.IsNull() {
			fieldConfig.Custom.Stacking.Mode = graph.StackSeries.ValueString()
		}
	}

	panel := &grafana.Panel{
		CommonPanel: grafana.CommonPanel{
			OfType:          grafana.TimeseriesType,
			Title:           data.Title.ValueString(),
			Type:            "timeseries",
			Span:            12,
			IsNew:           true,
			Transformations: transformations,
		},
		TimeseriesPanel: &grafana.TimeseriesPanel{
			Targets: targets,
			Options: grafana.TimeseriesOptions{
				Legend:  legendOptions,
				Tooltip: tooltipOptions,
			},
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

	var jsonData []byte
	var err error

	if data.CompactJson.ValueBool() || d.CompactJson {
		jsonData, err = json.Marshal(panel)
	} else {
		jsonData, err = json.MarshalIndent(panel, "", "  ")
	}

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
