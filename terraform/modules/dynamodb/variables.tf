variable "hash_key" {
    type = string
    description = "table partition key"
}

variable "name" {
    description = "table name"
    type = string
}

variable "attributes" {
    description = "table attributes"
    type = list(object({
        name = string
        type = string
    }))
}

variable "lsi" {
    description = "table local secondary index"
    type = list(object({
      name = string
      projection_type = string
      range_key = string
    }))

    default = []
}

variable "range_key" {
    description = "table sort key"
    type = string

    default = null
}

