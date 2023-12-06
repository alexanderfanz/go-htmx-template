package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"

	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type MainStackProps struct {
	awscdk.StackProps
}

func MainStack(scope constructs.Construct, id string, props *MainStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	lambdaHandler := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("GetHandler"), &awscdklambdagoalpha.GoFunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../src"),
		Bundling: &awscdklambdagoalpha.BundlingOptions{
			GoBuildFlags: jsii.Strings(`-ldflags "-s -w" -tags lambda.norpc`),
			// Environment: map[string]*string{
			// 	"HELLO": jsii.String("WORLD"),
			// },
		},
		LogRetention: awslogs.RetentionDays_THREE_MONTHS,
	})

	mainApi := awsapigateway.NewRestApi(stack, jsii.String("GoHtmxApi"), &awsapigateway.RestApiProps{
		RestApiName: jsii.String("GoHtmxApi"),
	})

	mainApi.Root().
		AddMethod(
			jsii.String("ANY"),
			awsapigateway.NewLambdaIntegration(lambdaHandler, &awsapigateway.LambdaIntegrationOptions{}),
			&awsapigateway.MethodOptions{},
		)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	mainStack := MainStack(app, "GoHtmxStack", &MainStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	awscdk.Tags_Of(mainStack).Add(jsii.String("Project"), jsii.String("GoHtmx"), &awscdk.TagProps{Priority: jsii.Number(90)})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
