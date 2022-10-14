data "gdashboard_timeseries" "jvm_memory" {
  title = "JVM memory"

  queries {
    prometheus {
      uid  = "prometheus"
      expr = "sum(increase(jvm_memory_total{container_name='container'}[$__rate_interval]))"
    }
  }
}
