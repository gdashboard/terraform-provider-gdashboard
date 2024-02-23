data "gdashboard_logs" "logs" {
  title       = "Logs example"
  description = "Text description"

  graph {
    show_time          = true
    show_unique_labels = true
    show_common_labels = true
    wrap_lines         = true
    prettify_json      = true
    enable_log_details = false
    deduplication      = "exact"
    order              = "oldest_first"
  }

}
