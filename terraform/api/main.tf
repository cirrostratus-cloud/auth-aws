terraform {
  backend "s3" {}
}

locals {
  api_gateway_name = "${var.function_name}-api-gateway"
}

resource "aws_apigatewayv2_api" "api" {
  name = local.api_gateway_name
  protocol_type = "HTTP"
  tags = var.common_tags
}

resource "aws_apigatewayv2_stage" "api" {
  api_id = aws_apigatewayv2_api.api.id
  name = var.stage_name
  auto_deploy = true
}

resource "aws_apigatewayv2_integration" "api" {
  api_id = aws_apigatewayv2_api.api.id
  integration_uri = var.invoke_arn
  integration_type = "AWS_PROXY"
}

resource "aws_apigatewayv2_route" "api" {
  count = length(var.routes)
  api_id = aws_apigatewayv2_api.api.id
  route_key = "${var.routes[count.index].method} ${var.routes[count.index].path}"
  target    = "integrations/${aws_apigatewayv2_integration.api.id}"
}

resource "aws_lambda_permission" "api" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action = "lambda:InvokeFunction"
  function_name = var.function_name
  principal = "apigateway.amazonaws.com"
  source_arn = "${aws_apigatewayv2_api.api.execution_arn}/*/*"
}