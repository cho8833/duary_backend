################## Lambda Options ##################

variable "function_name" {
  description = "lambda function name"
  type = string
}

variable "env_vars" {
  description = "lambda environment variables"
  type = map(string)
  default = {}
}

variable "file_path" {
  description = "lambda zip file path"
  type = string
}

variable "route_key" {
  type = string
  description = "route key of api, ex) PUT /example"
}

variable "attach_policy_arns_map" {
  description = "lambda additinal policy arns, key is logical name, value is arn"
  type = map(string)
  default = {}
}

variable "apigw_id" {
  type = string
  description = "ID of integration target api gateway"
}

variable "authorizer_id" {
  type = string
  description = "ID of jwt authorizer lambda"
  default = null
}

variable "api_type" {
  description = "type of the API, 'HTTP' or 'WEBSOCKET"

  validation {
    condition = contains(["HTTP", "WEBSOCKET"], var.api_type)
    error_message = "The api_type must be either 'HTTP' or 'WEBSOCKET'"
  }
}

variable "http_method" {
  description = "The HTTP method for the route (e.g., 'GET', 'POST'). Required for HTTP APIs."
  type        = string
  default     = null # WebSocket API에서는 사용하지 않으므로 nullable
}

variable "http_path" {
  description = "The path for the route (e.g., '/events'). Required for HTTP APIs."
  type        = string
  default = null
}

variable "ws_route_key" {
  description = "The route key for WebSocket APIs (e.g., '$connect', '$disconnect', 'message')."
  type        = string
  default     = null
}