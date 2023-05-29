package common

import (
	ecsutils "bott-the-pigeon/lib/aws/service/ecs"
	stsutils "bott-the-pigeon/lib/aws/service/sts"

	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/sts"
)

// Returns true if the minecraft server is currently running, and false otherwise.
// Note that it also returns false if an errors occurs.
func CheckMinecraftServerStatus(clusterNameOrArn string, serviceNameOrArn string, isCrossAccount bool) (bool, error) {
	serviceNames := []*string{&serviceNameOrArn}
	ecsDescribeServicesIn := &ecs.DescribeServicesInput{
		Cluster:  &clusterNameOrArn,
		Services: serviceNames,
	}
	var ecsOut *ecs.DescribeServicesOutput
	var err error
	if isCrossAccount {
		timestamp := time.Now().Unix() * 1000
		roleArn := os.Getenv("MC_VANILLA_STS_ROLE_ARN")
		roleSessionName := fmt.Sprintf("Bot-Assume-Minecraft-Account-Role-%d", timestamp)
		assumeRoleIn := &sts.AssumeRoleInput{
			RoleArn:         &roleArn,
			RoleSessionName: &roleSessionName,
		}
		assumeRoleOut, err := stsutils.GetCrossAccountCredentials(assumeRoleIn)
		if err != nil {
			return false, err
		}
		staticCreds := credentials.NewStaticCredentials(
			*assumeRoleOut.Credentials.AccessKeyId,
			*assumeRoleOut.Credentials.SecretAccessKey,
			*assumeRoleOut.Credentials.SessionToken,
		)
		cfg := &aws.Config{
			Credentials: staticCreds,
		}
		ecsOut, err = ecsutils.DescribeServicesWithCfg(ecsDescribeServicesIn, cfg)
	} else {
		ecsOut, err = ecsutils.DescribeServices(ecsDescribeServicesIn)
	}
	if err != nil {
		return false, err
	}
	mcService := *ecsOut.Services[0]
	if *mcService.RunningCount > 0 {
		return true, nil
	} else {
		return false, nil
	}
}
