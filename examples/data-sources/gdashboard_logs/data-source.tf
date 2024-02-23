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

  queries {
    cloudwatch {
      logs {
        uid    = "cloudwatch"
        region = "eu-west-2"

        expression = <<-EOT
        fields @timestamp, @message |
          sort @timestamp desc |
          limit 20
        EOT

        log_group {
          arn  = "arn:aws:logs:eu-west-2:123456789012:log-group:/ecs/production/service-1:*"
          name = "/ecs/production/service-1"
        }

        log_group {
          arn  = "arn:aws:logs:eu-west-2:123456789012:log-group:/ecs/production/service-2:*"
          name = "/ecs/production/service-2"
        }
      }
    }
  }
}
