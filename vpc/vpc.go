package vpc

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func PrivateSubnetsIds(vpcID string, region string) {
	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region), // Change to your desired region
	})
	if err != nil {
		fmt.Println("Error creating session:", err)
		return
	}

	// Create EC2 client
	svc := ec2.New(sess)

	// Describe subnets
	describeInput := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpcID)},
			},
		},
	}
	describeOutput, err := svc.DescribeSubnets(describeInput)
	if err != nil {
		fmt.Println("Error describing subnets:", err)
		return
	}

	// Extract private subnet IDs
	var privateSubnetIDs []string
	for _, subnet := range describeOutput.Subnets {
		for _, tag := range subnet.Tags {
			if *tag.Key == "Name" && *tag.Value == "private" {
				privateSubnetIDs = append(privateSubnetIDs, *subnet.SubnetId)
			}
		}
	}

	// Print subnet IDs
	fmt.Println("Private Subnet IDs:")
	for _, subnetID := range privateSubnetIDs {
		fmt.Println(subnetID)
	}
}
