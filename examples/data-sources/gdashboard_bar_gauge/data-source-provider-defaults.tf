provider "gdashboard" {
  defaults {
    bar_gauge {
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
    }
  }
}

data "gdashboard_bar_gauge" "jobs_processed" {
  title = "Jobs Processed"

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

data "gdashboard_bar_gauge" "mails_sent" {
  title = "Mails Sent"

  queries {
    prometheus {
      uid           = "prometheus"
      expr          = "sort_desc(sum(increase(mails_sent_total{container_name='container'}[$__range])) by (mail_type))"
      min_interval  = "30"
      legend_format = "{{mail_type}}"
      instant       = true
    }
  }
}
