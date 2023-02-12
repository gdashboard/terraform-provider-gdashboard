package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTextDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccTextDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_text.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_text.test", "json", testAccTextDataSourceConfigExpectedJson),
				),
			},
			{
				Config: testAccTextDataSourceProviderDefaultsConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gdashboard_text.test", "title", "Test"),
					resource.TestCheckResourceAttr("data.gdashboard_text.test", "json", testAccTextDataSourceProviderDefaultsConfigExpectedJson),
				),
			},
		},
	})
}

const testAccTextDataSourceConfig = `
data "gdashboard_text" "test" {
  title       = "Test"
  description = "Text description"

  graph {
    mode = "code"

	code {
	  language 			= "sql"
	  show_line_numbers = true
      show_mini_map     = true
	}
    
    content = "SELECT * FROM users;"
  }
}
`

const testAccTextDataSourceConfigExpectedJson = `{
  "editable": false,
  "error": false,
  "gridPos": {},
  "id": 0,
  "isNew": true,
  "span": 12,
  "title": "Test",
  "description": "Text description",
  "transparent": false,
  "type": "text",
  "options": {
    "mode": "code",
    "content": "SELECT * FROM users;",
    "code": {
      "language": "sql",
      "showLineNumbers": true,
      "showMiniMap": true
    }
  }
}`

const testAccTextDataSourceProviderDefaultsConfig = `
data "gdashboard_text" "test" {
  title = "Test"
}
`

const testAccTextDataSourceProviderDefaultsConfigExpectedJson = `{
  "editable": false,
  "error": false,
  "gridPos": {},
  "id": 0,
  "isNew": true,
  "span": 12,
  "title": "Test",
  "transparent": false,
  "type": "text",
  "options": {
    "mode": "markdown",
    "content": "",
    "code": {
      "language": "plaintext",
      "showLineNumbers": false,
      "showMiniMap": false
    }
  }
}`
