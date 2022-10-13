package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure GrafanaDashboardBuilderProvider satisfies various provider interfaces.
var _ provider.Provider = &GrafanaDashboardBuilderProvider{}
var _ provider.ProviderWithMetadata = &GrafanaDashboardBuilderProvider{}

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
}

type DashboardDefaultsModel struct {
	Editable     types.Bool   `tfsdk:"editable"`
	Style        types.String `tfsdk:"style"`
	GraphTooltip types.String `tfsdk:"graph_tooltip"`
	Time         []TimeModel  `tfsdk:"time"`
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

type TimeModel struct {
	From types.String `tfsdk:"from"`
	To   types.String `tfsdk:"to"`
}

func (p *GrafanaDashboardBuilderProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "gdashboard"
	resp.Version = p.version
}

func (p *GrafanaDashboardBuilderProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Blocks: map[string]tfsdk.Block{
			"defaults": {
				NestingMode: tfsdk.BlockNestingModeList,
				MinItems:    0,
				MaxItems:    1,
				Blocks: map[string]tfsdk.Block{
					"dashboard": {
						NestingMode: tfsdk.BlockNestingModeList,
						MinItems:    0,
						MaxItems:    1,
						Blocks: map[string]tfsdk.Block{
							"time": dashboardTimeBlock(),
						},
						Attributes: map[string]tfsdk.Attribute{
							"editable":      dashboardEditableAttribute(),
							"style":         dashboardStyleAttribute(),
							"graph_tooltip": dashboardGraphTooltipAttribute(),
						},
					},
					"timeseries": {
						NestingMode: tfsdk.BlockNestingModeList,
						MinItems:    0,
						MaxItems:    1,
						Blocks: map[string]tfsdk.Block{
							"legend":  timeseriesLegendBlock(),
							"tooltip": timeseriesTooltipBlock(),
							"field":   fieldBlock(),
							"axis":    axisBlock(),
							"graph":   timeseriesGraphBlock(),
						},
					},
					"bar_gauge": {
						NestingMode: tfsdk.BlockNestingModeList,
						MinItems:    0,
						MaxItems:    1,
						Blocks: map[string]tfsdk.Block{
							"field": fieldBlock(),
							"graph": barGaugeGraphBlock(),
						},
					},
					"stat": {
						NestingMode: tfsdk.BlockNestingModeList,
						MinItems:    0,
						MaxItems:    1,
						Blocks: map[string]tfsdk.Block{
							"field": fieldBlock(),
							"graph": statGraphBlock(),
						},
					},
					"gauge": {
						NestingMode: tfsdk.BlockNestingModeList,
						MinItems:    0,
						MaxItems:    1,
						Blocks: map[string]tfsdk.Block{
							"field": fieldBlock(),
							"graph": gaugeGraphBlock(),
						},
					},
				},
			},
		},
	}, nil
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
	}

	if len(data.Defaults) > 0 && len(data.Defaults[0].Dashboard) > 0 {
		opts := data.Defaults[0].Dashboard[0]

		if !opts.Style.Null {
			defaults.Dashboard.Style = opts.Style.Value
		}

		if !opts.Editable.Null {
			defaults.Dashboard.Editable = opts.Editable.Value
		}

		if !opts.GraphTooltip.Null {
			defaults.Dashboard.GraphTooltip = opts.GraphTooltip.Value
		}

		for _, time := range opts.Time {
			defaults.Dashboard.Time.From = time.From.Value
			defaults.Dashboard.Time.To = time.To.Value
		}
	}

	if len(data.Defaults) > 0 && len(data.Defaults[0].Timeseries) > 0 {
		opts := data.Defaults[0].Timeseries[0]

		updateFieldDefaults(&defaults.Timeseries.Field, opts.Field)

		for _, graph := range opts.Graph {
			if !graph.DrawStyle.Null {
				defaults.Timeseries.Graph.DrawStyle = graph.DrawStyle.Value
			}

			if !graph.LineInterpolation.Null {
				defaults.Timeseries.Graph.LineInterpolation = graph.LineInterpolation.Value
			}

			if !graph.LineWidth.Null {
				defaults.Timeseries.Graph.LineWidth = int(graph.LineWidth.Value)
			}

			if !graph.FillOpacity.Null {
				defaults.Timeseries.Graph.FillOpacity = int(graph.FillOpacity.Value)
			}

			if !graph.GradientMode.Null {
				defaults.Timeseries.Graph.GradientMode = graph.GradientMode.Value
			}

			if !graph.LineStyle.Null {
				defaults.Timeseries.Graph.LineStyle = graph.LineStyle.Value
			}

			if !graph.SpanNulls.Null {
				defaults.Timeseries.Graph.SpanNulls = graph.SpanNulls.Value
			}

			if !graph.ShowPoints.Null {
				defaults.Timeseries.Graph.ShowPoints = graph.ShowPoints.Value
			}

			if !graph.PointSize.Null {
				defaults.Timeseries.Graph.PointSize = int(graph.PointSize.Value)
			}

			if !graph.StackSeries.Null {
				defaults.Timeseries.Graph.StackSeries = graph.StackSeries.Value
			}
		}

		for _, legend := range opts.Legend {
			if len(legend.Calculations) > 0 {
				calculations := make([]string, len(legend.Calculations))

				for i, c := range legend.Calculations {
					calculations[i] = c.Value
				}

				defaults.Timeseries.Legend.Calculations = calculations
			}

			if !legend.DisplayMode.Null {
				defaults.Timeseries.Legend.DisplayMode = legend.DisplayMode.Value
			}

			if !legend.Placement.Null {
				defaults.Timeseries.Legend.Placement = legend.Placement.Value
			}
		}

		for _, tooltip := range opts.Tooltip {
			defaults.Timeseries.Tooltip.Mode = tooltip.Mode.Value
		}

		for _, axis := range opts.Axis {
			if !axis.Label.Null {
				defaults.Timeseries.Axis.Label = axis.Label.Value
			}

			if !axis.Placement.Null {
				defaults.Timeseries.Axis.Placement = axis.Placement.Value
			}

			if !axis.SoftMin.Null {
				min := int(axis.SoftMin.Value)
				defaults.Timeseries.Axis.SoftMin = &min
			}

			if !axis.SoftMax.Null {
				max := int(axis.SoftMax.Value)
				defaults.Timeseries.Axis.SoftMax = &max
			}

			for _, scale := range axis.Scale {
				if !scale.Type.Null {
					defaults.Timeseries.Axis.Scale.Type = scale.Type.Value
				}

				if !scale.Log.Null {
					defaults.Timeseries.Axis.Scale.Log = int(scale.Log.Value)
				}
			}
		}
	}

	if len(data.Defaults) > 0 && len(data.Defaults[0].BarGuage) > 0 {
		opts := data.Defaults[0].BarGuage[0]

		updateFieldDefaults(&defaults.BarGauge.Field, opts.Field)

		for _, graph := range opts.Graph {
			if !graph.Orientation.Null {
				defaults.BarGauge.Graph.Orientation = graph.Orientation.Value
			}

			if !graph.DisplayMode.Null {
				defaults.BarGauge.Graph.DisplayMode = graph.DisplayMode.Value
			}

			if !graph.TextAlignment.Null {
				defaults.BarGauge.Graph.TextAlignment = graph.TextAlignment.Value
			}

			updateTextSizeDefaults(&defaults.BarGauge.Graph.TextSize, graph.TextSize)
			updateReduceOptionsDefaults(&defaults.BarGauge.Graph.ReduceOptions, graph.ReduceOptions)
		}
	}

	if len(data.Defaults) > 0 && len(data.Defaults[0].Stat) > 0 {
		opts := data.Defaults[0].Stat[0]

		updateFieldDefaults(&defaults.Stat.Field, opts.Field)

		for _, graph := range opts.Graph {
			if !graph.Orientation.Null {
				defaults.Stat.Graph.Orientation = graph.Orientation.Value
			}

			if !graph.TextMode.Null {
				defaults.Stat.Graph.TextMode = graph.TextMode.Value
			}

			if !graph.ColorMode.Null {
				defaults.Stat.Graph.ColorMode = graph.ColorMode.Value
			}

			if !graph.GraphMode.Null {
				defaults.Stat.Graph.GraphMode = graph.GraphMode.Value
			}

			if !graph.TextAlignment.Null {
				defaults.Stat.Graph.TextAlignment = graph.TextAlignment.Value
			}

			updateTextSizeDefaults(&defaults.Stat.Graph.TextSize, graph.TextSize)
			updateReduceOptionsDefaults(&defaults.Stat.Graph.ReduceOptions, graph.ReduceOptions)
		}
	}

	if len(data.Defaults) > 0 && len(data.Defaults[0].Gauge) > 0 {
		opts := data.Defaults[0].Gauge[0]

		updateFieldDefaults(&defaults.Gauge.Field, opts.Field)

		for _, graph := range opts.Graph {
			if !graph.Orientation.Null {
				defaults.Gauge.Graph.Orientation = graph.Orientation.Value
			}

			if !graph.ShowThresholdLabels.Null {
				defaults.Gauge.Graph.ShowThresholdLabels = graph.ShowThresholdLabels.Value
			}

			if !graph.ShowThresholdMarkers.Null {
				defaults.Gauge.Graph.ShowThresholdMarkers = graph.ShowThresholdMarkers.Value
			}

			updateTextSizeDefaults(&defaults.Gauge.Graph.TextSize, graph.TextSize)
			updateReduceOptionsDefaults(&defaults.Gauge.Graph.ReduceOptions, graph.ReduceOptions)
		}
	}

	resp.DataSourceData = defaults
	resp.ResourceData = defaults
}

func updateFieldDefaults(defaults *FieldDefaults, opts []FieldOptions) {
	for _, field := range opts {
		if !field.Unit.Null {
			defaults.Unit = field.Unit.Value
		}

		if !field.Decimals.Null {
			decimals := int(field.Decimals.Value)
			defaults.Decimals = &decimals
		}

		if !field.Min.Null {
			defaults.Min = &field.Min.Value
		}

		if !field.Max.Null {
			defaults.Max = &field.Max.Value
		}

		if !field.NoValue.Null {
			defaults.NoValue = &field.NoValue.Value
		}

		for _, color := range field.Color {
			if !color.Mode.Null {
				defaults.Color.Mode = color.Mode.Value
			}

			if !color.FixedColor.Null {
				defaults.Color.FixedColor = color.FixedColor.Value
			}

			if !color.SeriesBy.Null {
				defaults.Color.SeriesBy = color.SeriesBy.Value
			}
		}

		for _, threshold := range field.Thresholds {
			steps := make([]ThresholdStepDefaults, len(threshold.Steps))

			if !threshold.Mode.Null {
				defaults.Thresholds.Mode = threshold.Mode.Value
			}

			for i, step := range threshold.Steps {
				s := ThresholdStepDefaults{
					Color: step.Color.Value,
				}

				if !step.Value.Null {
					value := step.Value.Value
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
		if !textSize.Title.Null {
			size := int(textSize.Title.Value)
			defaults.Title = &size
		}

		if !textSize.Value.Null {
			size := int(textSize.Value.Value)
			defaults.Value = &size
		}
	}
}

func updateReduceOptionsDefaults(defaults *ReduceOptionDefaults, opts []ReduceOptions) {
	for _, reducer := range opts {
		if !reducer.Values.Null {
			defaults.Values = reducer.Values.Value
		}

		if !reducer.Fields.Null {
			defaults.Fields = reducer.Fields.Value
		}

		if !reducer.Limit.Null {
			limit := int(reducer.Limit.Value)
			defaults.Limit = &limit
		}

		if !reducer.Calculation.Null {
			defaults.Calculation = reducer.Calculation.Value
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
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &GrafanaDashboardBuilderProvider{
			version: version,
		}
	}
}
