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

# Create a firewall with multiple rules
resource "idcloudhost_firewall" "web_firewall" {
  display_name        = "web-server-firewall"
  billing_account_id  = 1200132376  # Replace with your billing account ID
  description         = "Firewall for web servers"

  # Allow HTTP inbound traffic
  rules {
    direction          = "inbound"
    protocol           = "tcp"
    port_start         = 80
    port_end           = 80
    endpoint_spec_type = "any"
    endpoint_spec      = []
    description        = "Allow HTTP traffic"
  }

  # Allow HTTPS inbound traffic
  rules {
    direction          = "inbound"
    protocol           = "tcp"
    port_start         = 443
    port_end           = 443
    endpoint_spec_type = "any"
    endpoint_spec      = []
    description        = "Allow HTTPS traffic"
  }

  # Allow SSH from specific CIDR
  rules {
    direction          = "inbound"
    protocol           = "tcp"
    port_start         = 22
    port_end           = 22
    endpoint_spec_type = "ip_prefixes"
    endpoint_spec      = ["203.0.113.0/24"]
    description        = "Allow SSH from office network"
  }

  # Allow all outbound traffic
  rules {
    direction          = "outbound"
    protocol           = "any"
    endpoint_spec_type = "any"
    endpoint_spec      = []
    description        = "Allow all outbound traffic"
  }
}

# Create a database firewall
resource "idcloudhost_firewall" "database_firewall" {
  display_name        = "database-firewall"
  billing_account_id  = 1200132376  # Replace with your billing account ID
  description         = "Firewall for database servers"

  # Allow PostgreSQL from specific CIDR
  rules {
    direction          = "inbound"
    protocol           = "tcp"
    port_start         = 5432
    port_end           = 5432
    endpoint_spec_type = "ip_prefixes"
    endpoint_spec      = ["10.0.1.0/24"]
    description        = "Allow PostgreSQL from app network"
  }

  # Allow MySQL from specific CIDR
  rules {
    direction          = "inbound"
    protocol           = "tcp"
    port_start         = 3306
    port_end           = 3306
    endpoint_spec_type = "ip_prefixes"
    endpoint_spec      = ["10.0.1.0/24"]
    description        = "Allow MySQL from app network"
  }

  # Allow ICMP ping
  rules {
    direction          = "inbound"
    protocol           = "icmp"
    endpoint_spec_type = "any"
    endpoint_spec      = []
    description        = "Allow ICMP ping"
  }
}

# Output firewall details
output "web_firewall_uuid" {
  value       = idcloudhost_firewall.web_firewall.uuid
  description = "The UUID of the web firewall"
}

output "web_firewall_id" {
  value       = idcloudhost_firewall.web_firewall.id
  description = "The ID of the web firewall"
}

output "database_firewall_uuid" {
  value       = idcloudhost_firewall.database_firewall.uuid
  description = "The UUID of the database firewall"
}
