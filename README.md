# go-oci-linux-vip-routing
Go programs to manage VIP association to Oracle Cloud Instance VNICs and Subnet Routing Tables

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. 

### Prerequisites

This project requires 
* [Go programming language](https://golang.org/dl/) installed
* [Oracle Cloud Infrastructure Golang SDK](https://github.com/oracle/oci-go-sdk) installed

### Installing

* Install Go 
* Install oci-go-sdk - you can skip the SDK Configuration as we'll see later, these programs leverage the OCI IAM Feature called Instance Principals that will allow our VM instance to make API calls to OCI Services without configuring any user credentials
* Clone this repo under your Go workspace directory

```sh
$ cd $HOME/go/src
$ git clone https://github.com/daniel-pro/go-oci-linux-vip-routing.git
```

## Build

```sh
$ cd $HOME/go/src/go-oci-linux-vip-routing
$ go build movePrivateIp.go
$ go build moveRoutingRule.go
```

## Configuration
This projects to call OCI services uses Instance Principals https://docs.cloud.oracle.com/iaas/Content/Identity/Tasks/callingservicesfrominstances.htm to authenticate in OCI without any credentials . Therefore to make it working an additional configuration is required in OCI
* Create a Dynamic Group and call it IPSec-VRRP-DGroup or whatever name you like
* Create a rule which includes all instances where the script will be invoked :
```sh
    ANY {instance.id = 'ocid1.instance.oc1.eu-frankfurt-1.<instance-id>', instance.id = 'ocid1.instance.oc1.eu-frankfurt-1.<instance-id>', instance.id = 'ocid1.instance.oc1.eu-frankfurt-1.<instance-id>'}
```
* Creare a new Policy and call it IPSec-VRRP-DGroup-ManageNetwork-Policy or whatever suits you best
* Add the following Policy Statements :
```sh
   Allow dynamic-group IPSec-VRRP-DGroup to use private-ips in compartment <YOUR_COMPARTMENT>
   Allow dynamic-group IPSec-VRRP-DGroup to use subnets in compartment <YOUR_COMPARTMENT> 
   Allow dynamic-group IPSec-VRRP-DGroup to use vnics in compartment <YOUR_COMPARTMENT>
   Allow dynamic-group IPSec-VRRP-DGroup to manage virtual-network-family in compartment <YOUR_COMPARTMENT>
   Allow dynamic-group IPSec-VRRP-DGroup to use instances in compartment <YOUR_COMPARTMENT>
```
