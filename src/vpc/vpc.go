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

	svc := ec2.New(sess)

	vpcInput := &ec2.CreateVpcInput{
		CidrBlock: aws.String(cidr),
	}

	result, err := svc.CreateVpc(vpcInput)
	if err != nil {
		log.Fatalln("Error creating VPC:", err)
	}

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

	firstTwo := tagValues.Tags[:2]
	svc := ec2.New(sess)

	publicSubnets := []string{subnet_1, subnet_2}

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

		tagInput := &ec2.CreateTagsInput{
			Resources: []*string{result.Subnet.SubnetId},
			Tags: []*ec2.Tag{
				{
					Key:   aws.String(firstTwo[key].Key),
					Value: aws.String(firstTwo[key].Value),
				},
			},
		}
		_, err = svc.CreateTags(tagInput)
		if err != nil {
			fmt.Println("Error adding tags to subnet:", err)
			os.Exit(1)
		}

		subnets = append(subnets, *result.Subnet.SubnetId)
	}

	return subnets
}

func CreatePrivateSubnets(sess *session.Session, vpcID string, subnet_1 string, subnet_2 string, tagValues types.SubnetTags) (subnets []string) {

	lastTwo := tagValues.Tags[len(tagValues.Tags)-2:]
	svc := ec2.New(sess)

	privateSubnets := []string{subnet_1, subnet_2}
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

		tagInput := &ec2.CreateTagsInput{
			Resources: []*string{result.Subnet.SubnetId},
			Tags: []*ec2.Tag{
				{
					Key:   aws.String(lastTwo[key].Key),
					Value: aws.String(lastTwo[key].Value),
				},
			},
		}
		_, err = svc.CreateTags(tagInput)
		if err != nil {
			fmt.Println("Error adding tags to subnet:", err)
			os.Exit(1)
		}

		subnets = append(subnets, *result.Subnet.SubnetId)
	}

	return subnets
}

func CreateNatGateway(sess *session.Session, subnetID string) string {

	svc := ec2.New(sess)

	elasticIpInput := &ec2.AllocateAddressInput{
		Domain: aws.String("vpc"),
	}
	elasticIpResult, err := svc.AllocateAddress(elasticIpInput)
	if err != nil {
		log.Fatal(err)
	}

	natGatewayInput := &ec2.CreateNatGatewayInput{
		AllocationId: elasticIpResult.AllocationId,
		SubnetId:     aws.String(subnetID),
	}
	natGatewayResult, err := svc.CreateNatGateway(natGatewayInput)
	if err != nil {
		log.Fatal(err)
	}

	return *natGatewayResult.NatGateway.NatGatewayId
}

func CreateInternetGateway(sess *session.Session, vpcID string) string {

	svc := ec2.New(sess)

	internetGatewayInput := &ec2.CreateInternetGatewayInput{}
	internetGatewayResult, err := svc.CreateInternetGateway(internetGatewayInput)
	if err != nil {
		log.Fatal(err)
	}

	attachGatewayInput := &ec2.AttachInternetGatewayInput{
		InternetGatewayId: internetGatewayResult.InternetGateway.InternetGatewayId,
		VpcId:             aws.String(vpcID),
	}
	_, err = svc.AttachInternetGateway(attachGatewayInput)
	if err != nil {
		log.Fatal(err)
	}

	return *internetGatewayResult.InternetGateway.InternetGatewayId
}

func CreatePrivateRouteTable(sess *session.Session, vpcID string, ngwID string, subnets []string) {

	svc := ec2.New(sess)

	routeTableOutput, err := svc.CreateRouteTable(&ec2.CreateRouteTableInput{
		VpcId: aws.String(vpcID),
	})
	if err != nil {
		log.Fatalf("Error creating Route Table: %v", err)
	}

	routeTableID := routeTableOutput.RouteTable.RouteTableId

	_, err = svc.CreateRoute(&ec2.CreateRouteInput{
		RouteTableId:         aws.String(*routeTableID),
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		GatewayId:            aws.String(ngwID),
	})
	if err != nil {
		log.Fatalf("Error adding route to Route Table: %v", err)
	}

	for _, subnetID := range subnets {
		_, err = svc.AssociateRouteTable(&ec2.AssociateRouteTableInput{
			RouteTableId: aws.String(*routeTableID),
			SubnetId:     aws.String(subnetID),
		})
		if err != nil {
			log.Fatalf("Error associating subnet with Route Table: %v", err)
		}
	}
}

func UpdatePublicRouteTable(sess *session.Session, vpcID string, igwID string, subnets []string) {

	svc := ec2.New(sess)

	resp, err := svc.DescribeRouteTables(&ec2.DescribeRouteTablesInput{})
	if err != nil {
		fmt.Println("failed to describe route tables,", err)
		return
	}

	routeTableID := *resp.RouteTables[1].RouteTableId

	_, err = svc.CreateRoute(&ec2.CreateRouteInput{
		RouteTableId:         aws.String(routeTableID),
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		GatewayId:            aws.String(igwID),
	})

	if err != nil {
		log.Fatalln(err)
	}

	for _, subnetID := range subnets {

		_, err = svc.AssociateRouteTable(&ec2.AssociateRouteTableInput{
			RouteTableId: aws.String(routeTableID),
			SubnetId:     aws.String(subnetID),
		})

		if err != nil {
			log.Fatalln(err)
		}

	}

}
