package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDashboardDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccDashboardDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_dashboard.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_dashboard.test", "json", testAccDashboardDataSourceConfigExpectedJson),
				),
			},
			{
				Config: testAccDashboardDataSourceProviderDefaultsConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_dashboard.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_dashboard.test", "json", testAccDashboardDataSourceProviderDefaultsConfigExpectedJson),
				),
			},
			{
				Config: testAccDashboardDataSourceProviderCustomDefaultsConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_dashboard.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_dashboard.test", "json", testAccDashboardDataSourceProviderCustomDefaultsConfigExpectedJson),
				),
			},
			{
				Config: testAccDashboardDataSourceProvider_Variable_Custom_Valid,
				Check:  resource.TestCheckResourceAttr("data.gdashboard_dashboard.test", "json", testAccDashboardDataSourceProvider_Variable_Custom_Valid_ExpectedJson),
			},
			{
				Config:      testAccDashboardDataSourceProvider_Variable_TextBox_MissingFields,
				ExpectError: regexp.MustCompile("The argument \"name\" is required, but no definition was found"),
			},
			{
				Config: testAccDashboardDataSourceProvider_Variable_TextBox_Valid,
				Check:  resource.TestCheckResourceAttr("data.gdashboard_dashboard.test", "json", testAccDashboardDataSourceProvider_Variable_TextBox_Valid_ExpectedJson),
			},
			{
				Config:      testAccDashboardDataSourceProvider_Variable_Adhoc_MissingFields,
				ExpectError: regexp.MustCompile("Attribute \"variables\\[0]\\.adhoc\\[0]\\.datasource\" must be specified when"),
			},
			{
				Config: testAccDashboardDataSourceProvider_Variable_Adhoc_Valid,
				Check:  resource.TestCheckResourceAttr("data.gdashboard_dashboard.test", "json", testAccDashboardDataSourceProvider_Variable_Adhoc_Valid_ExpectedJson),
			},
			{
				Config:      testAccDashboardDataSourceProvider_Variable_Datasource_MissingFields,
				ExpectError: regexp.MustCompile("Attribute \"variables\\[0]\\.datasource\\[0]\\.source\" must be specified when"),
			},
			{
				Config: testAccDashboardDataSourceProvider_Variable_Datasource_Valid,
				Check:  resource.TestCheckResourceAttr("data.gdashboard_dashboard.test", "json", testAccDashboardDataSourceProvider_Variable_Datasource_Valid_ExpectedJson),
			},
			{
				Config:      testAccDashboardDataSourceProvider_Variable_Query_MissingFields,
				ExpectError: regexp.MustCompile("Attribute \"variables\\[0]\\.query\\[0]\\.target\" must be specified when"),
			},
			{
				Config: testAccDashboardDataSourceProvider_Variable_Query_Valid,
				Check:  resource.TestCheckResourceAttr("data.gdashboard_dashboard.test", "json", testAccDashboardDataSourceProvider_Variable_Query_Valid_ExpectedJson),
			},
			{
				Config:      testAccDashboardDataSourceProvider_Variable_Interval_MissingFields,
				ExpectError: regexp.MustCompile("The argument \"name\" is required, but no definition was found"),
			},
			{
				Config: testAccDashboardDataSourceProvider_Variable_Interval_Valid,
				Check:  resource.TestCheckResourceAttr("data.gdashboard_dashboard.test", "json", testAccDashboardDataSourceProvider_Variable_Interval_Valid_ExpectedJson),
			},
		},
	})
}

const testAccDashboardDataSourceConfig = `
data "gdashboard_dashboard" "test" {
  title         = "Test"
  uid             = "test-uid"
  editable      = false
  style         = "light"
  graph_tooltip = "shared-crosshair"
  
  time {
    from = "now-1h"
    to   = "now+1h"
  }

  variables {
	const {
      name  = "var"
      value = "const-value"
	}

    custom {
      name = "custom"
      hide = "label" 

      option {
        text  = "entry-1"
        value = "value"
      }

      option {
        text  	 = "entry-2"
        value 	 = "value"
        selected = true
      }
    }
  }

  layout {
	row {
      panel {
      	size = {
          height = 8
          width  = 10
        }
      	source = "{\"title\": \"Panel 1\"}"
      }

      panel {
      	size = {
          height = 3
          width  = 24
        }
      	source = "{\"title\": \"Panel 2\"}"
      } 
	}

	row {
      panel {
      	size = {
          height = 4
          width  = 24
        }
      	source = "{\"title\": \"Panel 3\"}"
      }

      panel {
      	size = {
          height = 3
          width  = 3
        }
      	source = "{\"title\": \"Panel 4\"}"
      } 
	}
  }
}
`

const testAccDashboardDataSourceConfigExpectedJson = `{
  "uid": "test-uid",
  "slug": "",
  "title": "Test",
  "originalTitle": "",
  "style": "light",
  "timezone": "",
  "editable": false,
  "hideControls": false,
  "panels": [
    {
      "editable": false,
      "error": false,
      "gridPos": {
        "h": 8,
        "w": 10,
        "x": 0,
        "y": 0
      },
      "id": 0,
      "isNew": false,
      "span": 0,
      "title": "Panel 1",
      "transparent": false,
      "type": "",
      "title": "Panel 1"
    },
    {
      "editable": false,
      "error": false,
      "gridPos": {
        "h": 3,
        "w": 24,
        "x": 10,
        "y": 0
      },
      "id": 0,
      "isNew": false,
      "span": 0,
      "title": "Panel 2",
      "transparent": false,
      "type": "",
      "title": "Panel 2"
    },
    {
      "editable": false,
      "error": false,
      "gridPos": {
        "h": 4,
        "w": 24,
        "x": 0,
        "y": 9
      },
      "id": 0,
      "isNew": false,
      "span": 0,
      "title": "Panel 3",
      "transparent": false,
      "type": "",
      "title": "Panel 3"
    },
    {
      "editable": false,
      "error": false,
      "gridPos": {
        "h": 3,
        "w": 3,
        "x": 24,
        "y": 9
      },
      "id": 0,
      "isNew": false,
      "span": 0,
      "title": "Panel 4",
      "transparent": false,
      "type": "",
      "title": "Panel 4"
    }
  ],
  "templating": {
    "list": [
      {
        "type": "custom",
        "name": "custom",
        "label": "",
        "hide": 1,
        "refresh": false,
        "options": [
          {
            "text": "entry-1",
            "value": "value",
            "selected": false
          },
          {
            "text": "entry-2",
            "value": "value",
            "selected": true
          }
        ],
        "includeAll": false,
        "allValue": "",
        "multi": false,
        "query": "entry-1 : value, entry-2 : value",
        "regex": "",
        "current": {
          "text": "entry-2",
          "value": "value",
          "selected": true
        },
        "sort": 0
      },
      {
        "type": "constant",
        "name": "var",
        "label": "",
        "hide": 0,
        "refresh": false,
        "options": [],
        "includeAll": false,
        "allValue": "",
        "multi": false,
        "query": "const-value",
        "regex": "",
        "current": {
          "text": null,
          "value": "",
          "selected": false
        },
        "sort": 0
      }
    ]
  },
  "annotations": {
    "list": null
  },
  "schemaVersion": 0,
  "version": 1,
  "links": null,
  "time": {
    "from": "now-1h",
    "to": "now+1h"
  },
  "timepicker": {
    "refresh_intervals": null,
    "time_options": null
  },
  "graphTooltip": 1
}`

const testAccDashboardDataSourceProviderCustomDefaultsConfig = `
provider "gdashboard" {
  defaults {
    dashboard {
      editable        = false
      graph_tooltip = "shared-tooltip"
       style         = "light"
      time {
        from = "now-12h"
        to   = "now-3h"
      }
    }
  }
}

data "gdashboard_dashboard" "test" {
  title = "Test"

  layout {

  }
}
`

const testAccDashboardDataSourceProviderCustomDefaultsConfigExpectedJson = `{
  "slug": "",
  "title": "Test",
  "originalTitle": "",
  "style": "light",
  "timezone": "",
  "editable": false,
  "hideControls": false,
  "panels": [],
  "templating": {
    "list": []
  },
  "annotations": {
    "list": null
  },
  "schemaVersion": 0,
  "version": 1,
  "links": null,
  "time": {
    "from": "now-12h",
    "to": "now-3h"
  },
  "timepicker": {
    "refresh_intervals": null,
    "time_options": null
  },
  "graphTooltip": 2
}`

const testAccDashboardDataSourceProviderDefaultsConfig = `
data "gdashboard_dashboard" "test" {
  title = "Test"

  layout {

  }
}
`

const testAccDashboardDataSourceProviderDefaultsConfigExpectedJson = `{
  "slug": "",
  "title": "Test",
  "originalTitle": "",
  "style": "dark",
  "timezone": "",
  "editable": true,
  "hideControls": false,
  "panels": [],
  "templating": {
    "list": []
  },
  "annotations": {
    "list": null
  },
  "schemaVersion": 0,
  "version": 1,
  "links": null,
  "time": {
    "from": "now-6h",
    "to": "now"
  },
  "timepicker": {
    "refresh_intervals": null,
    "time_options": null
  }
}`

//// Custom start

const testAccDashboardDataSourceProvider_Variable_Custom_Valid = `
data "gdashboard_dashboard" "test" {
  title = "Test"

  variables {
    custom {
      name          = "custom"
      label         = "Label"
      description   = "Description"
      hide          = "label"
      multi   = true
      
       include_all {
        enabled      = true
        custom_value = "*"
      }

      option {
        text  = "entry-1"
        value = "value"
      }

      option {
        text  	 = "entry-2"
        value 	 = "value"
        selected = true
      }
    }
  }

  layout {

  }
}
`

const testAccDashboardDataSourceProvider_Variable_Custom_Valid_ExpectedJson = `{
  "slug": "",
  "title": "Test",
  "originalTitle": "",
  "style": "dark",
  "timezone": "",
  "editable": true,
  "hideControls": false,
  "panels": [],
  "templating": {
    "list": [
      {
        "type": "custom",
        "name": "custom",
        "description": "Description",
        "label": "Label",
        "hide": 1,
        "refresh": false,
        "options": [
          {
            "text": "entry-1",
            "value": "value",
            "selected": false
          },
          {
            "text": "entry-2",
            "value": "value",
            "selected": true
          }
        ],
        "includeAll": true,
        "allValue": "*",
        "multi": true,
        "query": "entry-1 : value, entry-2 : value",
        "regex": "",
        "current": {
          "text": "entry-2",
          "value": "value",
          "selected": true
        },
        "sort": 0
      }
    ]
  },
  "annotations": {
    "list": null
  },
  "schemaVersion": 0,
  "version": 1,
  "links": null,
  "time": {
    "from": "now-6h",
    "to": "now"
  },
  "timepicker": {
    "refresh_intervals": null,
    "time_options": null
  }
}`

//// Custom end

//// Textbox start

const testAccDashboardDataSourceProvider_Variable_TextBox_Valid = `
data "gdashboard_dashboard" "test" {
  title = "Test"

  variables {
    textbox {
      name          = "custom"
      label         = "Label"
      description   = "Description"
      default_value = "*"
      hide          = "label"
    }
  }

  layout {

  }
}`

const testAccDashboardDataSourceProvider_Variable_TextBox_Valid_ExpectedJson = `{
  "slug": "",
  "title": "Test",
  "originalTitle": "",
  "style": "dark",
  "timezone": "",
  "editable": true,
  "hideControls": false,
  "panels": [],
  "templating": {
    "list": [
      {
        "type": "textbox",
        "name": "custom",
        "description": "Description",
        "label": "Label",
        "hide": 1,
        "refresh": false,
        "options": [],
        "includeAll": false,
        "allValue": "",
        "multi": false,
        "query": "*",
        "regex": "",
        "current": {
          "text": null,
          "value": "",
          "selected": false
        },
        "sort": 0
      }
    ]
  },
  "annotations": {
    "list": null
  },
  "schemaVersion": 0,
  "version": 1,
  "links": null,
  "time": {
    "from": "now-6h",
    "to": "now"
  },
  "timepicker": {
    "refresh_intervals": null,
    "time_options": null
  }
}`

const testAccDashboardDataSourceProvider_Variable_TextBox_MissingFields = `
data "gdashboard_dashboard" "test" {
  title = "Test"

  variables {
    textbox {

    }
  }

  layout { }
}`

//// Textbox end

//// Adhoc start

const testAccDashboardDataSourceProvider_Variable_Adhoc_Valid = `
data "gdashboard_dashboard" "test" {
  title = "Test"

  variables {
    adhoc {
      name          = "custom"
      label         = "Label"
      description   = "Description"
      hide          = "label"

      datasource {
        type = "prometheus"
        uid  = "uid"
      }

      filter {
        key      = "__name__"
        operator = "!="
        value    = "any"
      }

      filter {
        key 	 = "host"
        operator = "=~"
        value 	 = "^prod$"
      }
    }
  }

  layout { }
}`

const testAccDashboardDataSourceProvider_Variable_Adhoc_Valid_ExpectedJson = `{
  "slug": "",
  "title": "Test",
  "originalTitle": "",
  "style": "dark",
  "timezone": "",
  "editable": true,
  "hideControls": false,
  "panels": [],
  "templating": {
    "list": [
      {
        "type": "adhoc",
        "name": "custom",
        "description": "Description",
        "label": "Label",
        "hide": 1,
        "datasource": {
          "uid": "uid",
          "type": "prometheus"
        },
        "filters": [
          {
            "condition": "",
            "key": "__name__",
            "operator": "!=",
            "value": "any"
          },
          {
            "condition": "",
            "key": "host",
            "operator": "=~",
            "value": "^prod$"
          }
        ],
        "refresh": false,
        "options": [],
        "includeAll": false,
        "allValue": "",
        "multi": false,
        "query": null,
        "regex": "",
        "current": {
          "text": null,
          "value": "",
          "selected": false
        },
        "sort": 0
      }
    ]
  },
  "annotations": {
    "list": null
  },
  "schemaVersion": 0,
  "version": 1,
  "links": null,
  "time": {
    "from": "now-6h",
    "to": "now"
  },
  "timepicker": {
    "refresh_intervals": null,
    "time_options": null
  }
}`

const testAccDashboardDataSourceProvider_Variable_Adhoc_MissingFields = `
data "gdashboard_dashboard" "test" {
  title = "Test"

  variables {
    adhoc {
      name = "test"
    }
  }
 
  layout { }
}`

//// Adhoc end

//// Datasource start

const testAccDashboardDataSourceProvider_Variable_Datasource_Valid = `
data "gdashboard_dashboard" "test" {
  title = "Test"

  variables {
    datasource {
      name          = "custom"
      label         = "Label"
      description   = "Description"
      hide          = "label"
      multi   = true
      
       include_all {
        enabled      = true
        custom_value = "*"
      }

      source {
        type   = "prometheus"
        filter = "^prod$"
      }
    }
  }

  layout { }
}`

const testAccDashboardDataSourceProvider_Variable_Datasource_Valid_ExpectedJson = `{
  "slug": "",
  "title": "Test",
  "originalTitle": "",
  "style": "dark",
  "timezone": "",
  "editable": true,
  "hideControls": false,
  "panels": [],
  "templating": {
    "list": [
      {
        "type": "datasource",
        "name": "custom",
        "description": "Description",
        "label": "Label",
        "hide": 1,
        "refresh": 1,
        "options": [],
        "includeAll": true,
        "allValue": "*",
        "multi": true,
        "query": "prometheus",
        "regex": "^prod$",
        "current": {
          "text": null,
          "value": "",
          "selected": false
        },
        "sort": 0
      }
    ]
  },
  "annotations": {
    "list": null
  },
  "schemaVersion": 0,
  "version": 1,
  "links": null,
  "time": {
    "from": "now-6h",
    "to": "now"
  },
  "timepicker": {
    "refresh_intervals": null,
    "time_options": null
  }
}`

const testAccDashboardDataSourceProvider_Variable_Datasource_MissingFields = `
data "gdashboard_dashboard" "test" {
  title = "Test"

  variables {
    datasource {
      name = "test"
    }
  }
 
  layout { }
}`

//// Datasource end

//// Query start

const testAccDashboardDataSourceProvider_Variable_Query_Valid = `
data "gdashboard_dashboard" "test" {
  title = "Test"

  variables {
    query {
      name          = "custom"
      label         = "Label"
      description   = "Description"
      hide          = "label"
      multi   = true
      refresh         = "time-range-change"
       regex          = "^prod$"

       include_all {
        enabled      = true
        custom_value = "*"
      }

      sort {
        type  = "alphabetical-case-insensitive"
        order = "desc"
      }

      target {
        prometheus {
          uid  = "uid"
          expr = "up{service='test'}"
        }
      }
    }
  }

  layout { }
}`

const testAccDashboardDataSourceProvider_Variable_Query_Valid_ExpectedJson = `{
  "slug": "",
  "title": "Test",
  "originalTitle": "",
  "style": "dark",
  "timezone": "",
  "editable": true,
  "hideControls": false,
  "panels": [],
  "templating": {
    "list": [
      {
        "type": "query",
        "name": "custom",
        "description": "Description",
        "label": "Label",
        "hide": 1,
        "datasource": {
          "uid": "uid",
          "type": "prometheus"
        },
        "refresh": 2,
        "options": [],
        "includeAll": true,
        "allValue": "*",
        "multi": true,
        "query": {
          "query": "up{service='test'}",
          "refId": "StandardVariableQuery"
        },
        "regex": "^prod$",
        "current": {
          "text": null,
          "value": "",
          "selected": false
        },
        "sort": 6,
        "definition": "up{service='test'}"
      }
    ]
  },
  "annotations": {
    "list": null
  },
  "schemaVersion": 0,
  "version": 1,
  "links": null,
  "time": {
    "from": "now-6h",
    "to": "now"
  },
  "timepicker": {
    "refresh_intervals": null,
    "time_options": null
  }
}`

const testAccDashboardDataSourceProvider_Variable_Query_MissingFields = `
data "gdashboard_dashboard" "test" {
  title = "Test"

  variables {
    query {
      name = "test"
    }
  }
 
  layout { }
}`

//// Query end

//// Interval start

const testAccDashboardDataSourceProvider_Variable_Interval_Valid = `
data "gdashboard_dashboard" "test" {
  title = "Test"

  variables {
    interval {
      name          = "custom"
      label         = "Label"
      description   = "Description"
      hide          = "label"
      intervals 	= ["1m", "10m", "30m", "1h", "6h", "12h", "1d", "7d", "14d", "30d"]

      auto {
        enabled 	 = true
        step_count 	 = 30
        min_interval = "10s"
      }
    }
  }

  layout { }
}`

const testAccDashboardDataSourceProvider_Variable_Interval_Valid_ExpectedJson = `{
  "slug": "",
  "title": "Test",
  "originalTitle": "",
  "style": "dark",
  "timezone": "",
  "editable": true,
  "hideControls": false,
  "panels": [],
  "templating": {
    "list": [
      {
        "type": "interval",
        "name": "custom",
        "description": "Description",
        "label": "Label",
        "hide": 1,
        "auto": true,
        "auto_count": 30,
        "auto_min": "10s",
        "refresh": false,
        "options": [
          {
            "text": "auto",
            "value": "$__auto_interval_custom",
            "selected": false
          },
          {
            "text": "1m",
            "value": "1m",
            "selected": true
          },
          {
            "text": "10m",
            "value": "10m",
            "selected": false
          },
          {
            "text": "30m",
            "value": "30m",
            "selected": false
          },
          {
            "text": "1h",
            "value": "1h",
            "selected": false
          },
          {
            "text": "6h",
            "value": "6h",
            "selected": false
          },
          {
            "text": "12h",
            "value": "12h",
            "selected": false
          },
          {
            "text": "1d",
            "value": "1d",
            "selected": false
          },
          {
            "text": "7d",
            "value": "7d",
            "selected": false
          },
          {
            "text": "14d",
            "value": "14d",
            "selected": false
          },
          {
            "text": "30d",
            "value": "30d",
            "selected": false
          }
        ],
        "includeAll": false,
        "allValue": "",
        "multi": false,
        "query": "1m,10m,30m,1h,6h,12h,1d,7d,14d,30d",
        "regex": "",
        "current": {
          "text": "1m",
          "value": "1m",
          "selected": true
        },
        "sort": 0
      }
    ]
  },
  "annotations": {
    "list": null
  },
  "schemaVersion": 0,
  "version": 1,
  "links": null,
  "time": {
    "from": "now-6h",
    "to": "now"
  },
  "timepicker": {
    "refresh_intervals": null,
    "time_options": null
  }
}`

const testAccDashboardDataSourceProvider_Variable_Interval_MissingFields = `
data "gdashboard_dashboard" "test" {
  title = "Test"

  variables {
    interval { }
  }
 
  layout { }
}`

//// Interval end
