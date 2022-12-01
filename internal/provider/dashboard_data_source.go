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
	Custom   []VariableCustom   `tfsdk:"custom"`
	Constant []VariableConstant `tfsdk:"const"`
}

type VariableCustom struct {
	Name    types.String           `tfsdk:"name"`
	Hide    types.String           `tfsdk:"hide"`
	Options []VariableCustomOption `tfsdk:"option"`
}

type VariableCustomOption struct {
	Selected types.Bool   `tfsdk:"selected"`
	Text     types.String `tfsdk:"text"`
	Value    types.String `tfsdk:"value"`
}

type VariableConstant struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
	Hide  types.String `tfsdk:"hide"`
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
								},
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Required:    true,
										Description: "The name of the variable.",
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
								listvalidator.SizeAtMost(10),
							},
						},
						"const": schema.ListNestedBlock{
							Description: "The constant variable.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "The name of the variable.",
										Required:    true,
									},
									"value": schema.StringAttribute{
										Description: "The value of the variable.",
										Required:    true,
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
								listvalidator.SizeAtMost(5),
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
							listvalidator.SizeAtMost(40),
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
				Type:    "custom",
				Name:    custom.Name.ValueString(),
				Options: opts,
				Query:   query,
				Current: current,
				Hide:    hide,
			}

			vars = append(vars, v)
		}

		for _, c := range variable.Constant {
			hide := uint8(0)

			if !c.Hide.IsNull() {
				switch v := c.Hide.ValueString(); v {
				case "label":
					hide = 1
				case "variable":
					hide = 2
				default:
					hide = 0
				}
			}

			v := grafana.TemplateVar{
				Type:  "constant",
				Name:  c.Name.ValueString(),
				Query: c.Value.ValueString(),
				Hide:  hide,
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

	//resp.Diagnostics.AddError("Client Error", fmt.Sprintf("%s", string(jsonData)))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
