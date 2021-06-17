resource "dyn_traffic_director" "example" {
  label = "my-traffic-director"

  node {
    zone = "my-zone.example.net"
    fqdn = "example.my-zone.example.net"
  }
}

resource "dyn_dsf_response_pool" "response_pool" {
  label               = "my-response-pool"
  traffic_director_id = dyn_traffic_director.example.id
  # automation        = "auto"
}

resource "dyn_dsf_rsfc" "rsfc" {
  label               = "my-rsfc"
  traffic_director_id = dyn_traffic_director.example.id
  response_pool_id    = dyn_dsf_response_pool.response_pool.id
}

resource "dyn_dsf_record_set" "record_set" {
  label               = "my-a-record-set"
  traffic_director_id = dyn_traffic_director.example.id
  response_pool_id    = dyn_dsf_response_pool.response_pool.id
  dsf_rsfc_id         = dyn_dsf_rsfc.rsfc.id

  rdata_class = "A"
  ttl         = 150
  automation  = "auto"

  serve_count = 1
}

resource "dyn_dsf_record" "my-record" {
  traffic_director_id = dyn_traffic_director.example.id
  record_set_id       = dyn_dsf_record_set.record_set.id

  label       = "my-record"
  master_line = "1.2.3.4"
}
