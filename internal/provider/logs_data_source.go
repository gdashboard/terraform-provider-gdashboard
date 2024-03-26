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
var _ datasource.DataSource = &LogsDataSource{}

func NewLogsDataSource() datasource.DataSource {
	return &LogsDataSource{}
}

// LogsDataSource defines the data source implementation.
type LogsDataSource struct {
	CompactJson bool
}

// LogsDataSourceModel describes the data source data model.
type LogsDataSourceModel struct {
	Id          types.String  `tfsdk:"id"`
	Json        types.String  `tfsdk:"json"`
	CompactJson types.Bool    `tfsdk:"compact_json"`
	Title       types.String  `tfsdk:"title"`
	Description types.String  `tfsdk:"description"`
	Queries     []Query       `tfsdk:"queries"`
	Graph       []LogsOptions `tfsdk:"graph"`
}

type LogsOptions struct {
	ShowTime         types.Bool   `tfsdk:"show_time"`
	ShowUniqueLabels types.Bool   `tfsdk:"show_unique_labels"`
	ShowCommonLabels types.Bool   `tfsdk:"show_common_labels"`
	WrapLines        types.Bool   `tfsdk:"wrap_lines"`
	PrettifyJson     types.Bool   `tfsdk:"prettify_json"`
	EnableLogDetails types.Bool   `tfsdk:"enable_log_details"`
	Deduplication    types.String `tfsdk:"deduplication"`
	Order            types.String `tfsdk:"order"`
}

func (d *LogsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_logs"
}

func (d *LogsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description:         "Logs panel data source.",
		MarkdownDescription: "Logs panel data source. See Grafana [documentation](https://grafana.com/docs/grafana/latest/panels-visualizations/visualizations/logs/) for more details.",

		Blocks: map[string]schema.Block{
			"queries": queryBlock(),
			"graph": schema.ListNestedBlock{
				Description: "The visualization options.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"show_time": schema.BoolAttribute{
							Optional:    true,
							Description: "Whether to show the time column. This is the timestamp associated with the log line as reported from the data source.",
						},
						"show_unique_labels": schema.BoolAttribute{
							Optional:    true,
							Description: "Whether to show the unique labels column, which shows only non-common labels.",
						},
						"show_common_labels": schema.BoolAttribute{
							Optional:    true,
							Description: "Whether to show the common labels.",
						},
						"wrap_lines": schema.BoolAttribute{
							Optional:    true,
							Description: "Whether to wrap the lines.",
						},
						"prettify_json": schema.BoolAttribute{
							Optional:    true,
							Description: "Set this to true to pretty print all JSON logs. This setting does not affect logs in any format other than JSON.",
						},
						"enable_log_details": schema.BoolAttribute{
							Optional:    true,
							Description: "Whether to show log details view for each log row.",
						},
						"deduplication": schema.StringAttribute{
							Optional:            true,
							Description:         "The layout orientation. The choices are: none, exact, numbers, signature.",
							MarkdownDescription: "The layout orientation. The choices are: `none`, `exact`, `numbers`, `signature`.",
							Validators: []validator.String{
								stringvalidator.OneOf("none", "exact", "numbers", "signature"),
							},
						},
						"order": schema.StringAttribute{
							Optional:            true,
							Description:         "The order in which to show logs first. The choices are: newest_first, oldest_first.",
							MarkdownDescription: "The order in which to show logs first. The choices are: `newest_first`, `oldest_first`.",
							Validators: []validator.String{
								stringvalidator.OneOf("newest_first", "oldest_first"),
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
			"id":           idAttribute(),
			"json":         jsonAttribute(),
			"compact_json": compactJsonAttribute(),
			"title":        titleAttribute(),
			"description":  descriptionAttribute(),
		},
	}
}

func (d *LogsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
}

func (d *LogsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data LogsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	targets, minInterval := createTargets(data.Queries)

	falseVal := false
	trueVal := true
	dedupStrategy := "none"
	sortOrderAscending := "Ascending"
	sortOrderDescending := "Descending"

	options := grafana.LogsOptions{
		ShowTime:           &falseVal,
		ShowLabels:         &falseVal,
		ShowCommonLabels:   &falseVal,
		WrapLogMessage:     &falseVal,
		PrettifyLogMessage: &falseVal,
		EnableLogDetails:   &trueVal,
		DedupStrategy:      &dedupStrategy,
		SortOrder:          &sortOrderDescending,
	}

	for _, graph := range data.Graph {
		if !graph.ShowTime.IsNull() {
			options.ShowTime = graph.ShowTime.ValueBoolPointer()
		}

		if !graph.ShowUniqueLabels.IsNull() {
			options.ShowLabels = graph.ShowUniqueLabels.ValueBoolPointer()
		}

		if !graph.ShowCommonLabels.IsNull() {
			options.ShowCommonLabels = graph.ShowCommonLabels.ValueBoolPointer()
		}

		if !graph.WrapLines.IsNull() {
			options.WrapLogMessage = graph.WrapLines.ValueBoolPointer()
		}

		if !graph.PrettifyJson.IsNull() {
			options.PrettifyLogMessage = graph.PrettifyJson.ValueBoolPointer()
		}

		if !graph.EnableLogDetails.IsNull() {
			options.EnableLogDetails = graph.EnableLogDetails.ValueBoolPointer()
		}

		if !graph.Deduplication.IsNull() {
			options.DedupStrategy = graph.Deduplication.ValueStringPointer()
		}

		if !graph.Order.IsNull() {
			if graph.Order.ValueString() == "oldest_first" {
				options.SortOrder = &sortOrderAscending
			} else {
				options.SortOrder = &sortOrderDescending
			}
		}
	}

	panel := &grafana.Panel{
		CommonPanel: grafana.CommonPanel{
			OfType:   grafana.LogsType,
			Title:    data.Title.ValueString(),
			Type:     "logs",
			Span:     12,
			IsNew:    true,
			Interval: minInterval,
		},
		LogsPanel: &grafana.LogsPanel{
			Targets: targets,
			Options: options,
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
