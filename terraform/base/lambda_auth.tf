module "get_oidc_public_key" {
  source = "../modules/lambda"

  file_path = "${var.root_path}/build/package/get_oidc_public_key.zip"
  function_name = "get_oidc_public_key"

  attach_policy_arns_map = {
    "dynamodb_dev_cert_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_cert_crud_policy_arn
    "dynamodb_prod_cert_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_cert_crud_policy_arn
  }
}

module "kakao_sign_in_api" {
  source = "../modules/api_lambda_endpoint"

  file_path     = "${var.root_path}/build/package/kakao_sign_in.zip"
  function_name = "kakao_sign_in_api"

  env_vars = {
    aud = var.kakao_aud
    secretKey = var.jwt_secret_key
    nonce = var.oidc_nonce
  }

  attach_policy_arns_map = {
    "dev_couple_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_couple_crud_policy_arn
    "dev_member_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_member_crud_policy_arn
    "prod_couple_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_couple_crud_policy_arn
    "prod_member_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_member_crud_policy_arn
    "invoke_get_oidc_public_key" = aws_iam_policy.allow_invoke_lambda_get_oidc_public_key_policy.arn
  }

  apigw_id  = aws_apigatewayv2_api.duary_apigw.id
  api_type  = "HTTP"
  http_method = "POST"
  http_path = "/auth/signin/kakao"
}

module "google_sign_in_api" {
  source ="../modules/api_lambda_endpoint"

  function_name = "google_sign_in_api"
  file_path     = "${var.root_path}/build/package/google_sign_in.zip"

  api_type      = "HTTP"
  apigw_id      = aws_apigatewayv2_api.duary_apigw.id
  http_method = "POST"
  http_path = "/auth/signin/google"

  attach_policy_arns_map = {
    "dev_couple_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_couple_crud_policy_arn
    "dev_member_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_member_crud_policy_arn
    "prod_couple_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_couple_crud_policy_arn
    "prod_member_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_member_crud_policy_arn
    "invoke_get_oidc_public_key" = aws_iam_policy.allow_invoke_lambda_get_oidc_public_key_policy.arn
  }

  env_vars = {
    nonce = var.oidc_nonce
    secretKey = var.jwt_secret_key
    ANDROID_CLIENT_ID = var.gcp_android_client_id
    IOS_CLIENT_ID = var.gcp_ios_client_id
    WEB_CLIENT_ID = var.gcp_web_client_id
  }
}

module "apple_sign_in_api" {
  source = "../modules/api_lambda_endpoint"

  function_name = "apple_sign_in_api"
  file_path = "${var.root_path}/build/package/apple_sign_in.zip"

  api_type = "HTTP"
  apigw_id = aws_apigatewayv2_api.duary_apigw.id
  http_method = "POST"
  http_path = "/auth/signin/apple"

  attach_policy_arns_map = {
    "dev_couple_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_couple_crud_policy_arn
    "dev_member_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_member_crud_policy_arn
    "prod_couple_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_couple_crud_policy_arn
    "prod_member_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_member_crud_policy_arn
    "invoke_get_oidc_public_key" = aws_iam_policy.allow_invoke_lambda_get_oidc_public_key_policy.arn
  }

  env_vars = {
    nonce = var.oidc_nonce
    secretKey = var.jwt_secret_key
  }
}

module "sign_out_api" {
  source = "../modules/api_lambda_endpoint"

  function_name = "sign_out_api"
  file_path = "${var.root_path}/build/package/sign_out.zip"

  api_type = "HTTP"
  apigw_id = aws_apigatewayv2_api.duary_apigw.id
  http_method = "POST"
  http_path = "/auth/signout"

  attach_policy_arns_map = {
    "prod_member_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_member_crud_policy_arn
    "dev_member_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_member_crud_policy_arn
  }

  authorizer_id = aws_apigatewayv2_authorizer.http_jwt_authorizer.id
}

module "token_sign_in_api" {
  source = "../modules/api_lambda_endpoint"

  function_name = "token_sign_in_api"
  file_path = "${var.root_path}/build/package/token_sign_in.zip"

  api_type = "HTTP"
  apigw_id = aws_apigatewayv2_api.duary_apigw.id
  http_method = "POST"
  http_path = "/auth/signin/token"

  authorizer_id = aws_apigatewayv2_authorizer.http_jwt_authorizer.id

  env_vars = {
    secretKey = var.jwt_secret_key
  }

  attach_policy_arns_map = {
    "dev_couple_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_couple_crud_policy_arn
    "dev_member_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_member_crud_policy_arn
    "prod_couple_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_couple_crud_policy_arn
    "prod_member_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_member_crud_policy_arn
  }
}

module "reissue_token_api" {
  source = "../modules/api_lambda_endpoint"

  function_name = "reissue_token_api"
  file_path     = "${var.root_path}/build/package/reissue_token.zip"

  api_type      = "HTTP"
  apigw_id      = aws_apigatewayv2_api.duary_apigw.id
  http_method = "POST"
  http_path = "/auth/token"

  env_vars = {
    secretKey = var.jwt_secret_key
  }

  attach_policy_arns_map = {
    "dev_member_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_member_crud_policy_arn
    "prod_member_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_member_crud_policy_arn
  }

}

module "withdrawal_api" {
  source = "../modules/api_lambda_endpoint"

  function_name = "withdrawal_api"
  file_path = "${var.root_path}/build/package/withdrawal.zip"

  api_type = "HTTP"
  apigw_id = aws_apigatewayv2_api.duary_apigw.id
  http_method = "POST"
  http_path = "/auth/withdrawal"

  authorizer_id = aws_apigatewayv2_authorizer.http_jwt_authorizer.id

  attach_policy_arns_map = {
    "dev_couple_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_couple_crud_policy_arn
    "dev_member_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_member_crud_policy_arn
    "dev_event_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_event_crud_policy_arn
    "prod_couple_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_couple_crud_policy_arn
    "prod_member_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_member_crud_policy_arn
    "prod_event_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_event_crud_policy_arn
    "event_bridge_scheduler_crud" = aws_iam_policy.allow_crud_event_bridge_scheduler.arn
  }
}