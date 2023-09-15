package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/solacecommunity/solace-go-semp-client/client/all"
	"github.com/solacecommunity/solace-go-semp-client/models"
	"github.com/solacecommunity/solace-terraform-provider/provider"
)

// Main resource definition for Client Usernames
func ResourceClientCertAuthoritoy() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaClientCertAuthority(),
		// Provider CRUD functions
		CreateContext: createClientCertAuthority,
		ReadContext:   readClientCertAuthority,
		UpdateContext: updateClientCertAuthority,
		DeleteContext: deleteClientCertAuthority,
	}
}

func schemaClientCertAuthority() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		CERT_AUTHORITY_NAME: {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The name of the Certificate Authority.",
		},
		CERT_CONTENT: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The PEM formatted content for the trusted root certificate of a client Certificate Authority. The default value is `\"\"`",
		},
	}
}

func getClientCertAuthorityModelFromResource(d *schema.ResourceData) *models.ClientCertAuthority {
	q := &models.ClientCertAuthority{
		CertAuthorityName: d.Get(CERT_AUTHORITY_NAME).(string),
		CertContent:       d.Get(CERT_CONTENT).(string),
	}
	return q
}

func createClientCertAuthority(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	body := getClientCertAuthorityModelFromResource(d)

	params := all.NewCreateClientCertAuthorityParams().WithContext(ctx).WithBody(body)
	_, err := state.Client.All.CreateClientCertAuthority(params, state.Auth)
	if err != nil {
		if semp, err2 := provider.GetSempError(err); err2 == nil {
			if *semp.Status == "ALREADY_EXISTS" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Client Cert Authority %s already exists. Going to import state from Broker", body.CertAuthorityName),
				})
			} else {
				return provider.AppendError(diags, err)
			}
		} else {
			return provider.AppendError(diags, err)
		}
	}
	d.SetId("sol_cca_" + body.CertAuthorityName)

	// read after create to make sure tf state is in sync with broker
	return append(diags, readClientCertAuthority(ctx, d, meta)...)
}

func readClientCertAuthority(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	certname := d.Get(CERT_AUTHORITY_NAME).(string)
	params := all.NewGetClientCertAuthorityParamsWithContext(ctx).WithCertAuthorityName(certname)
	resp, err := state.Client.All.GetClientCertAuthority(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	p := resp.Payload.Data
	d.Set(CERT_AUTHORITY_NAME, p.CertAuthorityName)
	d.Set(CERT_CONTENT, p.CertContent)

	return diags
}

func updateClientCertAuthority(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	certname := d.Get(CERT_AUTHORITY_NAME).(string)
	body := getClientCertAuthorityModelFromResource(d)

	params := all.NewUpdateClientCertAuthorityParamsWithContext(ctx).WithCertAuthorityName(certname).WithBody(body)
	_, err := state.Client.All.UpdateClientCertAuthority(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}

	// read after update to make sure tf state is in sync with broker
	return readClientCertAuthority(ctx, d, meta)
}

func deleteClientCertAuthority(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	certname := d.Get(CERT_AUTHORITY_NAME).(string)

	params := all.NewDeleteClientCertAuthorityParamsWithContext(ctx).WithCertAuthorityName(certname)
	_, err := state.Client.All.DeleteClientCertAuthority(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	return diags
}
