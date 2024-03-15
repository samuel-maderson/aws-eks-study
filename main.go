package main

import (
	"fmt"

	"github.com/alexflint/go-arg"
)

var args struct {
	Region string `args:"-r, --region" help:"TEST"`
	VpcID  string `args: help:"TEST2"`
}

var (
	region = "us-east-1"
	vpcID  = "vpc-0d5ac2223eb54c91d"
)

func init() {

}

func main() {

	arg.MustParse(&args)
	fmt.Println(args.Region)
	//vpc.PrivateSubnetsIds(vpcID, region)
}
