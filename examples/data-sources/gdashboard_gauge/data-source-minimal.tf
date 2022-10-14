data "gdashboard_gauge" "jvm_memory" {
  title = "JVM Memory"

  targets {
    prometheus {
      uid     = "prometheus"
      expr    = "sum(increase(jvm_memory_total{container_name='container'}[$__rate_interval]))"
      instant = true
    }
  }
}
