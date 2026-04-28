package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/rds"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/secretsmanager"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

// RDSResources holds the RDS resources created by createRDS.
type RDSResources struct {
	Instance     *rds.Instance
	DBSecret     *secretsmanager.Secret
	ConnStringFn pulumi.StringOutput
}

// createRDS creates an RDS PostgreSQL instance (db.t4g.micro, single-AZ) for
// subject-data. The connection string is stored in Secrets Manager and passed
// to ECS as DATABASE_URL.
func createRDS(ctx *pulumi.Context, appEnv string, net *NetworkResources) (*RDSResources, error) {
	cfg := config.New(ctx, "")
	dbName := cfg.Get("dbName")
	if dbName == "" {
		dbName = "subjectdata"
	}
	dbUser := cfg.Get("dbUser")
	if dbUser == "" {
		dbUser = "subjectdata"
	}

	// --- DB subnet group ---

	subnetGroup, err := rds.NewSubnetGroup(ctx, "db-subnet-group", &rds.SubnetGroupArgs{
		Name:        pulumi.Sprintf("subject-data-%s", appEnv),
		Description: pulumi.Sprintf("subject-data %s DB subnets", appEnv),
		SubnetIds:   net.SubnetIDs,
		Tags:        tags("db-subnet-group", appEnv),
	})
	if err != nil {
		return nil, err
	}

	// --- Security group: inbound 5432 from ECS SG only ---

	dbSG, err := ec2.NewSecurityGroup(ctx, "db-sg", &ec2.SecurityGroupArgs{
		VpcId:       net.VPC.ID(),
		Description: pulumi.String("subject-data RDS PostgreSQL"),
		Ingress: ec2.SecurityGroupIngressArray{
			&ec2.SecurityGroupIngressArgs{
				Protocol:       pulumi.String("tcp"),
				FromPort:       pulumi.Int(5432),
				ToPort:         pulumi.Int(5432),
				SecurityGroups: pulumi.StringArray{net.ECSSG.ID()},
				Description:    pulumi.String("PostgreSQL from ECS tasks"),
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
		Tags: tags("db-sg", appEnv),
	})
	if err != nil {
		return nil, err
	}

	// --- Secrets Manager: DB password ---

	dbSecret, err := secretsmanager.NewSecret(ctx, "db-password", &secretsmanager.SecretArgs{
		Name: pulumi.Sprintf("subject-data/%s/db-password", appEnv),
		Tags: tags("db-password", appEnv),
	})
	if err != nil {
		return nil, err
	}

	// Generate a random password and store it.
	dbSecretVersion, err := secretsmanager.NewSecretVersion(ctx, "db-password-version", &secretsmanager.SecretVersionArgs{
		SecretId:     dbSecret.ID(),
		SecretString: pulumi.Sprintf("subject-data-%s-db-password", appEnv),
	})
	if err != nil {
		return nil, err
	}

	// --- RDS instance ---

	instance, err := rds.NewInstance(ctx, "db", &rds.InstanceArgs{
		Identifier:              pulumi.Sprintf("subject-data-%s", appEnv),
		Engine:                  pulumi.String("postgres"),
		EngineVersion:           pulumi.String("16"),
		InstanceClass:           pulumi.String("db.t4g.micro"),
		AllocatedStorage:        pulumi.Int(20),
		StorageType:             pulumi.String("gp3"),
		DbName:                  pulumi.String(dbName),
		Username:                pulumi.String(dbUser),
		Password:                dbSecretVersion.SecretString,
		DbSubnetGroupName:       subnetGroup.Name,
		VpcSecurityGroupIds:     pulumi.StringArray{dbSG.ID()},
		MultiAz:                 pulumi.Bool(false),
		PubliclyAccessible:      pulumi.Bool(false),
		StorageEncrypted:        pulumi.Bool(true),
		SkipFinalSnapshot:       pulumi.Bool(false),
		FinalSnapshotIdentifier: pulumi.Sprintf("subject-data-%s-final", appEnv),
		BackupRetentionPeriod:   pulumi.Int(7),
		Tags:                    tags("db", appEnv),
	})
	if err != nil {
		return nil, err
	}

	// Build the connection string: postgres://user:pass@host:port/dbname?sslmode=require
	// instance.Endpoint is already "host:port"
	connString := pulumi.All(
		instance.Endpoint, dbSecretVersion.SecretString,
	).ApplyT(func(args []any) string {
		endpoint := args[0].(string)
		password := *args[1].(*string)
		return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=require",
			dbUser, password, endpoint, dbName)
	}).(pulumi.StringOutput)

	return &RDSResources{
		Instance:     instance,
		DBSecret:     dbSecret,
		ConnStringFn: connString,
	}, nil
}
