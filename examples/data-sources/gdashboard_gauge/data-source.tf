data "gdashboard_gauge" "jobs_processed" {
  title = "Example"

  field {
    unit = "percent"
  }

  graph {
    orientation            = "horizontal"
    show_threshold_labels  = true
    show_threshold_markers = true

    options {
      calculation = "lastNotNull"
    }
  }

  targets {
    prometheus {
      uid           = "prometheus"
      expr          = "sum(increase(jvm_memory_total{container_name='container'}[$__rate_interval]))"
      min_interval  = "30"
      legend_format = "{{job_type}}"
      instant       = true
    }
  }
}
