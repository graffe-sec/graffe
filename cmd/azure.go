package cmd

import (
	"fmt"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	//"github.com/graffe-sec/graffe-azure"
	"github.com/spf13/cobra"
)

var azureCmd = &cobra.Command{
	Use:   "azure",
	Short: "Sets the scope to Azure",
}

var azureReviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Reviews the specified Azure service (this is purely READ ONLY)",
}

// Device token authentication flag vars
var azureCliAuth bool

// Client credential authentication flag vars
var (
	// Global configuration variables
	appID       string
	appSec      string
	cloudName   string = "AzurePublicCloud"
	environment *azure.Environment
	tenantID    string

	// Authorizers
	graphAuthorizer autorest.Authorizer
)

type OAuthGrantType int

const (
	// OAuthGrantTypeServicePrincipal for client credentials flow
	OAuthGrantTypeServicePrincipal OAuthGrantType = iota
	// OAuthGrantTypeDeviceFlow
	OAuthGrantTypeDeviceFlow
)

func init() {
	rootCmd.AddCommand(azureCmd)
	azureCmd.AddCommand(azureReviewCmd)

	// Device token authentication flag (uses Azure CLI)
	azureCmd.PersistentFlags().BoolVar(&azureCliAuth, "cli-auth", false, "")

	// Client credentials authentication flags (uses App Registrations)
	azureCmd.PersistentFlags().StringVar(&appID, "app-id", "false", "Provide the App ID of the AAD App Registration")
	azureCmd.PersistentFlags().StringVar(&appSec, "app-sec", "false", "Provide the App Password of the AAD App Registration")
	azureCmd.PersistentFlags().StringVar(&tenantID, "tenant-id", "false", "Provide the Azyre AD tenant ID")
}

func azureAuth() (autorest.Authorizer, error) {
	if graphAuthorizer != nil {
		return graphAuthorizer, nil
	}

	var a autorest.Authorizer
	var err error
	if azureCliAuth == true {
		a, err = getAuthorizerForResource(OAuthGrantTypeDeviceFlow, Environment().GraphEndpoint)

		if err == nil {
			graphAuthorizer = a
		} else {
			graphAuthorizer = nil
		}
	}

	return graphAuthorizer, err
}

func getAuthorizerForResource(grantType OAuthGrantType, resource string) (autorest.Authorizer, error) {
	var a autorest.Authorizer
	var err error

	switch grantType {

	case OAuthGrantTypeServicePrincipal:
		oauthConfig, err := adal.NewOAuthConfig(
			environment.ActiveDirectoryEndpoint, tenantID)
		if err != nil {
			return nil, err
		}

		token, err := adal.NewServicePrincipalToken(
			*oauthConfig, appID, appSec, resource)
		if err != nil {
			return nil, err
		}
		a = autorest.NewBearerAuthorizer(token)

	case OAuthGrantTypeDeviceFlow:
		deviceConfig := auth.NewDeviceFlowConfig(appID, tenantID)
		deviceConfig.Resource = resource
		a, err = deviceConfig.Authorizer()
		if err != nil {
			return nil, err
		}

	default:
		return a, fmt.Errorf("invalid grant type specified")
	}

	return a, err
}

func Environment() *azure.Environment {
	if environment != nil {
		return environment
	}
	env, err := azure.EnvironmentFromName(cloudName)
	if err != nil {
		panic(fmt.Sprintf("invalid cloud name '%s' specified, cannot continue\n", cloudName))
	}
	environment = &env
	return environment
}
