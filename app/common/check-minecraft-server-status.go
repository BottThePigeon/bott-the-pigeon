package common

import (
	ecsutils "bott-the-pigeon/lib/aws/service/ecs"
	"os"

	"github.com/aws/aws-sdk-go/service/ecs"
)

// Returns true if the minecraft server is currently running, and false otherwise.
// Note that it also returns false if an errors occurs.
func CheckMinecraftServerStatus() (bool, error) {
	clusterName := os.Getenv("MC_CLUSTER_NAME")
	serviceName := os.Getenv("MC_SERVICE_NAME")
	serviceNames := []*string{&serviceName}
	ecsDescribeServicesIn := &ecs.DescribeServicesInput{
		Cluster:  &clusterName,
		Services: serviceNames,
	}
	ecsOut, err := ecsutils.DescribeServices(ecsDescribeServicesIn)
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
