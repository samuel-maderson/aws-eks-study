package vpc

import (
	"fmt"
	"log"
	"os"

	"aws-eks-study/src/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func CreateVpc(sess *session.Session, cidr string) string {

	// Create EC2 service client
	svc := ec2.New(sess)

	// Create VPC
	vpcInput := &ec2.CreateVpcInput{
		CidrBlock: aws.String(cidr),
	}

	result, err := svc.CreateVpc(vpcInput)
	if err != nil {
		log.Fatalln("Error creating VPC:", err)
	}

	// Add tags to the subnet
	tagInput := &ec2.CreateTagsInput{
		Resources: []*string{result.Vpc.VpcId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(os.Getenv("APP_NAME")),
			},
		},
	}
	_, err = svc.CreateTags(tagInput)
	if err != nil {
		fmt.Println("Error adding tags to subnet:", err)
		os.Exit(1)
	}

	return *result.Vpc.VpcId

}

func CreatePublicSubnets(sess *session.Session, vpcID string, subnet_1 string, subnet_2 string, tagValues types.SubnetTags) (subnets []string) {

	svc := ec2.New(sess)
	// Create public subnets
	publicSubnets := []string{subnet_1, subnet_2} // Change CIDR blocks as needed

	for key, cidr := range publicSubnets {
		subnetInput := &ec2.CreateSubnetInput{
			CidrBlock: aws.String(cidr),
			VpcId:     aws.String(vpcID),
		}

		result, err := svc.CreateSubnet(subnetInput)
		if err != nil {
			fmt.Println("Error creating subnet:", err)
			os.Exit(1)
		}
		// Add tags to the subnet
		tagInput := &ec2.CreateTagsInput{
			Resources: []*string{result.Subnet.SubnetId},
			Tags: []*ec2.Tag{
				{
					Key:   aws.String(tagValues.Tags[key].Key),
					Value: aws.String(tagValues.Tags[key].Value),
				},
			},
		}
		_, err = svc.CreateTags(tagInput)
		if err != nil {
			fmt.Println("Error adding tags to subnet:", err)
			os.Exit(1)
		}
	}

	return subnets
}

func CreatePrivateSubnets(sess *session.Session, vpcID string, subnet_1 string, subnet_2 string, tagValues types.SubnetTags) (subnets []string) {

	svc := ec2.New(sess)
	// Create private subnets
	privateSubnets := []string{subnet_1, subnet_2} // Change CIDR blocks as needed
	for key, cidr := range privateSubnets {
		subnetInput := &ec2.CreateSubnetInput{
			CidrBlock: aws.String(cidr),
			VpcId:     aws.String(vpcID),
		}
		result, err := svc.CreateSubnet(subnetInput)
		if err != nil {
			fmt.Println("Error creating subnet:", err)
			os.Exit(1)
		}
		// Add tags to the subnet
		tagInput := &ec2.CreateTagsInput{
			Resources: []*string{result.Subnet.SubnetId},
			Tags: []*ec2.Tag{
				{
					Key:   aws.String(tagValues.Tags[key].Key),
					Value: aws.String(tagValues.Tags[key].Value),
				},
			},
		}
		_, err = svc.CreateTags(tagInput)
		if err != nil {
			fmt.Println("Error adding tags to subnet:", err)
			os.Exit(1)
		}

		fmt.Println(*result.Subnet.SubnetId)
	}

	return subnets
}
