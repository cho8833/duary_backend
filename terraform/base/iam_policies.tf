data "aws_iam_policy_document" "allow_invoke_lambda_get_oidc_public_key_policy_doc" {
  statement {
    sid ="AllowInvokeLambdaGetOidcPublicKey"
    effect = "Allow"
    actions = [
      "lambda:InvokeFunction"
    ]

    resources = [
      module.get_oidc_public_key.function_arn
    ]
  }

  version = "2012-10-17"
}

resource "aws_iam_policy" "allow_invoke_lambda_get_oidc_public_key_policy" {
  policy = data.aws_iam_policy_document.allow_invoke_lambda_get_oidc_public_key_policy_doc.json
  name = "invoke-lambda-get-oidc-public-key-policy"
}

data "aws_iam_policy_document" "allow_crud_event_bridge_scheduler_doc" {
  statement {
    sid = "AllowCRUDEventBridgeScheduler"
    effect = "Allow"
    actions = [
      "scheduler:GetSchedule",
      "scheduler:UpdateSchedule",
      "scheduler:CreateSchedule",
      "scheduler:DeleteSchedule"
    ]

    resources = [
      "arn:aws:scheduler:${data.aws_region.current.region}:${data.aws_caller_identity.current.account_id}:schedule/*/*"
    ]
  }

  version = "2012-10-17"
}

resource "aws_iam_policy" "allow_crud_event_bridge_scheduler" {
  policy = data.aws_iam_policy_document.allow_crud_event_bridge_scheduler_doc.json
  name = "full-access-event-bridge-scheduler"
}

data "aws_iam_policy_document" "allow_execute_ws_api_doc" {
  statement {
    sid = "AllowExecuteWSAPI"
    effect = "Allow"
    actions = [
      "apigateway:*",
      "lambda:InvokeFunction",
      "lambda:InvokeFunctionUrl",
      "execute-api:Invoke",
      "execute-api:ManageConnections"
    ]
    resources = [
      "*"
    ]
  }

  version = "2012-10-17"
}

resource "aws_iam_policy" "allow_execute_ws_api" {
  name = "AllowExecuteWSAPI"
  policy = data.aws_iam_policy_document.allow_execute_ws_api_doc.json
}