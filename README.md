# go-oci-linux-vip-routing
Go programs to manage VIP association to Oracle Cloud Instance VNICs and Subnet Routing Tables

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

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

