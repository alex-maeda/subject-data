package main

import (
	"fmt"
	"net"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// subnetCIDRs splits a VPC CIDR into two equal halves by incrementing the
// prefix length by 1. For example, 10.2.0.0/24 → [10.2.0.0/25, 10.2.0.128/25].
func subnetCIDRs(vpcCidr string, count int) ([]string, error) {
	ip, ipNet, err := net.ParseCIDR(vpcCidr)
	if err != nil {
		return nil, fmt.Errorf("parse VPC CIDR %q: %w", vpcCidr, err)
	}

	ones, bits := ipNet.Mask.Size()
	subnetBits := 0
	for (1 << subnetBits) < count {
		subnetBits++
	}
	newOnes := ones + subnetBits
	if newOnes > bits {
		return nil, fmt.Errorf("VPC CIDR %q too small to split into %d subnets", vpcCidr, count)
	}

	subnetSize := 1 << (bits - newOnes)
	base := ipToUint32(ip.Mask(ipNet.Mask).To4())

	cidrs := make([]string, count)
	for i := range count {
		subnetIP := uint32ToIP(base + uint32(i*subnetSize))
		cidrs[i] = fmt.Sprintf("%s/%d", subnetIP, newOnes)
	}
	return cidrs, nil
}

func ipToUint32(ip net.IP) uint32 {
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

func uint32ToIP(n uint32) net.IP {
	return net.IPv4(byte(n>>24), byte(n>>16&0xff), byte(n>>8&0xff), byte(n&0xff))
}

// NetworkResources holds the networking resources created by createNetworking.
type NetworkResources struct {
	VPC       *ec2.Vpc
	SubnetIDs pulumi.StringArray
	ECSSG     *ec2.SecurityGroup
}

var azs = []string{"us-east-2a", "us-east-2b"}

// createNetworking creates a VPC with public subnets and the ECS security group.
// Beta uses public subnets only; add private subnets + NAT gateway for prod.
func createNetworking(ctx *pulumi.Context, appEnv, vpcCidr string) (*NetworkResources, error) {
	vpc, err := ec2.NewVpc(ctx, "vpc", &ec2.VpcArgs{
		CidrBlock:          pulumi.String(vpcCidr),
		EnableDnsHostnames: pulumi.Bool(true),
		EnableDnsSupport:   pulumi.Bool(true),
		Tags:               tags("vpc", appEnv),
	})
	if err != nil {
		return nil, err
	}

	igw, err := ec2.NewInternetGateway(ctx, "igw", &ec2.InternetGatewayArgs{
		VpcId: vpc.ID(),
		Tags:  tags("igw", appEnv),
	})
	if err != nil {
		return nil, err
	}

	publicRT, err := ec2.NewRouteTable(ctx, "public-rt", &ec2.RouteTableArgs{
		VpcId: vpc.ID(),
		Routes: ec2.RouteTableRouteArray{
			&ec2.RouteTableRouteArgs{
				CidrBlock: pulumi.String("0.0.0.0/0"),
				GatewayId: igw.ID(),
			},
		},
		Tags: tags("public-rt", appEnv),
	}, pulumi.DependsOn([]pulumi.Resource{igw}))
	if err != nil {
		return nil, err
	}

	cidrs, err := subnetCIDRs(vpcCidr, len(azs))
	if err != nil {
		return nil, err
	}

	subnetIDs := make(pulumi.StringArray, len(azs))
	for i, az := range azs {
		sub, err := ec2.NewSubnet(ctx, fmt.Sprintf("public-%d", i), &ec2.SubnetArgs{
			VpcId:               vpc.ID(),
			CidrBlock:           pulumi.String(cidrs[i]),
			AvailabilityZone:    pulumi.String(az),
			MapPublicIpOnLaunch: pulumi.Bool(true),
			Tags:                pulumi.StringMap{"Name": pulumi.Sprintf("subject-data-public-%s-%s", az, appEnv)},
		}, pulumi.DependsOn([]pulumi.Resource{igw}))
		if err != nil {
			return nil, err
		}
		subnetIDs[i] = sub.ID()

		_, err = ec2.NewRouteTableAssociation(ctx, fmt.Sprintf("public-rta-%d", i), &ec2.RouteTableAssociationArgs{
			SubnetId:     sub.ID(),
			RouteTableId: publicRT.ID(),
		})
		if err != nil {
			return nil, err
		}
	}

	ecsSG, err := ec2.NewSecurityGroup(ctx, "ecs-sg", &ec2.SecurityGroupArgs{
		VpcId:       vpc.ID(),
		Description: pulumi.String("subject-data ECS tasks"),
		Ingress: ec2.SecurityGroupIngressArray{
			&ec2.SecurityGroupIngressArgs{
				Protocol:    pulumi.String("tcp"),
				FromPort:    pulumi.Int(8380),
				ToPort:      pulumi.Int(8380),
				CidrBlocks:  pulumi.StringArray{pulumi.String(vpcCidr)},
				Description: pulumi.String("API Gateway VPC Link"),
			},
		},
		Egress: ec2.SecurityGroupEgressArray{
			&ec2.SecurityGroupEgressArgs{
				Protocol:   pulumi.String("-1"),
				FromPort:   pulumi.Int(0),
				ToPort:     pulumi.Int(0),
				CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
			},
		},
		Tags: tags("ecs-sg", appEnv),
	})
	if err != nil {
		return nil, err
	}

	return &NetworkResources{
		VPC:       vpc,
		SubnetIDs: subnetIDs,
		ECSSG:     ecsSG,
	}, nil
}
