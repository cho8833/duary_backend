data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

data "terraform_remote_state" "dev" {
  backend = "local"

  config = {
    path = "${path.module}/../environments/dev/terraform.tfstate"
  }
}

data "terraform_remote_state" "prod" {
  backend = "local"

  config = {
    path = "${path.module}/../environments/prod/terraform.tfstate"
  }
}