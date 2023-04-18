package provider

import (
	"github.com/gdashboard/terraform-provider-gdashboard/internal/provider/grafana"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Transformations struct {
	Steps []TransformationsStep `tfsdk:"step"`
}

func transformationsBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: "The ",
		NestedObject: schema.NestedBlockObject{
			Blocks: map[string]schema.Block{
				"step": transformationsStepBlock(),
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

type TransformationsStep struct {
	FilterByName     []TransformationFilterFieldsByName `tfsdk:"filter_fields_by_name"`
	GroupBy          []TransformationGroupBy            `tfsdk:"group_by"`
	GroupingToMatrix []TransformationGroupingToMatrix   `tfsdk:"grouping_to_matrix"`
	Limit            []TransformationLimit              `tfsdk:"limit"`
	SeriesToRows     []TransformationSeriesToRows       `tfsdk:"series_to_rows"`
	SortBy           []TransformationSortBy             `tfsdk:"sort_by"`
}

func transformationsStepBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: "The transform step.",
		NestedObject: schema.NestedBlockObject{
			Blocks: map[string]schema.Block{
				"filter_fields_by_name": transformationFilterFieldsByNameBlock(),
				"group_by":              transformationGroupBy(),
				"grouping_to_matrix":    transformationGroupingToMatrixBlock(),
				"limit":                 transformationLimitBlock(),
				"series_to_rows":        transformationSeriesToRowsBlock(),
				"sort_by":               transformationSortByBlock(),
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(20),
		},
	}
}

type TransformationFilterFieldsByName struct {
	Names []types.String `tfsdk:"names"`
}

func transformationFilterFieldsByNameBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: "Remove portions of the query results.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"names": schema.ListAttribute{
					ElementType: types.StringType,
					Required:    true,
					Description: "The fields to keep.",
					Validators: []validator.List{
						listvalidator.SizeAtLeast(1),
					},
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("group_by")),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("grouping_to_matrix")),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("limit")),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("sort_by")),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("series_to_rows")),
		},
	}
}

type TransformationGroupBy struct {
	GroupBy   []types.String            `tfsdk:"by"`
	Aggregate map[string][]types.String `tfsdk:"aggregate"`
}

func transformationGroupBy() schema.Block {
	return schema.ListNestedBlock{
		Description: "Group the data by a specified field (column) value and processes calculations on each group.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"by": schema.ListAttribute{
					ElementType: types.StringType,
					Required:    true,
					Description: "Fields (columns) to group the records by.",
					Validators: []validator.List{
						listvalidator.SizeAtLeast(1),
					},
				},
				"aggregate": schema.MapAttribute{
					ElementType: types.ListType{
						ElemType: types.StringType,
					},
					Optional:    true,
					Description: "Choose the fields should appear in calculations.",
					Validators: []validator.Map{
						mapvalidator.ValueListsAre(listvalidator.ValueStringsAre(stringvalidator.OneOf(CalculationTypes()...))),
					},
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("filter_fields_by_name")),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("grouping_to_matrix")),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("limit")),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("sort_by")),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("series_to_rows")),
		},
	}
}

type TransformationGroupingToMatrix struct {
	Column types.String `tfsdk:"column"`
	Row    types.String `tfsdk:"row"`
	Cell   types.String `tfsdk:"cell"`
}

func transformationGroupingToMatrixBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: "Limit the number of rows displayed.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"column": schema.StringAttribute{
					Required:    true,
					Description: "The column to group the records by.",
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},
				"row": schema.StringAttribute{
					Required:    true,
					Description: "The row to group the records by.",
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},
				"cell": schema.StringAttribute{
					Required:    true,
					Description: "The value to display in a cell.",
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("filter_fields_by_name")),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("group_by")),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("limit")),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("sort_by")),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("series_to_rows")),
		},
	}
}

type TransformationLimit struct {
	Limit types.Int64 `tfsdk:"limit"`
}

func transformationLimitBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: "Limit the number of rows displayed.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"limit": schema.Int64Attribute{
					Required:    true,
					Description: "How many rows to display.",
					Validators: []validator.Int64{
						int64validator.AtLeast(0),
					},
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("filter_fields_by_name")),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("group_by")),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("grouping_to_matrix")),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("sort_by")),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("series_to_rows")),
		},
	}
}

type TransformationSeriesToRows struct {
}

func transformationSeriesToRowsBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: "Create a row for each field and a column for each calculation.",
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("filter_fields_by_name")),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("group_by")),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("grouping_to_matrix")),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("limit")),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("sort_by")),
		},
	}
}

type TransformationSortBy struct {
	Field   types.String `tfsdk:"field"`
	Reverse types.Bool   `tfsdk:"reverse"`
}

func transformationSortByBlock() schema.Block {
	return schema.ListNestedBlock{
		Description: "Sort each frame by the configured field.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"field": schema.StringAttribute{
					Required:    true,
					Description: "The field to sort the frame by.",
				},
				"reverse": schema.BoolAttribute{
					Optional:    true,
					Description: "Whether to sort frames in a reverse order.",
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("filter_fields_by_name")),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("group_by")),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("grouping_to_matrix")),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("limit")),
			listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("series_to_rows")),
		},
	}
}

// creators

func createTransformations(transformations []Transformations) []grafana.Transformation {
	targets := make([]grafana.Transformation, 0)

	type GroupByField struct {
		Operation    string   `json:"operation"`
		Aggregations []string `json:"aggregations,omitempty"`
	}

	for _, transformation := range transformations {
		for _, step := range transformation.Steps {
			for _, byName := range step.FilterByName {
				names := make([]string, len(byName.Names))

				for i, name := range byName.Names {
					names[i] = name.ValueString()
				}

				targets = append(targets, grafana.Transformation{
					Id: "filterFieldsByName",
					Options: map[string]interface{}{
						"include": map[string]interface{}{
							"names": names,
						},
					},
				})
			}

			for _, groupBy := range step.GroupBy {
				fields := make(map[string]interface{})

				for _, field := range groupBy.GroupBy {
					fields[field.ValueString()] = GroupByField{
						Operation: "groupby",
					}
				}

				for field, calculations := range groupBy.Aggregate {
					calc := make([]string, 0)

					for _, c := range calculations {
						calc = append(calc, c.ValueString())
					}

					fields[field] = GroupByField{
						Operation:    "aggregate",
						Aggregations: calc,
					}
				}

				targets = append(targets, grafana.Transformation{
					Id: "groupBy",
					Options: map[string]interface{}{
						"fields": fields,
					},
				})
			}

			for _, groupToMatrix := range step.GroupingToMatrix {
				targets = append(targets, grafana.Transformation{
					Id: "groupingToMatrix",
					Options: map[string]interface{}{
						"columnField": groupToMatrix.Column.ValueString(),
						"rowField":    groupToMatrix.Row.ValueString(),
						"valueField":  groupToMatrix.Cell.ValueString(),
					},
				})
			}

			for _, limit := range step.Limit {
				targets = append(targets, grafana.Transformation{
					Id: "limit",
					Options: map[string]interface{}{
						"limitField": limit.Limit.ValueInt64(),
					},
				})
			}

			for _ = range step.SeriesToRows {
				targets = append(targets, grafana.Transformation{
					Id: "seriesToRows",
				})
			}

			for _, sortBy := range step.SortBy {
				targets = append(targets, grafana.Transformation{
					Id: "sortBy",
					Options: map[string]interface{}{
						"sort": []map[string]interface{}{
							{
								"desc":  sortBy.Reverse.ValueBool(),
								"field": sortBy.Field.ValueString(),
							},
						},
					},
				})
			}
		}
	}

	return targets
}
