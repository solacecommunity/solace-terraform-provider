package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solacecommunity/solace-terraform-provider/provider"
)

// Main resource definition for Jndi Queue entities
func ResourceLdapCliGroup() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaLdapCliGroup(),
		// Provider CRUD functions
		CreateContext: createLdapCliGroup,
		ReadContext:   readLdapCliGroup,
		UpdateContext: updateLdapCliGroup,
		DeleteContext: deleteLdapCliGroup,
	}
}

func schemaLdapCliGroup() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		LDAP_CLI_GROUPNAME: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		LDAP_CLI_GLOBAL_ACCESS_LEVEL: {
			Type:         schema.TypeString,
			Default:      "none",
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"none", "admin", "read-only", "read-write"}, false),
		},
		LDAP_CLI_MSGVPN_DEFAULT_ACCESS_LEVEL: {
			Type:         schema.TypeString,
			Default:      "none",
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"none", "read-only", "read-write"}, false),
		},
		LDAP_CLI_MSGVPN_ACCESS_LEVEL_EXCEPTION: {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"none", "read-only", "read-write"}, false),
			},
		},
	}
}

func createLdapCliGroup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*provider.ProviderState).SempV1Client
	ldap_cli_group_name := d.Get(LDAP_CLI_GROUPNAME).(string)
	ldap_cli_group_global_access_level := d.Get(LDAP_CLI_GLOBAL_ACCESS_LEVEL).(string)
	ldap_cli_group_msgvpn_default_access_level := d.Get(LDAP_CLI_MSGVPN_DEFAULT_ACCESS_LEVEL).(string)
	ldap_cli_group_msgvpn_access_level_exception := d.Get(LDAP_CLI_MSGVPN_ACCESS_LEVEL_EXCEPTION)

	err := client.CreateLdapCliGroup(ctx, ldap_cli_group_name)
	if err != nil {
		return provider.AppendError(diags, err)
	}

	err = client.SetLdapCliGroupGlobalAccessLevel(ctx, ldap_cli_group_name, ldap_cli_group_global_access_level)
	if err != nil {
		return provider.AppendError(diags, fmt.Errorf("cannot set global access-level for ldap cli group %s: %s", ldap_cli_group_name, err))
	}
	err = client.SetLdapCliGroupMessageVpnDefaultAccessLevel(ctx, ldap_cli_group_name, ldap_cli_group_msgvpn_default_access_level)
	if err != nil {
		return provider.AppendError(diags, fmt.Errorf("cannot set msgvpn default access-level for ldap cli group %s: %s", ldap_cli_group_name, err))
	}

	if exceptions, ok := ldap_cli_group_msgvpn_access_level_exception.(map[string]interface{}); ok {
		for vpn, accessLevel := range exceptions {
			err = client.CreateLdapCliGroupMessageVpnAccessLevelException(ctx, ldap_cli_group_name, vpn)
			if err != nil {
				return provider.AppendError(diags, fmt.Errorf("cannot create msgvpn access-level exception for ldap cli group %s and vpn %s: %s", ldap_cli_group_name, vpn, err))
			}
			err = client.SetLdapCliGroupMessageVpnAccesslevelException(ctx, ldap_cli_group_name, vpn, accessLevel.(string))
			if err != nil {
				return provider.AppendError(diags, fmt.Errorf("cannot set msgvpn access-level exception for ldap cli group %s and vpn %s: %s", ldap_cli_group_name, vpn, err))
			}
		}
	} else {
		return provider.AppendError(diags, fmt.Errorf("cannot read CLI_MSGVPN_ACCESS_LEVEL_EXCEPTIONs"))
	}

	d.SetId("sol_ldapclig_" + ldap_cli_group_name)
	return append(diags, readLdapCliGroup(ctx, d, meta)...)
}

func readLdapCliGroup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*provider.ProviderState).SempV1Client
	ldap_cli_group_name := d.Get(LDAP_CLI_GROUPNAME).(string)

	result, err := client.GetLdapCliGroup(ctx, ldap_cli_group_name)
	if err != nil {
		return provider.AppendError(diags, err)
	}

	if len(result.Rpc.Show.Authentication.Authentications.Authentication.AccessLevelConfiguration.Ldap.Group) == 1 {
		ldapCliGroup := result.Rpc.Show.Authentication.Authentications.Authentication.AccessLevelConfiguration.Ldap.Group[0]
		d.Set(LDAP_CLI_GROUPNAME, ldapCliGroup.Name)
		d.Set(LDAP_CLI_GLOBAL_ACCESS_LEVEL, ldapCliGroup.GlobalAccessLevel)
		d.Set(LDAP_CLI_MSGVPN_DEFAULT_ACCESS_LEVEL, ldapCliGroup.DefaultVpnAccessLevel)
		exceptions := map[string]string{}
		for _, ex := range ldapCliGroup.VpnAccessLevelException {
			exceptions[ex.VpnName] = ex.AccessLevel
		}
		d.Set(LDAP_CLI_MSGVPN_ACCESS_LEVEL_EXCEPTION, exceptions)

	} else {
		return provider.AppendError(diags, fmt.Errorf("query requires exactly one returned ldap cli group"))
	}

	return diags
}

func updateLdapCliGroup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*provider.ProviderState).SempV1Client
	ldap_cli_group_name := d.Get(LDAP_CLI_GROUPNAME).(string)
	ldap_cli_group_global_access_level := d.Get(LDAP_CLI_GLOBAL_ACCESS_LEVEL).(string)
	ldap_cli_group_msgvpn_default_access_level := d.Get(LDAP_CLI_MSGVPN_DEFAULT_ACCESS_LEVEL).(string)
	ldap_cli_group_msgvpn_access_level_exception := d.Get(LDAP_CLI_MSGVPN_ACCESS_LEVEL_EXCEPTION)

	err := client.SetLdapCliGroupGlobalAccessLevel(ctx, ldap_cli_group_name, ldap_cli_group_global_access_level)
	if err != nil {
		return provider.AppendError(diags, fmt.Errorf("cannot set global access-level for ldap cli gorup %s: %s", ldap_cli_group_name, err))
	}
	err = client.SetLdapCliGroupMessageVpnDefaultAccessLevel(ctx, ldap_cli_group_name, ldap_cli_group_msgvpn_default_access_level)
	if err != nil {
		return provider.AppendError(diags, fmt.Errorf("cannot set msgvpn default access-level for ldap cli gorup %s: %s", ldap_cli_group_name, err))
	}

	if d.HasChange(LDAP_CLI_MSGVPN_ACCESS_LEVEL_EXCEPTION) {
		target_ex := ldap_cli_group_msgvpn_access_level_exception.(map[string]interface{})
		target := iSMapTosSMap(&target_ex)
		group, err := client.GetLdapCliGroup(ctx, ldap_cli_group_name)
		if err != nil {
			return provider.AppendError(diags, fmt.Errorf("cannot read current msgvpn exceptions for ldap cli group %s: %s", ldap_cli_group_name, err))
		}
		current_ex := group.LdapCliGroupExceptions()

		// Get the delta of specified exceptions and currently assigned exceptions (c=toCreate, u=toUpdatem, d=toDelete)
		c, u, d := StringMapDelta(target, current_ex)
		for k := range c {
			err = client.CreateLdapCliGroupMessageVpnAccessLevelException(ctx, ldap_cli_group_name, k)
			if err != nil {
				return provider.AppendError(diags, fmt.Errorf("cannot create msgvpn exception for for ldap cli group %s and msgvpn %s: %s", ldap_cli_group_name, k, err))
			}
		}
		for k, v := range u {
			err = client.SetLdapCliGroupMessageVpnAccesslevelException(ctx, ldap_cli_group_name, k, v)
			if err != nil {
				return provider.AppendError(diags, fmt.Errorf("cannot update msgvpn exception for user %s, msgvpn %s and access level %s: %s", ldap_cli_group_name, k, v, err))
			}
		}
		for k := range d {
			err = client.DeleteLdapCliGroupMessageVpnAccessLevelException(ctx, ldap_cli_group_name, k)
			if err != nil {
				return provider.AppendError(diags, fmt.Errorf("cannot delete msgvpn exception for user %s and msgvpn %s: %s", ldap_cli_group_name, k, err))
			}
		}
	}

	return append(diags, readLdapCliGroup(ctx, d, meta)...)
}

func deleteLdapCliGroup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*provider.ProviderState).SempV1Client
	ldap_cli_group_name := d.Get(LDAP_CLI_GROUPNAME).(string)

	err := client.DeleteLdapCliGroup(ctx, ldap_cli_group_name)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	d.SetId("")
	return diags
}
