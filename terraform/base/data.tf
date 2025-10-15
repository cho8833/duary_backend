data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

data "terraform_remote_state" "dev" {
  backend = "s3"

  config = {
    bucket = "duary-terraform"
    key = "dev-terraform.tfstate"
    region = "ap-northeast-2"
  }
}

data "terraform_remote_state" "prod" {
  backend = "s3"

  config = {
    bucket = "duary-terraform"
    key = "prod-terraform.tfstate"
    region = "ap-northeast-2"
  }
}

locals {
  dev_websocket_post_url = "https://${aws_apigatewayv2_api.duary_ws_apigw.id}.execute-api.${data.aws_region.current.region}.amazonaws.com/${aws_apigatewayv2_stage.duary_ws_apigw_dev_stage.name}"
  prod_websocket_post_url = "https://${aws_apigatewayv2_api.duary_ws_apigw.id}.execute-api.${data.aws_region.current.region}.amazonaws.com/${aws_apigatewayv2_stage.duary_ws_apigw_prod_stage.name}"
}