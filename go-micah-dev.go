package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type GoMicahDevStackProps struct {
	awscdk.StackProps
}

func NewGoMicahDevStack(scope constructs.Construct, id string, props *GoMicahDevStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// This defines a vpc with three private and three public subnets and one nat gateway
	// NOTE: there is a running cost for the NAT GW of $0.045 per hour in us-east-1
	vpc := awsec2.NewVpc(stack, jsii.String("GoMicahDevVPC"), &awsec2.VpcProps{
		VpcName:     jsii.String("micadev"),
		IpAddresses: awsec2.IpAddresses_Cidr(jsii.String("10.0.0.0/16")),
		MaxAzs:      jsii.Number(3),
		NatGateways: jsii.Number(1),
	})

	// Next we will define our EC2 instance properties

	// This creates a new role using the AmazonSSMManagedInstanceCore managed policy. We will use this to log into our EC2 instance
	ssmPolicy := awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonSSMManagedInstanceCore"))
	instanceRole := awsiam.NewRole(stack, jsii.String("micadevInstanceRole"),
		&awsiam.RoleProps{
			AssumedBy:       awsiam.NewServicePrincipal(jsii.String("ec2.amazonaws.com"), nil),
			Description:     jsii.String("Instance Role"),
			ManagedPolicies: &[]awsiam.IManagedPolicy{ssmPolicy},
		},
	)

	// This defines an EC2 instance we will use for remote development
	awsec2.NewInstance(stack, jsii.String("micadevInstance"),
		&awsec2.InstanceProps{
			InstanceType: awsec2.InstanceType_Of(awsec2.InstanceClass_C5, awsec2.InstanceSize_LARGE),
			MachineImage: awsec2.NewAmazonLinuxImage(nil),
			Vpc:          vpc,
			KeyName:      jsii.String(os.Getenv("AWS_INSTANCE_KEY_PAIR")),
			InstanceName: jsii.String(os.Getenv("INSTANCE_NAME")),
			Role:         instanceRole,
			VpcSubnets: &awsec2.SubnetSelection{
				SubnetType: awsec2.SubnetType_PRIVATE_WITH_EGRESS,
			},
		})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewGoMicahDevStack(app, "GoMicahDevStack", &GoMicahDevStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

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
