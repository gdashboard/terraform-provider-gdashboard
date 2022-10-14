data "gdashboard_stat" "status" {
  title = "Status"

  targets {
    prometheus {
      uid     = "prometheus"
      expr    = "up{container_name='container'}"
      instant = true
    }
  }
}