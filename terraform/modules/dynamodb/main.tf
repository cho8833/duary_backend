resource "aws_dynamodb_table" "this" {
    name = var.name

    hash_key = var.hash_key

    range_key = var.range_key

    dynamic "attribute" {
      for_each = var.attributes

      content {
        name = attribute.value.name
        type = attribute.value.type
      }
    }

    dynamic "local_secondary_index" {
        for_each = var.lsi

        content {
            name = local_secondary_index.value.name
            projection_type = local_secondary_index.value.projection_type
            range_key = local_secondary_index.value.range_key
        }
    }

    point_in_time_recovery {
      enabled = "false"
      recovery_period_in_days = "1"
    }

    read_capacity = "1"

    region = "ap-northeast-2"

    stream_enabled = "false"

    table_class = "STANDARD"

    write_capacity = "1"

    billing_mode = "PROVISIONED"

    deletion_protection_enabled = "false"
}