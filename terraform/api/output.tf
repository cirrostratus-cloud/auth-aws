output api_endpoint {
  value       = aws_apigatewayv2_stage.api.invoke_url
  description = "The API endpoint"
}

output api_name {
  value       = aws_apigatewayv2_api.api.name
  description = "The API name"
}
