package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/iRevive/terraform-provider-gdashboard/internal/provider/grafana"
	"math"
	"strconv"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &DashboardDataSource{}

func NewDashboardDataSource() datasource.DataSource {
	return &DashboardDataSource{}
}

// DashboardDataSource defines the data source implementation.
type DashboardDataSource struct {
	Defaults DashboardDefaults
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
	Id           types.String `tfsdk:"id"`
	Json         types.String `tfsdk:"json"`
	Title        types.String `tfsdk:"title"`
	UID          types.String `tfsdk:"uid"`
	Editable     types.Bool   `tfsdk:"editable"`
	Style        types.String `tfsdk:"style"`
	GraphTooltip types.String `tfsdk:"graph_tooltip"`
	Time         []TimeModel  `tfsdk:"time"`
	Layout       Layout       `tfsdk:"layout"`
	Variables    []Variable   `tfsdk:"variables"`
}

type Layout struct {
	Rows []Row `tfsdk:"row"`
}

type Row struct {
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
}

type VariableCustom struct {
	Name        types.String           `tfsdk:"name"`
	Label       types.String           `tfsdk:"label"`
	Description types.String           `tfsdk:"description"`
	Hide        types.String           `tfsdk:"hide"`
	Options     []VariableCustomOption `tfsdk:"option"`
	MultiValue  types.Bool             `tfsdk:"multi_value"`
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
	Name        types.String              `tfsdk:"name"`
	Label       types.String              `tfsdk:"label"`
	Description types.String              `tfsdk:"description"`
	Hide        types.String              `tfsdk:"hide"`
	DataSource  []VariableAdHocDataSource `tfsdk:"datasource"`
}

type VariableAdHocDataSource struct {
	Type types.String `tfsdk:"type"`
	UID  types.String `tfsdk:"uid"`
}

type VariableDataSource struct {
	Name        types.String                 `tfsdk:"name"`
	Label       types.String                 `tfsdk:"label"`
	Description types.String                 `tfsdk:"description"`
	Hide        types.String                 `tfsdk:"hide"`
	MultiValue  types.Bool                   `tfsdk:"multi_value"`
	IncludeAll  []VariableIncludeAll         `tfsdk:"include_all"`
	DataSource  []VariableDataSourceSelector `tfsdk:"source"`
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
	MultiValue  types.Bool            `tfsdk:"multi_value"`
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

func (d *DashboardDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description:         "Dashboard data source.",
		MarkdownDescription: "Dashboard data source. See Grafana [documentation](https://grafana.com/docs/grafana/latest/dashboards/use-dashboards/) for more details.",

		Blocks: map[string]schema.Block{
			"time": dashboardTimeBlock(),
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
									"name": schema.StringAttribute{
										Required:    true,
										Description: "The name of the variable.",
										Validators: []validator.String{
											stringvalidator.LengthAtMost(50),
										},
									},
									"label": schema.StringAttribute{
										Optional:    true,
										Description: "The optional display name.",
									},
									"description": schema.StringAttribute{
										Optional:    true,
										Description: "The description of the variable.",
									},
									"multi_value": schema.BoolAttribute{
										Optional:    true,
										Description: "Whether to allow selecting multiple values at the same time or not.",
									},
									"hide": schema.StringAttribute{
										Optional:            true,
										Description:         "Which variable information to hide. The choices are: label, variable.",
										MarkdownDescription: "Which variable information to hide. The choices are: `label`, `variable`.",
										Validators: []validator.String{
											stringvalidator.OneOf("label", "variable"),
										},
									},
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
									"name": schema.StringAttribute{
										Description: "The name of the variable.",
										Required:    true,
										Validators: []validator.String{
											stringvalidator.LengthAtMost(50),
										},
									},
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
									"name": schema.StringAttribute{
										Description: "The name of the variable.",
										Required:    true,
										Validators: []validator.String{
											stringvalidator.LengthAtMost(50),
										},
									},
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
									"hide": schema.StringAttribute{
										Optional:            true,
										Description:         "Which variable information to hide. The choices are: label, variable.",
										MarkdownDescription: "Which variable information to hide. The choices are: `label`, `variable`.",
										Validators: []validator.String{
											stringvalidator.OneOf("label", "variable"),
										},
									},
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
									"datasource": schema.ListNestedBlock{
										Description: "The datasource to use.",
										NestedObject: schema.NestedBlockObject{
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
										Validators: []validator.List{
											listvalidator.SizeAtLeast(1),
											listvalidator.SizeAtMost(1),
										},
									},
								},

								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "The name of the variable.",
										Required:    true,
										Validators: []validator.String{
											stringvalidator.LengthAtMost(50),
										},
									},
									"label": schema.StringAttribute{
										Optional:    true,
										Description: "The optional display name.",
									},
									"description": schema.StringAttribute{
										Optional:    true,
										Description: "The description of the variable.",
									},
									"hide": schema.StringAttribute{
										Optional:            true,
										Description:         "Which variable information to hide. The choices are: label, variable.",
										MarkdownDescription: "Which variable information to hide. The choices are: `label`, `variable`.",
										Validators: []validator.String{
											stringvalidator.OneOf("label", "variable"),
										},
									},
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
									"source": schema.ListNestedBlock{
										Description: "The datasource selector.",
										NestedObject: schema.NestedBlockObject{
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
										Validators: []validator.List{
											listvalidator.SizeAtLeast(1),
											listvalidator.SizeAtMost(1),
										},
									},
									"include_all": variableIncludeAllBlock(),
								},

								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "The name of the variable.",
										Required:    true,
										Validators: []validator.String{
											stringvalidator.LengthAtMost(50),
										},
									},
									"label": schema.StringAttribute{
										Optional:    true,
										Description: "The optional display name.",
									},
									"description": schema.StringAttribute{
										Optional:    true,
										Description: "The description of the variable.",
									},
									"multi_value": schema.BoolAttribute{
										Optional:    true,
										Description: "Whether to allow selecting multiple values at the same time or not.",
									},
									"hide": schema.StringAttribute{
										Optional:            true,
										Description:         "Which variable information to hide. The choices are: label, variable.",
										MarkdownDescription: "Which variable information to hide. The choices are: `label`, `variable`.",
										Validators: []validator.String{
											stringvalidator.OneOf("label", "variable"),
										},
									},
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
									"name": schema.StringAttribute{
										Description: "The name of the variable.",
										Required:    true,
										Validators: []validator.String{
											stringvalidator.LengthAtMost(50),
										},
									},
									"label": schema.StringAttribute{
										Optional:    true,
										Description: "The optional display name.",
									},
									"description": schema.StringAttribute{
										Optional:    true,
										Description: "The description of the variable.",
									},
									"multi_value": schema.BoolAttribute{
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
									"hide": schema.StringAttribute{
										Optional:            true,
										Description:         "Which variable information to hide. The choices are: label, variable.",
										MarkdownDescription: "Which variable information to hide. The choices are: `label`, `variable`.",
										Validators: []validator.String{
											stringvalidator.OneOf("label", "variable"),
										},
									},
								},

								Validators: []validator.Object{
									objectvalidator.AlsoRequires(path.MatchRelative().AtName("target")),
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
					"row": schema.ListNestedBlock{
						Description: "The row within the dashboard.",
						NestedObject: schema.NestedBlockObject{
							Blocks: map[string]schema.Block{
								"panel": schema.ListNestedBlock{
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
									Validators: []validator.List{
										listvalidator.SizeAtMost(20),
									},
								},
							},
						},
						Validators: []validator.List{
							listvalidator.SizeAtMost(200),
						},
					},
				},
				Validators: []validator.Object{
					objectvalidator.AtLeastOneOf(
						path.MatchRoot("layout"),
					),
				},
			},
		},

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"json": schema.StringAttribute{
				Computed:    true,
				Description: "The Grafana-API-compatible JSON of this panel.",
			},
			"title": schema.StringAttribute{
				Required:    true,
				Description: "The title of the dashboard.",
			},
			"uid": schema.StringAttribute{
				Optional:    true,
				Description: "The UID of the dashboard.",
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

	d.Defaults = defaults.Dashboard
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
			var current grafana.Current

			for i, opt := range custom.Options {
				opts[i] = grafana.Option{
					Text:     opt.Text.ValueString(),
					Value:    opt.Value.ValueString(),
					Selected: opt.Selected.ValueBool(),
				}

				if query != "" {
					query = query + ", "
				}

				query = query + opt.Text.ValueString() + " : " + opt.Value.ValueString()

				if opt.Selected.ValueBool() {
					current = grafana.Current{
						Text: &grafana.StringSliceString{
							Value: []string{opt.Text.ValueString()},
							Valid: true,
						},
						Value: opt.Value.ValueString(),
					}
				}
			}

			hide := uint8(0)

			if !custom.Hide.IsNull() {
				switch v := custom.Hide.ValueString(); v {
				case "label":
					hide = 1
				case "variable":
					hide = 2
				default:
					hide = 0
				}
			}

			v := grafana.TemplateVar{
				Type:        "custom",
				Name:        custom.Name.ValueString(),
				Label:       custom.Label.ValueString(),
				Description: custom.Description.ValueString(),
				Multi:       custom.MultiValue.ValueBool(),
				Options:     opts,
				Query:       query,
				Current:     current,
				Hide:        hide,
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
			}

			vars = append(vars, v)
		}

		for _, textbox := range variable.TextBox {
			hide := uint8(0)

			if !textbox.Hide.IsNull() {
				switch v := textbox.Hide.ValueString(); v {
				case "label":
					hide = 1
				case "variable":
					hide = 2
				default:
					hide = 0
				}
			}

			v := grafana.TemplateVar{
				Type:        "textbox",
				Options:     make([]grafana.Option, 0),
				Name:        textbox.Name.ValueString(),
				Label:       textbox.Label.ValueString(),
				Description: textbox.Description.ValueString(),
				Query:       textbox.DefaultValue.ValueString(),
				Hide:        hide,
			}

			vars = append(vars, v)
		}

		for _, adhoc := range variable.AdHoc {
			hide := uint8(0)

			if !adhoc.Hide.IsNull() {
				switch v := adhoc.Hide.ValueString(); v {
				case "label":
					hide = 1
				case "variable":
					hide = 2
				default:
					hide = 0
				}
			}

			v := grafana.TemplateVar{
				Type:        "adhoc",
				Options:     make([]grafana.Option, 0),
				Name:        adhoc.Name.ValueString(),
				Label:       adhoc.Label.ValueString(),
				Description: adhoc.Description.ValueString(),
				Datasource: &grafana.TemplateVarDataSource{
					UID:  adhoc.DataSource[0].UID.ValueString(),
					Type: adhoc.DataSource[0].Type.ValueString(),
				},
				Hide: hide,
			}

			vars = append(vars, v)
		}

		for _, ds := range variable.DataSource {
			hide := uint8(0)

			if !ds.Hide.IsNull() {
				switch v := ds.Hide.ValueString(); v {
				case "label":
					hide = 1
				case "variable":
					hide = 2
				default:
					hide = 0
				}
			}

			v := grafana.TemplateVar{
				Type:        "datasource",
				Options:     make([]grafana.Option, 0),
				Name:        ds.Name.ValueString(),
				Label:       ds.Label.ValueString(),
				Description: ds.Description.ValueString(),
				Multi:       ds.MultiValue.ValueBool(),
				Query:       ds.DataSource[0].Type.ValueString(),
				Regex:       ds.DataSource[0].Filter.ValueString(),
				Hide:        hide,
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
			hide := uint8(0)

			if !query.Hide.IsNull() {
				switch v := query.Hide.ValueString(); v {
				case "label":
					hide = 1
				case "variable":
					hide = 2
				default:
					hide = 0
				}
			}

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
				Multi:       query.MultiValue.ValueBool(),
				Hide:        hide,
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
	}

	panels := make([]*grafana.Panel, 0)

	for rowIdx, row := range data.Layout.Rows {
		for columnIdx, column := range row.Panels {
			var panel grafana.Panel

			err := json.Unmarshal([]byte(column.Source.ValueString()), &panel)
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Could not unmarshall json as Panel: %s", err))
				return
			}

			height := int(column.Size.Height.ValueInt64())
			width := int(column.Size.Width.ValueInt64())

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
				for _, r := range data.Layout.Rows[0:rowIdx] {
					max := 0
					for _, c := range r.Panels {
						max = int(math.Max(float64(max), float64(c.Size.Height.ValueInt64())))
					}
					total += max
				}
				y = total + rowIdx
			}

			panel.GridPos = struct {
				H *int `json:"h,omitempty"`
				W *int `json:"w,omitempty"`
				X *int `json:"x,omitempty"`
				Y *int `json:"y,omitempty"`
			}{
				H: &height,
				W: &width,
				X: &x,
				Y: &y,
			}

			panels = append(panels, &panel)
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

	if !data.UID.IsNull() {
		dashboard.UID = data.UID.ValueString()
	}

	if !data.Editable.IsNull() {
		dashboard.Editable = data.Editable.ValueBool()
	}

	if !data.Style.IsNull() {
		dashboard.Style = data.Style.ValueString()
	}

	for _, time := range data.Time {
		dashboard.Time.From = time.From.ValueString()
		dashboard.Time.To = time.To.ValueString()
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

	jsonData, err := json.MarshalIndent(dashboard, "", "  ")
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Could not marshal json: %s", err))
		return
	}

	data.Json = types.StringValue(string(jsonData))
	data.Id = types.StringValue(strconv.Itoa(hashcode(jsonData)))

	// resp.Diagnostics.AddError("Client Error", fmt.Sprintf("%s", string(jsonData)))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
