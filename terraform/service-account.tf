resource kubernetes_service_account residential_dns {
  metadata {
    name = "residential-dns"
    namespace = var.create_namespace ? kubernetes_namespace.ns["ns"].metadata[0].name : data.kubernetes_namespace.ns["ns"].metadata[0].name
    annotations = var.service_account_annotations
  }
}