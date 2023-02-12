data "gdashboard_text" "code" {
  title       = "Code example"
  description = "Text description"

  graph {
    mode = "code"

    code {
      language          = "sql"
      show_line_numbers = true
      show_mini_map     = true
    }

    content = "SELECT * FROM users;"
  }
}
