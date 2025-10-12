# TODO: 앱 도메인 변경 업데이트 시, stage 의 values 에서 dev 지울 것


data "aws_iam_policy_document" "dynamodb_prod_event_crud_policy_doc" {
  statement {
    sid    = "AllowCRUDProdEventTable"
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
      [module.prod_event_table.arn],
      values(module.prod_event_table.lsi_arns)
    )

    condition {
      test     = "StringEquals"
      values = ["dev", "prod"]
      variable = "apigateway:stage"
    }
  }
  version = "2012-10-17"
}

resource "aws_iam_policy" "dynamodb_prod_event_crud_policy" {
  name   = "prod-dynamodb-event-crud-policy"
  policy = data.aws_iam_policy_document.dynamodb_prod_event_crud_policy_doc.json
}

data "aws_iam_policy_document" "dynamodb_prod_cert_crud_policy_doc" {
  statement {
    sid    = "AllowCRUDProdCertTable"
    effect = "Allow"
    actions = [
      "dynamodb:GetItem",
      "dynamodb:Scan",
      "dynamodb:Query",
      "dynamodb:PutItem",
      "dynamodb:UpdateItem",
      "dynamodb:DeleteItem"
    ]

    resources = [module.prod_cert_table.arn]

    condition {
      test     = "StringEquals"
      values = ["dev", "prod"]
      variable = "apigateway:stage"
    }
  }
  version = "2012-10-17"
}

resource "aws_iam_policy" "dynamodb_prod_cert_crud_policy" {
  name = "prod-dynamodb-cert-crud-policy"
  policy = data.aws_iam_policy_document.dynamodb_prod_cert_crud_policy_doc.json
}

data "aws_iam_policy_document" "dynamodb_prod_member_crud_policy_doc" {
  statement {
    sid    = "AllowCRUDProdMemberTable"
    effect = "Allow"
    actions = [
      "dynamodb:GetItem",
      "dynamodb:Scan",
      "dynamodb:Query",
      "dynamodb:PutItem",
      "dynamodb:UpdateItem",
      "dynamodb:DeleteItem"
    ]

    resources = [module.prod_member_table.arn]

    condition {
      test     = "StringEquals"
      values = ["dev", "prod"]
      variable = "apigateway:stage"
    }
  }
  version = "2012-10-17"
}

resource "aws_iam_policy" "dynamodb_prod_member_crud_policy" {
  name = "prod-dynamodb-member-crud-policy"
  policy = data.aws_iam_policy_document.dynamodb_prod_member_crud_policy_doc.json
}

data "aws_iam_policy_document" "dynamodb_prod_couple_crud_policy_doc" {
  statement {
    sid    = "AllowCRUDProdCoupleTable"
    effect = "Allow"
    actions = [
      "dynamodb:GetItem",
      "dynamodb:Scan",
      "dynamodb:Query",
      "dynamodb:PutItem",
      "dynamodb:UpdateItem",
      "dynamodb:DeleteItem"
    ]

    resources = [module.prod_couple_table.arn]

    condition {
      test     = "StringEquals"
      values = ["dev", "prod"]
      variable = "apigateway:stage"
    }
  }
  version = "2012-10-17"
}

resource "aws_iam_policy" "dynamodb_prod_couple_crud_policy" {
  name = "prod-dynamodb-couple-crud-policy"
  policy = data.aws_iam_policy_document.dynamodb_prod_couple_crud_policy_doc.json
}

data "aws_iam_policy_document" "dynamodb_prod_ws_connection_crud_policy_doc" {
  statement {
    sid    = "AllowCRUDProdWSConectionTable"
    effect = "Allow"
    actions = [
      "dynamodb:GetItem",
      "dynamodb:Scan",
      "dynamodb:Query",
      "dynamodb:PutItem",
      "dynamodb:UpdateItem",
      "dynamodb:DeleteItem"
    ]

    resources = [module.prod_ws_connection_table.arn]

    condition {
      test     = "StringEquals"
      values = ["dev", "prod"]
      variable = "apigateway:stage"
    }
  }
  version = "2012-10-17"
}

resource "aws_iam_policy" "dynamodb_prod_ws_connection_crud_policy" {
  name = "prod-dynamodb-ws-connection-crud-policy"
  policy = data.aws_iam_policy_document.dynamodb_prod_ws_connection_crud_policy_doc.json
}