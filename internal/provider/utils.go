package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/iRevive/terraform-provider-gdashboard/internal/provider/grafana"
	"hash/crc32"
	"strings"
)

// constants

// CalculationTypes The following table contains a list of calculations you can perform in Grafana.
// You can find these calculations in the Transform tab and in the bar gauge, gauge, and stat visualizations.
func CalculationTypes() []string {
	return []string{
		"lastNotNull",   // Last (not null) - Last, not null value
		"last",          // Last - Last value
		"firstNotNull",  // First (not null) - First, not null value
		"first",         // First - First value
		"min",           // Min - Minimum value of a field
		"max",           // Max - Maximum value of a field
		"mean",          // Mean - Mean value of all values in a field
		"sum",           // Total - Sum of all values in a field
		"count",         // Count - Number of values in a field
		"range",         // Range - Difference between maximum and minimum values of a field
		"delta",         // Delta - Cumulative change in value, only counts increments
		"step",          // Step - Minimal interval between values of a field
		"diff",          // Difference - Difference between first and last value of a field
		"diffperc",      // Difference percent - Percentage change between first and last value of a field
		"logmin",        // Min (above zero) - Minimum, positive value of a field
		"allIsZero",     // All zeros - True when all values are 0
		"allIsNull",     // All nulls - True when all values are null
		"changeCount",   // Change count - Number of times the field’s value changes
		"distinctCount", // Distinct count - Number of unique values in a field
		"stdDev",        // StdDev - Standard deviation of all values in a field
		"variance",      // Variance - Variance of all values in a field
		"allValues",     // All values - Returns an array with all values
		"uniqueValues",  // All unique values - Returns an array with all unique values
	}
}

func CalculationTypesString() string {
	return strings.Join(CalculationTypes(), ", ")
}

func CalculationTypesMarkdown() string {
	res := make([]string, 0)
	for _, tpe := range CalculationTypes() {
		res = append(res, "`"+tpe+"`")
	}

	return strings.Join(res, ", ")
}

// defaults

type FieldDefaults struct {
	Unit       string
	Decimals   *int64
	Min        *float64
	Max        *float64
	NoValue    *float64
	Color      ColorDefaults
	Thresholds ThresholdDefaults
}

func NewFieldDefaults() FieldDefaults {
	return FieldDefaults{
		Unit:     "",
		Decimals: nil,
		Min:      nil,
		Max:      nil,
		Color: ColorDefaults{
			Mode:       "palette-classic",
			FixedColor: "green",
			SeriesBy:   "last",
		},
		Thresholds: ThresholdDefaults{
			Mode: "absolute",
			Steps: []ThresholdStepDefaults{
				{
					Color: "green",
					Value: nil,
				},
			},
		},
	}
}

type ColorDefaults struct {
	Mode       string
	FixedColor string
	SeriesBy   string
}

type ThresholdDefaults struct {
	Mode  string
	Steps []ThresholdStepDefaults
}

type ThresholdStepDefaults struct {
	Color string
	Value *float64
}

type ReduceOptionDefaults struct {
	Values      bool
	Fields      string
	Limit       *int64
	Calculation string
}

func NewReduceOptionDefaults() ReduceOptionDefaults {
	return ReduceOptionDefaults{
		Values:      false,
		Fields:      "",
		Calculation: "lastNotNull",
	}
}

type TextSizeDefaults struct {
	Title *int64
	Value *int64
}

type AxisDefaults struct {
	Label     string
	Placement string
	SoftMin   *int64
	SoftMax   *int64
	Scale     ScaleDefaults
}

type ScaleDefaults struct {
	Type string
	Log  int
}

// Terraform projections

type AxisOptions struct {
	Label     types.String   `tfsdk:"label"`
	Placement types.String   `tfsdk:"placement"`
	SoftMin   types.Int64    `tfsdk:"soft_min"`
	SoftMax   types.Int64    `tfsdk:"soft_max"`
	Scale     []ScaleOptions `tfsdk:"scale"`
}

type ScaleOptions struct {
	Type types.String `tfsdk:"type"`
	Log  types.Int64  `tfsdk:"log"`
}

type MappingOptions struct {
	Value   []ValueMappingOptions   `tfsdk:"value"`
	Range   []RangeMappingOptions   `tfsdk:"range"`
	Regex   []RegexMappingOptions   `tfsdk:"regex"`
	Special []SpecialMappingOptions `tfsdk:"special"`
}

type ValueMappingOptions struct {
	Value       types.String `tfsdk:"value"`
	DisplayText types.String `tfsdk:"display_text"`
	Color       types.String `tfsdk:"color"`
}

type RangeMappingOptions struct {
	From        types.Float64 `tfsdk:"from"`
	To          types.Float64 `tfsdk:"to"`
	DisplayText types.String  `tfsdk:"display_text"`
	Color       types.String  `tfsdk:"color"`
}

type RegexMappingOptions struct {
	Pattern     types.String `tfsdk:"pattern"`
	DisplayText types.String `tfsdk:"display_text"`
	Color       types.String `tfsdk:"color"`
}

type SpecialMappingOptions struct {
	Match       types.String `tfsdk:"match"`
	DisplayText types.String `tfsdk:"display_text"`
	Color       types.String `tfsdk:"color"`
}

type FieldOptions struct {
	Unit       types.String       `tfsdk:"unit"`
	Decimals   types.Int64        `tfsdk:"decimals"`
	Min        types.Float64      `tfsdk:"min"`
	Max        types.Float64      `tfsdk:"max"`
	NoValue    types.Float64      `tfsdk:"no_value"`
	Color      []ColorOptions     `tfsdk:"color"`
	Mappings   []MappingOptions   `tfsdk:"mappings"`
	Thresholds []ThresholdOptions `tfsdk:"thresholds"`
	// todo links
}

type ColorOptions struct {
	Mode       types.String `tfsdk:"mode"`
	FixedColor types.String `tfsdk:"fixed_color"`
	SeriesBy   types.String `tfsdk:"series_by"`
}

type ThresholdOptions struct {
	Mode  types.String    `tfsdk:"mode"`
	Steps []ThresholdStep `tfsdk:"step"`
}

type ThresholdStep struct {
	Color types.String  `tfsdk:"color"`
	Value types.Float64 `tfsdk:"value"`
}

type ReduceOptions struct {
	Values      types.Bool   `tfsdk:"values"`
	Fields      types.String `tfsdk:"fields"`
	Limit       types.Int64  `tfsdk:"limit"`
	Calculation types.String `tfsdk:"calculation"`
}

type TextSizeOptions struct {
	Title types.Int64 `tfsdk:"title"`
	Value types.Int64 `tfsdk:"value"`
}

type FieldOverrideOptions struct {
	ByName    []ByNameOverrideOptions    `tfsdk:"by_name"`
	ByRegex   []ByRegexOverrideOptions   `tfsdk:"by_regex"`
	ByType    []ByTypeOverrideOptions    `tfsdk:"by_type"`
	ByQueryID []ByQueryIDOverrideOptions `tfsdk:"by_query_id"`
}

type ByNameOverrideOptions struct {
	Name  types.String   `tfsdk:"name"`
	Field []FieldOptions `tfsdk:"field"`
}

type ByRegexOverrideOptions struct {
	Regex types.String   `tfsdk:"regex"`
	Field []FieldOptions `tfsdk:"field"`
}

type ByTypeOverrideOptions struct {
	Type  types.String   `tfsdk:"type"`
	Field []FieldOptions `tfsdk:"field"`
}

type ByQueryIDOverrideOptions struct {
	QueryID types.String   `tfsdk:"query_id"`
	Field   []FieldOptions `tfsdk:"field"`
}

type Query struct {
	Prometheus []PrometheusTarget `tfsdk:"prometheus"`
	CloudWatch []CloudWatchTarget `tfsdk:"cloudwatch"`
}

type PrometheusTarget struct {
	UID     types.String `tfsdk:"uid"`
	Expr    types.String `tfsdk:"expr"`
	Instant types.Bool   `tfsdk:"instant"`
	Format  types.String `tfsdk:"format"`
	// etc
	RefId        types.String `tfsdk:"ref_id"`
	MinInterval  types.String `tfsdk:"min_interval"`
	LegendFormat types.String `tfsdk:"legend_format"`
}

type CloudWatchTarget struct {
	UID        types.String          `tfsdk:"uid"`
	Namespace  types.String          `tfsdk:"namespace"`
	MetricName types.String          `tfsdk:"metric_name"`
	Statistic  types.String          `tfsdk:"statistic"`
	Dimensions []CloudWatchDimension `tfsdk:"dimension"`
	MatchExact types.Bool            `tfsdk:"match_exact"`
	Region     types.String          `tfsdk:"region"`
	// etc
	RefId  types.String `tfsdk:"ref_id"`
	Period types.String `tfsdk:"period"`
	Label  types.String `tfsdk:"label"`
}

type CloudWatchDimension struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

func axisBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: "Axis display options.",
		NestedObject: schema.NestedBlockObject{
			Blocks: map[string]schema.Block{
				"scale": schema.ListNestedBlock{
					Description: "Can be used to configure the scale of the y-axis.",
					MarkdownDescription: "Can be used to configure the scale of the y-axis. " +
						"Another way visualize series that differ by orders of magnitude is to use a logarithmic scales. " +
						"This is really useful for data usage or latency measurements. " +
						"The goal here is to avoid one series dominating and delegating all the others to the bottom of the graph.",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Required:            true,
								Description:         "The type of the scale. The choices are: linear, log.",
								MarkdownDescription: "The type of the scale. The choices are: `linear`, `log`.",
								Validators: []validator.String{
									stringvalidator.OneOf("linear", "log"),
								},
							},
							"log": schema.Int64Attribute{
								Optional:            true,
								Description:         "The power of the logarithmic scale. The choices are: 2, 10.",
								MarkdownDescription: "The power of the logarithmic scale. The choices are: `2`, `10`.",
								Validators: []validator.Int64{
									int64validator.OneOf(2, 10),
								},
							},
						},
					},
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
				},
			},
			Attributes: map[string]schema.Attribute{
				"label": schema.StringAttribute{
					Description: "The custom text label for the y-axis.",
					Optional:    true,
				},
				"placement": schema.StringAttribute{
					Optional:            true,
					Description:         "The placement of the y-axis. The choices are: auto, left, right, hidden.",
					MarkdownDescription: "The placement of the y-axis. The choices are: `auto`, `left`, `right`, `hidden`.",
					Validators: []validator.String{
						stringvalidator.OneOf("auto", "left", "right", "hidden"),
					},
				},
				"soft_min": schema.Int64Attribute{
					Optional:    true,
					Description: "The soft minimum of y-axis.",
					MarkdownDescription: "The soft minimum of y-axis. " +
						"By default, the Grafana workspace sets the range for the y-axis automatically based on the data." +
						"The `soft_min` setting can prevent blips from appearing as mountains when the data is mostly flat, " +
						"and hard min or max derived from standard min and max field options can prevent intermittent spikes " +
						"from flattening useful detail by clipping the spikes past a defined point.",
				},
				"soft_max": schema.Int64Attribute{
					Optional:    true,
					Description: "The soft maximum of y-axis.",
					MarkdownDescription: "The soft maximum of y-axis. " +
						"By default, the Grafana workspace sets the range for the y-axis automatically based on the data." +
						"The `soft_max` setting can prevent blips from appearing as mountains when the data is mostly flat, " +
						"and hard min or max derived from standard min and max field options can prevent intermittent spikes " +
						"from flattening useful detail by clipping the spikes past a defined point.",
				},
				/*"width": {
					Type:     types.Int64Type,
					Optional: true,
					Description: "The fixed width of the y-axis.",
					MarkdownDescription: "The fixed width of the y-axis. By default, the Grafana workspace dynamically calculates the axis width. " +
						"By setting the width of the axis, data whose axes types are different can share the same display proportions. " +
						"This makes it easier to compare more than one graph’s worth of data because the axes are not shifted or stretched within visual proximity of each other.",
				},*/
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func fieldBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: "The customization of field options.",
		NestedObject: schema.NestedBlockObject{

			Blocks: map[string]schema.Block{
				"color": schema.ListNestedBlock{
					Description:         "Defines how Grafana colors series or fields.",
					MarkdownDescription: "Defines how Grafana colors series or fields. There are multiple modes here that work differently, and their utility depends largely on the currently selected visualization.",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"mode": schema.StringAttribute{
								Optional:    true,
								Description: "The colorization mode.",
								MarkdownDescription: "The colorization mode. The most popular options:\n" +
									"1) `fixed` - specific color set by using the value of `fixed_color`.\n" +
									"2) `thresholds` - a color is derived from the matching threshold. This is useful for gauges, stat, and table visualizations.\n" +
									"3) `palette-classic` - a color is derived from the matching threshold using the classic color palette.",
								Validators: []validator.String{
									stringvalidator.OneOf(
										"fixed", "thresholds", "palette-classic",
										"continuous-GrYlRd", "continuous-RdYlGr", "continuous-BlYlRd", "continuous-YlRd", "continuous-BlPu", "continuous-YlBl",
										"continuous-blues", "continuous-reds", "continuous-greens", "continuous-purples",
									),
								},
							},
							"fixed_color": schema.StringAttribute{
								Optional:    true,
								Description: "The series to use to define the color. This is useful for graphs and pie charts, for example.",
							},
							"series_by": schema.StringAttribute{
								Optional:    true,
								Description: "The series to use to define the color. This is useful for graphs and pie charts, for example.",
							},
						},
						/* when fixed_color or series_by is present, mode must be present too

						Validators: []validator.Object{
							objectvalidator.
						},*/
					},
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
				},
				"thresholds": schema.ListNestedBlock{
					Description: "Thresholds set the color of the value text depending on conditions that you define.",
					NestedObject: schema.NestedBlockObject{
						Blocks: map[string]schema.Block{
							"step": schema.ListNestedBlock{
								Description: "The threshold steps.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"color": schema.StringAttribute{
											Required:    true,
											Description: "The color for the matching values.",
										},
										"value": schema.Float64Attribute{
											Optional:    true,
											Description: "The value to match. Either percentage or absolute. Depends on the mode.",
											MarkdownDescription: "The value to match. Either percentage or absolute. Depends on the mode. " +
												"The step without `value` indicates the base color. It is generally the good color.",
										},
									},
								},
								Validators: []validator.List{
									listvalidator.SizeAtMost(20),
								},
							},
						},
						Attributes: map[string]schema.Attribute{
							"mode": schema.StringAttribute{
								Optional:    true,
								Description: "The threshold mode. The choices are: absolute, percentage.",
								MarkdownDescription: "The threshold mode. The choices are:\n" +
									"1) `absolute` - defined based on a number; for example, 80 on a scale of 1 to 150. \n" +
									"2) `percentage` - defined relative to minimum or maximum; for example, 80 percent.",
								Validators: []validator.String{
									stringvalidator.OneOf("absolute", "percentage"),
								},
							},
						},
					},
					Validators: []validator.List{
						listvalidator.SizeAtMost(1),
					},
				},
				"mappings": mappingsBlock(),
			},
			Attributes: map[string]schema.Attribute{
				"unit": schema.StringAttribute{
					Optional:    true,
					Description: "The unit the field should use.",
				},
				"decimals": schema.Int64Attribute{
					Optional:            true,
					Description:         "The number of decimals to include when rendering a value. Must be between 0 and 20 (inclusive).",
					MarkdownDescription: "The number of decimals to include when rendering a value. Must be between `0` and `20` (inclusive).",
					Validators: []validator.Int64{
						int64validator.Between(0, 20),
					},
				},
				"min": schema.Float64Attribute{
					Optional:    true,
					Description: "The minimum value used in percentage threshold calculations.",
				},
				"max": schema.Float64Attribute{
					Optional:    true,
					Description: "The maximum value used in percentage threshold calculations.",
				},
				"no_value": schema.Float64Attribute{
					Optional:    true,
					Description: "The value to display if the field value is empty or null.",
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func reduceOptionsBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: "Reduction or calculation options for a value.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"values": schema.BoolAttribute{
					Optional:    true,
					Description: "Whether to calculate a single value per column or series or show each row.",
				},
				"fields": schema.StringAttribute{
					Optional:    true,
					Description: "The fields that should be included in the panel.", // todo schema validation `values = true`
				},
				"limit": schema.Int64Attribute{
					Optional:    true,
					Description: "The max number of rows to display.", // todo schema validation `values = true`
				},
				"calculation": schema.StringAttribute{ // todo schema validation `values = false`
					Optional:            true,
					Description:         "A reducer function or calculation. The choices are: " + CalculationTypesString() + ".",
					MarkdownDescription: "A reducer function or calculation. The choices are: " + CalculationTypesMarkdown() + ".",
					Validators: []validator.String{
						stringvalidator.OneOf(CalculationTypes()...),
					},
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func textSizeBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: "The size of the text elements on the panel.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"title": schema.Int64Attribute{
					Optional:            true,
					Description:         "The size of the title. Must be between 1 and 100 (inclusive).",
					MarkdownDescription: "The size of the title. Must be between `1` and `100` (inclusive).",
					Validators: []validator.Int64{
						int64validator.Between(1, 100),
					},
				},
				"value": schema.Int64Attribute{
					Optional:            true,
					Description:         "The size of the value. Must be between 1 and 100 (inclusive).",
					MarkdownDescription: "The size of the value. Must be between `1` and `100` (inclusive).",
					Validators: []validator.Int64{
						int64validator.Between(1, 100),
					},
				},
			},
			Validators: []validator.Object{
				objectvalidator.AtLeastOneOf(
					path.MatchRelative().AtParent().AtName("title"),
					path.MatchRelative().AtParent().AtName("value"),
				),
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func queryBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: "The queries to collect values from data sources.",
		NestedObject: schema.NestedBlockObject{
			Blocks: map[string]schema.Block{
				"prometheus": schema.ListNestedBlock{
					Description: "The Prometheus query.",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"uid": schema.StringAttribute{
								Description: "The UID of a Prometheus DataSource to use in this query.",
								Required:    true,
							},
							"expr": schema.StringAttribute{
								Required:    true,
								Description: "The query expression.",
							},
							"instant": schema.BoolAttribute{
								Optional:    true,
								Description: "Whether to return the latest value from the time series or not.",
							},
							"ref_id": schema.StringAttribute{
								Optional:    true,
								Description: "The ID of the query. The ID can be used to reference queries in math expressions.",
							},
							"format": schema.StringAttribute{
								Optional:            true,
								Description:         "The query format. The choices are: time_series, table, heatmap.",
								MarkdownDescription: "The query format. The choices are: `time_series`, `table`, `heatmap`.",
								Validators: []validator.String{
									stringvalidator.OneOf("time_series", "table", "heatmap"),
								},
							},
							"min_interval": schema.StringAttribute{
								Optional:    true,
								Description: "The lower bounds on the interval between data points.",
							},
							"legend_format": schema.StringAttribute{
								Optional:    true,
								Description: "The legend name.",
							},
						},
					},
					Validators: []validator.List{
						listvalidator.SizeAtMost(26),
					},
				},
				"cloudwatch": schema.ListNestedBlock{
					Description: "The CloudWatch query.",
					NestedObject: schema.NestedBlockObject{
						Blocks: map[string]schema.Block{
							"dimension": schema.ListNestedBlock{
								Description: "The dimension to filter the metric with.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"name": schema.StringAttribute{
											Required:    true,
											Description: "The name of the dimension.",
										},
										"value": schema.StringAttribute{
											Required:    true,
											Description: "The value of the dimension.",
										},
									},
								},
								Validators: []validator.List{
									listvalidator.SizeAtMost(5),
								},
							},
						},
						Attributes: map[string]schema.Attribute{
							"uid": schema.StringAttribute{
								Description: "The UID of a CloudWatch DataSource to use in this query.",
								Required:    true,
							},
							"namespace": schema.StringAttribute{
								Required:    true,
								Description: "The namespace to query the metrics from.",
							},
							"metric_name": schema.StringAttribute{
								Required:            true,
								Description:         "The name of the metric to query.",
								MarkdownDescription: "The name of the metric to query. Example: `CPUUtilization`",
							},
							"statistic": schema.StringAttribute{
								Required:    true,
								Description: "The calculation to apply to the time series.",
							},
							"match_exact": schema.BoolAttribute{
								Optional:            true,
								Description:         "If enabled you also need to specify all the dimensions of the metric you’re querying.",
								MarkdownDescription: "If enabled you also need to specify **all** the dimensions of the metric you’re querying.",
							},
							"region": schema.StringAttribute{
								Optional:    true,
								Description: "The AWS region to query the metrics from.",
							},
							"ref_id": schema.StringAttribute{
								Optional:    true,
								Description: "The ID of the query. The ID can be used to reference queries in math expressions.",
							},
							"period": schema.StringAttribute{
								Optional:    true,
								Description: "The minimum interval between points in seconds.",
							},
							"label": schema.StringAttribute{
								Optional:    true,
								Description: "The legend name.",
							},
						},
					},
					Validators: []validator.List{
						listvalidator.SizeAtMost(26),
					},
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(3),
		},
	}
}

func mappingsBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: "The set of rules that translate a field value or range of values into explicit text.",
		NestedObject: schema.NestedBlockObject{
			Blocks: map[string]schema.Block{
				"value": schema.ListNestedBlock{
					Description: "Match a specific text value.",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"value": schema.StringAttribute{
								Required:    true,
								Description: "The exact value to match.",
							},
							"display_text": schema.StringAttribute{
								Optional:    true,
								Description: "Text to display if the condition is met. This field accepts Grafana variables.",
							},
							"color": schema.StringAttribute{
								Optional:    true,
								Description: "The color to use if the condition is met.",
							},
						},
					},
					Validators: []validator.List{
						listvalidator.SizeAtMost(10),
					},
				},
				"range": schema.ListNestedBlock{
					Description: "Match a numerical range of values.",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"from": schema.Float64Attribute{
								Required:    true,
								Description: "The start of the range.",
							},
							"to": schema.Float64Attribute{
								Required:    true,
								Description: "The end of the range.",
							},
							"display_text": schema.StringAttribute{
								Optional:    true,
								Description: "Text to display if the condition is met. This field accepts Grafana variables.",
							},
							"color": schema.StringAttribute{
								Optional:    true,
								Description: "The color to use if the condition is met.",
							},
						},
					},
					Validators: []validator.List{
						listvalidator.SizeAtMost(10),
					},
				},
				"regex": schema.ListNestedBlock{
					Description: "Match a regular expression with replacement.",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"pattern": schema.StringAttribute{
								Required:    true,
								Description: "The regular expression to match.",
							},
							"display_text": schema.StringAttribute{
								Optional:    true,
								Description: "Text to display if the condition is met. This field accepts Grafana variables.",
							},
							"color": schema.StringAttribute{
								Optional:    true,
								Description: "The color to use if the condition is met.",
							},
						},
					},
					Validators: []validator.List{
						listvalidator.SizeAtMost(10),
					},
				},
				"special": schema.ListNestedBlock{
					Description: "Match on null, NaN, boolean and empty values.",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"match": schema.StringAttribute{
								Optional:            true,
								Description:         "The category to match. The choices are: null, nan, null+nan, true, false, empty.",
								MarkdownDescription: "The category to match. The choices are: `null`, `nan`, `null+nan`, `true`, `false`, `empty`.",
								Validators: []validator.String{
									stringvalidator.OneOf("null", "nan", "null+nan", "true", "false", "empty"),
								},
							},
							"display_text": schema.StringAttribute{
								Optional:    true,
								Description: "Text to display if the condition is met. This field accepts Grafana variables.",
							},
							"color": schema.StringAttribute{
								Optional:    true,
								Description: "The color to use if the condition is met.",
							},
						},
					},
					Validators: []validator.List{
						listvalidator.SizeAtMost(10),
					},
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func fieldOverrideBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: "The set of rules that override attributes of a field.",
		NestedObject: schema.NestedBlockObject{
			Blocks: map[string]schema.Block{
				"by_name": schema.ListNestedBlock{
					Description: "Override properties for a field with a specific name.",
					NestedObject: schema.NestedBlockObject{
						Blocks: map[string]schema.Block{
							"field": fieldBlock(),
						},
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Required:    true,
								Description: "The name of the field to override attributes for.",
							},
						},
					},
					Validators: []validator.List{
						listvalidator.SizeAtMost(100),
					},
				},
				"by_regex": schema.ListNestedBlock{
					Description: "Override properties for a field with a matching name.",
					NestedObject: schema.NestedBlockObject{
						Blocks: map[string]schema.Block{
							"field": fieldBlock(),
						},
						Attributes: map[string]schema.Attribute{
							"regex": schema.StringAttribute{
								Required:    true,
								Description: "The regex the field's name should match.",
							},
						},
					},
					Validators: []validator.List{
						listvalidator.SizeAtMost(100),
					},
				},
				"by_type": schema.ListNestedBlock{
					Description: "Override properties for a field with a specific type.",
					NestedObject: schema.NestedBlockObject{
						Blocks: map[string]schema.Block{
							"field": fieldBlock(),
						},
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Required:    true,
								Description: "The type of the field to override attributes for.",
							},
						},
					},
					Validators: []validator.List{
						listvalidator.SizeAtMost(100),
					},
				},
				"by_query_id": schema.ListNestedBlock{
					Description: "Override properties for a field returned by a specific query.",
					NestedObject: schema.NestedBlockObject{
						Blocks: map[string]schema.Block{
							"field": fieldBlock(),
						},
						Attributes: map[string]schema.Attribute{
							"query_id": schema.StringAttribute{
								Required:    true,
								Description: "The name of the field to override attributes for.",
							},
						},
					},
					Validators: []validator.List{
						listvalidator.SizeAtMost(100),
					},
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(10),
		},
	}
}

// attributes
func idAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Computed: true,
	}
}

func jsonAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Computed:    true,
		Description: "The Grafana-API-compatible JSON of this panel.",
	}
}

func titleAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Required:    true,
		Description: "The title of this panel.",
	}
}

func descriptionAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Optional:    true,
		Description: "The description of this panel.",
	}
}

// creators

func createTargets(queries []Query) []grafana.Target {
	targets := make([]grafana.Target, 0)

	for _, group := range queries {
		for _, target := range group.Prometheus {
			t := grafana.Target{
				Datasource: grafana.Datasource{
					UID:  target.UID.ValueString(),
					Type: "prometheus",
				},
				RefID:        target.RefId.ValueString(),
				Expr:         target.Expr.ValueString(),
				Interval:     target.MinInterval.ValueString(),
				LegendFormat: target.LegendFormat.ValueString(),
				Instant:      target.Instant.ValueBool(),
				Format:       target.Format.ValueString(),
			}

			targets = append(targets, t)
		}

		for _, target := range group.CloudWatch {
			dimensions := make(map[string]string)

			for _, dim := range target.Dimensions {
				dimensions[dim.Name.ValueString()] = dim.Value.ValueString()
			}

			t := grafana.Target{
				Datasource: grafana.Datasource{
					UID:  target.UID.ValueString(),
					Type: "cloudwatch",
				},
				RefID:      target.RefId.ValueString(),
				Namespace:  target.Namespace.ValueString(),
				MetricName: target.MetricName.ValueString(),
				Statistics: []string{target.Statistic.ValueString()},
				Dimensions: dimensions,
				Period:     target.Period.ValueString(),
				Region:     target.Region.ValueString(),
				Label:      target.Label.ValueString(),
			}

			targets = append(targets, t)
		}
	}

	return targets
}

type ValueMappingResult struct {
	Color string `json:"color,omitempty"`
	Text  string `json:"text,omitempty"`
	Index int    `json:"index"`
}

func createFieldConfig(defaults FieldDefaults, fieldOptions []FieldOptions) grafana.FieldConfigDefaults {
	thresholdStep := make([]grafana.ThresholdStep, len(defaults.Thresholds.Steps))

	for i, step := range defaults.Thresholds.Steps {
		thresholdStep[i] = grafana.ThresholdStep{
			Color: step.Color,
			Value: step.Value,
		}
	}

	fieldConfig := grafana.FieldConfigDefaults{
		Unit:     defaults.Unit,
		Decimals: defaults.Decimals,
		Min:      defaults.Min,
		Max:      defaults.Max,
		Color: grafana.FieldConfigColor{
			Mode:       defaults.Color.Mode,
			FixedColor: defaults.Color.FixedColor,
			SeriesBy:   defaults.Color.SeriesBy,
		},
		Thresholds: grafana.Thresholds{
			Mode:  defaults.Thresholds.Mode,
			Steps: thresholdStep,
		},
	}

	for _, field := range fieldOptions {
		if !field.Unit.IsNull() {
			fieldConfig.Unit = field.Unit.ValueString()
		}

		fieldConfig.Decimals = field.Decimals.ValueInt64Pointer()
		fieldConfig.Min = field.Min.ValueFloat64Pointer()
		fieldConfig.Max = field.Max.ValueFloat64Pointer()
		fieldConfig.NoValue = field.NoValue.ValueFloat64Pointer()

		for _, color := range field.Color {
			if !color.Mode.IsNull() {
				fieldConfig.Color.Mode = color.Mode.ValueString()
			}

			if !color.FixedColor.IsNull() {
				fieldConfig.Color.FixedColor = color.FixedColor.ValueString()
			}

			if !color.SeriesBy.IsNull() {
				fieldConfig.Color.SeriesBy = color.SeriesBy.ValueString()
			}
		}

		mappings := createMappings(field.Mappings)

		if len(mappings) > 0 {
			fieldConfig.Mappings = mappings
		}

		updateThresholds(&fieldConfig.Thresholds, field.Thresholds)
	}

	return fieldConfig
}

func createMappings(mappingOptions []MappingOptions) []grafana.FieldMapping {
	mappings := make([]grafana.FieldMapping, 0)

	for _, mapping := range mappingOptions {
		idx := 0
		valuesMap := make(map[string]interface{})

		for _, value := range mapping.Value {
			v := ValueMappingResult{
				Color: value.Color.ValueString(),
				Text:  value.DisplayText.ValueString(),
				Index: idx,
			}

			valuesMap[value.Value.ValueString()] = v
			idx += 1
		}

		if len(valuesMap) > 0 {
			mapping := grafana.FieldMapping{
				Type:    "value",
				Options: valuesMap,
			}

			mappings = append(mappings, mapping)
		}

		for _, range_ := range mapping.Range {
			mapping := grafana.FieldMapping{
				Type: "range",
				Options: map[string]interface{}{
					"from": range_.From.ValueFloat64(),
					"to":   range_.From.ValueFloat64(),
					"result": ValueMappingResult{
						Color: range_.Color.ValueString(),
						Text:  range_.DisplayText.ValueString(),
						Index: idx,
					},
				},
			}
			idx += 1

			mappings = append(mappings, mapping)
		}

		for _, regex := range mapping.Regex {
			mapping := grafana.FieldMapping{
				Type: "regex",
				Options: map[string]interface{}{
					"pattern": regex.Pattern.ValueString(),
					"result": ValueMappingResult{
						Color: regex.Color.ValueString(),
						Text:  regex.DisplayText.ValueString(),
						Index: idx,
					},
				},
			}
			idx += 1

			mappings = append(mappings, mapping)
		}

		for _, special := range mapping.Special {
			mapping := grafana.FieldMapping{
				Type: "special",
				Options: map[string]interface{}{
					"match": special.Match.ValueString(),
					"result": ValueMappingResult{
						Color: special.Color.ValueString(),
						Text:  special.DisplayText.ValueString(),
						Index: idx,
					},
				},
			}
			idx += 1

			mappings = append(mappings, mapping)
		}
	}

	return mappings
}

func createOverrides(overrides []FieldOverrideOptions) []grafana.FieldOverride {
	fieldOverrides := make([]grafana.FieldOverride, 0)

	for _, override := range overrides {
		for _, byName := range override.ByName {
			fieldOverride := grafana.FieldOverride{
				Matcher: grafana.FieldOverrideMatcher{
					Id:      "byName",
					Options: byName.Name.ValueString(),
				},
				Properties: createOverrideProperties(byName.Field),
			}
			fieldOverrides = append(fieldOverrides, fieldOverride)
		}

		for _, byRegex := range override.ByRegex {
			fieldOverride := grafana.FieldOverride{
				Matcher: grafana.FieldOverrideMatcher{
					Id:      "byRegexp",
					Options: byRegex.Regex.ValueString(),
				},
				Properties: createOverrideProperties(byRegex.Field),
			}
			fieldOverrides = append(fieldOverrides, fieldOverride)
		}

		for _, byType := range override.ByType {
			fieldOverride := grafana.FieldOverride{
				Matcher: grafana.FieldOverrideMatcher{
					Id:      "byType",
					Options: byType.Type.ValueString(),
				},
				Properties: createOverrideProperties(byType.Field),
			}
			fieldOverrides = append(fieldOverrides, fieldOverride)
		}

		for _, byQueryID := range override.ByQueryID {
			fieldOverride := grafana.FieldOverride{
				Matcher: grafana.FieldOverrideMatcher{
					Id:      "byFrameRefID",
					Options: byQueryID.QueryID.ValueString(),
				},
				Properties: createOverrideProperties(byQueryID.Field),
			}
			fieldOverrides = append(fieldOverrides, fieldOverride)
		}
	}

	return fieldOverrides
}

func createOverrideProperties(fieldOptions []FieldOptions) []grafana.FieldOverrideProperty {
	properties := make([]grafana.FieldOverrideProperty, 0)

	for _, field := range fieldOptions {
		if !field.Unit.IsNull() {
			properties = append(properties, grafana.FieldOverrideProperty{
				Id:    "unit",
				Value: field.Unit.ValueString(),
			})
		}

		if !field.Decimals.IsNull() {
			properties = append(properties, grafana.FieldOverrideProperty{
				Id:    "decimals",
				Value: field.Decimals.ValueInt64(),
			})
		}

		if !field.Min.IsNull() {
			properties = append(properties, grafana.FieldOverrideProperty{
				Id:    "min",
				Value: field.Min.ValueFloat64(),
			})
		}

		if !field.Max.IsNull() {
			properties = append(properties, grafana.FieldOverrideProperty{
				Id:    "max",
				Value: field.Max.ValueFloat64(),
			})
		}

		if !field.NoValue.IsNull() {
			properties = append(properties, grafana.FieldOverrideProperty{
				Id:    "noValue",
				Value: field.Decimals.ValueInt64(),
			})
		}

		for _, color := range field.Color {
			fieldColor := grafana.FieldConfigColor{}
			if !color.Mode.IsNull() {
				fieldColor.Mode = color.Mode.ValueString()
			}

			if !color.FixedColor.IsNull() {
				fieldColor.FixedColor = color.FixedColor.ValueString()
			}

			if !color.SeriesBy.IsNull() {
				fieldColor.SeriesBy = color.SeriesBy.ValueString()
			}

			properties = append(properties, grafana.FieldOverrideProperty{
				Id:    "color",
				Value: fieldColor,
			})
		}

		mappings := createMappings(field.Mappings)

		if len(mappings) > 0 {
			properties = append(properties, grafana.FieldOverrideProperty{
				Id:    "mappings",
				Value: mappings,
			})
		}

		thresholds := grafana.Thresholds{}
		updateThresholds(&thresholds, field.Thresholds)

		if len(field.Thresholds) > 0 {
			properties = append(properties, grafana.FieldOverrideProperty{
				Id:    "thresholds",
				Value: thresholds,
			})
		}
	}

	return properties
}

// updaters
func updateThresholds(thresholds *grafana.Thresholds, thresholdOptions []ThresholdOptions) {
	for _, threshold := range thresholdOptions {
		steps := make([]grafana.ThresholdStep, len(threshold.Steps))

		if !threshold.Mode.IsNull() {
			thresholds.Mode = threshold.Mode.ValueString()
		}

		for i, step := range threshold.Steps {
			s := grafana.ThresholdStep{
				Color: step.Color.ValueString(),
			}

			if !step.Value.IsNull() {
				s.Value = step.Value.ValueFloat64Pointer()
			}

			steps[i] = s
		}

		if len(steps) > 0 {
			thresholds.Steps = steps
		}
	}
}

func updateTextSize(options *grafana.TextSize, opts []TextSizeOptions) {
	for _, textSize := range opts {
		options.TitleSize = textSize.Title.ValueInt64Pointer()
		options.ValueSize = textSize.Value.ValueInt64Pointer()
	}
}

func updateReduceOptions(options *grafana.ReduceOptions, opts []ReduceOptions) {
	for _, reducer := range opts {
		if !reducer.Values.IsNull() {
			options.Values = reducer.Values.ValueBool()
		}

		if !reducer.Fields.IsNull() {
			options.Fields = reducer.Fields.ValueString()
		}

		if !reducer.Limit.IsNull() {
			options.Limit = reducer.Limit.ValueInt64Pointer()
		}

		if !reducer.Calculation.IsNull() {
			options.Calcs = []string{reducer.Calculation.ValueString()}
		}
	}
}

// etc
func hashcode(s []byte) int {
	v := int(crc32.ChecksumIEEE(s))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}
