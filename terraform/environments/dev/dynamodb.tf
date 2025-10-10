module "dev_event_table" {
    source = "../../modules/dynamodb"
    name = "dev_Event"

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

module "dev_member_table" {
    source = "../../modules/dynamodb"

    name = "dev_Member"

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

module "dev_ws_conneciton_table" {
    source = "../../modules/dynamodb"

    name = "dev_WSConnection"

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

module "dev_cert_table" {
    source = "../../modules/dynamodb"

    name = "dev_Cert"

    hash_key = "provider"

    attributes = [
        {
            name = "provider"
            type = "S"
        }
    ]
}

module "dev_couple_table" {
    source = "../../modules/dynamodb"

    name = "dev_Couple"

    hash_key = "id"

    attributes = [
        {
            name = "id"
            type = "S"
        }
    ]
}