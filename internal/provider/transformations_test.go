package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTransformationsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTransformationTableDataSource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_table.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_table.test", "json", testAccTransformationTableDataSourceExpectedJson),
				),
			},
		},
	})
}

const testAccTransformationTableDataSource = `
data "gdashboard_table" "test" {
  title = "Test"

  transform {
	step {
	  series_to_rows { }
	}

	step {
      sort_by {
	    field = "test 2"
      }
	}

	step {
	  limit {
	    limit = 10
	  }
	}

	step {
      group_by {
		by 		  = ["a", "b"]
		aggregate = {
		  "c" = ["min", "max"]
		  "d" = ["mean", "first"]
		}
	  }
	}

	step {
	  grouping_to_matrix {
		column = "a"
		row    = "b"
		cell   = "c"
	  }
    }

    step {
      filter_fields_by_name {
        names = ["a", "b", "c"]
      }
    }
  }
}
`

const testAccTransformationTableDataSourceExpectedJson = `{
  "editable": false,
  "error": false,
  "gridPos": {},
  "id": 0,
  "isNew": true,
  "span": 12,
  "title": "Test",
  "transparent": false,
  "type": "table",
  "transformations": [
    {
      "id": "seriesToRows",
      "options": null
    },
    {
      "id": "sortBy",
      "options": {
        "sort": [
          {
            "desk": false,
            "field": "test 2"
          }
        ]
      }
    },
    {
      "id": "limit",
      "options": {
        "limitField": 10
      }
    },
    {
      "id": "groupBy",
      "options": {
        "fields": {
          "a": {
            "operation": "groupby"
          },
          "b": {
            "operation": "groupby"
          },
          "c": {
            "operation": "aggregate",
            "aggregations": [
              "min",
              "max"
            ]
          },
          "d": {
            "operation": "aggregate",
            "aggregations": [
              "mean",
              "first"
            ]
          }
        }
      }
    },
    {
      "id": "groupingToMatrix",
      "options": {
        "columnField": "a",
        "rowField": "b",
        "valueField": "c"
      }
    },
    {
      "id": "filterFieldsByName",
      "options": {
        "include": {
          "names": [
            "a",
            "b",
            "c"
          ]
        }
      }
    }
  ],
  "options": {
    "showHeader": true,
    "footer": {
      "show": false,
      "enablePagination": false
    }
  },
  "fieldConfig": {
    "defaults": {
      "unit": "",
      "color": {
        "mode": "palette-classic",
        "fixedColor": "green",
        "seriesBy": "last"
      },
      "thresholds": {
        "mode": "absolute",
        "steps": [
          {
            "color": "green",
            "value": null
          }
        ]
      },
      "custom": {
        "axisPlacement": "",
        "barAlignment": 0,
        "drawStyle": "",
        "fillOpacity": 0,
        "gradientMode": "",
        "lineInterpolation": "",
        "lineWidth": 0,
        "pointSize": 0,
        "showPoints": "",
        "spanNulls": false,
        "hideFrom": {
          "legend": false,
          "tooltip": false,
          "viz": false
        },
        "lineStyle": {
          "fill": ""
        },
        "scaleDistribution": {
          "type": ""
        },
        "stacking": {
          "group": "",
          "mode": ""
        },
        "thresholdsStyle": {
          "mode": ""
        }
      }
    }
  }
}`
