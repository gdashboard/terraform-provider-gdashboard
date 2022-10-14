data "gdashboard_bar_gauge" "jobs_processed" {
  title = "Jobs Processed"

  targets {
    prometheus {
      uid     = "prometheus"
      expr    = "sort_desc(sum(increase(jobs_processed_total{container_name='container'}[$__range])) by (job_type))"
      instant = true
    }
  }
}
