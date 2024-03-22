package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"aws-eks-study/src/types"
	"aws-eks-study/src/vpc"

	"github.com/alexflint/go-arg"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/joho/godotenv"
)

var args struct {
	Region string `args:"" help:"AWS Region"`
	//VpcID  string `args:"" help:"AWS VPC ID"`
}

var (
	appName          string
	cidr             string
	sess             *session.Session
	err              error
	tagsFields       types.SubnetsTagsFields
	tags             types.SubnetTags
	private_subnet_1 string
	private_subnet_2 string
	public_subnet_1  string
	public_subnet_2  string
	tagValues        types.SubnetTags
)

func init() {

	err = godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	arg.MustParse(&args)

	if args.Region == "" {
		log.Fatalln("\033[1;31m[+]\033[0m Please provide a region, --help for options")
	}

	appName = os.Getenv("APP_NAME")
	cidr = os.Getenv("AWS_VPC_CIDR")
	private_subnet_1 = os.Getenv("AWS_VPC_PRIVATE_SUBNET_1")
	private_subnet_2 = os.Getenv("AWS_VPC_PRIVATE_SUBNET_2")
	public_subnet_1 = os.Getenv("AWS_VPC_PUBLIC_SUBNET_1")
	public_subnet_2 = os.Getenv("AWS_VPC_PUBLIC_SUBNET_2")

	tags = types.SubnetTags{
		Tags: []types.SubnetsTagsFields{},
	}

	for i := 1; i < 5; i++ {

		if i < 3 {
			tagsFields = types.SubnetsTagsFields{
				Key:   "Name",
				Value: fmt.Sprintf("%s%s-%d", "public-", appName, i),
			}

		} else {
			tagsFields = types.SubnetsTagsFields{
				Key:   "Name",
				Value: fmt.Sprintf("%s%s-%d", "private-", appName, i),
			}
		}

		tags = types.SubnetTags{
			Tags: append(tags.Tags, tagsFields),
		}

	}

	var bytes bytes.Buffer
	enc := json.NewEncoder(&bytes)
	enc.Encode(tags)
	json.Unmarshal(bytes.Bytes(), &tagValues)

	sess, err = session.NewSession(&aws.Config{
		Region: aws.String(args.Region),
	})
	if err != nil {
		fmt.Println("Error creating session:", err)
		return
	}
}

func main() {

	log.Println("\033[1;32m[+]\033[0m Creating VPC...")
	vpcID := vpc.CreateVpc(sess, cidr)
	log.Println("\033[1;32m[+]\033[0m Creating Private Subnets on vpcID:", vpcID)
	vpc.CreatePrivateSubnets(sess, vpcID, private_subnet_1, private_subnet_2, tagValues)
	log.Println("\033[1;32m[+]\033[0m Creating Public Subnets on vpcID:", vpcID)
	publicSubnets := vpc.CreatePublicSubnets(sess, vpcID, public_subnet_1, public_subnet_2, tagValues)
	log.Println("\033[1;32m[+]\033[0m Creating InternetGateway...")
	igwID := vpc.CreateInternetGateway(sess, vpcID)
	log.Println("\033[1;32m[+]\033[0m Creating NatGateway...")
	vpc.CreateNatGateway(sess, publicSubnets[0])
	log.Println("\033[1;32m[+]\033[0m Creating Private RouteTable...")
	vpc.CreatePrivateRouteTable(sess, vpcID, igwID, publicSubnets)
	log.Println("\033[1;32m[+]\033[0m Updating Main Public RouteTable...")
	vpc.UpdatePublicRouteTable(sess, vpcID, igwID, publicSubnets)

}
