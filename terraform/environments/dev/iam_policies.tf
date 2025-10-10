data "aws_iam_policy_document" "dynamodb_event_crud_policy" {
    statement {
        sid = "AllowCRUDDevEventTable"
        effect = "Allow"
        actions = [
            "dynamodb:GetItem",
            "dynamodb:Scan",
            "dynamodb:Query",
            "dynamodb:PutItem",
            "dynamodb:UpdateItem",
            "dynamodb:DeleteItem"
        ]
        resources = concat(
            [module.dev_event_table.arn],
            values(module.dev_event_table.lsi_arns)
        )
    }
    version = "2012-10-17"
}

resource "aws_iam_policy" "dynamodb_event_crud_policy" {
    name = "dev-dynamodb-event-crud-policy"
    policy = data.aws_iam_policy_document.dynamodb_event_crud_policy.json
}