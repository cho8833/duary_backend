module "cert_table" {
  source = "../modules/dynamodb"

  name = "Cert"

  hash_key = "provider"

  attributes = [
    {
      name = "provider"
      type = "S"
    }
  ]
}