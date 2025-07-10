output "supabase_secret_name" {
  value = kubernetes_secret.supabase_credentials.metadata[0].name
}
