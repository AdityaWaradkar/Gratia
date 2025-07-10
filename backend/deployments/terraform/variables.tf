variable "kubeconfig_path" {
  type        = string
  default     = "~/.kube/config"
}

variable "namespace_auth" {
  type    = string
  default = "auth"
}

variable "namespace_monitoring" {
  type    = string
  default = "monitoring"
}

variable "supabase_url" {
  type    = string
  default = "https://uxtvvrzqygvjmpgqlcgz.supabase.co"
}

variable "supabase_anon_key" {
  type    = string
  default = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InV4dHZ2cnpxeWd2am1wZ3FsY2d6Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3NTE4Njg2MDgsImV4cCI6MjA2NzQ0NDYwOH0.859r4hJ86T47QyjI0peuFIf5zy91mJ9-Qjf-2as5iQc"
}

variable "db_name" {
  type    = string
  default = "gratia-database"
}

variable "nginx_ingress_version" {
  type    = string
  default = "4.12.3"
}

variable "prometheus_version" {
  type    = string
  default = "56.6.1"
}
