package main

import (
	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/awslambdago"
	"github.com/aws/aws-cdk-go/awscdk/awss3assets"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type Props struct {
	awscdk.StackProps
}

func NewStack(scope constructs.Construct, id string, props *Props) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// wkhtmltopdf layer, from:
	// https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6-4/wkhtmltox-0.12.6-4.amazonlinux2_lambda.zip
	wkhtmltopdfLayer := awslambda.NewLayerVersion(stack, jsii.String("wkhtmltopdfLayer"), &awslambda.LayerVersionProps{
		Code:                    awslambda.AssetCode_FromAsset(jsii.String("../wkhtmltox-0.12.6-4.amazonlinux2_lambda"), &awss3assets.AssetOptions{}),
		CompatibleArchitectures: &[]awslambda.Architecture{awslambda.Architecture_X86_64()},
	})

	// POST /documents handler.
	documentsPost := awslambdago.NewGoFunction(stack, jsii.String("documentsPost"), &awslambdago.GoFunctionProps{
		MemorySize: jsii.Number(1024),
		Tracing:    awslambda.Tracing_ACTIVE,
		Layers:     &[]awslambda.ILayerVersion{wkhtmltopdfLayer},
		Entry:      jsii.String("../api/documents/post"),
		Timeout:    awscdk.Duration_Seconds(jsii.Number(30)),
		Environment: &map[string]*string{
			"FONTCONFIG_PATH": jsii.String("/opt/fonts"),
		},
	})

	// Create API Gateway.
	api := awsapigatewayv2.NewHttpApi(stack, jsii.String("documentsApi"), &awsapigatewayv2.HttpApiProps{})

	// POST /documents
	api.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path: jsii.String("/documents"),
		Integration: awsapigatewayv2integrations.NewLambdaProxyIntegration(&awsapigatewayv2integrations.LambdaProxyIntegrationProps{
			Handler:              documentsPost,
			PayloadFormatVersion: awsapigatewayv2.PayloadFormatVersion_VERSION_2_0(),
		}),
	})

	// Output the endpoint address.
	awscdk.NewCfnOutput(stack, jsii.String("documentApiEndpoint"), &awscdk.CfnOutputProps{
		Value: api.Url(),
	})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)
	NewStack(app, "wkhtmltopdfLambdaStack", &Props{
		awscdk.StackProps{
			Env: env(),
		},
	})
	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	return nil
}
