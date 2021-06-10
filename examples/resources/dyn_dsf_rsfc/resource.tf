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
