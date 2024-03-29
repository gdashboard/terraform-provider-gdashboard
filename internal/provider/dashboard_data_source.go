package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gdashboard/terraform-provider-gdashboard/internal/provider/grafana"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"math"
	"regexp"
	"strconv"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &DashboardDataSource{}

func NewDashboardDataSource() datasource.DataSource {
	return &DashboardDataSource{}
}

// DashboardDataSource defines the data source implementation.
type DashboardDataSource struct {
	CompactJson bool
	Defaults    DashboardDefaults
}

type DashboardDefaults struct {
	Editable     bool
	Style        string
	GraphTooltip string
	Time         Time
}

type Time struct {
	From string
	To   string
}

// DashboardDataSourceModel describes the data source data model.
type DashboardDataSourceModel struct {
	Id           types.String           `tfsdk:"id"`
	Json         types.String           `tfsdk:"json"`
	CompactJson  types.Bool             `tfsdk:"compact_json"`
	Title        types.String           `tfsdk:"title"`
	Description  types.String           `tfsdk:"description"`
	Version      types.Int64            `tfsdk:"version"`
	UID          types.String           `tfsdk:"uid"`
	Editable     types.Bool             `tfsdk:"editable"`
	Style        types.String           `tfsdk:"style"`
	GraphTooltip types.String           `tfsdk:"graph_tooltip"`
	Tags         []types.String         `tfsdk:"tags"`
	TimeOptions  []DashboardTimeOptions `tfsdk:"time"`
	Layout       Layout                 `tfsdk:"layout"`
	Variables    []Variable             `tfsdk:"variables"`
	Annotations  []Annotation           `tfsdk:"annotations"`
	Links        []Link                 `tfsdk:"links"`
}

type DashboardTimeOptions struct {
	Timezone              types.String                 `tfsdk:"timezone"`
	WeekStart             types.String                 `tfsdk:"week_start"`
	RefreshLiveDashboards types.Bool                   `tfsdk:"refresh_live_dashboards"`
	Range                 []TimeModel                  `tfsdk:"default_range"`
	TimePicker            []DashboardTimePickerOptions `tfsdk:"picker"`
}

type DashboardTimePickerOptions struct {
	Hide             types.Bool     `tfsdk:"hide"`
	NowDelay         types.String   `tfsdk:"now_delay"`
	RefreshIntervals []types.String `tfsdk:"refresh_intervals"`
	TimeOptions      []types.String `tfsdk:"time_options"`
}

type Layout struct {
	Sections []Section `tfsdk:"section"`
}

// Section has two modes: auto layout and explicit rows
type Section struct {
	Title     types.String `tfsdk:"title"`
	Collapsed types.Bool   `tfsdk:"collapsed"`
	Panels    []Panel      `tfsdk:"panel"`
	Rows      []SectionRow `tfsdk:"row"`
}

type SectionRow struct {
	Panels []Panel `tfsdk:"panel"`
}

type Panel struct {
	Size   Size         `tfsdk:"size"`
	Source types.String `tfsdk:"source"`
}

type Size struct {
	Width  types.Int64 `tfsdk:"width"`
	Height types.Int64 `tfsdk:"height"`
}

type Variable struct {
	Custom     []VariableCustom     `tfsdk:"custom"`
	Constant   []VariableConstant   `tfsdk:"const"`
	TextBox    []VariableTextBox    `tfsdk:"textbox"`
	AdHoc      []VariableAdHoc      `tfsdk:"adhoc"`
	DataSource []VariableDataSource `tfsdk:"datasource"`
	Query      []VariableQuery      `tfsdk:"query"`
	Interval   []VariableInterval   `tfsdk:"interval"`
}

type VariableCustom struct {
	Name        types.String           `tfsdk:"name"`
	Label       types.String           `tfsdk:"label"`
	Description types.String           `tfsdk:"description"`
	Hide        types.String           `tfsdk:"hide"`
	Options     []VariableCustomOption `tfsdk:"option"`
	Multi       types.Bool             `tfsdk:"multi"`
	IncludeAll  []VariableIncludeAll   `tfsdk:"include_all"`
}

type VariableCustomOption struct {
	Selected types.Bool   `tfsdk:"selected"`
	Text     types.String `tfsdk:"text"`
	Value    types.String `tfsdk:"value"`
}

type VariableIncludeAll struct {
	Enabled     types.Bool   `tfsdk:"enabled"`
	CustomValue types.String `tfsdk:"custom_value"`
}

type VariableConstant struct {
	Name        types.String `tfsdk:"name"`
	Label       types.String `tfsdk:"label"`
	Description types.String `tfsdk:"description"`
	Value       types.String `tfsdk:"value"`
}

type VariableTextBox struct {
	Name         types.String `tfsdk:"name"`
	Label        types.String `tfsdk:"label"`
	Description  types.String `tfsdk:"description"`
	DefaultValue types.String `tfsdk:"default_value"`
	Hide         types.String `tfsdk:"hide"`
}

type VariableAdHoc struct {
	Name        types.String            `tfsdk:"name"`
	Label       types.String            `tfsdk:"label"`
	Description types.String            `tfsdk:"description"`
	Hide        types.String            `tfsdk:"hide"`
	DataSource  VariableAdHocDataSource `tfsdk:"datasource"`
	Filters     []VariableAdHocFilter   `tfsdk:"filter"`
}

type VariableAdHocDataSource struct {
	Type types.String `tfsdk:"type"`
	UID  types.String `tfsdk:"uid"`
}

type VariableAdHocFilter struct {
	Key      types.String `tfsdk:"key"`
	Operator types.String `tfsdk:"operator"`
	Value    types.String `tfsdk:"value"`
}

type VariableDataSource struct {
	Name        types.String               `tfsdk:"name"`
	Label       types.String               `tfsdk:"label"`
	Description types.String               `tfsdk:"description"`
	Hide        types.String               `tfsdk:"hide"`
	Multi       types.Bool                 `tfsdk:"multi"`
	IncludeAll  []VariableIncludeAll       `tfsdk:"include_all"`
	DataSource  VariableDataSourceSelector `tfsdk:"source"`
}

type VariableDataSourceSelector struct {
	Type   types.String `tfsdk:"type"`
	Filter types.String `tfsdk:"filter"`
}

type VariableQuery struct {
	Name        types.String          `tfsdk:"name"`
	Label       types.String          `tfsdk:"label"`
	Description types.String          `tfsdk:"description"`
	Hide        types.String          `tfsdk:"hide"`
	Refresh     types.String          `tfsdk:"refresh"`
	Multi       types.Bool            `tfsdk:"multi"`
	IncludeAll  []VariableIncludeAll  `tfsdk:"include_all"`
	Sort        []VariableQuerySort   `tfsdk:"sort"`
	Regex       types.String          `tfsdk:"regex"`
	Target      []VariableQueryTarget `tfsdk:"target"`
}

type VariableQuerySort struct {
	Type  types.String `tfsdk:"type"`
	Order types.String `tfsdk:"order"`
}

type VariableQueryTarget struct {
	Prometheus []VariableQueryTargetPrometheus `tfsdk:"prometheus"`
}

type VariableQueryTargetPrometheus struct {
	UID  types.String `tfsdk:"uid"`
	Expr types.String `tfsdk:"expr"`
}

type VariableInterval struct {
	Name        types.String           `tfsdk:"name"`
	Label       types.String           `tfsdk:"label"`
	Description types.String           `tfsdk:"description"`
	Hide        types.String           `tfsdk:"hide"`
	Intervals   []types.String         `tfsdk:"intervals"`
	Auto        []VariableIntervalAuto `tfsdk:"auto"`
}

type VariableIntervalAuto struct {
	Enabled     types.Bool   `tfsdk:"enabled"`
	StepCount   types.Int64  `tfsdk:"step_count"`
	MinInterval types.String `tfsdk:"min_interval"`
}

type Annotation struct {
	Grafana    []AnnotationGrafana    `tfsdk:"grafana"`
	Prometheus []AnnotationPrometheus `tfsdk:"prometheus"`
}

type AnnotationGrafana struct {
	Name        types.String                        `tfsdk:"name"`
	Enabled     types.Bool                          `tfsdk:"enabled"`
	Hidden      types.Bool                          `tfsdk:"hidden"`
	Color       types.String                        `tfsdk:"color"`
	ByDashboard []AnnotationGrafanaQueryByDashboard `tfsdk:"by_dashboard"`
	ByTags      []AnnotationGrafanaQueryByTags      `tfsdk:"by_tags"`
}

type AnnotationGrafanaQueryByDashboard struct {
	Limit types.Int64 `tfsdk:"limit"`
}

type AnnotationGrafanaQueryByTags struct {
	Limit    types.Int64    `tfsdk:"limit"`
	MatchAny types.Bool     `tfsdk:"match_any"`
	Tags     []types.String `tfsdk:"tags"`
}

type AnnotationPrometheus struct {
	Name    types.String              `tfsdk:"name"`
	Enabled types.Bool                `tfsdk:"enabled"`
	Hidden  types.Bool                `tfsdk:"hidden"`
	Color   types.String              `tfsdk:"color"`
	Query   AnnotationPrometheusQuery `tfsdk:"query"`
}

type AnnotationPrometheusQuery struct {
	UID                 types.String `tfsdk:"datasource_uid"`
	Expr                types.String `tfsdk:"expr"`
	Step                types.String `tfsdk:"min_step"`
	Title               types.String `tfsdk:"title_format"`
	Text                types.String `tfsdk:"text_format"`
	UseValueAsTimestamp types.Bool   `tfsdk:"use_value_as_timestamp"`
	TagKeys             types.String `tfsdk:"tag_keys"`
}

type Link struct {
	Dashboards []LinkDashboards `tfsdk:"dashboards"`
	External   []LinkExternal   `tfsdk:"external"`
}

type LinkDashboards struct {
	Title                    types.String   `tfsdk:"title"`
	Tags                     []types.String `tfsdk:"tags"`
	AsDropdown               types.Bool     `tfsdk:"as_dropdown"`
	IncludeTimeRange         types.Bool     `tfsdk:"include_time_range"`
	IncludeTemplateVariables types.Bool     `tfsdk:"include_template_variables"`
	NewTab                   types.Bool     `tfsdk:"new_tab"`
}

type LinkExternal struct {
	Title                    types.String `tfsdk:"title"`
	Url                      types.String `tfsdk:"url"`
	Tooltip                  types.String `tfsdk:"tooltip"`
	Icon                     types.String `tfsdk:"icon"`
	IncludeTimeRange         types.Bool   `tfsdk:"include_time_range"`
	IncludeTemplateVariables types.Bool   `tfsdk:"include_template_variables"`
	NewTab                   types.Bool   `tfsdk:"new_tab"`
}

func (d *DashboardDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dashboard"
}

func dashboardTimeBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: "The default query time range.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"from": schema.StringAttribute{
					Required:    true,
					Description: "The start of the range.",
				},
				"to": schema.StringAttribute{
					Required:    true,
					Description: "The end of the range.",
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func dashboardEditableAttribute() schema.Attribute {
	return schema.BoolAttribute{
		Optional:    true,
		Description: "Whether to make the dashboard editable or not.",
	}
}

func dashboardStyleAttribute() schema.Attribute {
	return schema.StringAttribute{
		Optional:            true,
		Description:         "The dashboard style. The choices are: dark, light.",
		MarkdownDescription: "The dashboard style. The choices are: `dark`, `light`.",
		Validators: []validator.String{
			stringvalidator.OneOf("dark", "light"),
		},
	}
}

func dashboardGraphTooltipAttribute() schema.Attribute {
	return schema.StringAttribute{
		Optional:            true,
		Description:         "Controls tooltip and hover highlight behavior across different panels: default, shared-crosshair, shared-tooltip.",
		MarkdownDescription: "Controls tooltip and hover highlight behavior across different panels: `default`, `shared-crosshair`, `shared-tooltip`.",
		Validators: []validator.String{
			stringvalidator.OneOf("default", "shared-crosshair", "shared-tooltip"),
		},
	}
}

func variableIncludeAllBlock() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description:         "An option to include all variables. If 'custom_value' is blank, then the Grafana concatenates (adds together) all the values in the query.",
		MarkdownDescription: "An option to include all variables. If `custom_value` is blank, then the Grafana concatenates (adds together) all the values in the query.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"enabled": schema.BoolAttribute{
					Required:    true,
					Description: "Whether to enable the option to include all variables or not.",
				},
				"custom_value": schema.StringAttribute{
					Optional:            true,
					Description:         "The value to use when 'include_all' is enabled",
					MarkdownDescription: "The value to use when `include_all` is enabled. Example: `*`, `all`, etc.",
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func variableNameAttribute() schema.Attribute {
	return schema.StringAttribute{
		Description: "The name of the variable.",
		Required:    true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
			stringvalidator.LengthAtMost(50),
			stringvalidator.RegexMatches(regexp.MustCompile("^\\w+$"), "Only word and digit characters are allowed in variable names"),
		},
	}
}

func variableHideAttribute() schema.Attribute {
	return schema.StringAttribute{
		Optional:            true,
		Description:         "Which variable information to hide from the dashboard. The choices are: label, variable.",
		MarkdownDescription: "Which variable information to hide from the dashboard. The choices are: `label`, `variable`.",
		Validators: []validator.String{
			stringvalidator.OneOf("label", "variable"),
		},
	}
}

func annotationCommonAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Required:    true,
			Description: "The name of the annotation.",
		},
		"enabled": schema.BoolAttribute{
			Optional:    true,
			Description: "When enabled the annotation query is issued every dashboard refresh.",
		},
		"hidden": schema.BoolAttribute{
			Optional:    true,
			Description: "Whether the annotation can be toggled on or off at the top of the dashboard. With this option checked this toggle will be hidden.",
		},
		"color": schema.StringAttribute{
			Optional:    true,
			Description: "The color to use for the annotation event markers.",
		},
	}
}

func panelBlock() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "The definition of the panel within the row.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"size": schema.SingleNestedAttribute{
					Description: "The size of the panel.",
					Required:    true,
					Attributes: map[string]schema.Attribute{
						"height": schema.Int64Attribute{
							Required:    true,
							Description: "The height of the panel.",
						},
						"width": schema.Int64Attribute{
							Required:    true,
							Description: "The width of the panel.",
						},
					},
				},
				"source": schema.StringAttribute{
					Description: "The JSON source of the panel.",
					Required:    true,
				},
			},
		},
	}
}

func (d *DashboardDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description:         "Dashboard data source.",
		MarkdownDescription: "Dashboard data source. See Grafana [documentation](https://grafana.com/docs/grafana/latest/dashboards/use-dashboards/) for more details.",

		Blocks: map[string]schema.Block{
			"time": schema.ListNestedBlock{
				Description: "The time-specific options.",
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"default_range": dashboardTimeBlock(),
						"picker": schema.ListNestedBlock{
							Description: "The time picker options.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"refresh_intervals": schema.ListAttribute{
										ElementType: types.StringType,
										Optional:    true,
										Description: "The auto refresh intervals that should be available in the auto refresh list. " +
											"The following time units are supported: s (seconds), m (minutes), h (hours), d (days), w (weeks), M (months), and y (years). " +
											"Example: 5s, 10s, 30s, 1m, 5m, 15m, 30m, 1h, 2h, 1d.",
										MarkdownDescription: "The auto refresh intervals that should be available in the auto refresh list. " +
											"The following time units are supported: `s (seconds)`, `m (minutes)`, `h (hours)`, `d (days)`, `w (weeks)`, `M (months)`, and `y (years)`. " +
											"Example: `1m`, `10s`, `30s`, `1m`, `5m`, `15m`, `30m`, `1h`, `2h`, `1d`.",
									},
									"time_options": schema.ListAttribute{
										ElementType: types.StringType,
										Optional:    true,
										Description: "The time options (Last X) that should be available in the range selection list. " +
											"The following time units are supported: s (seconds), m (minutes), h (hours), d (days), w (weeks), M (months), and y (years). " +
											"Example: 5s, 10s, 30s, 1m, 5m, 15m, 30m, 1h, 2h, 1d.",
										MarkdownDescription: "The time options (Last X) that should be available in the range selection list. " +
											"The following time units are supported: `s (seconds)`, `m (minutes)`, `h (hours)`, `d (days)`, `w (weeks)`, `M (months)`, and `y (years)`. " +
											"Example: `1m`, `10s`, `30s`, `1m`, `5m`, `15m`, `30m`, `1h`, `2h`, `1d`.",
									},
									"now_delay": schema.StringAttribute{
										Optional:            true,
										Description:         "Exclude recent data that may be incomplete. Example: 10s, 1m, 15m.",
										MarkdownDescription: "Exclude recent data that may be incomplete. Example: `10s`, `1m`, `15m.`",
									},
									"hide": schema.BoolAttribute{
										Optional:    true,
										Description: "Whether to hide time picker or not.",
									},
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(1),
							},
						},
					},

					Attributes: map[string]schema.Attribute{
						"timezone": schema.StringAttribute{
							Optional:            true,
							Description:         "The timezone to use. Predefined: utc, browser. Custom: Europe/Kyiv.",
							MarkdownDescription: "The timezone to use. Predefined: `utc`, `browser`. Custom: `Europe/Kyiv`.",
						},
						"week_start": schema.StringAttribute{
							Optional:            true,
							Description:         "The custom week start. The choices are: saturday, sunday, monday.",
							MarkdownDescription: "The custom week start. The choices are: `saturday`, `sunday`, `monday`.",
							Validators: []validator.String{
								stringvalidator.OneOf("saturday", "sunday", "monday"),
							},
						},
						"refresh_live_dashboards": schema.BoolAttribute{
							Optional:    true,
							Description: "Continuously re-draw panels where the time range references 'now'.",
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},

			"variables": schema.ListNestedBlock{
				Description: "The variables.",
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"custom": schema.ListNestedBlock{
							Description: "The variable options defined as a comma-separated list.",

							NestedObject: schema.NestedBlockObject{
								Blocks: map[string]schema.Block{
									"option": schema.ListNestedBlock{
										Description: "The option entry.",
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"text": schema.StringAttribute{
													Required:    true,
													Description: "The text (label) of the entry.",
												},
												"value": schema.StringAttribute{
													Required:    true,
													Description: "The value of the entry.",
												},
												"selected": schema.BoolAttribute{
													Optional:    true,
													Description: "Whether to mark the option as selected or not.",
												},
											},
										},
										Validators: []validator.List{
											listvalidator.SizeAtLeast(1),
											listvalidator.SizeAtMost(10),
										},
									},
									"include_all": variableIncludeAllBlock(),
								},
								Attributes: map[string]schema.Attribute{
									"name": variableNameAttribute(),
									"label": schema.StringAttribute{
										Optional:    true,
										Description: "The optional display name.",
									},
									"description": schema.StringAttribute{
										Optional:    true,
										Description: "The description of the variable.",
									},
									"multi": schema.BoolAttribute{
										Optional:    true,
										Description: "Whether to allow selecting multiple values at the same time or not.",
									},
									"hide": variableHideAttribute(),
								},
							},

							Validators: []validator.List{
								listvalidator.SizeAtLeast(1),
								listvalidator.SizeAtMost(20),
							},
						},

						"const": schema.ListNestedBlock{
							Description: "The constant variable. Defines a hidden constant.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"name": variableNameAttribute(),
									"label": schema.StringAttribute{
										Optional:    true,
										Description: "The optional display name.",
									},
									"description": schema.StringAttribute{
										Optional:    true,
										Description: "The description of the variable.",
									},
									"value": schema.StringAttribute{
										Description: "The value of the variable.",
										Required:    true,
									},
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(20),
							},
						},

						"textbox": schema.ListNestedBlock{
							Description: "The textbox variable. Displays a free text input field with an optional default value.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"name": variableNameAttribute(),
									"label": schema.StringAttribute{
										Optional:    true,
										Description: "The optional display name.",
									},
									"description": schema.StringAttribute{
										Optional:    true,
										Description: "The description of the variable.",
									},
									"default_value": schema.StringAttribute{
										Description: "The default value of the variable to use when the textbox is empty.",
										Optional:    true,
									},
									"hide": variableHideAttribute(),
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(20),
							},
						},

						"adhoc": schema.ListNestedBlock{
							Description: "The adhoc variable. " +
								"Allows adding key/value filters that are automatically added to all metric queries that use the specified data source. " +
								"Unlike other variables, you do not use ad hoc filters in queries. Instead, you use ad hoc filters to write filters for existing queries.",

							NestedObject: schema.NestedBlockObject{
								Blocks: map[string]schema.Block{
									"datasource": schema.SingleNestedBlock{
										Description: "The datasource to use.",
										Attributes: map[string]schema.Attribute{
											"type": schema.StringAttribute{
												Required:            true,
												Description:         "The type of the datasource. The choices are: prometheus, loki, influxdb, elasticsearch.",
												MarkdownDescription: "The type of the datasource. The choices are: `prometheus`, `loki`, `influxdb`, `elasticsearch`.",
												Validators: []validator.String{
													stringvalidator.OneOf("prometheus", "loki", "influxdb", "elasticsearch"),
												},
											},
											"uid": schema.StringAttribute{
												Required:    true,
												Description: "The uid of the datasource.",
											},
										},
									},
									"filter": schema.ListNestedBlock{
										Description: "The predefined filters.",
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"key": schema.StringAttribute{
													Required:    true,
													Description: "The name of the dimensional label to filter by.",
												},
												"operator": schema.StringAttribute{
													Required:            true,
													Description:         "The operator to use for comparison. The choices are: =, !=, >, <, =~, !~.",
													MarkdownDescription: "The operator to use for comparison. The choices are: `=`, `!=`, `>`, `<`, `=~`, `!~`.",
													Validators: []validator.String{
														stringvalidator.OneOf("=", "!=", ">", "<", "=~", "!~"),
													},
												},
												"value": schema.StringAttribute{
													Required:    true,
													Description: "The expected value of the dimensional label.",
												},
											},
										},
										Validators: []validator.List{
											listvalidator.SizeAtMost(30),
										},
									},
								},

								Attributes: map[string]schema.Attribute{
									"name": variableNameAttribute(),
									"label": schema.StringAttribute{
										Optional:    true,
										Description: "The optional display name.",
									},
									"description": schema.StringAttribute{
										Optional:    true,
										Description: "The description of the variable.",
									},
									"hide": variableHideAttribute(),
								},

								Validators: []validator.Object{
									objectvalidator.AlsoRequires(path.MatchRelative().AtName("datasource")),
								},
							},

							Validators: []validator.List{
								listvalidator.SizeAtMost(20),
							},
						},

						"datasource": schema.ListNestedBlock{
							Description: "The datasource variable. Quickly change the data source for an entire dashboard.",

							NestedObject: schema.NestedBlockObject{
								Blocks: map[string]schema.Block{
									"source": schema.SingleNestedBlock{
										Description: "The datasource selector.",
										Attributes: map[string]schema.Attribute{
											"type": schema.StringAttribute{
												Required:            true,
												Description:         "The type of the datasource. Example: prometheus, loki, influxdb, elasticsearch, cloudwatch, etc.",
												MarkdownDescription: "The type of the datasource. Example: `prometheus`, `loki`, `influxdb`, `elasticsearch`, `cloudwatch`, etc.",
											},
											"filter": schema.StringAttribute{
												Optional:            true,
												Description:         "Regex filter for which data source instances to choose from in the variable value list. Leave empty for all. Example: /^prod/.",
												MarkdownDescription: "Regex filter for which data source instances to choose from in the variable value list. Leave empty for all. Example: `/^prod/`.",
											},
										},
									},
									"include_all": variableIncludeAllBlock(),
								},

								Attributes: map[string]schema.Attribute{
									"name": variableNameAttribute(),
									"label": schema.StringAttribute{
										Optional:    true,
										Description: "The optional display name.",
									},
									"description": schema.StringAttribute{
										Optional:    true,
										Description: "The description of the variable.",
									},
									"multi": schema.BoolAttribute{
										Optional:    true,
										Description: "Whether to allow selecting multiple values at the same time or not.",
									},
									"hide": variableHideAttribute(),
								},

								Validators: []validator.Object{
									objectvalidator.AlsoRequires(path.MatchRelative().AtName("source")),
								},
							},

							Validators: []validator.List{
								listvalidator.SizeAtLeast(1),
								listvalidator.SizeAtMost(20),
							},
						},

						"query": schema.ListNestedBlock{
							Description: "The query variable. Allows adding a query that can return a list of metric names, tag values, or keys.",

							NestedObject: schema.NestedBlockObject{
								Blocks: map[string]schema.Block{
									"sort": schema.ListNestedBlock{
										Description: "The sort order for values to be displayed in the dropdown list.",
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"type": schema.StringAttribute{
													Required:            true,
													Description:         "The type of sorting. The choices are: disabled, alphabetical, numerical, alphabetical-case-insensitive.",
													MarkdownDescription: "The type of sorting. The choices are: `disabled`, `alphabetical`, `numerical`, `alphabetical-case-insensitive`.",
													Validators: []validator.String{
														stringvalidator.OneOf("disabled", "alphabetical", "numerical", "alphabetical-case-insensitive"),
													},
												},
												"order": schema.StringAttribute{
													Optional:            true,
													Description:         "The order of the sorting. The choices are: asc, desc.",
													MarkdownDescription: "The order of the sorting. The choices are: `asc`, `desc`.",
													Validators: []validator.String{
														stringvalidator.OneOf("asc", "desc"),
													},
												},
											},
										},
										Validators: []validator.List{
											listvalidator.SizeAtLeast(1),
											listvalidator.SizeAtMost(1),
										},
									},
									"target": schema.ListNestedBlock{
										Description: "The datasource-specific query.",
										NestedObject: schema.NestedBlockObject{
											Blocks: map[string]schema.Block{
												"prometheus": schema.ListNestedBlock{
													NestedObject: schema.NestedBlockObject{
														Attributes: map[string]schema.Attribute{
															"uid": schema.StringAttribute{
																Required:    true,
																Description: "The uid of the datasource.",
															},
															"expr": schema.StringAttribute{
																Required:    true,
																Description: "The query expression.",
															},
														},
													},
													Validators: []validator.List{
														listvalidator.SizeAtLeast(1),
														listvalidator.SizeAtMost(1),
													},
												},
											},
											Validators: []validator.Object{
												objectvalidator.AlsoRequires(path.MatchRelative().AtName("prometheus")),
											},
										},
										Validators: []validator.List{
											listvalidator.SizeAtLeast(1),
											listvalidator.SizeAtMost(1),
										},
									},
									"include_all": variableIncludeAllBlock(),
								},

								Attributes: map[string]schema.Attribute{
									"name": variableNameAttribute(),
									"label": schema.StringAttribute{
										Optional:    true,
										Description: "The optional display name.",
									},
									"description": schema.StringAttribute{
										Optional:    true,
										Description: "The description of the variable.",
									},
									"multi": schema.BoolAttribute{
										Optional:    true,
										Description: "Whether to allow selecting multiple values at the same time or not.",
									},
									"refresh": schema.StringAttribute{
										Optional:            true,
										Description:         "When to update the values of this variable. The choices are: dashboard-load, time-range-change.",
										MarkdownDescription: "When to update the values of this variable. The choices are: `dashboard-load`, `time-range-change`.",
										Validators: []validator.String{
											stringvalidator.OneOf("dashboard-load", "time-range-change"),
										},
									},
									"regex": schema.StringAttribute{
										Optional:            true,
										Description:         "The regex expression to filter or capture specific parts of the names returned by your data source query. Example: /.*instance=\"([^\"]*).*/.",
										MarkdownDescription: "The regex expression to filter or capture specific parts of the names returned by your data source query. Example: `/.*instance=\"([^\"]*).*/`.",
									},
									"hide": variableHideAttribute(),
								},

								Validators: []validator.Object{
									objectvalidator.AlsoRequires(path.MatchRelative().AtName("target")),
								},
							},

							Validators: []validator.List{
								listvalidator.SizeAtMost(20),
							},
						},

						"interval": schema.ListNestedBlock{
							Description: "The interval variable. Represents time spans such as 1m, 1h, 1d. You can think of them as a dashboard-wide 'group by time' command. " +
								"Interval variables change how the data is grouped in the visualization.",

							MarkdownDescription: "The interval variable. Represents time spans such as `1m`, `1h`, `1d`. You can think of them as a dashboard-wide 'group by time' command. " +
								"Interval variables change how the data is grouped in the visualization.",

							NestedObject: schema.NestedBlockObject{
								Blocks: map[string]schema.Block{
									"auto": schema.ListNestedBlock{
										Description: "Defines how many times the current time range should be divided to calculate the current auto time span.",
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"enabled": schema.BoolAttribute{
													Required:    true,
													Description: "Whether to enable calculation of auto time spans or not.",
												},
												"step_count": schema.Int64Attribute{
													Required: true,
													Description: "How many times the current time range should be divided to calculate the value. " +
														"The choices are: 1, 2, 3, 4, 5, 10, 20, 30, 40, 50, 100, 200, 300, 400, 500.",
													MarkdownDescription: "How many times the current time range should be divided to calculate the value. " +
														"The choices are: `1`, `2`, `3`, `4`, `5`, `10`, `20`, `30`, `40`, `50`, `100`, `200`, `300`, `400`, `500`.",
													Validators: []validator.Int64{
														int64validator.OneOf(1, 2, 3, 4, 5, 10, 20, 30, 40, 50, 100, 200, 300, 400, 500),
													},
												},
												"min_interval": schema.StringAttribute{
													Required:    true,
													Description: "The calculated value will not go below this threshold.",
												},
											},
										},
										Validators: []validator.List{
											listvalidator.SizeAtLeast(1),
											listvalidator.SizeAtMost(1),
										},
									},
								},

								Attributes: map[string]schema.Attribute{
									"name": variableNameAttribute(),
									"label": schema.StringAttribute{
										Optional:    true,
										Description: "The optional display name.",
									},
									"description": schema.StringAttribute{
										Optional:    true,
										Description: "The description of the variable.",
									},
									"hide": variableHideAttribute(),
									"intervals": schema.ListAttribute{
										ElementType: types.StringType,
										Required:    true,
										Description: "The time range intervals that you want to appear in the variable drop-down list. " +
											"The following time units are supported: s (seconds), m (minutes), h (hours), d (days), w (weeks), M (months), and y (years). " +
											"Example: 1m, 10m, 30m, 1h, 6h, 12h, 1d, 7d, 14d, 30d.",
										MarkdownDescription: "The time range intervals that you want to appear in the variable drop-down list. " +
											"The following time units are supported: `s (seconds)`, `m (minutes)`, `h (hours)`, `d (days)`, `w (weeks)`, `M (months)`, and `y (years)`. " +
											"Example: `1m`, `10m`, `30m`, `1h`, `6h`, `12h`, `1d`, `7d`, `14d`, `30d`.",
									},
								},
							},

							Validators: []validator.List{
								listvalidator.SizeAtMost(20),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(10),
				},
			},

			"layout": schema.SingleNestedBlock{
				Description: "The layout of the dashboard.",
				Blocks: map[string]schema.Block{
					"section": schema.ListNestedBlock{
						Description: "The row within the dashboard.",
						NestedObject: schema.NestedBlockObject{
							Blocks: map[string]schema.Block{
								"panel": panelBlock(),
								"row": schema.ListNestedBlock{
									Description: "The new row to align the nested panels.",
									NestedObject: schema.NestedBlockObject{
										Blocks: map[string]schema.Block{
											"panel": panelBlock(),
										},
									},
									Validators: []validator.List{
										listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("panel")),
									},
								},
							},
							Attributes: map[string]schema.Attribute{
								"title": schema.StringAttribute{
									Optional:    true,
									Description: "The title of the row. If the title is defined the row is treated as collapsible.",
								},
								"collapsed": schema.BoolAttribute{
									Optional:    true,
									Description: "Whether the row is collapsed or not.",
								},
							},
						},
					},
				},
				Validators: []validator.Object{
					objectvalidator.AtLeastOneOf(
						path.MatchRoot("layout"),
					),
				},
			},

			"annotations": schema.ListNestedBlock{
				Description: "The annotations to add to the dashboard.",

				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"grafana": schema.ListNestedBlock{
							Description: "The Grafana annotation query.",
							NestedObject: schema.NestedBlockObject{
								Blocks: map[string]schema.Block{
									"by_dashboard": schema.ListNestedBlock{
										Description: "Query for events created on this dashboard and show them in the panels where they were created.",
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"limit": schema.Int64Attribute{
													Required:    true,
													Description: "The limit of events.",
												},
											},
										},
										Validators: []validator.List{
											listvalidator.SizeAtMost(1),
											listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("by_tags")),
										},
									},
									"by_tags": schema.ListNestedBlock{
										Description: "This will fetch any annotation events that match the tags filter.",
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"limit": schema.Int64Attribute{
													Required:    true,
													Description: "The limit of events.",
												},
												"match_any": schema.BoolAttribute{
													Optional:    true,
													Description: "Enabling this returns annotations that match any of the tags specified below.",
												},
												"tags": schema.ListAttribute{
													ElementType: types.StringType,
													Optional:    true,
													Description: "Specify a list of tags to match. To specify a key and value tag use `key:value` syntax.",
												},
											},
										},
										Validators: []validator.List{
											listvalidator.SizeAtMost(1),
											listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("by_dashboard")),
										},
									},
								},
								Attributes: annotationCommonAttributes(),
							},
						},

						"prometheus": schema.ListNestedBlock{
							Description: "The Prometheus annotation query.",
							NestedObject: schema.NestedBlockObject{
								Blocks: map[string]schema.Block{
									"query": schema.SingleNestedBlock{
										Attributes: map[string]schema.Attribute{
											"datasource_uid": schema.StringAttribute{
												Required:    true,
												Description: "The uid of the datasource.",
											},
											"expr": schema.StringAttribute{
												Required:    true,
												Description: "The query expression.",
											},
											"min_step": schema.StringAttribute{
												Optional:    true,
												Description: "The minimum step interval to use when evaluating the query. ",
											},
											"title_format": schema.StringAttribute{
												Optional:            true,
												Description:         "Use either the name or a pattern. For example, {{instance}} is replaced with the label value for the label instance.",
												MarkdownDescription: "Use either the name or a pattern. For example, `{{instance}}` is replaced with the label value for the label instance.",
											},
											"text_format": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "Use either the name or a pattern. For example, `{{instance}}` is replaced with the label value for the label instance.",
											},
											"tag_keys": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "The tags to use.",
											},
											"use_value_as_timestamp": schema.BoolAttribute{
												Optional:    true,
												Description: "Treat the value of the series as a timestamp.",
											},
										},
									},
								},
								Attributes: annotationCommonAttributes(),
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(20),
							},
						},
					},
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(50),
				},
			},

			"links": schema.ListNestedBlock{
				Description: "The links to add to the dashboard.",

				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"dashboards": schema.ListNestedBlock{
							Description: "The dashboards link.",

							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"title": schema.StringAttribute{
										Optional:    true,
										Description: "The title you want the link to display.",
									},
									"as_dropdown": schema.BoolAttribute{
										Optional:    true,
										Description: "Whether to show links as a dropdown. Otherwise, Grafana displays the dashboard links side by side across the top of your dashboard.",
									},
									"tags": schema.ListAttribute{
										ElementType: types.StringType,
										Optional:    true,
										Description: "Include dashboards only with these tags. Otherwise, Grafana includes links to all other dashboards.",
									},
									"include_time_range": schema.BoolAttribute{
										Optional:    true,
										Description: "Whether to include the dashboard's time range in the link.",
									},
									"include_template_variables": schema.BoolAttribute{
										Optional:    true,
										Description: "Whether to include the dashboard's template variables in the link.",
									},
									"new_tab": schema.BoolAttribute{
										Optional:    true,
										Description: "Whether to open a link in a new tab or window.",
									},
								},
							},
						},

						"external": schema.ListNestedBlock{
							Description: "The external link.",

							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"title": schema.StringAttribute{
										Optional:    true,
										Description: "The title you want the link to display.",
									},
									"url": schema.StringAttribute{
										Optional:    true,
										Description: "The url you want to link to.",
									},
									"tooltip": schema.StringAttribute{
										Optional:    true,
										Description: "The tooltip to display when the user hovers their mouse over the link.",
									},
									"icon": schema.StringAttribute{
										Optional:            true,
										Description:         "The icon you want displayed with the link. The choices are: external link, dashboard, question, info, bolt, doc, cloud.",
										MarkdownDescription: "The icon you want displayed with the link. The choices are: `external link`, `dashboard`, `question`, `info`, `bolt`, `doc`, `cloud`.",
										Validators: []validator.String{
											stringvalidator.OneOf("external link", "dashboard", "question", "info", "bolt", "doc", "cloud"),
										},
									},
									"include_time_range": schema.BoolAttribute{
										Optional:    true,
										Description: "Whether to include the dashboard's time range in the link.",
									},
									"include_template_variables": schema.BoolAttribute{
										Optional:    true,
										Description: "Whether to include the dashboard's template variables in the link.",
									},
									"new_tab": schema.BoolAttribute{
										Optional:    true,
										Description: "Whether to open a link in a new tab or window.",
									},
								},
							},
						},
					},
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(50),
				},
			},
		},

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"json":         jsonAttribute(),
			"compact_json": compactJsonAttribute(),
			"title": schema.StringAttribute{
				Required:    true,
				Description: "The title of the dashboard.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the dashboard.",
			},
			"uid": schema.StringAttribute{
				Optional:    true,
				Description: "The UID of the dashboard.",
			},
			"tags": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "The set of tags to associate with the dashboard.",
			},
			"version": schema.Int64Attribute{
				Optional:    true,
				Description: "The version of the dashboard.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"editable":      dashboardEditableAttribute(),
			"style":         dashboardStyleAttribute(),
			"graph_tooltip": dashboardGraphTooltipAttribute(),
		},
	}
}

func (d *DashboardDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.Defaults = defaults.Dashboard
}

func decodeHide(hideValue types.String) uint8 {
	hide := uint8(0)

	if !hideValue.IsNull() {
		switch v := hideValue.ValueString(); v {
		case "label":
			hide = 1
		case "variable":
			hide = 2
		default:
			hide = 0
		}
	}

	return hide
}

func findFreeBlock(matrix [][]uint8, H int, W int) (int, int) {
	rows := len(matrix)

	for i := 0; i <= rows-H; i++ {
		cols := len(matrix[i])

		for j := 0; j <= cols-W; j++ {
			allZero := true
			for k := i; k < i+H; k++ {
				for l := j; l < j+W; l++ {
					if matrix[k][l] != 0 {
						allZero = false
						break
					}
				}
				if !allZero {
					break
				}
			}
			if allZero {
				return i, j
			}
		}
	}
	return -1, -1
}

// todo add verification of bounds?
func calculateAutoLayout(panels []Panel, startY int) ([]grafana.Panel, error) {
	matrix := make([][]uint8, 0)
	result := make([]grafana.Panel, 0)

	for _, panel := range panels {
		var grafanaPanel grafana.Panel

		err := json.Unmarshal([]byte(panel.Source.ValueString()), &grafanaPanel)
		if err != nil {
			return result, err
		}

		height := int(panel.Size.Height.ValueInt64())
		width := int(panel.Size.Width.ValueInt64())

		var y, x = findFreeBlock(matrix, height, width)

		if y == -1 {
			y = len(matrix)
			x = 0

			for i := 0; i < height; i++ {
				matrix = append(matrix, make([]uint8, 24))
			}
		}

		for i := 0; i < height; i++ {
			for j := 0; j < width; j++ {
				matrix[y+i][x+j] = 1
			}
		}

		/*fmt.Println("New matrix:") // Move to the next line after printing each row
		for i := 0; i < len(matrix); i++ {
			for j := 0; j < len(matrix[i]); j++ {
				fmt.Printf("%4d", matrix[i][j]) // Use a fixed width of 4 characters
			}
			fmt.Println() // Move to the next line after printing each row
		}*/

		posX := x
		posY := startY + y
		grafanaPanel.GridPos = grafana.GridPos{
			H: &height,
			W: &width,
			X: &posX,
			Y: &posY,
		}

		result = append(result, grafanaPanel)
	}

	return result, nil
}

func calculateManualLayout(rows []SectionRow, startY int) ([]grafana.Panel, error) {
	result := make([]grafana.Panel, 0)

	for rowIdx, row := range rows {
		for columnIdx, panel := range row.Panels {
			var grafanaPanel grafana.Panel

			err := json.Unmarshal([]byte(panel.Source.ValueString()), &grafanaPanel)
			if err != nil {
				return result, err
			}

			height := int(panel.Size.Height.ValueInt64())
			width := int(panel.Size.Width.ValueInt64())

			var x int
			var y int

			if columnIdx == 0 {
				x = 0
			} else {
				total := 0
				for _, item := range row.Panels[0:columnIdx] {
					total += int(item.Size.Width.ValueInt64())
				}
				x = total
			}

			if rowIdx == 0 {
				y = 0
			} else {
				total := 0
				for _, r := range rows[0:rowIdx] {
					max := 0
					for _, c := range r.Panels {
						max = int(math.Max(float64(max), float64(c.Size.Height.ValueInt64())))
					}
					total += max
				}
				y = total + rowIdx
			}

			posX := x
			posY := startY + y

			grafanaPanel.GridPos = grafana.GridPos{
				H: &height,
				W: &width,
				X: &posX,
				Y: &posY,
			}

			result = append(result, grafanaPanel)
		}
	}

	return result, nil
}

func (d *DashboardDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DashboardDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	vars := make([]grafana.TemplateVar, 0)
	for _, variable := range data.Variables {
		for _, custom := range variable.Custom {
			opts := make([]grafana.Option, len(custom.Options))
			query := ""
			var current grafana.Option

			for i, opt := range custom.Options {
				text := opt.Text.ValueString()
				opts[i] = grafana.Option{
					Text:     &text,
					Value:    opt.Value.ValueString(),
					Selected: opt.Selected.ValueBool(),
				}

				if query != "" {
					query = query + ", "
				}

				query = query + opt.Text.ValueString() + " : " + opt.Value.ValueString()

				if opt.Selected.ValueBool() {
					current = opts[i]
				}
			}

			v := grafana.TemplateVar{
				Type:        "custom",
				Name:        custom.Name.ValueString(),
				Label:       custom.Label.ValueString(),
				Description: custom.Description.ValueString(),
				Multi:       custom.Multi.ValueBool(),
				Options:     opts,
				Query:       query,
				Current:     current,
				Hide:        decodeHide(custom.Hide),
			}

			for _, all := range custom.IncludeAll {
				if !all.Enabled.IsNull() {
					v.IncludeAll = all.Enabled.ValueBool()
				}

				if !all.CustomValue.IsNull() {
					v.AllValue = all.CustomValue.ValueString()
				}
			}

			vars = append(vars, v)
		}

		for _, c := range variable.Constant {
			v := grafana.TemplateVar{
				Type:        "constant",
				Options:     make([]grafana.Option, 0),
				Name:        c.Name.ValueString(),
				Label:       c.Label.ValueString(),
				Description: c.Description.ValueString(),
				Query:       c.Value.ValueString(),
				Hide:        uint8(2),
			}

			vars = append(vars, v)
		}

		for _, textbox := range variable.TextBox {
			v := grafana.TemplateVar{
				Type:        "textbox",
				Options:     make([]grafana.Option, 0),
				Name:        textbox.Name.ValueString(),
				Label:       textbox.Label.ValueString(),
				Description: textbox.Description.ValueString(),
				Query:       textbox.DefaultValue.ValueString(),
				Hide:        decodeHide(textbox.Hide),
			}

			vars = append(vars, v)
		}

		for _, adhoc := range variable.AdHoc {
			filters := make([]grafana.TemplateVarAdHocFilter, len(adhoc.Filters))

			for i, filter := range adhoc.Filters {
				filters[i] = grafana.TemplateVarAdHocFilter{
					Condition: "",
					Key:       filter.Key.ValueString(),
					Operator:  filter.Operator.ValueString(),
					Value:     filter.Value.ValueString(),
				}
			}

			v := grafana.TemplateVar{
				Type:        "adhoc",
				Options:     make([]grafana.Option, 0),
				Name:        adhoc.Name.ValueString(),
				Label:       adhoc.Label.ValueString(),
				Description: adhoc.Description.ValueString(),
				Datasource: &grafana.TemplateVarDataSource{
					UID:  adhoc.DataSource.UID.ValueString(),
					Type: adhoc.DataSource.Type.ValueString(),
				},
				Filters: filters,
				Hide:    decodeHide(adhoc.Hide),
			}

			vars = append(vars, v)
		}

		for _, ds := range variable.DataSource {
			refresh := int64(1)
			v := grafana.TemplateVar{
				Type:        "datasource",
				Options:     make([]grafana.Option, 0),
				Name:        ds.Name.ValueString(),
				Label:       ds.Label.ValueString(),
				Description: ds.Description.ValueString(),
				Multi:       ds.Multi.ValueBool(),
				Query:       ds.DataSource.Type.ValueString(),
				Regex:       ds.DataSource.Filter.ValueString(),
				Hide:        decodeHide(ds.Hide),
				Refresh:     grafana.BoolInt{Value: &refresh},
			}

			for _, all := range ds.IncludeAll {
				if !all.Enabled.IsNull() {
					v.IncludeAll = all.Enabled.ValueBool()
				}

				if !all.CustomValue.IsNull() {
					v.AllValue = all.CustomValue.ValueString()
				}
			}

			vars = append(vars, v)
		}

		for _, query := range variable.Query {
			refresh := int64(1)

			if !query.Refresh.IsNull() {
				switch v := query.Refresh.ValueString(); v {
				case "dashboard-load":
					refresh = 1
				case "time-range-change":
					refresh = 2
				default:
					refresh = 1
				}
			}

			sort := 1

			for _, s := range query.Sort {
				asc := s.Order.ValueString() != "desc"

				if !s.Type.IsNull() {
					switch v := s.Type.ValueString(); v {
					case "disabled":
						sort = 0
					case "alphabetical":
						if asc {
							sort = 1
						} else {
							sort = 2
						}
					case "numerical":
						if asc {
							sort = 3
						} else {
							sort = 4
						}
					case "alphabetical-case-insensitive":
						if asc {
							sort = 5
						} else {
							sort = 6
						}
					default:
						sort = 1
					}
				}
			}

			v := grafana.TemplateVar{
				Type:        "query",
				Options:     make([]grafana.Option, 0),
				Name:        query.Name.ValueString(),
				Label:       query.Label.ValueString(),
				Description: query.Description.ValueString(),
				Multi:       query.Multi.ValueBool(),
				Hide:        decodeHide(query.Hide),
				Sort:        sort,
				Refresh:     grafana.BoolInt{Value: &refresh},
				Regex:       query.Regex.ValueString(),
			}

			for _, target := range query.Target {
				for _, prometheus := range target.Prometheus {
					v.Datasource = &grafana.TemplateVarDataSource{
						UID:  prometheus.UID.ValueString(),
						Type: "prometheus",
					}
					v.Definition = prometheus.Expr.ValueString()
					v.Query = grafana.TemplateVarQueryPrometheus{
						Query: prometheus.Expr.ValueString(),
						RefID: "StandardVariableQuery",
					}
				}
			}

			for _, all := range query.IncludeAll {
				if !all.Enabled.IsNull() {
					v.IncludeAll = all.Enabled.ValueBool()
				}

				if !all.CustomValue.IsNull() {
					v.AllValue = all.CustomValue.ValueString()
				}
			}

			vars = append(vars, v)
		}

		for _, interval := range variable.Interval {
			query := ""
			totalIntervals := len(interval.Intervals)
			opts := make([]grafana.Option, totalIntervals)

			for i, intervalValue := range interval.Intervals {
				text := intervalValue.ValueString()
				opts[i] = grafana.Option{
					Text:     &text,
					Value:    intervalValue.ValueString(),
					Selected: i == 0,
				}

				query = query + intervalValue.ValueString()

				if i < (totalIntervals - 1) {
					query = query + ","
				}
			}

			v := grafana.TemplateVar{
				Type:        "interval",
				Options:     opts,
				Current:     opts[0],
				Name:        interval.Name.ValueString(),
				Label:       interval.Label.ValueString(),
				Description: interval.Description.ValueString(),
				Hide:        decodeHide(interval.Hide),
				Query:       query,
			}

			for _, auto := range interval.Auto {
				v.Auto = auto.Enabled.ValueBool()

				if !auto.StepCount.IsNull() {
					v.AutoCount = auto.StepCount.ValueInt64Pointer()
				}

				if !auto.MinInterval.IsNull() {
					v.AutoMin = auto.MinInterval.ValueStringPointer()
				}

				if v.Auto {
					text := "auto"
					autoOpt := grafana.Option{
						Text:  &text,
						Value: "$__auto_interval_" + interval.Name.ValueString(),
					}
					v.Options = append([]grafana.Option{autoOpt}, opts...)
				}
			}

			vars = append(vars, v)
		}
	}

	annotations := make([]grafana.Annotation, 0)
	for _, annotation := range data.Annotations {
		for _, grafanaQuery := range annotation.Grafana {
			hide := true
			result := grafana.Annotation{
				Name:      grafanaQuery.Name.ValueString(),
				Enable:    true,
				Hide:      &hide,
				IconColor: "rgba(0, 211, 255, 1)",
				Datasource: grafana.AnnotationDataSource{
					UID:  "-- Grafana --",
					Type: "prometheus",
				},
			}

			if !grafanaQuery.Enabled.IsNull() {
				result.Enable = grafanaQuery.Enabled.ValueBool()
			}

			if !grafanaQuery.Hidden.IsNull() {
				result.Enable = grafanaQuery.Enabled.ValueBool()
			}

			if !grafanaQuery.Color.IsNull() {
				result.IconColor = grafanaQuery.Color.ValueString()
			}

			for _, byDashboard := range grafanaQuery.ByDashboard {
				result.Target = &grafana.AnnotationGrafanaTarget{
					Limit: 100,
					Type:  "dashboard",
				}

				if !byDashboard.Limit.IsNull() {
					result.Target.Limit = byDashboard.Limit.ValueInt64()
				}
			}

			for _, byTags := range grafanaQuery.ByTags {
				tags := make([]string, 0)

				for _, tag := range byTags.Tags {
					tags = append(tags, tag.ValueString())
				}

				result.Target = &grafana.AnnotationGrafanaTarget{
					Limit:    100,
					MatchAny: byTags.MatchAny.ValueBool(),
					Tags:     tags,
					Type:     "tags",
				}

				if !byTags.Limit.IsNull() {
					result.Target.Limit = byTags.Limit.ValueInt64()
				}
			}

			annotations = append(annotations, result)
		}

		for _, prometheus := range annotation.Prometheus {
			result := grafana.Annotation{
				Name:      prometheus.Name.ValueString(),
				Enable:    true,
				Hide:      prometheus.Hidden.ValueBoolPointer(),
				IconColor: "red",
				Datasource: grafana.AnnotationDataSource{
					UID:  prometheus.Query.UID.ValueString(),
					Type: "prometheus",
				},
				Expr:            prometheus.Query.Expr.ValueStringPointer(),
				Step:            prometheus.Query.Step.ValueStringPointer(),
				UseValueForTime: prometheus.Query.UseValueAsTimestamp.ValueBoolPointer(),
				TitleFormat:     prometheus.Query.Title.ValueStringPointer(),
				TextFormat:      prometheus.Query.Text.ValueStringPointer(),
				TagKeys:         prometheus.Query.TagKeys.ValueStringPointer(),
			}

			if !prometheus.Enabled.IsNull() {
				result.Enable = prometheus.Enabled.ValueBool()
			}

			if !prometheus.Color.IsNull() {
				result.IconColor = prometheus.Color.ValueString()
			}

			annotations = append(annotations, result)
		}
	}

	links := make([]grafana.Link, 0)
	for _, link := range data.Links {
		for _, dashboards := range link.Dashboards {
			tags := make([]string, 0)
			for _, tag := range dashboards.Tags {
				if !tag.IsNull() {
					tags = append(tags, tag.ValueString())
				}
			}

			result := grafana.Link{
				Title:       dashboards.Title.ValueString(),
				Type:        "dashboards",
				AsDropdown:  dashboards.AsDropdown.ValueBoolPointer(),
				Tags:        tags,
				IncludeVars: dashboards.IncludeTemplateVariables.ValueBool(),
				KeepTime:    dashboards.IncludeTimeRange.ValueBoolPointer(),
				TargetBlank: dashboards.NewTab.ValueBoolPointer(),
			}

			links = append(links, result)
		}

		for _, external := range link.External {
			result := grafana.Link{
				Title:       external.Title.ValueString(),
				Type:        "link",
				URL:         external.Url.ValueStringPointer(),
				Tooltip:     external.Tooltip.ValueStringPointer(),
				Icon:        external.Icon.ValueStringPointer(),
				IncludeVars: external.IncludeTemplateVariables.ValueBool(),
				KeepTime:    external.IncludeTimeRange.ValueBoolPointer(),
				TargetBlank: external.NewTab.ValueBoolPointer(),
			}

			links = append(links, result)
		}
	}

	panels := make([]grafana.Panel, 0)

	for _, section := range data.Layout.Sections {
		isCollapsibleRow := !section.Title.IsNull() || section.Collapsed.ValueBool()
		sectionPanels := make([]grafana.Panel, 0)

		startY := 0

		if len(panels) == 0 {
			startY = 0
		} else {
			maxY := 0
			maxHeight := 0

			for _, panel := range panels {
				if *panel.GridPos.Y > maxY {
					maxY = *panel.GridPos.Y
					maxHeight = *panel.GridPos.H
				} else if *panel.GridPos.Y == maxY && *panel.GridPos.H > maxHeight {
					maxHeight = *panel.GridPos.H
				}

				if panel.Type == "row" {
					for _, collapsedPanel := range panel.Panels {
						if *collapsedPanel.GridPos.Y > maxY {
							maxY = *collapsedPanel.GridPos.Y
							maxHeight = *collapsedPanel.GridPos.H
						} else if *collapsedPanel.GridPos.Y == maxY && *collapsedPanel.GridPos.H > maxHeight {
							maxHeight = *collapsedPanel.GridPos.H
						}
					}
				}
			}

			startY = maxY + maxHeight + 1
		}

		if isCollapsibleRow {
			startY = startY + 1
		}

		// auto layout
		if len(section.Panels) > 0 {
			grafanaPanels, err := calculateAutoLayout(section.Panels, startY)
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Could not unmarshall json as Panel: %s", err))
				return
			}

			sectionPanels = append(sectionPanels, grafanaPanels...)
		} else { // manual layout
			grafanaPanels, err := calculateManualLayout(section.Rows, startY)
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Could not unmarshall json as Panel: %s", err))
				return
			}

			sectionPanels = append(sectionPanels, grafanaPanels...)
		}

		if isCollapsibleRow {
			x := 0
			y := startY - 1
			height := 1
			width := 24

			rowPanel := grafana.Panel{
				CommonPanel: grafana.CommonPanel{
					OfType: grafana.RowType,
					Title:  section.Title.ValueString(),
					Type:   "row",
					Span:   12,
					IsNew:  true,
					GridPos: grafana.GridPos{
						H: &height,
						W: &width,
						X: &x,
						Y: &y,
					},
				},
				RowPanel: &grafana.RowPanel{Collapsed: section.Collapsed.ValueBool()},
			}

			if section.Collapsed.ValueBool() {
				rowPanel.RowPanel.Panels = sectionPanels
				panels = append(panels, rowPanel)
			} else {
				panels = append(panels, rowPanel)
				panels = append(panels, sectionPanels...)
			}
		} else {
			panels = append(panels, sectionPanels...)
		}
	}

	dashboard := &grafana.Board{
		Title:         data.Title.ValueString(),
		Editable:      d.Defaults.Editable,
		Style:         d.Defaults.Style,
		SchemaVersion: 0,
		Version:       1,
		Panels:        panels,
		Templating: grafana.Templating{
			List: vars,
		},
		Time: grafana.Time{
			From: d.Defaults.Time.From,
			To:   d.Defaults.Time.To,
		},
	}

	if !data.Description.IsNull() {
		dashboard.Description = data.Description.ValueString()
	}

	if !data.UID.IsNull() {
		dashboard.UID = data.UID.ValueString()
	}

	if !data.Editable.IsNull() {
		dashboard.Editable = data.Editable.ValueBool()
	}

	if !data.Style.IsNull() {
		dashboard.Style = data.Style.ValueString()
	}

	if !data.Version.IsNull() {
		dashboard.Version = data.Version.ValueInt64()
	}

	if len(data.Tags) > 0 {
		tags := make([]string, len(data.Tags))
		for i, tag := range data.Tags {
			tags[i] = tag.ValueString()
		}

		dashboard.Tags = tags
	}

	if len(annotations) > 0 {
		dashboard.Annotations = grafana.Annotations{
			List: annotations,
		}
	}

	if len(links) > 0 {
		dashboard.Links = links
	}

	for _, timeOptions := range data.TimeOptions {
		if !timeOptions.Timezone.IsNull() {
			dashboard.Timezone = timeOptions.Timezone.ValueString()
		}

		if !timeOptions.WeekStart.IsNull() {
			dashboard.WeekStart = timeOptions.WeekStart.ValueString()
		}

		if !timeOptions.RefreshLiveDashboards.IsNull() {
			dashboard.LiveNow = timeOptions.RefreshLiveDashboards.ValueBool()
		}

		for _, time := range timeOptions.Range {
			dashboard.Time.From = time.From.ValueString()
			dashboard.Time.To = time.To.ValueString()
		}

		for _, picker := range timeOptions.TimePicker {
			if !picker.Hide.IsNull() {
				dashboard.Timepicker.Hidden = picker.Hide.ValueBoolPointer()
			}

			if !picker.NowDelay.IsNull() {
				dashboard.Timepicker.NowDelay = picker.NowDelay.ValueString()
			}

			if len(picker.RefreshIntervals) > 0 {
				intervals := make([]string, len(picker.RefreshIntervals))
				for i, interval := range picker.RefreshIntervals {
					intervals[i] = interval.ValueString()
				}
				dashboard.Timepicker.RefreshIntervals = intervals
			}

			if len(picker.TimeOptions) > 0 {
				options := make([]string, len(picker.TimeOptions))
				for i, option := range picker.TimeOptions {
					options[i] = option.ValueString()
				}
				dashboard.Timepicker.TimeOptions = options
			}
		}
	}

	tooltip := ""
	if !data.GraphTooltip.IsNull() {
		tooltip = data.GraphTooltip.ValueString()
	} else {
		tooltip = d.Defaults.GraphTooltip
	}

	if tooltip == "shared-crosshair" {
		dashboard.GraphTooltip = 1
	} else if tooltip == "shared-tooltip" {
		dashboard.GraphTooltip = 2
	}

	var jsonData []byte
	var err error

	if data.CompactJson.ValueBool() || d.CompactJson {
		jsonData, err = json.Marshal(dashboard)
	} else {
		jsonData, err = json.MarshalIndent(dashboard, "", "  ")
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Could not marshal json: %s", err))
		return
	}

	data.Json = types.StringValue(string(jsonData))
	data.Id = types.StringValue(strconv.Itoa(hashcode(jsonData)))

	// resp.Diagnostics.AddError("Client Error", fmt.Sprintf("%s", string(jsonData)))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
