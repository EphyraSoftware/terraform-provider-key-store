resource "tls_private_key" "example" {
  algorithm = "RSA"
  rsa_bits = 2048
}

resource "tls_self_signed_cert" "example_certificate" {
  key_algorithm   = "${tls_private_key.example.algorithm}"
  private_key_pem = "${tls_private_key.example.private_key_pem}"

  subject {
    common_name  = "example.com"
    organization = "EphyraSoftware"
  }

  dns_names = ["example.com"]

  validity_period_hours = 2190 // Three months

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth"
  ]
}

resource "keystore_pkcs12_bundle" "my-bundle" {
  name = "my-bundle-name"
  cert_pem = "${tls_self_signed_cert.example_certificate.cert_pem}"
  key_pem = "${tls_private_key.example.private_key_pem}"
}

resource "local_file" "my-bundle" {
    content_base64 = "${keystore_pkcs12_bundle.my-bundle.bundle}"
    filename       = "${path.module}/my-bundle.p12"
}
