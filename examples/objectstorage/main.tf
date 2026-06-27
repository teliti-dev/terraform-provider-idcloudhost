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

# Create an object storage bucket for application data
resource "idcloudhost_objectstorage" "app_bucket" {
  name               = "my-app-data-bucket"
  billing_account_id = 1200132376  # Replace with your billing account ID
}

# Create an object storage bucket for backups
resource "idcloudhost_objectstorage" "backup_bucket" {
  name               = "my-backup-bucket"
  billing_account_id = 1200132376  # Replace with your billing account ID
}

# Create an object storage bucket for static assets
resource "idcloudhost_objectstorage" "static_assets_bucket" {
  name               = "my-static-assets"
  billing_account_id = 1200132376  # Replace with your billing account ID
}

# Output bucket details
output "app_bucket_name" {
  value       = idcloudhost_objectstorage.app_bucket.name
  description = "The name of the application data bucket"
}

output "app_bucket_owner" {
  value       = idcloudhost_objectstorage.app_bucket.owner
  description = "The owner of the application data bucket"
}

output "app_bucket_size" {
  value       = idcloudhost_objectstorage.app_bucket.size_bytes
  description = "The size of the application data bucket in bytes"
}

output "backup_bucket_name" {
  value       = idcloudhost_objectstorage.backup_bucket.name
  description = "The name of the backup bucket"
}

output "static_assets_bucket_name" {
  value       = idcloudhost_objectstorage.static_assets_bucket.name
  description = "The name of the static assets bucket"
}

output "static_assets_num_objects" {
  value       = idcloudhost_objectstorage.static_assets_bucket.num_objects
  description = "The number of objects in the static assets bucket"
}
