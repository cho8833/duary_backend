module "prod_event_table" {
  source = "../../modules/dynamodb"
  name = "Event"

  attributes = [
    {
      name = "coupleId"
      type = "S"
    },
    {
      name = "eventType"
      type = "S"
    },
    {
      name = "startDateTime"
      type = "S"
    }
  ]

  hash_key = "coupleId"

  range_key = "startDateTime"

  lsi = [
    {
      name = "eventType-index"
      projection_type = "KEYS_ONLY"
      range_key = "eventType"
    }
  ]
}

module "prod_member_table" {
  source = "../../modules/dynamodb"

  name = "Member"

  hash_key = "socialId"

  range_key = "provider"

  attributes = [
    {
      name = "provider"
      type = "S"
    },
    {
      name = "socialId"
      type = "S"
    }
  ]

}

module "prod_ws_connection_table" {
  source = "../../modules/dynamodb"

  name = "WSConnection"

  hash_key = "socialId"
  range_key = "provider"

  attributes = [
    {
      name = "socialId"
      type = "S"
    },
    {
      name = "provider"
      type = "S"
    }
  ]
}

module "prod_cert_table" {
  source = "../../modules/dynamodb"

  name = "Cert"

  hash_key = "provider"

  attributes = [
    {
      name = "provider"
      type = "S"
    }
  ]
}

module "prod_couple_table" {
  source = "../../modules/dynamodb"

  name = "Couple"

  hash_key = "id"

  attributes = [
    {
      name = "id"
      type = "S"
    }
  ]
}