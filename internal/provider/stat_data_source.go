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
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &StatDataSource{}

func NewStatDataSource() datasource.DataSource {
	return &StatDataSource{}
}

// StatDataSource defines the data source implementation.
type StatDataSource struct {
	CompactJson bool
	Defaults    StatDefaults
}

type StatDefaults struct {
	Field FieldDefaults
	Graph StatGraphDefaults
}

type StatGraphDefaults struct {
	Orientation   string
	TextMode      string
	ColorMode     string
	GraphMode     string
	TextAlignment string
	ReduceOptions ReduceOptionDefaults
	TextSize      TextSizeDefaults
}

// StatDataSourceModel describes the data source data model.
type StatDataSourceModel struct {
	Id              types.String           `tfsdk:"id"`
	Json            types.String           `tfsdk:"json"`
	CompactJson     types.Bool             `tfsdk:"compact_json"`
	Title           types.String           `tfsdk:"title"`
	Description     types.String           `tfsdk:"description"`
	Queries         []Query                `tfsdk:"queries"`
	Field           []FieldOptions         `tfsdk:"field"`
	Graph           []StatOptions          `tfsdk:"graph"`
	Overrides       []FieldOverrideOptions `tfsdk:"overrides"`
	Transformations []Transformations      `tfsdk:"transform"`
}

type StatOptions struct {
	Orientation   types.String      `tfsdk:"orientation"`
	TextMode      types.String      `tfsdk:"text_mode"`
	ColorMode     types.String      `tfsdk:"color_mode"`
	GraphMode     types.String      `tfsdk:"graph_mode"`
	TextAlignment types.String      `tfsdk:"text_alignment"`
	TextSize      []TextSizeOptions `tfsdk:"text_size"`
	ReduceOptions []ReduceOptions   `tfsdk:"options"`
}

func (d *StatDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stat"
}

func statGraphBlock() schema.Block {
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
				"text_mode": schema.StringAttribute{
					Optional:            true,
					Description:         "What show on panel. The choices are: auto, value, value_and_name, name, none.",
					MarkdownDescription: "What show on panel. The choices are: `auto`, `value`, `value_and_name`, `name`, `none`.",
					Validators: []validator.String{
						stringvalidator.OneOf("auto", "value", "value_and_name", "name", "none"),
					},
				},
				"color_mode": schema.StringAttribute{
					Optional:            true,
					Description:         "The color mode. The choices are: none, value, background.",
					MarkdownDescription: "The color mode. The choices are: `none`, `value`, `background`.",
					Validators: []validator.String{
						stringvalidator.OneOf("none", "value", "background"),
					},
				},
				"graph_mode": schema.StringAttribute{
					Optional:            true,
					Description:         "The graph mode. The choices are: none, area.",
					MarkdownDescription: "The graph mode. The choices are: `none`, `area`.",
					Validators: []validator.String{
						stringvalidator.OneOf("none", "area"),
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

func (d *StatDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description:         "Stat panel data source.",
		MarkdownDescription: "Stat panel data source. See Grafana [documentation](https://grafana.com/docs/grafana/latest/panels-visualizations/visualizations/stat/) for more details.",

		Blocks: map[string]schema.Block{
			"queries":   queryBlock(),
			"field":     fieldBlock(false),
			"graph":     statGraphBlock(),
			"overrides": fieldOverrideBlock(false),
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

func (d *StatDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.Defaults = defaults.Stat
}

func (d *StatDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data StatDataSourceModel

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
		TextMode:    d.Defaults.Graph.TextMode,
		ColorMode:   d.Defaults.Graph.ColorMode,
		GraphMode:   d.Defaults.Graph.GraphMode,
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

		if !graph.TextMode.IsNull() {
			options.TextMode = graph.TextMode.ValueString()
		}

		if !graph.ColorMode.IsNull() {
			options.ColorMode = graph.ColorMode.ValueString()
		}

		if !graph.GraphMode.IsNull() {
			options.GraphMode = graph.GraphMode.ValueString()
		}

		if !graph.TextAlignment.IsNull() {
			options.JustifyMode = graph.TextAlignment.ValueString()
		}

		updateTextSize(&options.TextSize, graph.TextSize)
		updateReduceOptions(&options.ReduceOptions, graph.ReduceOptions)
	}

	panel := &grafana.Panel{
		CommonPanel: grafana.CommonPanel{
			OfType:          grafana.StatType,
			Title:           data.Title.ValueString(),
			Type:            "stat",
			Span:            12,
			IsNew:           true,
			Transformations: transformations,
		},
		StatPanel: &grafana.StatPanel{
			Targets: targets,
			Options: options,
			FieldConfig: grafana.FieldConfig{
				Defaults:  fieldConfig,
				Overrides: createOverrides(data.Overrides),
			},
		},
	}

	if !data.Description.IsNull() {
		panel.CommonPanel.Description = data.Description.ValueStringPointer()
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
