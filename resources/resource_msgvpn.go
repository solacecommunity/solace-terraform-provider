package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solacecommunity/solace-go-semp-client/client/all"
	"github.com/solacecommunity/solace-go-semp-client/models"
	"github.com/solacecommunity/solace-terraform-provider/provider"
)

// Main resource definition for Client Usernames
func ResourceMsgVpn() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaMsgVpn(),
		// Provider CRUD functions
		CreateContext: createMsgVpn,
		ReadContext:   readMsgVpn,
		UpdateContext: updateMsgVpn,
		DeleteContext: deleteMsgVpn,
	}
}

func schemaMsgVpn() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		MSG_VPN_NAME: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		JNDI_ENABLED: {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Enable or disable JNDI access for clients in the Message VPN.",
		},
		MAX_CONNECTION_COUNT: {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "The maximum number of client connections to the Message VPN. The default is the maximum value supported by the platform.",
		},
		MAX_MSG_SPOOL_USAGE: {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "The maximum message spool usage by the Message VPN, in megabytes. The default value is 0.",
		},
		MSG_VPN_ENABLED: {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Enables or disables the Message VPN.",
		},
		AUTHENTICATION_BASIC_ENABLED: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
			Description: "Enable or disable basic authentication for clients connecting to the Message VPN." +
				"Basic authentication is authentication that involves the use of a username and password to prove identity." +
				"If a user provides credentials for a different authentication scheme, this setting is not applicable." +
				"Changes to this attribute are synchronized to HA mates and replication sites via config-sync.",
		},
		AUTHENTICATION_BASIC_TYPE: {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "none",
			Description: "The type of basic authentication to use for clients connecting to the Message VPN." +
				"Changes to this attribute are synchronized to HA mates and replication sites via config-sync." +
				"The default value is radius. The allowed values and their meaning are:" +
				"internal - Internal database. Authentication is against Client Usernames." +
				"ldap - LDAP authentication. An LDAP profile name must be provided." +
				"radius - RADIUS authentication. A RADIUS profile name must be provided." +
				"none - No authentication. Anonymous login allowed.",
			ValidateFunc: validation.StringInSlice([]string{"internal", "ldap", "radius", "none"}, false),
		},
		AUTHENTICATION_BASIC_PROFILE_NAME: {
			Type:     schema.TypeString,
			Optional: true,
			Description: "The name of the RADIUS or LDAP Profile to use for basic authentication." +
				"Changes to this attribute are synchronized to HA mates and replication sites via config-sync.",
		},
		SERVICE_MQTT_PLAIN_TEXT_ENABLED: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			Description: "Enable or disable the plain-text MQTT service in the Message VPN." +
				"Disabling causes clients currently connected to be disconnected." +
				"Changes to this attribute are synchronized to HA mates and replication sites via config-sync.",
		},
		SERVICE_MQTT_WEB_SOCKET_ENABLED: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			Description: "Enable or disable the plain-text MQTT WSS service in the Message VPN." +
				"Disabling causes clients currently connected to be disconnected." +
				"Changes to this attribute are synchronized to HA mates and replication sites via config-sync.",
		},
		SERVICE_AMQP_PLAIN_TEXT_ENABLED: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			Description: "Enable or disable the plain-text AMQP service in the Message VPN." +
				"Disabling causes clients currently connected to be disconnected." +
				"Changes to this attribute are synchronized to HA mates and replication sites via config-sync.",
		},
		SERVICE_REST_INCOMING_PLAIN_TEXT_ENABLED: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			Description: "Enable or disable the plain-text REST service in the Message VPN." +
				"Disabling causes clients currently connected to be disconnected." +
				"Changes to this attribute are synchronized to HA mates and replication sites via config-sync.",
		},
		MSG_VPN_DMR_ENABLED: {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Enables or disables the Message VPN DMR Function.",
		},
		SERVICE_SMF_PLAIN_TEXT_ENABLED: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			Description: "Enable or disable the plain-text SMF service in the Message VPN." +
				"Disabling causes clients currently connected to be disconnected." +
				"Changes to this attribute are synchronized to HA mates and replication sites via config-sync.",
		},
		SERVICE_WEB_PLAIN_TEXT_ENABLED: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			Description: "Enable or disable the plain-text web service in the Message VPN." +
				"Disabling causes clients currently connected to be disconnected." +
				"Changes to this attribute are synchronized to HA mates and replication sites via config-sync.",
		},
		SERVICE_MQTT_TLS_LISTEN_PORT: {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  8883,
			Description: "The port number for TLS MQTT clients that connect to the Message VPN." +
				"The port must be unique across the message backbone." +
				"A value of 0 means that the listen-port is unassigned and cannot be enabled.",
		},
		SERVICE_MQTT_TLS_ENABLED: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
			Description: "Enable or disable the use of encryption (TLS) for the MQTT service in the Message VPN." +
				"Disabling causes clients currently connected over TLS to be disconnected.",
		},
		SERVICE_AMQP_TLS_ENABLED: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
			Description: "Enable or disable the use of encryption (TLS) for the AMQP service in the Message VPN." +
				"Disabling causes clients currently connected over TLS to be disconnected.",
		},
		SERVICE_AMQP_TLS_LISTEN_PORT: {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  5671,
			Description: "The port number for AMQP clients that connect to the Message VPN over TLS." +
				"The port must be unique across the message backbone." +
				"A value of 0 means that the listen-port is unassigned and cannot be enabled.",
		},
		EVENT_LARGE_MSG_THRESHOLD: {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  5120,
			Description: "The threshold, in kilobytes, after which a message is considered to be large for the Message VPN." +
				"Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is 5120",
		},
		SERVICE_REST_TLS_ENABLED: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
			Description: "Enable or disable the use of encryption (TLS) for the REST service in the Message VPN." +
				"Disabling causes clients currently connected over TLS to be disconnected.",
		},
		SERVICE_REST_TLS_LISTEN_PORT: {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  9444,
			Description: "The port number for REST clients that connect to the Message VPN over TLS." +
				"The port must be unique across the message backbone." +
				"A value of 0 means that the listen-port is unassigned and cannot be enabled.",
		},
	}
}
func getMsgVpnModelFromResource(d *schema.ResourceData) *models.MsgVpn {
	m := &models.MsgVpn{
		MsgVpnName: d.Get(MSG_VPN_NAME).(string),
	}
	if v, s := d.GetOk(JNDI_ENABLED); s {
		m.JndiEnabled = v.(bool)
	}
	if v, s := d.GetOk(MAX_CONNECTION_COUNT); s {
		m.MaxConnectionCount = int64(v.(int))
	}
	if v, s := d.GetOk(MAX_MSG_SPOOL_USAGE); s {
		m.MaxMsgSpoolUsage = int64(v.(int))
	}
	if v, s := d.GetOk(MSG_VPN_ENABLED); s {
		m.Enabled = v.(bool)
	}
	if v, s := d.GetOk(MSG_VPN_DMR_ENABLED); s {
		m.DmrEnabled = v.(bool)
	}
	if v, s := d.GetOk(AUTHENTICATION_BASIC_ENABLED); s {
		m.AuthenticationBasicEnabled = v.(bool)
	}
	if v, s := d.GetOk(AUTHENTICATION_BASIC_TYPE); s {
		m.AuthenticationBasicType = v.(string)
	}
	if v, s := d.GetOk(AUTHENTICATION_BASIC_PROFILE_NAME); s {
		if d.Get(AUTHENTICATION_BASIC_TYPE).(string) != "internal" {
			m.AuthenticationBasicProfileName = v.(string)
		}
	}
	if v, s := d.GetOk(SERVICE_MQTT_PLAIN_TEXT_ENABLED); s {
		m.ServiceMqttPlainTextEnabled = v.(bool)
	}
	if v, s := d.GetOk(SERVICE_AMQP_PLAIN_TEXT_ENABLED); s {
		m.ServiceAmqpPlainTextEnabled = v.(bool)
	}
	if v, s := d.GetOk(SERVICE_MQTT_WEB_SOCKET_ENABLED); s {
		m.ServiceMqttWebSocketEnabled = v.(bool)
	}
	if v, s := d.GetOk(SERVICE_REST_INCOMING_PLAIN_TEXT_ENABLED); s {
		m.ServiceRestIncomingPlainTextEnabled = v.(bool)
	}
	if v, s := d.GetOk(SERVICE_SMF_PLAIN_TEXT_ENABLED); s {
		m.ServiceSmfPlainTextEnabled = v.(bool)
	}
	if v, s := d.GetOk(SERVICE_WEB_PLAIN_TEXT_ENABLED); s {
		m.ServiceWebPlainTextEnabled = v.(bool)
	}
	if v, s := d.GetOk(SERVICE_MQTT_TLS_LISTEN_PORT); s {
		m.ServiceMqttTLSListenPort = int64(v.(int))
	}
	if v, s := d.GetOk(SERVICE_MQTT_TLS_ENABLED); s {
		m.ServiceMqttTLSEnabled = v.(bool)
	}
	if v, s := d.GetOk(SERVICE_AMQP_TLS_ENABLED); s {
		m.ServiceAmqpTLSEnabled = v.(bool)
	}
	if v, s := d.GetOk(SERVICE_AMQP_TLS_LISTEN_PORT); s {
		m.ServiceAmqpTLSListenPort = int64(v.(int))
	}
	if v, s := d.GetOk(EVENT_LARGE_MSG_THRESHOLD); s {
		m.EventLargeMsgThreshold = int64(v.(int))
	}
	if v, s := d.GetOk(SERVICE_REST_TLS_ENABLED); s {
		m.ServiceRestIncomingTLSEnabled = v.(bool)
	}
	if v, s := d.GetOk(SERVICE_REST_TLS_LISTEN_PORT); s {
		m.ServiceRestIncomingTLSListenPort = int64(v.(int))
	}

	return m
}

func createMsgVpn(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)

	body := getMsgVpnModelFromResource(d)
	params := all.NewCreateMsgVpnParamsWithContext(ctx).WithBody(body)
	_, err := state.Client.All.CreateMsgVpn(params, state.Auth)
	if err != nil {
		if semp, err2 := provider.GetSempError(err); err2 == nil {
			if *semp.Status == "ALREADY_EXISTS" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("MsgVpn %s already exists. Going to import state from Broker", body.MsgVpnName),
				})
			} else {
				return provider.AppendError(diags, err)
			}
		} else {
			return provider.AppendError(diags, err)
		}
	}
	d.SetId("sol_" + body.MsgVpnName)

	// read after create to make sure tf state is in sync with broker
	return append(diags, readMsgVpn(ctx, d, meta)...)
}

func readMsgVpn(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	msgvpn := d.Get(MSG_VPN_NAME).(string)

	params := all.NewGetMsgVpnParamsWithContext(ctx).WithMsgVpnName(msgvpn)
	resp, err := state.Client.All.GetMsgVpn(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	c := resp.Payload.Data
	d.Set(JNDI_ENABLED, c.JndiEnabled)
	d.Set(MAX_CONNECTION_COUNT, c.MaxConnectionCount)
	d.Set(MAX_MSG_SPOOL_USAGE, c.MaxMsgSpoolUsage)
	d.Set(MSG_VPN_ENABLED, c.Enabled)
	d.Set(MSG_VPN_DMR_ENABLED, c.DmrEnabled)
	d.Set(AUTHENTICATION_BASIC_ENABLED, c.AuthenticationBasicEnabled)
	d.Set(AUTHENTICATION_BASIC_PROFILE_NAME, c.AuthenticationBasicProfileName)
	d.Set(AUTHENTICATION_BASIC_TYPE, c.AuthenticationBasicType)
	d.Set(SERVICE_MQTT_PLAIN_TEXT_ENABLED, c.ServiceMqttPlainTextEnabled)
	d.Set(SERVICE_AMQP_PLAIN_TEXT_ENABLED, c.ServiceAmqpPlainTextEnabled)
	d.Set(SERVICE_MQTT_WEB_SOCKET_ENABLED, c.ServiceMqttWebSocketEnabled)
	d.Set(SERVICE_REST_INCOMING_PLAIN_TEXT_ENABLED, c.ServiceRestIncomingPlainTextEnabled)
	d.Set(SERVICE_SMF_PLAIN_TEXT_ENABLED, c.ServiceSmfPlainTextEnabled)
	d.Set(SERVICE_WEB_PLAIN_TEXT_ENABLED, c.ServiceWebPlainTextEnabled)
	d.Set(SERVICE_MQTT_TLS_LISTEN_PORT, c.ServiceMqttTLSListenPort)
	d.Set(SERVICE_MQTT_TLS_ENABLED, c.ServiceMqttTLSEnabled)
	d.Set(SERVICE_AMQP_TLS_ENABLED, c.ServiceAmqpTLSEnabled)
	d.Set(SERVICE_AMQP_TLS_LISTEN_PORT, c.ServiceAmqpTLSListenPort)
	d.Set(EVENT_LARGE_MSG_THRESHOLD, c.EventLargeMsgThreshold)
	d.Set(SERVICE_REST_TLS_ENABLED, c.ServiceRestIncomingTLSEnabled)
	d.Set(SERVICE_REST_TLS_LISTEN_PORT, c.ServiceRestIncomingTLSListenPort)
	return diags
}

func updateMsgVpn(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	body := getMsgVpnModelFromResource(d)
	// we have to use this special function the default msgvpn to disable / enable it.
	if msgvpn == "default" {
		client := meta.(*provider.ProviderState).SempV1Client
		if d.Get(MSG_VPN_ENABLED) == false {
			err := client.DisableDefaultMessageVPN(ctx)
			if err != nil {
				return provider.AppendError(diags, err)
			}
		}
		if d.Get(MSG_VPN_ENABLED) == true {
			err := client.EnableDefaultMessageVPN(ctx)
			if err != nil {
				return provider.AppendError(diags, err)
			}
		}

	}
	params := all.NewUpdateMsgVpnParamsWithContext(ctx).WithMsgVpnName(msgvpn).WithBody(body)
	_, err := state.Client.All.UpdateMsgVpn(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	// read after update to make sure tf state is in sync with broker
	return readMsgVpn(ctx, d, meta)
}

func deleteMsgVpn(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	msgvpn := d.Get(MSG_VPN_NAME).(string)

	params := all.NewDeleteMsgVpnParamsWithContext(ctx).WithMsgVpnName(msgvpn)
	_, err := state.Client.All.DeleteMsgVpn(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	return diags
}
