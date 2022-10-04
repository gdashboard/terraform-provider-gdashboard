data "gdashboard_timeseries" "test" {
  title = "Test"

  legend {
    calculations = ["min", "max", "mean"]
    display_mode = "table"
    placement    = "bottom"
  }

  tooltip {
    mode = "multi"
  }

  field {
    unit     = "bytes"
    decimals = 1
    min      = 0
    max      = 10000

    color {
      mode        = "palette-classic"
      fixed_color = "green"
      series_by   = "last"
    }
  }

  graph {
    fill_opacity = 10
    show_points  = "always"
    span_nulls   = true
  }

  targets {
    prometheus {
      uid           = "prometheus"
      expr          = "sum(increase(jvm_memory_total{container_name='container'}[$__rate_interval]))"
      instant       = false
      ref_id        = "Prometheus_Query"
      min_interval  = "30"
      legend_format = "Memory total"
    }

    cloudwatch {
      uid         = "cloudwatch"
      namespace   = "AWS/ApplicationELB"
      metric_name = "HTTPCode_Target_2XX_Count"
      statistic   = "Sum"
      match_exact = true
      region      = "af-south-1"

      dimension {
        name  = "LoadBalancer"
        value = "lb_arn_suffix"
      }

      dimension {
        name  = "TargetGroup"
        value = "target_group"
      }

      ref_id        = "CW_Query"
      period        = "30"
      legend_format = "Request Count"
    }
  }

}
