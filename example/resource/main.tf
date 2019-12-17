resource "tls_private_key" "example-ca" {
  algorithm = "RSA"
  rsa_bits = 2048
}

resource "tls_self_signed_cert" "example-ca" {
  key_algorithm   = tls_private_key.example-ca.algorithm
  private_key_pem = tls_private_key.example-ca.private_key_pem

  subject {
    common_name  = "example CA"
    organization = "EphyraSoftware"
  }

  dns_names = ["exampleca"]

  validity_period_hours = 2190 // Three months

  is_ca_certificate = true

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth"
  ]
}

resource "tls_private_key" "example" {
  algorithm = "RSA"
  rsa_bits = 2048
}

resource "tls_cert_request" "example" {
  key_algorithm   = tls_private_key.example-ca.algorithm
  private_key_pem = tls_private_key.example.private_key_pem

  subject {
    common_name  = "example.com"
    organization = "ACME Examples, Inc"
  }
}

resource "tls_locally_signed_cert" "example" {
  cert_request_pem   = tls_cert_request.example.cert_request_pem
  ca_key_algorithm   = tls_private_key.example-ca.algorithm
  ca_private_key_pem = tls_private_key.example-ca.private_key_pem
  ca_cert_pem        = tls_self_signed_cert.example-ca.cert_pem

  validity_period_hours = 12

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]
}

provider "keystore" {
  path = "${path.module}/../out"
}

output "cert_pem" {
  value = tls_locally_signed_cert.example.cert_pem
}

output "key_pem" {
  value = tls_private_key.example.private_key_pem
}

output "ca_certs" {
  value = tls_self_signed_cert.example-ca.cert_pem
}

resource "keystore_pkcs12_bundle" "my-bundle" {
  name = "my-bundle-name"
  cert_pem = tls_locally_signed_cert.example.cert_pem
  key_pem = tls_private_key.example.private_key_pem
  ca_certs = [
    tls_self_signed_cert.example-ca.cert_pem
  ]
}
