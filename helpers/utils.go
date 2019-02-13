package helpers

import (
        "context"
	"log"
	"time"
	"errors"
	"io/ioutil"
	"net/http"
	"encoding/json"

        "github.com/oracle/oci-go-sdk/common/auth"
        "github.com/oracle/oci-go-sdk/core"

)
type RtRule struct {
	Destination string `json:"destination"`
	Gateway string `json:"gateway"`
}

type RouteTable struct {
	RtId string `json:"rtId"`
	RtRules []RtRule `json:"rtRules"`
}

func FatalIfError(err error) {
        if err != nil {
                log.Fatalln(err.Error())
        }
}

func DeletePrivateIP(privateIPId *string) {
        provider, err := auth.InstancePrincipalConfigurationProvider()
        FatalIfError(err)

        client, err := core.NewVirtualNetworkClientWithConfigurationProvider(provider)

        request := core.DeletePrivateIpRequest {
                PrivateIpId : privateIPId,
        }
	ctx := context.Background()
        r, err := client.DeletePrivateIp(ctx, request)
        FatalIfError(err)
	log.Println("DeletePrivateIP :",r)

}

// CreatePrivateIP creates a new private IP assigned to a specific VNIC
// it requires core.CreatePrivateIpDetails with at least VnicId and IpAddress information
func CreatePrivateIP(createPrivateIPDetails core.CreatePrivateIpDetails) {
        provider, err := auth.InstancePrincipalConfigurationProvider()
        FatalIfError(err)

        client, err := core.NewVirtualNetworkClientWithConfigurationProvider(provider)

        request := core.CreatePrivateIpRequest {
                CreatePrivateIpDetails: createPrivateIPDetails,
        }

	ctx := context.Background()
        r, err := client.CreatePrivateIp(ctx, request)
        FatalIfError(err)
	log.Println("CreatePrivateIP :",r)


}

// GetComputeInstanceInfo retrieves instance ID and its compartment from http://169.254.169.254/opc/v1/instance/id
func GetComputeInstanceInfo() (*string,*string,error) {

        var myClient = &http.Client{Timeout: 10 * time.Second}
	var compartmentId,computeId *string

        r, err := myClient.Get("http://169.254.169.254/opc/v1/instance/")
        FatalIfError(err)

        body, err := ioutil.ReadAll(r.Body)
	var f interface{}
        err = json.Unmarshal(body, &f)
	FatalIfError(err)

	m:=f.(map[string]interface{})
    	for k, v := range m {
		switch {
		case k == "compartmentId":
			s, _ := v.(string)
			compartmentId = &s
		case  k == "id":
			s, _  := v.(string)
			computeId = &s
		}
    	}
	if (compartmentId == nil || computeId ==nil) {
	 	return nil,nil,errors.New("Can't find compartment-id or compute instance id in http://169.254.169.254/opc/v1/instance/ .")
	}
	return compartmentId,computeId,nil
}

// GetRouteInfo retrieves routes information from http://169.254.169.254/opc/v1/instance/metadata/route_tables
func GetRouteInfoFromMetadata()([]RouteTable) {
	
        var myClient = &http.Client{Timeout: 10 * time.Second}
	//var updateRouteTableRequests []core.UpdateRouteTableRequest 

        r, err := myClient.Get("http://169.254.169.254/opc/v1/instance/metadata/route_tables")
        FatalIfError(err)

        body, err := ioutil.ReadAll(r.Body)
	var routeTables []RouteTable
        err = json.Unmarshal(body, &routeTables)
	FatalIfError(err)
	log.Printf("routeItem: %+v",routeTables)
	return routeTables
}

func BuildRouteUpdateStructs(client core.VirtualNetworkClient,ctx context.Context,attachedVNICs []core.VnicAttachment,routeTables []RouteTable) ([]core.UpdateRouteTableRequest) {
	var newRouteRules []core.RouteRule
	var updateRouteTableRequests []core.UpdateRouteTableRequest
	for _, rt := range routeTables {
		for _, rR := range rt.RtRules {
			privateIpID,_,err := GetPrivateIPID(client,ctx,attachedVNICs,&rR.Gateway)
			FatalIfError(err)
			getRouteTableRequest   := core.GetRouteTableRequest { RtId : &rt.RtId, } 
			getRouteTableResponse,err  := client.GetRouteTable(ctx,getRouteTableRequest)
			FatalIfError(err)
			currentRouteRules  := getRouteTableResponse.RouteTable.RouteRules 
			newRouteRules      = getRouteTableResponse.RouteTable.RouteRules 
			log.Println("  currentRouteRules: ",currentRouteRules)
			i		       := 0
		        for _, cRR := range currentRouteRules {
				if (*cRR.Destination == rR.Destination) {
					log.Println("====> Removing ",*cRR.Destination)
					newRouteRules = append(newRouteRules[:i],newRouteRules[i+1:]...)
					 log.Println("====> Removed, New values are :",newRouteRules)
				}
				i = i + 1
			}		

			routeRule := core.RouteRule{
						    NetworkEntityId: privateIpID,
						    Destination: &rR.Destination,
						    DestinationType: core.RouteRuleDestinationTypeCidrBlock,       
						   } 
			newRouteRules = append(newRouteRules,routeRule) 
			log.Println("====> Addeded, New values are :",newRouteRules)
		}
		updateRouteTableDetails := core.UpdateRouteTableDetails {
						RouteRules : newRouteRules,
		}
	 	updateRouteTableRequest	:= core.UpdateRouteTableRequest{
						RtId : &rt.RtId,
						UpdateRouteTableDetails : updateRouteTableDetails,
						
		}
		updateRouteTableRequests = append(updateRouteTableRequests,updateRouteTableRequest)
	}
	return updateRouteTableRequests
}

func ListAttachedVNICs(client core.ComputeClient,ctx context.Context,compartmentId,computeId *string) ([]core.VnicAttachment) {
        request := core.ListVnicAttachmentsRequest {
                CompartmentId: compartmentId,
		InstanceId:computeId,
        }

        r, err := client.ListVnicAttachments(ctx, request)
        FatalIfError(err)

	return r.Items
}

func GetPrivateIPID(client core.VirtualNetworkClient,ctx context.Context,attachedVNICs []core.VnicAttachment, IpAddress *string) (*string, *string, error) {
	for _,vnic := range attachedVNICs {
	        request := core.ListPrivateIpsRequest  {
	                IpAddress: IpAddress,
	                SubnetId: vnic.SubnetId,
	        }

		ctx := context.Background()
		r, err := client.ListPrivateIps(ctx, request)
		FatalIfError(err)

		for _, v := range r.Items {
			if *v.IpAddress == *IpAddress {
				return v.Id, vnic.VnicId, nil
	
			}
		}
	}
	return nil,nil,errors.New("Can't find :"+*IpAddress)
}

// ReassignPrivateIP reassigns a private IP to the specified VnicId
func ReassignPrivateIP(client core.VirtualNetworkClient,ctx context.Context,vnicId,privateIPId *string, ) (*core.PrivateIp){
	IpDetails := core.UpdatePrivateIpDetails {
		VnicId : vnicId,
		}	
        request := core.UpdatePrivateIpRequest {
                PrivateIpId : privateIPId,
		UpdatePrivateIpDetails : IpDetails,
        }
        r, err := client.UpdatePrivateIp(ctx, request)
        FatalIfError(err)
	return &r.PrivateIp
}

// ReassignPrivateIP reassigns a private IP to the specified VnicId
func ChangeRouteTable(client core.VirtualNetworkClient,ctx context.Context,updateRouteTableRequest core.UpdateRouteTableRequest) {

        _, err := client.UpdateRouteTable(ctx, updateRouteTableRequest)
        FatalIfError(err)
}

