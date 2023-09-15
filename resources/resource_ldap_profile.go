package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/solacecommunity/solace-terraform-provider/provider"
	"github.com/solacecommunity/solace-terraform-provider/sempv1"
)

func ResourceLdapProfile() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Schema:        schemaLdapProfile(),
		CreateContext: createLdapProfile,
		ReadContext:   readLdapProfile,
		UpdateContext: updateLdapProfile,
		DeleteContext: deleteLdapProfile,
	}
}

func schemaLdapProfile() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		LDAP_PROFILE_NAME: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		LDAP_PROFILE_HOST: {
			Type:     schema.TypeString,
			Required: true,
		},
		LDAP_PROFILE_ENABLED: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		LDAP_PROFILE_INDEX: {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  1,
		},
		LDAP_PROFILE_ADMIN_DN: {
			Type:     schema.TypeString,
			Optional: true,
		},
		LDAP_PROFILE_ADMIN_PWD: {
			Type:      schema.TypeString,
			Sensitive: true,
			Optional:  true,
		},
		LDAP_PROFILE_BASE_DN: {
			Type:     schema.TypeString,
			Required: true,
		},
		LDAP_PROFILE_SEARCH_FILER: {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

func getLdapHelperFromSchema(d *schema.ResourceData) *sempv1.LdapProfileHelper {
	ldap := sempv1.LdapProfileHelper{
		Name:          d.Get(LDAP_PROFILE_NAME).(string),
		Enabled:       d.Get(LDAP_PROFILE_ENABLED).(bool),
		AdminDN:       d.Get(LDAP_PROFILE_ADMIN_DN).(string),
		AdminPassword: d.Get(LDAP_PROFILE_ADMIN_PWD).(string),
		Index:         d.Get(LDAP_PROFILE_INDEX).(int),
		Host:          d.Get(LDAP_PROFILE_HOST).(string),
		BaseDN:        d.Get(LDAP_PROFILE_BASE_DN).(string),
		SearchFilter:  d.Get(LDAP_PROFILE_SEARCH_FILER).(string),
	}

	return &ldap
}

func createLdapProfile(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)

	ldap := getLdapHelperFromSchema(d)
	ldap.Client = state.SempV1Client
	ldap.Context = ctx

	err := ldap.CreateLdapProfile()
	if err != nil {
		return provider.AppendError(diags, err)
	}

	d.SetId("sol_ldapp_" + ldap.Name)
	return append(diags, readLdapProfile(ctx, d, meta)...)
}

func readLdapProfile(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)

	ldap := sempv1.LdapProfileHelper{
		Client:  state.SempV1Client,
		Context: ctx,
		Name:    d.Get(LDAP_PROFILE_NAME).(string),
	}
	err := ldap.ReadLdapProfile()
	if err != nil {
		return provider.AppendError(diags, err)
	}

	d.Set(LDAP_PROFILE_NAME, ldap.Name)
	d.Set(LDAP_PROFILE_ENABLED, ldap.Enabled)
	d.Set(LDAP_PROFILE_ADMIN_DN, ldap.AdminDN)
	d.Set(LDAP_PROFILE_INDEX, ldap.Index)
	d.Set(LDAP_PROFILE_HOST, ldap.Host)
	d.Set(LDAP_PROFILE_BASE_DN, ldap.BaseDN)
	d.Set(LDAP_PROFILE_SEARCH_FILER, ldap.SearchFilter)

	return diags
}

func updateLdapProfile(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)

	ldap := getLdapHelperFromSchema(d)
	ldap.Client = state.SempV1Client
	ldap.Context = ctx

	err := ldap.UpdateLdapProfile()
	if err != nil {
		return provider.AppendError(diags, err)
	}

	return append(diags, readLdapProfile(ctx, d, meta)...)
}

func deleteLdapProfile(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)

	ldap := getLdapHelperFromSchema(d)
	ldap.Client = state.SempV1Client
	ldap.Context = ctx

	err := ldap.DeleteLdapProfile()
	if err != nil {
		return provider.AppendError(diags, err)
	}

	return diags
}
