data "gdashboard_stat" "status" {
  title = "Status"

  queries {
    prometheus {
      uid     = "prometheus"
      expr    = "up{container_name='container'}"
      instant = true
    }
  }
}
