data "gdashboard_text" "text" {
  title       = "Markdown example"
  description = "Some text"

  graph {
    mode = "markdown"

    content = <<-EOT
      # Header

      ### Body

      Hello World!
    EOT
  }

}
