package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/iRevive/terraform-provider-gdashboard/internal/provider/grafana"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &TableDataSource{}

func NewTableDataSource() datasource.DataSource {
	return &TableDataSource{}
}

// TableDataSource defines the data source implementation.
type TableDataSource struct {
	Defaults TableDefaults
}

type TableDefaults struct {
	Field FieldDefaults
}

// TableDataSourceModel describes the data source data model.
type TableDataSourceModel struct {
	Id              types.String           `tfsdk:"id"`
	Json            types.String           `tfsdk:"json"`
	Title           types.String           `tfsdk:"title"`
	Description     types.String           `tfsdk:"description"`
	Queries         []Query                `tfsdk:"queries"`
	Field           []FieldOptions         `tfsdk:"field"`
	Graph           []TableOptions         `tfsdk:"graph"`
	Overrides       []FieldOverrideOptions `tfsdk:"overrides"`
	Transformations []Transformations      `tfsdk:"transform"`
}

type TableColumnOptions struct {
	Align      types.String `tfsdk:"align"`
	Filterable types.Bool   `tfsdk:"filterable"`
	Width      types.Int64  `tfsdk:"width"`
	MinWidth   types.Int64  `tfsdk:"min_width"`
}

type TableCellOptions struct {
	Inspectable types.Bool   `tfsdk:"inspectable"`
	DisplayMode types.String `tfsdk:"display_mode"`
}

type TableFooterOptions struct {
	Pagination   types.Bool     `tfsdk:"pagination"`
	Fields       []types.String `tfsdk:"fields"`
	Calculations []types.String `tfsdk:"calculations"`
}

type TableOptions struct {
	Column     []TableColumnOptions `tfsdk:"column"`
	Cell       []TableCellOptions   `tfsdk:"cell"`
	Footer     []TableFooterOptions `tfsdk:"footer"`
	ShowHeader types.Bool           `tfsdk:"show_header"`
}

func (d *TableDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_table"
}

func footerOptionsBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: "Table footer options.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"pagination": schema.BoolAttribute{
					Optional:    true,
					Description: "Whether to enable pagination or not.",
				},
				"fields": schema.ListAttribute{
					ElementType: types.StringType,
					Optional:    true,
					Description: "Choose the fields should appear in calculations.",
				},
				"calculations": schema.ListAttribute{
					ElementType:         types.StringType,
					Optional:            true,
					Description:         "A reducer function or calculation. The choices are: " + CalculationTypesString() + ".",
					MarkdownDescription: "A reducer function or calculation. The choices are: " + CalculationTypesMarkdown() + ".",
					Validators: []validator.List{
						listvalidator.ValueStringsAre(stringvalidator.OneOf(CalculationTypes()...)),
					},
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func columnOptionsBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: "Table column options.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"align": schema.StringAttribute{
					Optional:            true,
					Description:         "The alignment of cell content: auto, center, left, right.",
					MarkdownDescription: "The alignment of cell content: `auto`, `center`, `left`, `right`.",
					Validators: []validator.String{
						stringvalidator.OneOf("auto", "center", "left", "right"),
					},
				},
				"filterable": schema.BoolAttribute{
					Optional:    true,
					Description: "Whether to make table filterable or not.",
				},
				"width": schema.Int64Attribute{
					Optional:    true,
					Description: "The width for columns in pixels.",
					Validators: []validator.Int64{
						int64validator.AtLeast(20),
						int64validator.AtMost(300),
					},
				},
				"min_width": schema.Int64Attribute{
					Optional:    true,
					Description: "The minimum width for columns in pixels for auto resizing.",
					Validators: []validator.Int64{
						int64validator.AtLeast(50),
						int64validator.AtMost(500),
					},
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func cellOptionsBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: "Table column options.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"inspectable": schema.BoolAttribute{
					Optional:    true,
					Description: "Whether to make cells inspectable or not.",
				},
				"display_mode": schema.StringAttribute{
					Optional:            true,
					Description:         "The alignment of cell content: auto, color-text, color-background, color-background-solid, gradient-gauge, lcd-gauge, basic, json-view, image.",
					MarkdownDescription: "The alignment of cell content: `auto`, `color-text`, `color-background`, `color-background-solid`, `gradient-gauge`, `lcd-gauge`, `basic`, `json-view`, `image`.",
					Validators: []validator.String{
						stringvalidator.OneOf(
							"auto", "color-text", "color-background", "color-background-solid",
							"gradient-gauge", "lcd-gauge", "basic", "json-view", "image",
						),
					},
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func tableGraphBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: "The visualization options.",
		NestedObject: schema.NestedBlockObject{
			Blocks: map[string]schema.Block{
				"column": columnOptionsBlock(),
				"cell":   cellOptionsBlock(),
				"footer": footerOptionsBlock(),
			},
			Attributes: map[string]schema.Attribute{
				"show_header": schema.BoolAttribute{
					Optional:    true,
					Description: "Whether to show table header or not.",
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func (d *TableDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description:         "Table panel data source.",
		MarkdownDescription: "Table panel data source. See Grafana [documentation](https://grafana.com/docs/grafana/latest/panels-visualizations/visualizations/table/) for more details.",

		Blocks: map[string]schema.Block{
			"queries":   queryBlock(),
			"field":     fieldBlock(),
			"graph":     tableGraphBlock(),
			"overrides": fieldOverrideBlock(),
			"transform": transformationsBlock(),
		},

		Attributes: map[string]schema.Attribute{
			"id":          idAttribute(),
			"json":        jsonAttribute(),
			"title":       titleAttribute(),
			"description": descriptionAttribute(),
		},
	}
}

func (d *TableDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.Defaults = defaults.Table
}

func (d *TableDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TableDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	targets := createTargets(data.Queries)
	fieldConfig := createFieldConfig(d.Defaults.Field, data.Field)
	transformations := createTransformations(data.Transformations)

	options := grafana.TableOptions{
		ShowHeader: true,
		Footer: grafana.TableFooter{
			Show:             false,
			EnablePagination: false,
			Fields:           nil,
			Reducer:          nil,
		},
	}

	for _, graph := range data.Graph {
		for _, cell := range graph.Cell {
			if !cell.Inspectable.IsNull() {
				fieldConfig.Custom.Inspect = cell.Inspectable.ValueBool()
			}

			if !cell.DisplayMode.IsNull() {
				fieldConfig.Custom.DisplayMode = cell.DisplayMode.ValueString()
			}
		}

		for _, column := range graph.Column {
			if !column.Align.IsNull() {
				fieldConfig.Custom.Align = column.Align.ValueString()
			}

			if !column.Filterable.IsNull() {
				fieldConfig.Custom.Filterable = column.Filterable.ValueBool()
			}

			if !column.Width.IsNull() {
				fieldConfig.Custom.Width = column.Width.ValueInt64()
			}

			if !column.MinWidth.IsNull() {
				fieldConfig.Custom.MinWidth = column.MinWidth.ValueInt64()
			}
		}

		for _, footer := range graph.Footer {
			if !footer.Pagination.IsNull() {
				options.Footer.EnablePagination = footer.Pagination.ValueBool()
			}

			if len(footer.Fields) > 0 {
				fields := make([]string, 0)

				for _, field := range footer.Fields {
					if !field.IsNull() {
						fields = append(fields, field.ValueString())
					}
				}

				options.Footer.Fields = fields
			}

			if len(footer.Calculations) > 0 {
				reducer := make([]string, 0)

				for _, field := range footer.Fields {
					if !field.IsNull() {
						reducer = append(reducer, field.ValueString())
					}
				}

				options.Footer.Reducer = reducer
			}

			options.Footer.Show = options.Footer.EnablePagination || len(options.Footer.Reducer) > 0 || len(options.Footer.Fields) > 0
		}

		if !graph.ShowHeader.IsNull() {
			options.ShowHeader = graph.ShowHeader.ValueBool()
		}
	}

	panel := &grafana.Panel{
		CommonPanel: grafana.CommonPanel{
			OfType:          grafana.TableType,
			Title:           data.Title.ValueString(),
			Type:            "table",
			Span:            12,
			IsNew:           true,
			Transformations: transformations,
		},
		TablePanel: &grafana.TablePanel{
			Targets: targets,
			Options: options,
			FieldConfig: grafana.FieldConfig{
				Defaults:  fieldConfig,
				Overrides: createOverrides(data.Overrides),
			},
		},
	}

	if !data.Description.IsNull() {
		description := data.Description.ValueString()
		panel.CommonPanel.Description = &description
	}

	jsonData, err := json.MarshalIndent(panel, "", "  ")
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Could not marshal json: %s", err))
		return
	}

	data.Json = types.StringValue(string(jsonData))
	data.Id = types.StringValue(strconv.Itoa(hashcode(jsonData)))

	// resp.Diagnostics.AddError("Client Error", fmt.Sprintf("%s", string(jsonData)))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
