terraform {
  required_providers {
    idcloudhost = {
      version = "0.2.0"
      source  = "github.com/bapung/idcloudhost"
    }
  }
}

provider "idcloudhost" {
  # auth_token is read from IDCLOUDHOST_AUTH_TOKEN environment variable
  # region defaults to "jkt01"
}

# Create a network for the load balancer
resource "idcloudhost_network" "lb_network" {
  name = "loadbalancer-network"
}

# Create VMs to use as load balancer targets
resource "idcloudhost_vm" "web_server_1" {
  name               = "web-server-1"
  os_name            = "ubuntu"
  os_version         = "20.04"
  disks              = 20
  vcpu               = 2
  memory             = 2048
  username           = "ubuntu"
  initial_password   = "SecurePassword123!"
  billing_account_id = 1200132376  # Replace with your billing account ID
  backup             = false
}

resource "idcloudhost_vm" "web_server_2" {
  name               = "web-server-2"
  os_name            = "ubuntu"
  os_version         = "20.04"
  disks              = 20
  vcpu               = 2
  memory             = 2048
  username           = "ubuntu"
  initial_password   = "SecurePassword123!"
  billing_account_id = 1200132376  # Replace with your billing account ID
  backup             = false
}

# Create a load balancer with targets and forwarding rules
resource "idcloudhost_loadbalancer" "web_lb" {
  display_name        = "web-load-balancer"
  billing_account_id  = 1200132376  # Replace with your billing account ID
  network_uuid        = idcloudhost_network.lb_network.uuid
  reserve_public_ip   = true

  # Add VM targets
  targets {
    target_type = "vm"
    target_uuid = idcloudhost_vm.web_server_1.uuid
  }

  targets {
    target_type = "vm"
    target_uuid = idcloudhost_vm.web_server_2.uuid
  }

  # HTTP forwarding rule
  forwarding_rules {
    source_port = 80
    target_port = 80
  }

  # HTTPS forwarding rule
  forwarding_rules {
    source_port = 443
    target_port = 443
  }
}

# Create a simple load balancer without initial targets
resource "idcloudhost_loadbalancer" "simple_lb" {
  display_name        = "simple-load-balancer"
  billing_account_id  = 1200132376  # Replace with your billing account ID
  network_uuid        = idcloudhost_network.lb_network.uuid
  reserve_public_ip   = false
}

# Output load balancer details
output "web_lb_uuid" {
  value       = idcloudhost_loadbalancer.web_lb.uuid
  description = "The UUID of the web load balancer"
}

output "web_lb_private_address" {
  value       = idcloudhost_loadbalancer.web_lb.private_address
  description = "The private IP address of the web load balancer"
}

output "web_lb_targets" {
  value       = idcloudhost_loadbalancer.web_lb.targets
  description = "The targets attached to the web load balancer"
}

output "web_lb_forwarding_rules" {
  value       = idcloudhost_loadbalancer.web_lb.forwarding_rules
  description = "The forwarding rules of the web load balancer"
}

output "simple_lb_uuid" {
  value       = idcloudhost_loadbalancer.simple_lb.uuid
  description = "The UUID of the simple load balancer"
}
