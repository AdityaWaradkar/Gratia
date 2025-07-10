resource "kubernetes_secret" "supabase_credentials" {
  metadata {
    name      = "supabase-credentials"
    namespace = kubernetes_namespace.auth.metadata[0].name
  }

  data = {
    SUPABASE_URL      = base64encode(var.supabase_url)
    SUPABASE_ANON_KEY = base64encode(var.supabase_anon_key)
    DB_NAME           = base64encode(var.db_name)
  }

  type = "Opaque"
}
