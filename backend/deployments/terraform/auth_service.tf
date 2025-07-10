resource "kubernetes_deployment" "auth" {
  metadata {
    name      = "auth-service"
    namespace = kubernetes_namespace.auth.metadata[0].name
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "auth-service"
      }
    }

    template {
      metadata {
        labels = {
          app = "auth-service"
        }
      }

      spec {
        container {
          name  = "auth"
          image = "adityawaradkar/gratia-auth:latest"

          env {
            name = "SUPABASE_URL"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.supabase_credentials.metadata[0].name
                key  = "SUPABASE_URL"
              }
            }
          }

          env {
            name = "SUPABASE_ANON_KEY"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.supabase_credentials.metadata[0].name
                key  = "SUPABASE_ANON_KEY"
              }
            }
          }

          env {
            name = "DB_NAME"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.supabase_credentials.metadata[0].name
                key  = "DB_NAME"
              }
            }
          }

          port {
            container_port = 8081
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "auth" {
  metadata {
    name      = "auth-service"
    namespace = kubernetes_namespace.auth.metadata[0].name
  }

  spec {
    selector = {
      app = "auth-service"
    }

    port {
      port        = 80
      target_port = 8081
    }

    type = "ClusterIP"
  }
}
