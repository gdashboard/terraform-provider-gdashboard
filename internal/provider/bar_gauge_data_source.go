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
var _ datasource.DataSource = &BarGaugeDataSource{}

func NewBarGaugeDataSource() datasource.DataSource {
	return &BarGaugeDataSource{}
}

// BarGaugeDataSource defines the data source implementation.
type BarGaugeDataSource struct {
	Defaults BarGaugeDefaults
}

type BarGaugeDefaults struct {
	Field FieldDefaults
	Graph BarGaugeGraphDefault
}

type BarGaugeGraphDefault struct {
	Orientation   string
	DisplayMode   string
	TextAlignment string
	TextSize      TextSizeDefaults
	ReduceOptions ReduceOptionDefaults
}

// BarGaugeDataSourceModel describes the data source data model.
type BarGaugeDataSourceModel struct {
	Id      types.String      `tfsdk:"id"`
	Json    types.String      `tfsdk:"json"`
	Title   types.String      `tfsdk:"title"`
	Targets []Target          `tfsdk:"targets"`
	Field   []FieldOptions    `tfsdk:"field"`
	Graph   []BarGaugeOptions `tfsdk:"graph"`
}

type BarGaugeOptions struct {
	Orientation   types.String      `tfsdk:"orientation"`
	DisplayMode   types.String      `tfsdk:"display_mode"`
	TextAlignment types.String      `tfsdk:"text_alignment"`
	TextSize      []TextSizeOptions `tfsdk:"text_size"`
	ReduceOptions []ReduceOptions   `tfsdk:"options"`
}

func (d *BarGaugeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bar_gauge"
}

func barGaugeGraphBlock() tfsdk.Block {
	return tfsdk.Block{
		NestingMode: tfsdk.BlockNestingModeList,
		MinItems:    0,
		MaxItems:    1,
		Blocks: map[string]tfsdk.Block{
			"options":   reduceOptionsBlock(),
			"text_size": textSizeBlock(),
		},
		Attributes: map[string]tfsdk.Attribute{
			"orientation": {
				Type:     types.StringType,
				Optional: true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf("auto", "horizontal", "vertical"),
				},
			},
			"display_mode": {
				Type:     types.StringType,
				Optional: true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf("gradient", "lcd", "basic"),
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

func (d *BarGaugeDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Bar Gauge panel data source",

		Blocks: map[string]tfsdk.Block{
			"targets": targetBlock(),
			"field":   fieldBlock(),
			"graph":   barGaugeGraphBlock(),
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

func (d *BarGaugeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.Defaults = defaults.BarGauge
}

func (d *BarGaugeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data BarGaugeDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	targets := createTargets(data.Targets)
	fieldConfig := createFieldConfig(d.Defaults.Field, data.Field)

	options := grafana.Options{
		Orientation: d.Defaults.Graph.Orientation,
		DisplayMode: d.Defaults.Graph.DisplayMode,
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

		if !graph.DisplayMode.Null {
			options.DisplayMode = graph.DisplayMode.Value
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
			OfType: grafana.BarGaugeType,
			Title:  data.Title.Value,
			Type:   "bargauge",
			Span:   12,
			IsNew:  true,
		},
		BarGaugePanel: &grafana.BarGaugePanel{
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
