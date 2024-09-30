package sfn_manager

import (
	"encoding/json"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/cloudwatch"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/sfn"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func CreateSfnManager(ctx *pulumi.Context, buildOutputLambdaFunction *lambda.Function, filterLabelsLambdaFunction *lambda.Function) error {
	var (
		sfnInfo SfnInfoInterface
	)
	cfg := config.New(ctx, "")
	if err := cfg.TryObject("pulumiAWSSfn", &sfnInfo); err != nil {
		ctx.Log.Error("pulumiAWSSfn Not Defined", nil)
		return err
	}

	stateMachineTrustPolicy, err := iam.GetPolicyDocument(ctx, &iam.GetPolicyDocumentArgs{
		Statements: []iam.GetPolicyDocumentStatement{
			{
				Effect: pulumi.StringRef("Allow"),
				Principals: []iam.GetPolicyDocumentStatementPrincipal{
					{
						Type: "Service",
						Identifiers: []string{
							"states.amazonaws.com",
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
		ctx.Log.Error("stateMachineTrustPolicy creation failed ", nil)
		return err
	}
	lambdaArns := pulumi.StringArray{
		pulumi.Sprintf("%s:%s", buildOutputLambdaFunction.Arn, "$LATEST"),
		pulumi.Sprintf("%s:%s", filterLabelsLambdaFunction.Arn, "$LATEST"),
	}
	s3Arns := pulumi.StringArray{
		pulumi.Sprintf("arn:aws:s3:::%s/*", sfnInfo.InputBucket),
		pulumi.Sprintf("arn:aws:s3:::%s/*", sfnInfo.OutputBucket),
	}

	stateMachinePolicyData := map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Effect":   "Allow",
				"Action":   []string{"lambda:InvokeFunction"},
				"Resource": lambdaArns,
			},
			{
				"Effect":   "Allow",
				"Action":   []string{"s3:GetObject", "s3:DeleteObject", "s3:PutObject"},
				"Resource": s3Arns,
			},
			{
				"Effect":   "Allow",
				"Action":   []string{"rekognition:DetectLabels"},
				"Resource": "*",
			},
			{
				"Effect":   "Allow",
				"Action":   []string{"bedrock:InvokeModel"},
				"Resource": "*",
			},
		},
	}

	stateMachinePolicy, err := iam.NewPolicy(ctx, "TatumImageSummarizationSfnPolicy", &iam.PolicyArgs{
		Name:        pulumi.String("TatumImageSummarizationSfnPolicy"),
		Path:        pulumi.String("/"),
		Description: pulumi.String("Permission policy for Image Summarization state machine"),
		Policy:      pulumi.JSONMarshal(stateMachinePolicyData),
	})
	if err != nil {
		ctx.Log.Error("stateMachinePolicy creation failed ", nil)
		return err
	}

	stateMachineRole, err := iam.NewRole(ctx, "TatumImageSummarizationSfnRole", &iam.RoleArgs{
		Name:              pulumi.String("TatumImageSummarizationSfnRole"),
		AssumeRolePolicy:  pulumi.String(stateMachineTrustPolicy.Json),
		ManagedPolicyArns: pulumi.StringArray{stateMachinePolicy.Arn},
	})
	if err != nil {
		ctx.Log.Error("stateMachineRole creation failed ", nil)
		return err
	}

	stateMachineDefinition := map[string]interface{}{
		"Comment": "A description of my state machine",
		"StartAt": "Detect Labels",
		"States": map[string]interface{}{
			"Detect Labels": map[string]interface{}{
				"Type": "Task",
				"Parameters": map[string]interface{}{
					"Image": map[string]interface{}{
						"S3Object": map[string]interface{}{
							"Bucket.$": "$.detail.bucket.name",
							"Name.$":   "$.detail.object.key",
						},
					},
				},
				"Resource":   "arn:aws:states:::aws-sdk:rekognition:detectLabels",
				"Next":       "Filter Labels",
				"ResultPath": "$.Rekognition",
				"ResultSelector": map[string]interface{}{
					"Labels.$": "$.Labels",
				},
				"Comment": "Uses Rekognition to detect the labels in the image. Combines it's input and output into a single state.",
			},
			"Filter Labels": map[string]interface{}{
				"Type":     "Task",
				"Resource": "arn:aws:states:::lambda:invoke",
				"Parameters": map[string]interface{}{
					"Payload.$":    "$",
					"FunctionName": pulumi.Sprintf("%s:$LATEST", filterLabelsLambdaFunction.Arn),
				},
				"Comment": "Filters labels based on confidence and provides a count of occurencer per-label",
				"Retry": []map[string]interface{}{
					{
						"ErrorEquals":     []string{"Lambda.ServiceException", "Lambda.AWSLambdaException", "Lambda.SdkClientException", "Lambda.TooManyRequestsException"},
						"IntervalSeconds": 1,
						"MaxAttempts":     3,
						"BackoffRate":     2,
					},
				},
				"ResultPath": "$.Lambda",
				"ResultSelector": map[string]interface{}{
					"FilteredLabels.$": "$.Payload.labels",
				},
				"Next": "Bedrock InvokeModel",
			},
			"Bedrock InvokeModel": map[string]interface{}{
				"Type":     "Task",
				"Resource": "arn:aws:states:::bedrock:invokeModel",
				"Parameters": map[string]interface{}{
					"ModelId": "arn:aws:bedrock:us-east-1::foundation-model/amazon.titan-text-premier-v1:0",
					"Body": map[string]interface{}{
						"inputText.$": "States.Format('Human: Here is a comma seperated list of labels/objects seen in an image\n<labels>{}</labels>\n\nPlease provide a human readible and understandable summary based on these labels\n\nAssistant:', $.Lambda.FilteredLabels)",
						"textGenerationConfig": map[string]interface{}{
							"temperature":   0.7,
							"topP":          0.9,
							"maxTokenCount": 512,
						},
					},
				},
				"ResultPath": "$.Bedrock",
				"Next":       "Build Output",
			},
			"Build Output": map[string]interface{}{
				"Type":       "Task",
				"Resource":   "arn:aws:states:::lambda:invoke",
				"OutputPath": "$.Payload",
				"Parameters": map[string]interface{}{
					"Payload.$":    "$",
					"FunctionName": pulumi.Sprintf("%s:$LATEST", buildOutputLambdaFunction.Arn),
				},
				"Retry": []map[string]interface{}{
					{
						"ErrorEquals":     []string{"Lambda.ServiceException", "Lambda.AWSLambdaException", "Lambda.SdkClientException", "Lambda.TooManyRequestsException"},
						"IntervalSeconds": 1,
						"MaxAttempts":     3,
						"BackoffRate":     2,
					},
				},
				"Next": "Save Output",
			},
			"Save Output": map[string]interface{}{
				"Type": "Task",
				"End":  true,
				"Parameters": map[string]interface{}{
					"Body.$": "$",
					"Bucket": sfnInfo.OutputBucket,
					"Key.$":  "States.Format('{}.json', $.source.file)",
				},
				"Resource": "arn:aws:states:::aws-sdk:s3:putObject",
			},
		},
	}

	stateMachine, err := sfn.NewStateMachine(ctx, "TatumImageSummarizationStateMachine", &sfn.StateMachineArgs{
		Name:       pulumi.String("TatumImageSummarizationStateMachine"),
		RoleArn:    stateMachineRole.Arn,
		Definition: pulumi.JSONMarshal(stateMachineDefinition),
	})
	if err != nil {
		ctx.Log.Error("stateMachine creation failed ", nil)
		return err
	}

	inputRuleTrustPolicy, err := iam.GetPolicyDocument(ctx, &iam.GetPolicyDocumentArgs{
		Statements: []iam.GetPolicyDocumentStatement{
			{
				Effect: pulumi.StringRef("Allow"),
				Principals: []iam.GetPolicyDocumentStatementPrincipal{
					{
						Type: "Service",
						Identifiers: []string{
							"events.amazonaws.com",
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
		ctx.Log.Error("inputRuleTrustPolicy creation failed ", nil)
		return err
	}

	stateMachineArns := pulumi.StringArray{
		stateMachine.Arn,
	}

	ImageSummarizationRulePolicyData := map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Effect":   "Allow",
				"Action":   []string{"states:StartExecution"},
				"Resource": stateMachineArns,
			},
		},
	}

	inputRulePolicy, err := iam.NewPolicy(ctx, "TatumImageSummarizationRulePolicy", &iam.PolicyArgs{
		Name:   pulumi.String("TatumImageSummarizationRulePolicy"),
		Policy: pulumi.JSONMarshal(ImageSummarizationRulePolicyData),
	})
	if err != nil {
		ctx.Log.Error("inputRulePolicy creation failed ", nil)
		return err
	}

	inputRuleRole, err := iam.NewRole(ctx, "TatumImageSummarizationRuleRole", &iam.RoleArgs{
		Name:              pulumi.String("TatumImageSummarizationRuleRole"),
		AssumeRolePolicy:  pulumi.String(inputRuleTrustPolicy.Json),
		ManagedPolicyArns: pulumi.StringArray{inputRulePolicy.Arn},
	})
	if err != nil {
		ctx.Log.Error("inputRuleRole creation failed ", nil)
		return err
	}

	eventPattern, err := json.Marshal(map[string]interface{}{
		"source":      []string{"aws.s3"},
		"detail-type": []string{"Object Created"},
		"detail": map[string]interface{}{
			"bucket": map[string]interface{}{
				"name": []string{
					sfnInfo.InputBucket,
				},
			},
		},
	})
	if err != nil {
		ctx.Log.Error("eventPattern creation failed ", nil)
		return err
	}

	inputRule, err := cloudwatch.NewEventRule(ctx, "TatumInputBucketRule", &cloudwatch.EventRuleArgs{
		Name:         pulumi.String("TatumInputBucketRule"),
		EventPattern: pulumi.String(eventPattern),
		ForceDestroy: pulumi.Bool(true),
	})
	if err != nil {
		return err
	}

	_, err = cloudwatch.NewEventTarget(ctx, "TatumInputRuleTarget", &cloudwatch.EventTargetArgs{
		TargetId: pulumi.String("TatumInputRuleTarget"),
		Rule:     inputRule.Name,
		Arn:      stateMachine.Arn,
		RoleArn:  inputRuleRole.Arn,
	})

	if err != nil {
		ctx.Log.Error("role creation failed ", nil)
		return err
	}

	return nil
}
