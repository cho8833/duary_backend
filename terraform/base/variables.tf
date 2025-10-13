variable "jwt_secret_key" {
  type = string
  description = "secret key for jwt authorizer"
}

variable "kakao_aud" {
  type = string
  description = "audience of kakao auth jwt"
}

variable "oidc_nonce" {
  type = string
  description = "oidc nonce, pair with front end"
}

variable "gcp_android_client_id" {
  type = string
  description = "android client id of google cloud platform"
}

variable "gcp_web_client_id" {
  type = string
  description = "web client id of google cloud platform"
}

variable "gcp_ios_client_id" {
  type = string
  description = "ios client id of google cloud platform"
}

