package main

import (
	"aws-eks-study/vpc"
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/joho/godotenv"
)

var args struct {
	Region string `args:"" help:"AWS Region"`
	VpcID  string `args:"" help:"AWS VPC ID"`
}

var (
	appName string
	tagName []string
)

func init() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	arg.MustParse(&args)

	if args.Region == "" && args.VpcID == "" {
		fmt.Println("\033[1;31m[+]\033[0m Please provide a region and VPC ID, --help for options")
		return
	}

	appName = os.Getenv("APP_NAME")
	tagName = []string{
		fmt.Sprintf("%s%s%s%s", appName, "-subnet-private1-", args.Region, "a"),
		fmt.Sprintf("%s%s%s%s", appName, "-subnet-private2-", args.Region, "b"),
	}

}

func main() {

	vpc.PrivateSubnetsIds(args.VpcID, args.Region, tagName)
}
