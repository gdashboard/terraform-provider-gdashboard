data "gdashboard_text" "text" {
  title       = "Text"
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
