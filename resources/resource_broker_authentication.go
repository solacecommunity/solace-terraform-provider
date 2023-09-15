package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solacecommunity/solace-terraform-provider/provider"
	semp "github.com/solacecommunity/solace-terraform-provider/sempv1"
)

// Main resource definition for Broker
func ResourceBrokerAuthentication() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaBrokerAuthentication(),
		// Provider CRUD functions
		CreateContext: createBrokerAuthentication,
		ReadContext:   readBrokerAuthentication,
		UpdateContext: updateBrokerAuthentication,
		DeleteContext: deleteBrokerAuthentication,
	}
}

func schemaBrokerAuthentication() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		CLI_USER_DEFAULT_GLOBAL_ACCESS_LEVEL: {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      semp.SEMPv1_ACCESS_LEVEL_NONE,
			Description:  "Defines the global access-level for the default access-level on the broker for the user-class cli",
			ValidateFunc: validation.StringInSlice([]string{semp.SEMPv1_ACCESS_LEVEL_NONE, semp.SEMPv1_ACCESS_LEVEL_ADMIN, semp.SEMPv1_ACCESS_LEVEL_READ_ONLY, semp.SEMPv1_ACCESS_LEVEL_READ_WRITE}, false),
		},
		CLI_USER_DEFAULT_MESSAGEVPN_ACCESS_LEVEL: {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      semp.SEMPv1_ACCESS_LEVEL_NONE,
			Description:  "Defines the messagevpn access-level for the default access-level on the broker for the user-class cli",
			ValidateFunc: validation.StringInSlice([]string{semp.SEMPv1_ACCESS_LEVEL_NONE, semp.SEMPv1_ACCESS_LEVEL_READ_ONLY, semp.SEMPv1_ACCESS_LEVEL_READ_WRITE}, false),
		},
		CLI_USER_AUTH_TYPE: {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      semp.CLI_AUTH_TYPE_INTERNAL,
			Description:  "Method used to authenticate CLI users. If ldap or radius is used, make sure to set a proper profile name.",
			ValidateFunc: validation.StringInSlice([]string{semp.CLI_AUTH_TYPE_INTERNAL, semp.CLI_AUTH_TYPE_LDAP, semp.CLI_AUTH_TYPE_RADIUS}, false),
		},
		CLI_USER_AUTH_TYPE_PROFILE: {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Name of the profile used for ldap or radius authentication. Can be ignored for internal authentication.",
		},
		LDAP_ACCESS_LEVEL_GROUP_MEMBERSHIP_ATTRIBUTE_NAME: {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "groupMembership",
		},
	}
}
func createBrokerAuthentication(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("broker_default_authentication")
	diags = updateBrokerAuthentication(ctx, d, meta)
	return append(diags, readBrokerAuthentication(ctx, d, meta)...)
}

func readBrokerAuthentication(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*provider.ProviderState).SempV1Client
	cli_user_auth_type := d.Get(CLI_USER_AUTH_TYPE).(string)

	result, err := client.GetBrokerAuthentication(ctx)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	authType := result.Rpc.Show.Authentication.Authentications.Authentication.AuthType
	if authType == semp.CLI_AUTH_TYPE_INTERNAL_VAL {
		authType = semp.CLI_AUTH_TYPE_INTERNAL
	}
	if authType == semp.CLI_AUTH_TYPE_LDAP_VAL {
		authType = semp.CLI_AUTH_TYPE_LDAP
	}
	d.Set(CLI_USER_AUTH_TYPE, authType)
	d.Set(CLI_USER_DEFAULT_GLOBAL_ACCESS_LEVEL, result.Rpc.Show.Authentication.Authentications.Authentication.AccessLevelConfiguration.Default.GlobalAccessLevel)
	d.Set(CLI_USER_DEFAULT_MESSAGEVPN_ACCESS_LEVEL, result.Rpc.Show.Authentication.Authentications.Authentication.AccessLevelConfiguration.Default.DefaultVpnAccessLevel)
	d.Set(LDAP_ACCESS_LEVEL_GROUP_MEMBERSHIP_ATTRIBUTE_NAME, result.Rpc.Show.Authentication.Authentications.Authentication.AccessLevelConfiguration.Ldap.GroupMembershipAttributeName)
	if cli_user_auth_type == semp.CLI_AUTH_TYPE_LDAP || cli_user_auth_type == semp.CLI_AUTH_TYPE_RADIUS {
		d.Set(CLI_USER_AUTH_TYPE_PROFILE, result.Rpc.Show.Authentication.Authentications.Authentication.ProfileName)
	}

	return diags
}

func updateBrokerAuthentication(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*provider.ProviderState).SempV1Client

	// Change CLI user auth type
	if d.HasChanges(CLI_USER_AUTH_TYPE, CLI_USER_AUTH_TYPE_PROFILE) {
		authType := d.Get(CLI_USER_AUTH_TYPE).(string)
		authTypeProfile := d.Get(CLI_USER_AUTH_TYPE_PROFILE).(string)

		// make sure that a profile is set when using ldap or radius
		if authType == semp.CLI_AUTH_TYPE_LDAP || authType == semp.CLI_AUTH_TYPE_RADIUS {
			if authTypeProfile == "" {
				return provider.AppendError(diags, fmt.Errorf("no auth-type profile defined"))
			}
		}
		err := client.SetCliUserAuthType(ctx, authType, authTypeProfile)
		if err != nil {
			return provider.AppendError(diags, err)
		}
	}
	if d.HasChanges(CLI_USER_DEFAULT_GLOBAL_ACCESS_LEVEL) {
		access_level := d.Get(CLI_USER_DEFAULT_GLOBAL_ACCESS_LEVEL).(string)
		err := client.SetBrokerAuthentication(ctx, semp.CLI_AUTH_ACCESS_LEVEL_DEFAULT_TYPE_GLOBAL, access_level)
		if err != nil {
			return provider.AppendError(diags, err)
		}
	}
	if d.HasChanges(CLI_USER_DEFAULT_MESSAGEVPN_ACCESS_LEVEL) {
		access_level := d.Get(CLI_USER_DEFAULT_MESSAGEVPN_ACCESS_LEVEL).(string)
		err := client.SetBrokerAuthentication(ctx, semp.CLI_AUTH_ACCESS_LEVEL_DEFAULT_TYPE_MSGVPN, access_level)
		if err != nil {
			return provider.AppendError(diags, err)
		}
	}
	if d.HasChanges(LDAP_ACCESS_LEVEL_GROUP_MEMBERSHIP_ATTRIBUTE_NAME) {
		attribute := d.Get(LDAP_ACCESS_LEVEL_GROUP_MEMBERSHIP_ATTRIBUTE_NAME).(string)
		err := client.DeleteLdapCliGroupMembershipAttributeName(ctx)
		if err != nil {
			return provider.AppendError(diags, err)
		}
		err = client.SetLdapCliGroupMembershipAttributeName(ctx, attribute)
		if err != nil {
			return provider.AppendError(diags, err)
		}
	}
	return append(diags, readBrokerAuthentication(ctx, d, meta)...)
}

func deleteBrokerAuthentication(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	// Empty for obvious reasons
	return diags
}
