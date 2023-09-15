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
func ResourceClientUsername() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaClientUsername(),
		// Provider CRUD functions
		CreateContext: createClientUsername,
		ReadContext:   readClientUsername,
		UpdateContext: updateClientUsername,
		DeleteContext: deleteClientUsername,
	}
}

func schemaClientUsername() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		MSG_VPN_NAME: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		CLIENT_USERNAME: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		ENABLED: {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Enable or disable the Client Username. When disabled, all clients currently connected as the Client Username are disconnected. The default value is true",
		},
		ACL_PROFILE_NAME: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The ACL Profile of the Client Username.",
		},
		CLIENT_PROFILE_NAME: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The Client Profile of the Client Username.",
		},
		PASSWORD: {
			Type:        schema.TypeString,
			Sensitive:   true,
			Optional:    true,
			Description: "The password for the Client Username. Should not be set when the Broker is set for LDAP authentication.",
		},
		SUBSCRIPTION_MANAGER_ENABLED: {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable or disable the subscription management capability of the Client Username. This is the ability to manage subscriptions on behalf of other Client Usernames. The default value is false.",
		},
	}
}

func getClientUsernameModelFromResource(d *schema.ResourceData) *models.MsgVpnClientUsername {
	q := &models.MsgVpnClientUsername{
		MsgVpnName:                 d.Get(MSG_VPN_NAME).(string),
		ClientUsername:             d.Get(CLIENT_USERNAME).(string),
		Enabled:                    d.Get(ENABLED).(bool),
		ACLProfileName:             d.Get(ACL_PROFILE_NAME).(string),
		ClientProfileName:          d.Get(CLIENT_PROFILE_NAME).(string),
		SubscriptionManagerEnabled: d.Get(SUBSCRIPTION_MANAGER_ENABLED).(bool),
	}
	if v, s := d.GetOk(PASSWORD); s {
		q.Password = v.(string)
	}
	return q
}

func createClientUsername(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	msgvpn := d.Get(MSG_VPN_NAME).(string)

	body := getClientUsernameModelFromResource(d)
	params := all.NewCreateMsgVpnClientUsernameParams().WithContext(ctx).WithMsgVpnName(msgvpn).WithBody(body)
	_, err := state.Client.All.CreateMsgVpnClientUsername(params, state.Auth)
	if err != nil {
		if semp, err2 := provider.GetSempError(err); err2 == nil {
			if *semp.Status == provider.SEMP_STATUS_ALREADY_EXISTS {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Client username %s already exists. Going to import state from Broker", body.ClientUsername),
				})
			} else {
				return provider.AppendError(diags, err)
			}
		} else {
			return provider.AppendError(diags, err)
		}
	}
	d.SetId("sol_" + msgvpn + "_cu_" + body.ClientUsername)
	return diags
}

func readClientUsername(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	user := d.Get(CLIENT_USERNAME).(string)

	params := all.NewGetMsgVpnClientUsernameParams().WithContext(ctx).WithMsgVpnName(msgvpn).WithClientUsername(user)
	resp, err := state.Client.All.GetMsgVpnClientUsername(params, state.Auth)
	if err != nil {
		if semp, err2 := provider.GetSempError(err); err2 == nil {
			if *semp.Status == provider.SEMP_STATUS_NOT_FOUND {
				d.SetId("")
				return append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Client Username %s does not exist. Going to remove ressource from tfstate", user),
				})
			} else {
				return provider.AppendError(diags, err)
			}
		} else {
			return provider.AppendError(diags, err)
		}
	}
	c := resp.Payload.Data
	d.Set(ENABLED, c.Enabled)
	d.Set(CLIENT_PROFILE_NAME, c.ClientProfileName)
	d.Set(ACL_PROFILE_NAME, c.ACLProfileName)
	d.Set(SUBSCRIPTION_MANAGER_ENABLED, c.SubscriptionManagerEnabled)
	return diags
}

func updateClientUsername(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	user := d.Get(CLIENT_USERNAME).(string)

	body := getClientUsernameModelFromResource(d)

	params := all.NewUpdateMsgVpnClientUsernameParams().WithContext(ctx).WithMsgVpnName(msgvpn).WithClientUsername(user).WithBody(body)
	_, err := state.Client.All.UpdateMsgVpnClientUsername(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	return diags
}

func deleteClientUsername(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	user := d.Get(CLIENT_USERNAME).(string)

	params := all.NewDeleteMsgVpnClientUsernameParams().WithContext(ctx).WithMsgVpnName(msgvpn).WithClientUsername(user)
	_, err := state.Client.All.DeleteMsgVpnClientUsername(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	return diags
}
