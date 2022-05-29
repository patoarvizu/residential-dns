resource kubernetes_deployment residential_dns {
  metadata {
    name = "residential-dns"
    namespace = var.create_namespace ? kubernetes_namespace.ns["ns"].metadata[0].name : data.kubernetes_namespace.ns["ns"].metadata[0].name
  }
  spec {
    selector {
      match_labels = {
        app = "residential-dns"
      }
    }
    template {
      metadata {
        labels = {
          app = "residential-dns"
        }
        annotations = var.deployment_annotations
      }
      spec {
        service_account_name = kubernetes_service_account.residential_dns.metadata[0].name
        container {
          name = "residential-dns"
          command = [
            "/residential-dns",
          ]
          args = var.deployment_args
          image = format("patoarvizu/residential-dns:%s", var.image_version)
          dynamic "env" {
            for_each = var.env_vars
            content {
              name = env.value["name"]
              value = env.value["value"]
            }
          }
        }
      }
    }
  }
}