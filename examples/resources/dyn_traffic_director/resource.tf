resource "dyn_traffic_director" "example" {
  label = "my-traffic-director"
  # ttl = 300

  node {
    zone = "my-zone.example.net"
    fqdn = "example.my-zone.example.net"
  }

  node {
    zone = "my-zone.example.net"
    fqdn = "example2.my-zone.example.net"
  }
}
