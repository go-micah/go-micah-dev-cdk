package main

import (
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	"github.com/aws/jsii-runtime-go"
)

func TestGoMicahDevStack(t *testing.T) {
	// GIVEN
	app := awscdk.NewApp(nil)

	// WHEN
	stack := NewGoMicahDevStack(app, "MyStack", nil)

	// THEN
	template := assertions.Template_FromStack(stack, nil)

	template.HasResourceProperties(jsii.String("AWS::SQS::Queue"), map[string]interface{}{
		"VisibilityTimeout": 300,
	})
}

func TestGoMicahDevStackVPC(t *testing.T) {
	app := awscdk.NewApp(nil)

	stack := NewGoMicahDevStack(app, "MyStack", nil)

	template := assertions.Template_FromStack(stack, nil)

	template.HasResourceProperties(jsii.String("AWS::EC2::VPC"), map[string]interface{}{
		"CidrBlock": "10.0.0.0/16",
	})
}
