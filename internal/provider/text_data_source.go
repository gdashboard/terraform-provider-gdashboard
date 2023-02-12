package provider

import (
	"context"
	"encoding/json"
	"fmt"
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
var _ datasource.DataSource = &TextDataSource{}

func NewTextDataSource() datasource.DataSource {
	return &TextDataSource{}
}

// TextDataSource defines the data source implementation.
type TextDataSource struct {
}

// TextDataSourceModel describes the data source data model.
type TextDataSourceModel struct {
	Id          types.String  `tfsdk:"id"`
	Json        types.String  `tfsdk:"json"`
	Title       types.String  `tfsdk:"title"`
	Description types.String  `tfsdk:"description"`
	Graph       []TextOptions `tfsdk:"graph"`
}

type CodeOptions struct {
	Language        types.String `tfsdk:"language"`
	ShowLineNumbers types.Bool   `tfsdk:"show_line_numbers"`
	ShowMiniMap     types.Bool   `tfsdk:"show_mini_map"`
}

type TextOptions struct {
	Mode        types.String  `tfsdk:"mode"`
	Content     types.String  `tfsdk:"content"`
	CodeOptions []CodeOptions `tfsdk:"code"`
}

func (d *TextDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_text"
}

func textGraphBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: "The visualization options.",
		NestedObject: schema.NestedBlockObject{
			Blocks: map[string]schema.Block{
				"code": schema.ListNestedBlock{
					Description:         "The configuration for the 'code' mode.",
					MarkdownDescription: "The configuration for the `code` mode.",
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"language": schema.StringAttribute{
								Optional:            true,
								Description:         "The syntax highlighting language. The choices are: go, html, json, markdown, plaintext, sql, typescript, xml, yaml.",
								MarkdownDescription: "The syntax highlighting language. The choices are: `go`, `html`, `json`, `markdown`, `plaintext`, `sql`, `typescript`, `xml`, `yaml`",
								Validators: []validator.String{
									stringvalidator.OneOf("go", "html", "json", "markdown", "plaintext", "sql", "typescript", "xml", "yaml"),
								},
							},
							"show_line_numbers": schema.BoolAttribute{
								Optional:    true,
								Description: "Whether to show line numbers or not.",
							},
							"show_mini_map": schema.BoolAttribute{
								Optional:    true,
								Description: "Whether to show a VSCode-like code-navigation mini map or not.",
							},
						},
						/*Validators: []validator.Object{
							objectvalidator.ConflictsWith(
								path.MatchRoot("graph").AtAnyListIndex().AtName("mode").AtSetValue(types.StringValue("markdown")),
							),
						},*/
					},
				},
			},
			Attributes: map[string]schema.Attribute{
				"mode": schema.StringAttribute{
					Optional:            true,
					Description:         "The display mode. The choices are: markdown, html, code.",
					MarkdownDescription: "The display mode. The choices are: `markdown`, `html`, `code`.",
					Validators: []validator.String{
						stringvalidator.OneOf("markdown", "html", "code"),
					},
				},
				"content": schema.StringAttribute{
					Required:            true,
					Description:         "The content to display.",
					MarkdownDescription: "The content to display.",
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func (d *TextDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description:         "Text panel data source.",
		MarkdownDescription: "Text panel data source. See Grafana [documentation](https://grafana.com/docs/grafana/latest/panels-visualizations/visualizations/text/) for more details.",

		Blocks: map[string]schema.Block{
			"graph": textGraphBlock(),
		},

		Attributes: map[string]schema.Attribute{
			"id":          idAttribute(),
			"json":        jsonAttribute(),
			"title":       titleAttribute(),
			"description": descriptionAttribute(),
		},
	}
}

func (d *TextDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
}

func (d *TextDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TextDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	mode := "markdown"
	content := ""

	codeOptions := grafana.CodeOptions{
		Language:        "plaintext",
		ShowLineNumbers: false,
		ShowMiniMap:     false,
	}

	for _, graph := range data.Graph {
		if !graph.Mode.IsNull() {
			mode = graph.Mode.ValueString()
		}

		if !graph.Content.IsNull() {
			content = graph.Content.ValueString()
		}

		for _, options := range graph.CodeOptions {
			if !options.Language.IsNull() {
				codeOptions.Language = options.Language.ValueString()
			}

			if !options.ShowLineNumbers.IsNull() {
				codeOptions.ShowLineNumbers = options.ShowLineNumbers.ValueBool()
			}

			if !options.ShowMiniMap.IsNull() {
				codeOptions.ShowMiniMap = options.ShowMiniMap.ValueBool()
			}
		}
	}

	panel := &grafana.Panel{
		CommonPanel: grafana.CommonPanel{
			OfType: grafana.TextType,
			Title:  data.Title.ValueString(),
			Type:   "text",
			Span:   12,
			IsNew:  true,
		},
		TextPanel: &grafana.TextPanel{
			Options: grafana.TextPanelOptions{
				Mode:    mode,
				Content: content,
				Code:    codeOptions,
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
