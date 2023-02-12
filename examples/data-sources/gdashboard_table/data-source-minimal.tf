data "gdashboard_table" "status" {
  title = "Example"

  graph {
    column {
      filterable = true
    }

    cell {
      inspectable  = true
      display_mode = "basic"
    }

    footer {
      pagination = true
    }
  }

  queries {
    prometheus {
      uid     = "prometheus"
      expr    = "up{container_name='container'}"
      instant = true
    }
  }
}
