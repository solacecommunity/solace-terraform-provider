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
func ResourceClientProfile() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaClientProfile(),
		// Provider CRUD functions
		CreateContext: createClientProfile,
		ReadContext:   readClientProfile,
		UpdateContext: updateClientProfile,
		DeleteContext: deleteClientProfile,
	}
}

func schemaClientProfile() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		MSG_VPN_NAME: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		CLIENT_PROFILE_NAME: {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The name of the Client Profile.",
		},
		ALLOW_BRIDGE_CONNECTIONS: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			Description: "Enable or disable allowing Bridge clients using the Client Profile to connect. " +
				"Changing this setting does not affect existing Bridge client connections",
		},
		ALLOW_CUT_THROUGH_FORWARDING: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			Description: "Enable or disable allowing clients using the Client Profile to bind to endpoints with " +
				"the cut-through forwarding delivery mode. Changing this value does not affect existing client connections",
		},
		ALLOW_GUARANTEED_ENDPOINT_CREATE: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			Description: "Enable or disable allowing clients using the Client Profile to create topic endponts or queues. " +
				"Changing this value does not affect existing client connections",
		},
		ALLOW_GUARANTEED_ENDPOINT_CREATE_DURABILITY: {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "all",
			Description: "	// The types of Queues and Topic Endpoints that clients using the client-profile can create. " +
				"Changing this value does not affect existing client connections. The default value is all. " +
				"The allowed values and their meaning are: " +
				"all - Client can create any type of endpoint. " +
				"durable - Client can create only durable endpoints. " +
				"non-durable - Client can create only non-durable endpoints. ",
		},
		ALLOW_SHARED_SUBSCRIPTIONS: {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable or disable allowing shared subscriptions. Changing this setting does not affect existing subscriptions",
		},
		ALLOW_GUARANTEED_MSG_SEND: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
			Description: "Enable or disable allowing clients using the Client Profile to send guaranteed messages. " +
				"Changing this setting does not affect existing client connections. The default value is true.",
		},
		ALLOW_GUARANTEED_MSG_RECEIVE: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
			Description: "Enable or disable allowing clients using the Client Profile to receive guaranteed messages. " +
				"Changing this setting does not affect existing client connections. The default value is true.",
		},
		ALLOW_DOWNGRADE_TO_PLAIN_TEXT: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			Description: "Enable or disable allowing a client using the Client Profile to downgrade an encrypted " +
				"connection to plain text. The default value is false",
		},
		ALLOW_TRANSACTED_SESSIONS_ENABLED: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			Description: "Enable or disable allowing clients using the Client Profile to establish transacted sessions. " +
				"Changing this setting does not affect existing client connections. The default value is false",
		},
		COMPRESSION_ENABLED: {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Enable or disable allowing clients using the Client Profile to use compression. The default value is `true`. Available since 2.10.",
		},
		SERVICE_SMF_MAX_CONNECTION_COUNT_PER_CLIENT_USERNAME: {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     1000,
			Description: "The maximum number of SMF client connections per Client Username using the Client Profile. The default is the maximum value supported by the platform",
		},
		MAX_TRANSACTED_SESSION_COUNT: {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     10,
			Description: "The maximum number of transacted sessions that can be created by one client using the Client Profile. The default value is `10`.",
		},
	}
}

func getClientProfileModelFromResource(d *schema.ResourceData) *models.MsgVpnClientProfile {
	q := &models.MsgVpnClientProfile{
		MsgVpnName:        d.Get(MSG_VPN_NAME).(string),
		ClientProfileName: d.Get(CLIENT_PROFILE_NAME).(string),
	}
	if v, s := d.GetOk(ALLOW_TRANSACTED_SESSIONS_ENABLED); s {
		q.AllowTransactedSessionsEnabled = v.(bool)
	}
	if v, s := d.GetOk(ALLOW_GUARANTEED_MSG_SEND); s {
		q.AllowGuaranteedMsgSendEnabled = v.(bool)
	}
	if v, s := d.GetOk(ALLOW_GUARANTEED_MSG_RECEIVE); s {
		q.AllowGuaranteedMsgReceiveEnabled = v.(bool)
	}
	if v, s := d.GetOk(ALLOW_DOWNGRADE_TO_PLAIN_TEXT); s {
		q.TLSAllowDowngradeToPlainTextEnabled = v.(bool)
	}
	if v, s := d.GetOk(ALLOW_BRIDGE_CONNECTIONS); s {
		q.AllowBridgeConnectionsEnabled = v.(bool)
	}
	if v, s := d.GetOk(ALLOW_CUT_THROUGH_FORWARDING); s {
		q.AllowCutThroughForwardingEnabled = v.(bool)
	}
	if v, s := d.GetOk(ALLOW_GUARANTEED_ENDPOINT_CREATE); s {
		q.AllowGuaranteedEndpointCreateEnabled = v.(bool)
	}
	if v, s := d.GetOk(ALLOW_GUARANTEED_ENDPOINT_CREATE_DURABILITY); s {
		q.AllowGuaranteedEndpointCreateDurability = v.(string)
	}
	if v, s := d.GetOk(ALLOW_SHARED_SUBSCRIPTIONS); s {
		q.AllowSharedSubscriptionsEnabled = v.(bool)
	}
	if v, s := d.GetOk(COMPRESSION_ENABLED); s {
		q.CompressionEnabled = v.(bool)
	}
	if v, s := d.GetOk(SERVICE_SMF_MAX_CONNECTION_COUNT_PER_CLIENT_USERNAME); s {
		q.ServiceSmfMaxConnectionCountPerClientUsername = int64(v.(int))
	}
	if v, s := d.GetOk(MAX_TRANSACTED_SESSION_COUNT); s {
		q.MaxTransactedSessionCount = int64(v.(int))
	}

	return q
}

func createClientProfile(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	msgvpn := d.Get(MSG_VPN_NAME).(string)

	body := getClientProfileModelFromResource(d)
	params := all.NewCreateMsgVpnClientProfileParamsWithContext(ctx).WithMsgVpnName(msgvpn).WithBody(body)
	_, err := state.Client.All.CreateMsgVpnClientProfile(params, state.Auth)
	if err != nil {
		if semp, err2 := provider.GetSempError(err); err2 == nil {
			if *semp.Status == "ALREADY_EXISTS" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Client Profile %s already exists. Going to import state from Broker", body.ClientProfileName),
				})
			} else {
				return provider.AppendError(diags, err)
			}
		} else {
			return provider.AppendError(diags, err)
		}
	}
	d.SetId("sol_" + msgvpn + "_cp_" + body.ClientProfileName)

	// read after create to make sure tf state is in sync with broker
	return append(diags, readClientProfile(ctx, d, meta)...)
}

func readClientProfile(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	profile := d.Get(CLIENT_PROFILE_NAME).(string)

	params := all.NewGetMsgVpnClientProfileParamsWithContext(ctx).WithMsgVpnName(msgvpn).WithClientProfileName(profile)
	resp, err := state.Client.All.GetMsgVpnClientProfile(params, state.Auth)
	if err != nil {
		if semp, err2 := provider.GetSempError(err); err2 == nil {
			if *semp.Status == provider.SEMP_STATUS_NOT_FOUND {
				// the ressource must have been deleted on the broker
				d.SetId("")
				return append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Client Profile %s does not exist. Going to remove ressource from tfstate", profile),
				})
			} else {
				return provider.AppendError(diags, err)
			}
		} else {
			return provider.AppendError(diags, err)
		}
	}
	c := resp.Payload.Data
	d.Set(ALLOW_GUARANTEED_MSG_RECEIVE, c.AllowGuaranteedMsgReceiveEnabled)
	d.Set(ALLOW_GUARANTEED_MSG_SEND, c.AllowGuaranteedMsgSendEnabled)
	d.Set(ALLOW_DOWNGRADE_TO_PLAIN_TEXT, c.TLSAllowDowngradeToPlainTextEnabled)
	d.Set(ALLOW_BRIDGE_CONNECTIONS, c.AllowBridgeConnectionsEnabled)
	d.Set(ALLOW_CUT_THROUGH_FORWARDING, c.AllowCutThroughForwardingEnabled)
	d.Set(ALLOW_GUARANTEED_ENDPOINT_CREATE, c.AllowGuaranteedEndpointCreateEnabled)
	d.Set(ALLOW_GUARANTEED_ENDPOINT_CREATE_DURABILITY, c.AllowTransactedSessionsEnabled)
	d.Set(ALLOW_TRANSACTED_SESSIONS_ENABLED, c.AllowGuaranteedEndpointCreateDurability)
	d.Set(ALLOW_SHARED_SUBSCRIPTIONS, c.AllowSharedSubscriptionsEnabled)
	d.Set(COMPRESSION_ENABLED, c.CompressionEnabled)
	d.Set(SERVICE_SMF_MAX_CONNECTION_COUNT_PER_CLIENT_USERNAME, c.ServiceSmfMaxConnectionCountPerClientUsername)
	d.Set(MAX_TRANSACTED_SESSION_COUNT, c.MaxTransactedSessionCount)
	return diags
}

func updateClientProfile(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	profile := d.Get(CLIENT_PROFILE_NAME).(string)
	body := getClientProfileModelFromResource(d)

	params := all.NewUpdateMsgVpnClientProfileParamsWithContext(ctx).WithMsgVpnName(msgvpn).WithClientProfileName(profile).WithBody(body)
	_, err := state.Client.All.UpdateMsgVpnClientProfile(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	// read after update to make sure tf state is in sync with broker
	return readClientProfile(ctx, d, meta)
}

func deleteClientProfile(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	profile := d.Get(CLIENT_PROFILE_NAME).(string)

	params := all.NewDeleteMsgVpnClientProfileParamsWithContext(ctx).WithMsgVpnName(msgvpn).WithClientProfileName(profile)
	_, err := state.Client.All.DeleteMsgVpnClientProfile(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	return diags
}
