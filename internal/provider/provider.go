package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure GrafanaDashboardBuilderProvider satisfies various provider interfaces.
var _ provider.Provider = &GrafanaDashboardBuilderProvider{}

// GrafanaDashboardBuilderProvider defines the provider implementation.
type GrafanaDashboardBuilderProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type Defaults struct {
	Dashboard  DashboardDefaults
	Timeseries TimeseriesDefaults
	BarGauge   BarGaugeDefaults
	Stat       StatDefaults
	Gauge      GaugeDefaults
	Table      TableDefaults
}

// GrafanaDashboardBuilderProviderModel describes the provider data model.
type GrafanaDashboardBuilderProviderModel struct {
	Defaults []DefaultsModel `tfsdk:"defaults"`
}

type DefaultsModel struct {
	Dashboard  []DashboardDefaultsModel  `tfsdk:"dashboard"`
	Timeseries []TimeseriesDefaultsModel `tfsdk:"timeseries"`
	BarGuage   []BarGaugeDefaultsModel   `tfsdk:"bar_gauge"`
	Stat       []StatDefaultsModel       `tfsdk:"stat"`
	Gauge      []GaugeDefaultsModel      `tfsdk:"gauge"`
	Table      []TableDefaultsModel      `tfsdk:"table"`
}

type DashboardDefaultsModel struct {
	Editable     types.Bool   `tfsdk:"editable"`
	Style        types.String `tfsdk:"style"`
	GraphTooltip types.String `tfsdk:"graph_tooltip"`
	Time         []TimeModel  `tfsdk:"default_time_range"`
}

type TimeseriesDefaultsModel struct {
	Legend  []TimeseriesLegendOptions  `tfsdk:"legend"`
	Tooltip []TimeseriesTooltipOptions `tfsdk:"tooltip"`
	Field   []FieldOptions             `tfsdk:"field"`
	Axis    []AxisOptions              `tfsdk:"axis"`
	Graph   []TimeseriesGraphOptions   `tfsdk:"graph"`
}

type BarGaugeDefaultsModel struct {
	Field []FieldOptions    `tfsdk:"field"`
	Graph []BarGaugeOptions `tfsdk:"graph"`
}

type StatDefaultsModel struct {
	Field []FieldOptions `tfsdk:"field"`
	Graph []StatOptions  `tfsdk:"graph"`
}

type GaugeDefaultsModel struct {
	Field []FieldOptions `tfsdk:"field"`
	Graph []GaugeOptions `tfsdk:"graph"`
}

type TableDefaultsModel struct {
	Field []FieldOptions `tfsdk:"field"`
}

type TimeModel struct {
	From types.String `tfsdk:"from"`
	To   types.String `tfsdk:"to"`
}

func (p *GrafanaDashboardBuilderProvider) Metadata(_ context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "gdashboard"
	resp.Version = p.version
}

func (p *GrafanaDashboardBuilderProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The provider offers a handy syntax to define Grafana dashboards: time series, gauge, bar gauge, stat, etc.",
		Blocks: map[string]schema.Block{
			"defaults": schema.ListNestedBlock{
				Description: "The default values to use with when an attribute is missing in the data source definition.",
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"dashboard": schema.ListNestedBlock{
							Description: "Dashboard defaults.",
							NestedObject: schema.NestedBlockObject{
								Blocks: map[string]schema.Block{
									"default_time_range": dashboardTimeBlock(),
								},
								Attributes: map[string]schema.Attribute{
									"editable":      dashboardEditableAttribute(),
									"style":         dashboardStyleAttribute(),
									"graph_tooltip": dashboardGraphTooltipAttribute(),
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},
						"timeseries": schema.ListNestedBlock{
							Description: "Timeseries defaults.",
							NestedObject: schema.NestedBlockObject{
								Blocks: map[string]schema.Block{
									"legend":  timeseriesLegendBlock(),
									"tooltip": timeseriesTooltipBlock(),
									"field":   fieldBlock(),
									"axis":    axisBlock(),
									"graph":   timeseriesGraphBlock(),
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},
						"bar_gauge": schema.ListNestedBlock{
							Description: "Bar gauge defaults.",
							NestedObject: schema.NestedBlockObject{
								Blocks: map[string]schema.Block{
									"field": fieldBlock(),
									"graph": barGaugeGraphBlock(),
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},
						"stat": schema.ListNestedBlock{
							Description: "Stat defaults.",
							NestedObject: schema.NestedBlockObject{
								Blocks: map[string]schema.Block{
									"field": fieldBlock(),
									"graph": statGraphBlock(),
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},
						"gauge": schema.ListNestedBlock{
							Description: "Gauge defaults.",
							NestedObject: schema.NestedBlockObject{
								Blocks: map[string]schema.Block{
									"field": fieldBlock(),
									"graph": gaugeGraphBlock(),
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},
						"table": schema.ListNestedBlock{
							Description: "Table defaults.",
							NestedObject: schema.NestedBlockObject{
								Blocks: map[string]schema.Block{
									"field": fieldBlock(),
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
		},
	}
}

func (p *GrafanaDashboardBuilderProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data GrafanaDashboardBuilderProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	defaults := Defaults{
		Dashboard: DashboardDefaults{
			Editable:     true,
			Style:        "dark",
			GraphTooltip: "default",
			Time: Time{
				From: "now-6h",
				To:   "now",
			},
		},
		Timeseries: TimeseriesDefaults{
			Legend: TimeseriesLegendDefault{
				Calculations: nil,
				DisplayMode:  "list",
				Placement:    "bottom",
			},
			Tooltip: TimeseriesTooltipDefaults{
				Mode: "single",
			},
			Field: NewFieldDefaults(),
			Axis: AxisDefaults{
				Label:     "",
				Placement: "auto",
				SoftMin:   nil,
				SoftMax:   nil,
				Scale: ScaleDefaults{
					Type: "linear",
					Log:  0,
				},
			},
			Graph: TimeseriesGraphDefault{
				DrawStyle:         "line",
				LineInterpolation: "linear",
				LineWidth:         1,
				FillOpacity:       0,
				GradientMode:      "none",
				LineStyle:         "solid",
				SpanNulls:         false,
				ShowPoints:        "auto",
				PointSize:         5,
				StackSeries:       "none",
			},
		},
		BarGauge: BarGaugeDefaults{
			Field: NewFieldDefaults(),
			Graph: BarGaugeGraphDefault{
				Orientation:   "auto",
				DisplayMode:   "gradient",
				ReduceOptions: NewReduceOptionDefaults(),
			},
		},
		Stat: StatDefaults{
			Field: NewFieldDefaults(),
			Graph: StatGraphDefaults{
				Orientation:   "auto",
				TextMode:      "auto",
				ColorMode:     "value",
				GraphMode:     "area",
				ReduceOptions: NewReduceOptionDefaults(),
			},
		},
		Gauge: GaugeDefaults{
			Field: NewFieldDefaults(),
			Graph: GaugeGraphDefault{
				Orientation:          "auto",
				ShowThresholdLabels:  false,
				ShowThresholdMarkers: true,
				ReduceOptions:        NewReduceOptionDefaults(),
			},
		},
		Table: TableDefaults{
			Field: NewFieldDefaults(),
		},
	}

	if len(data.Defaults) > 0 && len(data.Defaults[0].Dashboard) > 0 {
		opts := data.Defaults[0].Dashboard[0]

		if !opts.Style.IsNull() {
			defaults.Dashboard.Style = opts.Style.ValueString()
		}

		if !opts.Editable.IsNull() {
			defaults.Dashboard.Editable = opts.Editable.ValueBool()
		}

		if !opts.GraphTooltip.IsNull() {
			defaults.Dashboard.GraphTooltip = opts.GraphTooltip.ValueString()
		}

		for _, time := range opts.Time {
			defaults.Dashboard.Time.From = time.From.ValueString()
			defaults.Dashboard.Time.To = time.To.ValueString()
		}
	}

	if len(data.Defaults) > 0 && len(data.Defaults[0].Timeseries) > 0 {
		opts := data.Defaults[0].Timeseries[0]

		updateFieldDefaults(&defaults.Timeseries.Field, opts.Field)

		for _, graph := range opts.Graph {
			if !graph.DrawStyle.IsNull() {
				defaults.Timeseries.Graph.DrawStyle = graph.DrawStyle.ValueString()
			}

			if !graph.LineInterpolation.IsNull() {
				defaults.Timeseries.Graph.LineInterpolation = graph.LineInterpolation.ValueString()
			}

			if !graph.LineWidth.IsNull() {
				defaults.Timeseries.Graph.LineWidth = int(graph.LineWidth.ValueInt64())
			}

			if !graph.FillOpacity.IsNull() {
				defaults.Timeseries.Graph.FillOpacity = int(graph.FillOpacity.ValueInt64())
			}

			if !graph.GradientMode.IsNull() {
				defaults.Timeseries.Graph.GradientMode = graph.GradientMode.ValueString()
			}

			if !graph.LineStyle.IsNull() {
				defaults.Timeseries.Graph.LineStyle = graph.LineStyle.ValueString()
			}

			if !graph.SpanNulls.IsNull() {
				defaults.Timeseries.Graph.SpanNulls = graph.SpanNulls.ValueBool()
			}

			if !graph.ShowPoints.IsNull() {
				defaults.Timeseries.Graph.ShowPoints = graph.ShowPoints.ValueString()
			}

			if !graph.PointSize.IsNull() {
				defaults.Timeseries.Graph.PointSize = int(graph.PointSize.ValueInt64())
			}

			if !graph.StackSeries.IsNull() {
				defaults.Timeseries.Graph.StackSeries = graph.StackSeries.ValueString()
			}
		}

		for _, legend := range opts.Legend {
			if len(legend.Calculations) > 0 {
				calculations := make([]string, len(legend.Calculations))

				for i, c := range legend.Calculations {
					calculations[i] = c.ValueString()
				}

				defaults.Timeseries.Legend.Calculations = calculations
			}

			if !legend.DisplayMode.IsNull() {
				defaults.Timeseries.Legend.DisplayMode = legend.DisplayMode.ValueString()
			}

			if !legend.Placement.IsNull() {
				defaults.Timeseries.Legend.Placement = legend.Placement.ValueString()
			}
		}

		for _, tooltip := range opts.Tooltip {
			defaults.Timeseries.Tooltip.Mode = tooltip.Mode.ValueString()
		}

		for _, axis := range opts.Axis {
			if !axis.Label.IsNull() {
				defaults.Timeseries.Axis.Label = axis.Label.ValueString()
			}

			if !axis.Placement.IsNull() {
				defaults.Timeseries.Axis.Placement = axis.Placement.ValueString()
			}

			if !axis.SoftMin.IsNull() {
				min := int(axis.SoftMin.ValueInt64())
				defaults.Timeseries.Axis.SoftMin = &min
			}

			if !axis.SoftMax.IsNull() {
				max := int(axis.SoftMax.ValueInt64())
				defaults.Timeseries.Axis.SoftMax = &max
			}

			for _, scale := range axis.Scale {
				if !scale.Type.IsNull() {
					defaults.Timeseries.Axis.Scale.Type = scale.Type.ValueString()
				}

				if !scale.Log.IsNull() {
					defaults.Timeseries.Axis.Scale.Log = int(scale.Log.ValueInt64())
				}
			}
		}
	}

	if len(data.Defaults) > 0 && len(data.Defaults[0].BarGuage) > 0 {
		opts := data.Defaults[0].BarGuage[0]

		updateFieldDefaults(&defaults.BarGauge.Field, opts.Field)

		for _, graph := range opts.Graph {
			if !graph.Orientation.IsNull() {
				defaults.BarGauge.Graph.Orientation = graph.Orientation.ValueString()
			}

			if !graph.DisplayMode.IsNull() {
				defaults.BarGauge.Graph.DisplayMode = graph.DisplayMode.ValueString()
			}

			if !graph.TextAlignment.IsNull() {
				defaults.BarGauge.Graph.TextAlignment = graph.TextAlignment.ValueString()
			}

			updateTextSizeDefaults(&defaults.BarGauge.Graph.TextSize, graph.TextSize)
			updateReduceOptionsDefaults(&defaults.BarGauge.Graph.ReduceOptions, graph.ReduceOptions)
		}
	}

	if len(data.Defaults) > 0 && len(data.Defaults[0].Stat) > 0 {
		opts := data.Defaults[0].Stat[0]

		updateFieldDefaults(&defaults.Stat.Field, opts.Field)

		for _, graph := range opts.Graph {
			if !graph.Orientation.IsNull() {
				defaults.Stat.Graph.Orientation = graph.Orientation.ValueString()
			}

			if !graph.TextMode.IsNull() {
				defaults.Stat.Graph.TextMode = graph.TextMode.ValueString()
			}

			if !graph.ColorMode.IsNull() {
				defaults.Stat.Graph.ColorMode = graph.ColorMode.ValueString()
			}

			if !graph.GraphMode.IsNull() {
				defaults.Stat.Graph.GraphMode = graph.GraphMode.ValueString()
			}

			if !graph.TextAlignment.IsNull() {
				defaults.Stat.Graph.TextAlignment = graph.TextAlignment.ValueString()
			}

			updateTextSizeDefaults(&defaults.Stat.Graph.TextSize, graph.TextSize)
			updateReduceOptionsDefaults(&defaults.Stat.Graph.ReduceOptions, graph.ReduceOptions)
		}
	}

	if len(data.Defaults) > 0 && len(data.Defaults[0].Gauge) > 0 {
		opts := data.Defaults[0].Gauge[0]

		updateFieldDefaults(&defaults.Gauge.Field, opts.Field)

		for _, graph := range opts.Graph {
			if !graph.Orientation.IsNull() {
				defaults.Gauge.Graph.Orientation = graph.Orientation.ValueString()
			}

			if !graph.ShowThresholdLabels.IsNull() {
				defaults.Gauge.Graph.ShowThresholdLabels = graph.ShowThresholdLabels.ValueBool()
			}

			if !graph.ShowThresholdMarkers.IsNull() {
				defaults.Gauge.Graph.ShowThresholdMarkers = graph.ShowThresholdMarkers.ValueBool()
			}

			updateTextSizeDefaults(&defaults.Gauge.Graph.TextSize, graph.TextSize)
			updateReduceOptionsDefaults(&defaults.Gauge.Graph.ReduceOptions, graph.ReduceOptions)
		}
	}

	if len(data.Defaults) > 0 && len(data.Defaults[0].Table) > 0 {
		opts := data.Defaults[0].Table[0]

		updateFieldDefaults(&defaults.Table.Field, opts.Field)
	}

	resp.DataSourceData = defaults
	resp.ResourceData = defaults
}

func updateFieldDefaults(defaults *FieldDefaults, opts []FieldOptions) {
	for _, field := range opts {
		if !field.Unit.IsNull() {
			defaults.Unit = field.Unit.ValueString()
		}

		if !field.Decimals.IsNull() {
			decimals := int(field.Decimals.ValueInt64())
			defaults.Decimals = &decimals
		}

		if !field.Min.IsNull() {
			min := field.Min.ValueFloat64()
			defaults.Min = &min
		}

		if !field.Max.IsNull() {
			max := field.Max.ValueFloat64()
			defaults.Max = &max
		}

		if !field.NoValue.IsNull() {
			noValue := field.NoValue.ValueFloat64()
			defaults.NoValue = &noValue
		}

		for _, color := range field.Color {
			if !color.Mode.IsNull() {
				defaults.Color.Mode = color.Mode.ValueString()
			}

			if !color.FixedColor.IsNull() {
				defaults.Color.FixedColor = color.FixedColor.ValueString()
			}

			if !color.SeriesBy.IsNull() {
				defaults.Color.SeriesBy = color.SeriesBy.ValueString()
			}
		}

		for _, threshold := range field.Thresholds {
			steps := make([]ThresholdStepDefaults, len(threshold.Steps))

			if !threshold.Mode.IsNull() {
				defaults.Thresholds.Mode = threshold.Mode.ValueString()
			}

			for i, step := range threshold.Steps {
				s := ThresholdStepDefaults{
					Color: step.Color.ValueString(),
				}

				if !step.Value.IsNull() {
					value := step.Value.ValueFloat64()
					s.Value = &value
				}

				steps[i] = s
			}

			defaults.Thresholds.Steps = steps
		}
	}
}

func updateTextSizeDefaults(defaults *TextSizeDefaults, opts []TextSizeOptions) {
	for _, textSize := range opts {
		if !textSize.Title.IsNull() {
			size := int(textSize.Title.ValueInt64())
			defaults.Title = &size
		}

		if !textSize.Value.IsNull() {
			size := int(textSize.Value.ValueInt64())
			defaults.Value = &size
		}
	}
}

func updateReduceOptionsDefaults(defaults *ReduceOptionDefaults, opts []ReduceOptions) {
	for _, reducer := range opts {
		if !reducer.Values.IsNull() {
			defaults.Values = reducer.Values.ValueBool()
		}

		if !reducer.Fields.IsNull() {
			defaults.Fields = reducer.Fields.ValueString()
		}

		if !reducer.Limit.IsNull() {
			limit := int(reducer.Limit.ValueInt64())
			defaults.Limit = &limit
		}

		if !reducer.Calculation.IsNull() {
			defaults.Calculation = reducer.Calculation.ValueString()
		}
	}
}

func (p *GrafanaDashboardBuilderProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}

func (p *GrafanaDashboardBuilderProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewTimeseriesDataSource,
		NewDashboardDataSource,
		NewBarGaugeDataSource,
		NewStatDataSource,
		NewRowDataSource,
		NewGaugeDataSource,
		NewTextDataSource,
		NewTableDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &GrafanaDashboardBuilderProvider{
			version: version,
		}
	}
}
