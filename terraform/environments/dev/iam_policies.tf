data "aws_iam_policy_document" "dynamodb_dev_event_crud_policy_doc" {
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

resource "aws_iam_policy" "dynamodb_dev_event_crud_policy" {
    name = "dev-dynamodb-event-crud-policy"
    policy = data.aws_iam_policy_document.dynamodb_dev_event_crud_policy_doc.json
}

data "aws_iam_policy_document" "dynamodb_dev_cert_crud_policy_doc" {
    statement {
        sid = "AllowCRUDDevCertTable"
        effect = "Allow"
        actions = [
            "dynamodb:GetItem",
            "dynamodb:Scan",
            "dynamodb:Query",
            "dynamodb:PutItem",
            "dynamodb:UpdateItem",
            "dynamodb:DeleteItem"
        ]

        resources = [
            module.dev_cert_table.arn
        ]

    }

    version = "2012-10-17"
}

resource "aws_iam_policy" "dynamodb_dev_cert_crud_policy" {
    name = "dev-dynamodb-cert-crud-policy"
    policy = data.aws_iam_policy_document.dynamodb_dev_cert_crud_policy_doc.json
}

data "aws_iam_policy_document" "dynamodb_dev_couple_crud_policy_doc" {
    statement {
        sid = "AllowCRUDDevCoupleTable"
        effect = "Allow"
        actions = [
            "dynamodb:GetItem",
            "dynamodb:Scan",
            "dynamodb:Query",
            "dynamodb:PutItem",
            "dynamodb:UpdateItem",
            "dynamodb:DeleteItem"
        ]

        resources = [
            module.dev_couple_table.arn
        ]

    }

    version = "2012-10-17"
}

resource "aws_iam_policy" "dynamodb_dev_couple_crud_policy" {
    name = "dev-dynamodb-couple-crud-policy"
    policy = data.aws_iam_policy_document.dynamodb_dev_couple_crud_policy_doc.json
}

data "aws_iam_policy_document" "dynamodb_dev_member_crud_policy_doc" {
    statement {
        sid = "AllowCRUDDevMemberTable"
        effect = "Allow"
        actions = [
            "dynamodb:GetItem",
            "dynamodb:Scan",
            "dynamodb:Query",
            "dynamodb:PutItem",
            "dynamodb:UpdateItem",
            "dynamodb:DeleteItem"
        ]

        resources = [
            module.dev_member_table.arn
        ]

    }

    version = "2012-10-17"
}

resource "aws_iam_policy" "dynamodb_dev_member_crud_policy" {
    name = "dev-dynamodb-member-crud-policy"
    policy = data.aws_iam_policy_document.dynamodb_dev_member_crud_policy_doc.json
}

data "aws_iam_policy_document" "dynamodb_dev_ws_connection_crud_policy_doc" {
    statement {
        sid = "AllowCRUDDevWSConnectionTable"
        effect = "Allow"
        actions = [
            "dynamodb:GetItem",
            "dynamodb:Scan",
            "dynamodb:Query",
            "dynamodb:PutItem",
            "dynamodb:UpdateItem",
            "dynamodb:DeleteItem"
        ]

        resources = [
            module.dev_ws_conneciton_table.arn
        ]

    }

    version = "2012-10-17"
}

resource "aws_iam_policy" "dynamodb_dev_ws_connection_crud_policy" {
    name = "dev-dynamodb-ws-connection-crud-policy"
    policy = data.aws_iam_policy_document.dynamodb_dev_ws_connection_crud_policy_doc.json
}