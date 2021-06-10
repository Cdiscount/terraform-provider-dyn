resource "dyn_dsf_monitor" "monitor" {
  label          = "my-monitor"
  protocol       = "HTTP"
  response_count = 1
  probe_interval = 60
  retries        = 1
  active         = true

  options {
    timeout = 2
    port    = 80
    path    = "/check"
    host    = "check.test"
  }
}
