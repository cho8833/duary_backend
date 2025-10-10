output "arn" {
    description = "arn of created dynamodb table"
    value = aws_dynamodb_table.this.arn
}



output "lsi_arns" {
    description = "arn of created dynamodb LSIs"

    value = {
        for lsi in aws_dynamodb_table.this.local_secondary_index : lsi.name => "${aws_dynamodb_table.this.arn}/index/${lsi.name}"
    }
}