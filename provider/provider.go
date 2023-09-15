package provider

import (
	"errors"
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	apiclient "github.com/solacecommunity/solace-go-semp-client/client"
	"github.com/solacecommunity/solace-go-semp-client/models"
	"github.com/solacecommunity/solace-terraform-provider/sempv1"
)

type ProviderState struct {
	Auth         runtime.ClientAuthInfoWriter
	Client       *apiclient.SEMPSolaceElementManagementProtocol
	SempV1Client *sempv1.SempV1Client
}

// List of supported configuration fields for the Solace provider.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		HOST: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The address of the Solace msg broker",
		},
		PORT: {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     443,
			Description: "The port of the Solace msg broker",
		},
		SCHEMA: {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "https",
			Description:  "The schema used to connect with the msg broker, like https (Default) or http",
			ValidateFunc: validation.StringInSlice([]string{"http", "https"}, false),
		},
		ADMIN_USER: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Admin identity to login to the Solace VMR.",
		},
		ADMIN_PASSWD: {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Password of the admin identity used to login to the Solace VMR.",
		},
	}
}

type ErrorResponse interface {
	Error() string
	GetPayload() *models.SempMetaOnlyResponse
}

func GetSempError(err interface{}) (*models.SempError, error) {
	if er, ok := err.(ErrorResponse); ok {
		return er.GetPayload().Meta.Error, nil
	}
	return nil, errors.New("err is not of type SempError")
}

// AppendError extracts an error from an HTTP Request to the Solace SEMP API and appends it to the diags collection.
func AppendError(diags diag.Diagnostics, err interface{}) diag.Diagnostics {
	if er, ok := err.(ErrorResponse); ok {
		meta := er.GetPayload().Meta
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("HTTP %d %s - %s", *meta.Error.Code, *meta.Error.Status, *meta.Error.Description),
			Detail:   fmt.Sprintf("Error calling %s %s\n%s", *meta.Request.Method, *meta.Request.URI, *meta.Error.Description),
		})
	} else {
		return diag.FromErr(err.(error))
	}
	return diags
}
