variable "lambda_arn" {
    type = string
    description = "arn of integration target lambda"
}

variable "lambda_invoke_arn" {
    type = string
    description = "invoke arn of integration target lambda"
}

variable "function_name" {
    type = string
    description = "function name of integration target lambda"
}

variable "authorizer_id" {
    type = string
    description = "ID of jwt authorizer lambda"
    default = null
}

variable "apigw_id" {
    type = string
    description = "id of integration target api gateway"
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
}

variable "ws_route_key" {
    description = "The route key for WebSocket APIs (e.g., '$connect', '$disconnect', 'message')."
    type        = string
    default     = null
}