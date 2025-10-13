module "lambda" {
  source = "../lambda"

  file_path     = var.file_path
  function_name = var.function_name

  env_vars = var.env_vars

  attach_policy_arns_map = var.attach_policy_arns_map
}

module "apigw_integration" {
  source = "../apigw_lambda_integration"


  apigw_id          = var.apigw_id
  function_name     = var.function_name
  lambda_arn        = module.lambda.function_arn
  lambda_invoke_arn = module.lambda.function_invoke_arn
  api_type = var.api_type
  http_path = var.http_path
  http_method = var.http_method
  ws_route_key = var.ws_route_key
}