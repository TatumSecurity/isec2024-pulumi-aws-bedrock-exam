package lambda_manager

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateLambdaManager(ctx *pulumi.Context) (*lambda.Function, *lambda.Function, error) {

	lambdaTrustPolicy, err := iam.GetPolicyDocument(ctx, &iam.GetPolicyDocumentArgs{
		Statements: []iam.GetPolicyDocumentStatement{
			{
				Effect: pulumi.StringRef("Allow"),
				Principals: []iam.GetPolicyDocumentStatementPrincipal{
					{
						Type: "Service",
						Identifiers: []string{
							"lambda.amazonaws.com",
						},
					},
				},
				Actions: []string{
					"sts:AssumeRole",
				},
			},
		},
	})
	if err != nil {
		return nil, nil, err
	}
	lambdaRole, err := iam.NewRole(ctx, "TatumImageSummarizationLambdaRole", &iam.RoleArgs{
		Name:             pulumi.String("TatumImageSummarizationLambdaRole"),
		AssumeRolePolicy: pulumi.String(lambdaTrustPolicy.Json),
	})
	if err != nil {
		ctx.Log.Warn("lambda role creation failed ", nil)
		return nil, nil, err
	}

	filterLabelsArchive := pulumi.NewFileArchive("../../static/lambda/dist/filterLabels.zip")
	buildOutputArchive := pulumi.NewFileArchive("../../static/lambda/dist/buildOutput.zip")

	filterLabelsLambdaFunction, err := lambda.NewFunction(ctx, "TatumImageSummarizationFilterLabels", &lambda.FunctionArgs{
		Code:    filterLabelsArchive,
		Role:    lambdaRole.Arn,
		Runtime: pulumi.String(lambda.RuntimeNodeJS20dX),
		Handler: pulumi.String("filterLabels.handler"),
		Name:    pulumi.String("TatumImageSummarizationFilterLabels"),
		Environment: lambda.FunctionEnvironmentArgs{
			Variables: pulumi.StringMap{
				"CONFIDENCE_LEVEL": pulumi.String("90"),
			},
		},
	})
	if err != nil {
		ctx.Log.Warn("buildOutputLambda creation failed ", nil)
		return nil, nil, err
	}

	buildOutputLambdaFunction, err := lambda.NewFunction(ctx, "TatumImageSummarizationBuildOutput", &lambda.FunctionArgs{
		Code:    buildOutputArchive,
		Role:    lambdaRole.Arn,
		Runtime: pulumi.String(lambda.RuntimeNodeJS20dX),
		Handler: pulumi.String("buildOutput.handler"),
		Name:    pulumi.String("TatumImageSummarizationBuildOutput"),
		Environment: lambda.FunctionEnvironmentArgs{
			Variables: pulumi.StringMap{
				"CONFIDENCE_LEVEL": pulumi.String("90"),
			},
		},
	})
	if err != nil {
		ctx.Log.Warn("buildOutputLambda creation failed ", nil)
		return nil, nil, err
	}
	return buildOutputLambdaFunction, filterLabelsLambdaFunction, nil
}
