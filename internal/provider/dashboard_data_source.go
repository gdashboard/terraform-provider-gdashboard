package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/schemavalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/iRevive/terraform-provider-gdashboard/internal/provider/grafana"
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

func (d *DashboardDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dashboard"
}

func dashboardTimeBlock() tfsdk.Block {
	return tfsdk.Block{
		NestingMode: tfsdk.BlockNestingModeList,
		MinItems:    0,
		MaxItems:    1,
		Description: "The default query time range.",
		Attributes: map[string]tfsdk.Attribute{
			"from": {
				Type:        types.StringType,
				Required:    true,
				Description: "The start of the range.",
			},
			"to": {
				Type:        types.StringType,
				Required:    true,
				Description: "The end of the range.",
			},
		},
	}
}

func dashboardEditableAttribute() tfsdk.Attribute {
	return tfsdk.Attribute{
		Type:        types.BoolType,
		Optional:    true,
		Description: "Whether to make the dashboard editable or not.",
	}
}

func dashboardStyleAttribute() tfsdk.Attribute {
	return tfsdk.Attribute{
		Type:                types.StringType,
		Optional:            true,
		Description:         "The dashboard style. The choices are: dark, light.",
		MarkdownDescription: "The dashboard style. The choices are: `dark`, `light`.",
		Validators: []tfsdk.AttributeValidator{
			stringvalidator.OneOf("dark", "light"),
		},
	}
}

func dashboardGraphTooltipAttribute() tfsdk.Attribute {
	return tfsdk.Attribute{
		Type:                types.StringType,
		Optional:            true,
		Description:         "Controls tooltip and hover highlight behavior across different panels: default, shared-crosshair, shared-tooltip.",
		MarkdownDescription: "Controls tooltip and hover highlight behavior across different panels: `default`, `shared-crosshair`, `shared-tooltip`.",
		Validators: []tfsdk.AttributeValidator{
			stringvalidator.OneOf("default", "shared-crosshair", "shared-tooltip"),
		},
	}
}

func (d *DashboardDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Dashboard data source",

		Blocks: map[string]tfsdk.Block{
			"time": dashboardTimeBlock(),
			"variables": {
				NestingMode: tfsdk.BlockNestingModeList,
				MaxItems:    10,
				Description: "The variables.",
				Blocks: map[string]tfsdk.Block{
					"custom": {
						NestingMode: tfsdk.BlockNestingModeList,
						MinItems:    1,
						MaxItems:    10,
						Description: "The variable options defined as a comma-separated list.",
						Blocks: map[string]tfsdk.Block{
							"option": {
								NestingMode: tfsdk.BlockNestingModeList,
								MinItems:    1,
								MaxItems:    10,
								Description: "The option entry.",
								Attributes: map[string]tfsdk.Attribute{
									"text": {
										Type:        types.StringType,
										Required:    true,
										Description: "The text (label) of the entry.",
									},
									"value": {
										Type:        types.StringType,
										Required:    true,
										Description: "The value of the entry.",
									},
									"selected": {
										Type:        types.BoolType,
										Optional:    true,
										Description: "Whether to mark the option as selected or not.",
									},
								},
							},
						},
						Attributes: map[string]tfsdk.Attribute{
							"name": {
								Type:        types.StringType,
								Required:    true,
								Description: "The name of the variable.",
							},
						},
					},
					"const": {
						NestingMode: tfsdk.BlockNestingModeList,
						MaxItems:    5,
						Description: "The constant variable.",
						Attributes: map[string]tfsdk.Attribute{
							"name": {
								Type:        types.StringType,
								Description: "The name of the variable.",
								Required:    true,
							},
							"value": {
								Type:        types.StringType,
								Description: "The value of the variable.",
								Required:    true,
							},
							"hide": {
								Type:                types.StringType,
								Optional:            true,
								Description:         "Which variable information to hide. The choices are: label, variable.",
								MarkdownDescription: "Which variable information to hide. The choices are: `label`, `variable`.",
								Validators: []tfsdk.AttributeValidator{
									stringvalidator.OneOf("label", "variable"),
								},
							},
						},
					},
				},
			},
			"layout": {
				NestingMode: tfsdk.BlockNestingModeSingle,
				Description: "The layout of the dashboard.",
				Blocks: map[string]tfsdk.Block{
					"row": {
						NestingMode: tfsdk.BlockNestingModeList,
						MaxItems:    40,
						Description: "The row within the dashboard.",
						Blocks: map[string]tfsdk.Block{
							"panel": {
								MaxItems:    20,
								NestingMode: tfsdk.BlockNestingModeList,
								Description: "The definition of the panel within the row.",
								Attributes: map[string]tfsdk.Attribute{
									"size": {
										Description: "The size of the panel.",
										Required:    true,
										Attributes: tfsdk.SingleNestedAttributes(
											map[string]tfsdk.Attribute{
												"height": {
													Type:        types.Int64Type,
													Required:    true,
													Description: "The height of the panel.",
												},
												"width": {
													Type:        types.Int64Type,
													Required:    true,
													Description: "The width of the panel.",
												},
											},
										),
									},
									"source": {
										Type:        types.StringType,
										Description: "The JSON source of the panel.",
										Required:    true,
									},
								},
							},
						},
					},
				},
				Validators: []tfsdk.AttributeValidator{
					schemavalidator.AtLeastOneOf(
						path.MatchRoot("layout"),
					),
				},
			},
		},

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"json": {
				Type:        types.StringType,
				Computed:    true,
				Description: "The Grafana-API-compatible JSON of this panel.",
			},
			"title": {
				Type:        types.StringType,
				Required:    true,
				Description: "The title of the dashboard.",
			},
			"uid": {
				Type:        types.StringType,
				Optional:    true,
				Description: "The UID of the dashboard.",
			},
			"editable":      dashboardEditableAttribute(),
			"style":         dashboardStyleAttribute(),
			"graph_tooltip": dashboardGraphTooltipAttribute(),
		},
	}, nil
}

func (d *DashboardDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
					Text:     opt.Text.Value,
					Value:    opt.Value.Value,
					Selected: opt.Selected.Value,
				}

				if query != "" {
					query = query + ", "
				}

				query = query + opt.Text.Value + " : " + opt.Value.Value

				if opt.Selected.Value {
					current = grafana.Current{
						Text: &grafana.StringSliceString{
							Value: []string{opt.Text.Value},
							Valid: true,
						},
						Value: opt.Value.Value,
					}
				}
			}

			v := grafana.TemplateVar{
				Type:    "custom",
				Name:    custom.Name.Value,
				Options: opts,
				Query:   query,
				Current: current,
			}

			vars = append(vars, v)
		}

		for _, c := range variable.Constant {
			hide := uint8(0)

			if !c.Hide.Null {
				switch v := c.Hide.Value; v {
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
				Name:  c.Name.Value,
				Query: c.Value.Value,
				Hide:  hide,
			}

			vars = append(vars, v)
		}
	}

	panels := make([]*grafana.Panel, 0)

	for rowIdx, row := range data.Layout.Rows {
		for columnIdx, column := range row.Panels {
			var panel grafana.Panel

			err := json.Unmarshal([]byte(column.Source.Value), &panel)
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Could not unmarshall json as Panel: %s", err))
				return
			}

			height := int(column.Size.Height.Value)
			width := int(column.Size.Width.Value)

			var x int
			var y int

			if columnIdx == 0 {
				x = 0
			} else {
				total := 0
				for _, item := range row.Panels[0:columnIdx] {
					total += int(item.Size.Width.Value)
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
						max = int(math.Max(float64(max), float64(c.Size.Height.Value)))
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
		Title:         data.Title.Value,
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

	if !data.UID.Null {
		dashboard.UID = data.UID.Value
	}

	if !data.Editable.Null {
		dashboard.Editable = data.Editable.Value
	}

	if !data.Style.Null {
		dashboard.Style = data.Style.Value
	}

	for _, time := range data.Time {
		dashboard.Time.From = time.From.Value
		dashboard.Time.To = time.To.Value
	}

	tooltip := ""
	if !data.GraphTooltip.Null {
		tooltip = data.GraphTooltip.Value
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

	data.Json = types.String{Value: string(jsonData)}
	data.Id = types.String{Value: strconv.Itoa(hashcode(jsonData))}

	//resp.Diagnostics.AddError("Client Error", fmt.Sprintf("%s", string(jsonData)))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
