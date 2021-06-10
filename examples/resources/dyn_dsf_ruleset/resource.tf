resource "dyn_traffic_director" "example" {
  label = "my-traffic-director"

  node {
    zone = "my-zone.example.net"
    fqdn = "example.my-zone.example.net"
  }
}

resource "dyn_dsf_response_pool" "response_pool1" {
  label               = "my-response-pool"
  traffic_director_id = dyn_traffic_director.example.id
}

resource "dyn_dsf_response_pool" "response_pool2" {
  label               = "my-second-response-pool"
  traffic_director_id = dyn_traffic_director.example.id
}

resource "dyn_dsf_ruleset" "ruleset" {
  label               = "my-ruleset"
  traffic_director_id = dyn_traffic_director.example.id

  response_pool_ids = [
    dyn_dsf_response_pool.response_pool1.id,
    dyn_dsf_response_pool.response_pool2.id,
  ]
}
