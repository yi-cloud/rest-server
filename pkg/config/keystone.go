package config

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/spf13/viper"
)

type KeystoneOptions struct {
	EndPoint    string
	DomainName  string
	ProjectName string
	ProjectId   string
	Region      string
}

var KeystoneOpt KeystoneOptions
var KeystoneService = "identity"

func initKeystoneOptions() {
	KeystoneOpt.EndPoint = viper.GetString("keystone.endpoint")
	if KeystoneOpt.EndPoint == "" {
		KeystoneOpt.EndPoint = "http://localhost:5000/v3"
	}

	KeystoneOpt.DomainName = viper.GetString("keystone.domainname")
	if KeystoneOpt.DomainName == "" {
		KeystoneOpt.DomainName = "Default"
	}

	KeystoneOpt.ProjectName = viper.GetString("keystone.projectname")
	KeystoneOpt.ProjectId = viper.GetString("keystone.projectid")

	KeystoneOpt.Region = viper.GetString("keystone.region")
	if KeystoneOpt.Region == "" {
		KeystoneOpt.Region = "RegionOne"
	}
}

func GetAuthScope() *gophercloud.AuthScope {
	if KeystoneOpt.ProjectName != "" {
		return &gophercloud.AuthScope{
			DomainName:  KeystoneOpt.DomainName,
			ProjectName: KeystoneOpt.ProjectName,
		}
	}

	return &gophercloud.AuthScope{
		ProjectID: KeystoneOpt.ProjectId,
	}
}
