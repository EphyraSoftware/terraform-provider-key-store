
provider "keystore" {
  path = "${path.module}/../out"
}

data "keystore_pkcs12_bundle" "my-bundle" {
  name = "my-bundle-name"
}

output "my-bundle" {
  value = "${data.keystore_pkcs12_bundle.my-bundle.bundle}"
}
