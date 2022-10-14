data "gdashboard_bar_gauge" "jobs_processed" {
  title = "Jobs Processed"

  field {
    decimals = 0
  }

  graph {
    orientation  = "horizontal"
    display_mode = "basic"

    options {
      calculation = "lastNotNull"
    }
  }

  queries {
    prometheus {
      uid           = "prometheus"
      expr          = "sort_desc(sum(increase(jobs_processed_total{container_name='container'}[$__range])) by (job_type))"
      min_interval  = "30"
      legend_format = "{{job_type}}"
      instant       = true
    }
  }
}
