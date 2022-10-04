package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/iRevive/terraform-provider-gdashboard/internal/provider/grafana"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &StatDataSource{}

func NewStatDataSource() datasource.DataSource {
	return &StatDataSource{}
}

// StatDataSource defines the data source implementation.
type StatDataSource struct {
	Defaults StatDefaults
}

type StatDefaults struct {
	Field FieldDefaults
	Graph StatGraphDefault
}

type StatGraphDefault struct {
	Orientation   string
	TextMode      string
	ColorMode     string
	GraphMode     string
	TextAlignment string
	ReduceOptions StatReduceOptionDefault
	TextSize      TextSizeDefault
}

type StatReduceOptionDefault struct {
	Values      bool
	Fields      string
	Limit       int64
	Calculation string
}

type TextSizeDefault struct {
	Title *int
	Value *int
}

// StatDataSourceModel describes the data source data model.
type StatDataSourceModel struct {
	Id      types.String   `tfsdk:"id"`
	Json    types.String   `tfsdk:"json"`
	Title   types.String   `tfsdk:"title"`
	Targets []Target       `tfsdk:"targets"`
	Field   []FieldOptions `tfsdk:"field"`
	Graph   []StatOptions  `tfsdk:"graph"`
}

type StatOptions struct {
	Orientation   types.String        `tfsdk:"orientation"`
	TextMode      types.String        `tfsdk:"text_mode"`
	ColorMode     types.String        `tfsdk:"color_mode"`
	GraphMode     types.String        `tfsdk:"graph_mode"`
	TextAlignment types.String        `tfsdk:"text_alignment"`
	TextSize      []TextSizeOptions   `tfsdk:"text_size"`
	ReduceOptions []StatReduceOptions `tfsdk:"options"`
}

type TextSizeOptions struct {
	Title types.Int64 `tfsdk:"title"`
	Value types.Int64 `tfsdk:"value"`
}

type StatReduceOptions struct {
	Values      types.Bool   `tfsdk:"values"`
	Fields      types.String `tfsdk:"fields"`
	Limit       types.Int64  `tfsdk:"limit"`
	Calculation types.String `tfsdk:"calculation"`
}

func (d *StatDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stat"
}

func statGraphBlock() tfsdk.Block {
	return tfsdk.Block{
		NestingMode: tfsdk.BlockNestingModeList,
		MinItems:    0,
		MaxItems:    1,
		Blocks: map[string]tfsdk.Block{
			"options": {
				NestingMode: tfsdk.BlockNestingModeList,
				MinItems:    0,
				MaxItems:    1,
				Attributes: map[string]tfsdk.Attribute{
					"values": {
						Type:     types.BoolType,
						Optional: true,
					},
					"fields": {
						Type:     types.StringType,
						Optional: true,
					},
					"limit": {
						Type:     types.Int64Type,
						Optional: true,
					},
					"calculation": {
						Type:     types.StringType,
						Optional: true,
					},
				},
			},
			"text_size": {
				NestingMode: tfsdk.BlockNestingModeList,
				MinItems:    0,
				MaxItems:    1,
				Attributes: map[string]tfsdk.Attribute{
					"title": {
						Type:     types.Int64Type,
						Optional: true,
					},
					"value": {
						Type:     types.Int64Type,
						Optional: true,
					},
				},
			},
		},
		Attributes: map[string]tfsdk.Attribute{
			"orientation": {
				Type:     types.StringType,
				Optional: true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf("auto", "horizontal", "vertical"),
				},
			},
			"text_mode": {
				Type:     types.StringType,
				Optional: true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf("auto", "value", "value_and_name", "name", "none"),
				},
			},
			"color_mode": {
				Type:     types.StringType,
				Optional: true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf("none", "value", "background"),
				},
			},
			"graph_mode": {
				Type:     types.StringType,
				Optional: true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf("none", "area"),
				},
			},
			"text_alignment": {
				Type:     types.StringType,
				Optional: true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf("auto", "center"),
				},
			},
		},
	}
}

func (d *StatDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Stat panel data source",

		Blocks: map[string]tfsdk.Block{
			"targets": targetBlock(),
			"field":   fieldBlock(),
			"graph":   statGraphBlock(),
		},

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"json": {
				Type:     types.StringType,
				Computed: true,
			},
			"title": {
				Type:        types.StringType,
				Required:    true,
				Description: "The title of the panel",
			},
		},
	}, nil
}

func (d *StatDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.Defaults = defaults.Stat
}

func (d *StatDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data StatDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	targets := createTargets(data.Targets)
	fieldConfig := createFieldConfig(d.Defaults.Field, data.Field)

	options := grafana.Options{
		Orientation: d.Defaults.Graph.Orientation,
		TextMode:    d.Defaults.Graph.TextMode,
		ColorMode:   d.Defaults.Graph.ColorMode,
		GraphMode:   d.Defaults.Graph.GraphMode,
		JustifyMode: d.Defaults.Graph.TextAlignment,
	}

	options.ReduceOptions.Values = d.Defaults.Graph.ReduceOptions.Values
	options.ReduceOptions.Fields = d.Defaults.Graph.ReduceOptions.Fields
	options.ReduceOptions.Calcs = []string{d.Defaults.Graph.ReduceOptions.Calculation}

	options.Text.TitleSize = d.Defaults.Graph.TextSize.Title
	options.Text.ValueSize = d.Defaults.Graph.TextSize.Value

	for _, graph := range data.Graph {
		if !graph.Orientation.Null {
			options.Orientation = graph.Orientation.Value
		}

		if !graph.TextMode.Null {
			options.TextMode = graph.TextMode.Value
		}

		if !graph.ColorMode.Null {
			options.ColorMode = graph.ColorMode.Value
		}

		if !graph.GraphMode.Null {
			options.GraphMode = graph.GraphMode.Value
		}

		if !graph.TextAlignment.Null {
			options.JustifyMode = graph.TextAlignment.Value
		}

		for _, textSize := range graph.TextSize {
			if !textSize.Title.Null {
				size := int(textSize.Title.Value)
				options.Text.TitleSize = &size
			}

			if !textSize.Value.Null {
				size := int(textSize.Value.Value)
				options.Text.ValueSize = &size
			}
		}

		for _, reducer := range graph.ReduceOptions {
			if !reducer.Values.Null {
				options.ReduceOptions.Values = reducer.Values.Value
			}

			if !reducer.Fields.Null {
				options.ReduceOptions.Fields = reducer.Fields.Value
			}

			if !reducer.Calculation.Null {
				options.ReduceOptions.Calcs = []string{reducer.Calculation.Value}
			}
		}
	}

	panel := &grafana.Panel{
		CommonPanel: grafana.CommonPanel{
			OfType: grafana.StatType,
			Title:  data.Title.Value,
			Type:   "stat",
			Span:   12,
			IsNew:  true,
		},
		StatPanel: &grafana.StatPanel{
			Targets: targets,
			Options: options,
			FieldConfig: grafana.FieldConfig{
				Defaults: fieldConfig,
			},
		},
	}

	jsonData, err := json.MarshalIndent(panel, "", "  ")
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Could not marshal json: %s", err))
		return
	}

	data.Json = types.String{Value: string(jsonData)}
	data.Id = types.String{Value: strconv.Itoa(hashcode(jsonData))}

	//resp.Diagnostics.AddError("Client Error", fmt.Sprintf("%s", string(jsonData)))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
