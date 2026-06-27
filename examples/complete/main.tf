terraform {
  required_providers {
    idcloudhost = {
      version = "0.2.0"
      source  = "teliti-dev/idcloudhost"
    }
  }
}

provider "idcloudhost" {
  # auth_token is read from IDCLOUDHOST_AUTH_TOKEN environment variable
  # region defaults to "jkt01"
}

# Variables for configuration
locals {
  billing_account_id = 1200177265  # Replace with your billing account ID
  project_name       = "example-app"
}

# ============================================================================
# Network Resources
# ============================================================================

# Create a private network for the application
resource "idcloudhost_network" "app_network" {
  name    = "${local.project_name}-network"
  default = false
}

# ============================================================================
# Firewall Resources
# ============================================================================

# Create a firewall for web servers
resource "idcloudhost_firewall" "web_firewall" {
  display_name       = "${local.project_name}-web-firewall"
  billing_account_id = local.billing_account_id
  description        = "Firewall rules for web servers"

  # Allow HTTP
  rules {
    direction          = "inbound"
    protocol           = "tcp"
    port_start         = 80
    port_end           = 80
    endpoint_spec_type = "any"
    endpoint_spec      = []
    description        = "Allow HTTP"
  }

  # Allow HTTPS
  rules {
    direction          = "inbound"
    protocol           = "tcp"
    port_start         = 443
    port_end           = 443
    endpoint_spec_type = "any"
    endpoint_spec      = []
    description        = "Allow HTTPS"
  }

  # Allow SSH from specific CIDR
  rules {
    direction          = "inbound"
    protocol           = "tcp"
    port_start         = 22
    port_end           = 22
    endpoint_spec_type = "cidr"
    endpoint_spec      = ["0.0.0.0/0"]  # Change to your IP range
    description        = "Allow SSH"
  }

  # Allow all outbound
  rules {
    direction          = "outbound"
    protocol           = "any"
    endpoint_spec_type = "any"
    endpoint_spec      = []
    description        = "Allow all outbound"
  }
}

# Create a firewall for database servers
resource "idcloudhost_firewall" "db_firewall" {
  display_name       = "${local.project_name}-db-firewall"
  billing_account_id = local.billing_account_id
  description        = "Firewall rules for database servers"

  # Allow PostgreSQL from app network
  rules {
    direction          = "inbound"
    protocol           = "tcp"
    port_start         = 5432
    port_end           = 5432
    endpoint_spec_type = "cidr"
    endpoint_spec      = ["10.0.0.0/8"]
    description        = "Allow PostgreSQL"
  }

  # Allow SSH
  rules {
    direction          = "inbound"
    protocol           = "tcp"
    port_start         = 22
    port_end           = 22
    endpoint_spec_type = "cidr"
    endpoint_spec      = ["0.0.0.0/0"]  # Change to your IP range
    description        = "Allow SSH"
  }

  # Allow all outbound
  rules {
    direction          = "outbound"
    protocol           = "any"
    endpoint_spec_type = "any"
    endpoint_spec      = []
    description        = "Allow all outbound"
  }
}

# ============================================================================
# Virtual Machine Resources
# ============================================================================

# Create web server VMs
resource "idcloudhost_vm" "web_server_1" {
  name               = "${local.project_name}-web-1"
  os_name            = "ubuntu"
  os_version         = "20.04"
  disks              = 20
  vcpu               = 2
  memory             = 2048
  username           = "ubuntu"
  initial_password   = "ChangeMe123!"  # Change this!
  billing_account_id = local.billing_account_id
  backup             = false
}

resource "idcloudhost_vm" "web_server_2" {
  name               = "${local.project_name}-web-2"
  os_name            = "ubuntu"
  os_version         = "20.04"
  disks              = 20
  vcpu               = 2
  memory             = 2048
  username           = "ubuntu"
  initial_password   = "ChangeMe123!"  # Change this!
  billing_account_id = local.billing_account_id
  backup             = false
}

# Create database server VM
resource "idcloudhost_vm" "db_server" {
  name               = "${local.project_name}-db"
  os_name            = "ubuntu"
  os_version         = "20.04"
  disks              = 40
  vcpu               = 4
  memory             = 4096
  username           = "ubuntu"
  initial_password   = "ChangeMe123!"  # Change this!
  billing_account_id = local.billing_account_id
  backup             = true
}

# ============================================================================
# Load Balancer Resources
# ============================================================================

# Create a load balancer for web servers
resource "idcloudhost_loadbalancer" "web_lb" {
  display_name       = "${local.project_name}-web-lb"
  billing_account_id = local.billing_account_id
  network_uuid       = idcloudhost_network.app_network.uuid
  reserve_public_ip  = true

  # Add web servers as targets
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
    target_port = 8080
  }

  # HTTPS forwarding rule
  forwarding_rules {
    source_port = 443
    target_port = 8443
  }
}

# ============================================================================
# Object Storage Resources
# ============================================================================

# Create object storage bucket for application data
resource "idcloudhost_objectstorage" "app_data" {
  name               = "${local.project_name}-app-data"
  billing_account_id = local.billing_account_id
}

# Create object storage bucket for backups
resource "idcloudhost_objectstorage" "backups" {
  name               = "${local.project_name}-backups"
  billing_account_id = local.billing_account_id
}

# Create object storage bucket for static assets
resource "idcloudhost_objectstorage" "static_assets" {
  name               = "${local.project_name}-static"
  billing_account_id = local.billing_account_id
}

# ============================================================================
# Floating IP Resources
# ============================================================================

# Create a floating IP for the database server
resource "idcloudhost_floating_ip" "db_floating_ip" {
  name               = "${local.project_name}-db-fip"
  billing_account_id = local.billing_account_id
}

# ============================================================================
# Outputs
# ============================================================================

output "network_uuid" {
  value       = idcloudhost_network.app_network.uuid
  description = "The UUID of the application network"
}

output "web_firewall_uuid" {
  value       = idcloudhost_firewall.web_firewall.uuid
  description = "The UUID of the web firewall"
}

output "db_firewall_uuid" {
  value       = idcloudhost_firewall.db_firewall.uuid
  description = "The UUID of the database firewall"
}

output "web_server_1_uuid" {
  value       = idcloudhost_vm.web_server_1.uuid
  description = "The UUID of web server 1"
}

output "web_server_1_private_ip" {
  value       = idcloudhost_vm.web_server_1.private_ipv4
  description = "The private IP of web server 1"
}

output "web_server_2_uuid" {
  value       = idcloudhost_vm.web_server_2.uuid
  description = "The UUID of web server 2"
}

output "web_server_2_private_ip" {
  value       = idcloudhost_vm.web_server_2.private_ipv4
  description = "The private IP of web server 2"
}

output "db_server_uuid" {
  value       = idcloudhost_vm.db_server.uuid
  description = "The UUID of the database server"
}

output "db_server_private_ip" {
  value       = idcloudhost_vm.db_server.private_ipv4
  description = "The private IP of the database server"
}

output "load_balancer_uuid" {
  value       = idcloudhost_loadbalancer.web_lb.uuid
  description = "The UUID of the load balancer"
}

output "load_balancer_private_address" {
  value       = idcloudhost_loadbalancer.web_lb.private_address
  description = "The private IP address of the load balancer"
}

output "load_balancer_targets" {
  value       = idcloudhost_loadbalancer.web_lb.targets
  description = "The targets of the load balancer"
}

output "app_data_bucket" {
  value       = idcloudhost_objectstorage.app_data.name
  description = "The name of the app data bucket"
}

output "backups_bucket" {
  value       = idcloudhost_objectstorage.backups.name
  description = "The name of the backups bucket"
}

output "static_assets_bucket" {
  value       = idcloudhost_objectstorage.static_assets.name
  description = "The name of the static assets bucket"
}

output "db_floating_ip_address" {
  value       = idcloudhost_floating_ip.db_floating_ip.address
  description = "The floating IP address for the database server"
}
