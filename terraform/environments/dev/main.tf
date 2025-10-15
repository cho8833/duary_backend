terraform {
  required_providers {
    aws = {
      version = "~> 6.16.0"
    }
  }
}

terraform {
  backend "s3" {
    bucket = "duary-terraform"
    key = "dev-terraform.tfstate"
    region = "ap-northeast-2"
    encrypt = true
    acl = "private"
  }
}