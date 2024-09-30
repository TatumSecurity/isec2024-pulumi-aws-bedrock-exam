package main

import (
	s3_manager "pulumi-cloud-ai-exam/resources/s3"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		err := s3_manager.CreateS3Manager(ctx)
		if err != nil {
			ctx.Log.Error("Failed to create S3 manager: "+err.Error(), nil)
			return err
		}
		// buildOutputLambdaFunction, filterLabelsLambdaFunction, err := lambda_manager.CreateLambdaManager(ctx)
		// if err != nil {
		// 	ctx.Log.Error("Failed to create lambda manager: "+err.Error(), nil)
		// 	return err
		// }

		// err = sfn_manager.CreateSfnManager(ctx, buildOutputLambdaFunction, filterLabelsLambdaFunction)
		// if err != nil {
		// 	ctx.Log.Error("Failed to create sfn manager: "+err.Error(), nil)
		// 	return err
		// }
		return nil
	})
}
