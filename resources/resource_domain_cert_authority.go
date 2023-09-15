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
func ResourceDomainCertAuthoritoy() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaDomainCertAuthority(),
		// Provider CRUD functions
		CreateContext: createDomainCertAuthority,
		ReadContext:   readDomainCertAuthority,
		UpdateContext: updateDomainCertAuthority,
		DeleteContext: deleteDomainCertAuthority,
	}
}

func schemaDomainCertAuthority() map[string]*schema.Schema {
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
			Description: "The PEM formatted content for the trusted root certificate of a domain Certificate Authority. The default value is `\"\"`",
		},
	}
}

func getDomainCertAuthorityModelFromResource(d *schema.ResourceData) *models.DomainCertAuthority {
	q := &models.DomainCertAuthority{
		CertAuthorityName: d.Get(CERT_AUTHORITY_NAME).(string),
		CertContent:       d.Get(CERT_CONTENT).(string),
	}
	return q
}

func createDomainCertAuthority(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	body := getDomainCertAuthorityModelFromResource(d)

	params := all.NewCreateDomainCertAuthorityParams().WithContext(ctx).WithBody(body)
	_, err := state.Client.All.CreateDomainCertAuthority(params, state.Auth)
	if err != nil {
		if semp, err2 := provider.GetSempError(err); err2 == nil {
			if *semp.Status == "ALREADY_EXISTS" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Domain Cert Authority %s already exists. Going to import state from Broker", body.CertAuthorityName),
				})
			} else {
				return provider.AppendError(diags, err)
			}
		} else {
			return provider.AppendError(diags, err)
		}
	}
	d.SetId("sol_dca_" + body.CertAuthorityName)

	// read after create to make sure tf state is in sync with broker
	return append(diags, readDomainCertAuthority(ctx, d, meta)...)
}

func readDomainCertAuthority(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	certname := d.Get(CERT_AUTHORITY_NAME).(string)
	params := all.NewGetDomainCertAuthorityParamsWithContext(ctx).WithCertAuthorityName(certname)
	resp, err := state.Client.All.GetDomainCertAuthority(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	p := resp.Payload.Data
	d.Set(CERT_AUTHORITY_NAME, p.CertAuthorityName)
	d.Set(CERT_CONTENT, p.CertContent)

	return diags
}

func updateDomainCertAuthority(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	certname := d.Get(CERT_AUTHORITY_NAME).(string)
	body := getDomainCertAuthorityModelFromResource(d)

	params := all.NewUpdateDomainCertAuthorityParamsWithContext(ctx).WithCertAuthorityName(certname).WithBody(body)
	_, err := state.Client.All.UpdateDomainCertAuthority(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}

	// read after update to make sure tf state is in sync with broker
	return readDomainCertAuthority(ctx, d, meta)
}

func deleteDomainCertAuthority(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	certname := d.Get(CERT_AUTHORITY_NAME).(string)

	params := all.NewDeleteDomainCertAuthorityParamsWithContext(ctx).WithCertAuthorityName(certname)
	_, err := state.Client.All.DeleteDomainCertAuthority(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	return diags
}
