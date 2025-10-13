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