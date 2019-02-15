package main

import (
	
        "context"
	"log"

        "github.com/oracle/oci-go-sdk/common/auth"
	"github.com/oracle/oci-go-sdk/core"


	"dp/helpers"
)

func main() {

	log.Println("[ START ]")

	//
	// Defining common stuff
	//
	log.Println("         Initializing ...")
        provider, err := auth.InstancePrincipalConfigurationProvider()
        helpers.FatalIfError(err)

        virtualNetworkClient, err := core.NewVirtualNetworkClientWithConfigurationProvider(provider)
        helpers.FatalIfError(err)

        computeInstanceClient, err := core.NewComputeClientWithConfigurationProvider(provider)
        helpers.FatalIfError(err)

	ctx := context.Background()
	
	helpers.ParseEnvironmentVariables()

	log.Println("         Done")

	//
	// Getting my OCID by calling http://169.254.169.254/opc/v1/instance/
	//
	compartmentId,computeId,err := helpers.GetComputeInstanceInfo()
	helpers.FatalIfError(err)
	log.Println("         ...My Instance ID :",*computeId)
	log.Println("         My Compartment ID :",*compartmentId)

	//
	// Getting all my attached VNICs and their SubnetIDs
	//
	attachedVNICs := helpers.ListAttachedVNICs(computeInstanceClient,ctx,compartmentId,computeId)
	log.Println("         My Attached VNICs :[")
	for _, vnic := range attachedVNICs {
		log.Printf("                             %v,",*vnic.VnicId)

	}
	log.Println("                            ]")

	//
	// Moving the routing
	//
	routeTables 			:= helpers.GetRouteInfoFromMetadata()
	updateRouteTableRequests	:= helpers.BuildRouteUpdateStructs(virtualNetworkClient,ctx,attachedVNICs,routeTables)
	helpers.FatalIfError(err)
	for _, v := range updateRouteTableRequests {
		helpers.ChangeRouteTable(virtualNetworkClient,ctx,v)
	}
	log.Println("         Routing Rules Updated")

	// Done
	log.Println("[  END  ]")
}
