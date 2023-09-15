package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/solacecommunity/solace-go-semp-client/client/all"
	"github.com/solacecommunity/solace-terraform-provider/provider"
	res "github.com/solacecommunity/solace-terraform-provider/resources"
)

// Main resource definition for Client Usernames
func ResourceMsgVpn() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: map[string]*schema.Schema{
			provider.ENABLED: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable or disable the Message VPN.",
			},
			provider.DMR_ENABLED: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable or disable Dynamic Message Routing (DMR) for the Message VPN.",
			},
			provider.EVENT_LOG_TAG: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A prefix applied to all published Events in the Message VPN.",
			},
			res.MSG_VPN_NAME: {
				Type:     schema.TypeString,
				Required: true,
			},
			res.JNDI_ENABLED: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable or disable JNDI access for clients in the Message VPN.",
			},
			res.MAX_CONNECTION_COUNT: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of client connections to the Message VPN. The default is the maximum value supported by the platform.",
			},
			res.MAX_MSG_SPOOL_USAGE: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum message spool usage by the Message VPN, in megabytes. The default value is 0.",
			},
			provider.MAX_ENDPOINT_COUNT: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of Queues and Topic Endpoints that can be created in the Message VPN. The default value is 16000.",
			},
			provider.AUTHENTICATION_BASIC_ENABLED: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable or disable basic authentication for clients connecting to the Message VPN.",
			},
			provider.AUTHENTICATION_BASIC_TYPE: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of basic authentication to use for clients connecting to the Message VPN.",
			},
			provider.AUTHENTICATION_CLIENT_CERT_ENABLED: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable or disable client certificate authentication in the Message VPN.",
			},
			provider.AUTHORIZATION_TYPE: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of authorization to use for clients connecting to the Message VPN. ",
			},
			provider.SERVICE_AMQP_MAX_CONNECTION_COUNT: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of AMQP client connections that can be simultaneously connected to the Message VPN.",
			},
			provider.SERVICE_AMQP_TLS_ENABLED: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable or disable the use of encryption (TLS) for the AMQP service in the Message VPN.",
			},
			provider.SERVICE_AMQP_TLS_LISTEN_PORT: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The port number for AMQP clients that connect to the Message VPN over TLS. The port must be unique across the message backbone. A value of 0 means that the listen-port is unassigned and cannot be enabled.",
			},
			provider.SERVICE_MQTT_MAX_CONNECTION_COUNT: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of MQTT client connections that can be simultaneously connected to the Message VPN.",
			},
			provider.SERVICE_MQTT_TLS_ENABLED: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable or disable the use of encryption (TLS) for the MQTT service in the Message VPN.",
			},
			provider.SERVICE_MQTT_TLS_LISTEN_PORT: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The port number for MQTT clients that connect to the Message VPN over TLS.",
			},
		},
		ReadContext: readMsgVpn,
		Description: "Datasource for Solace Message VPN",
	}
}

func readMsgVpn(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	msgvpn := d.Get(res.MSG_VPN_NAME).(string)

	params := all.NewGetMsgVpnParamsWithContext(ctx).WithMsgVpnName(msgvpn)
	resp, err := state.Client.All.GetMsgVpn(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}

	c := resp.Payload.Data
	d.SetId("msgvpn_" + c.MsgVpnName)
	d.Set(res.JNDI_ENABLED, c.JndiEnabled)
	d.Set(res.MAX_CONNECTION_COUNT, c.MaxConnectionCount)
	d.Set(res.MAX_MSG_SPOOL_USAGE, c.MaxMsgSpoolUsage)
	d.Set(provider.MAX_ENDPOINT_COUNT, c.MaxEndpointCount)
	d.Set(provider.AUTHENTICATION_BASIC_ENABLED, c.AuthenticationBasicEnabled)
	d.Set(provider.AUTHENTICATION_BASIC_TYPE, c.AuthenticationBasicType)
	d.Set(provider.AUTHENTICATION_CLIENT_CERT_ENABLED, c.AuthenticationClientCertEnabled)
	d.Set(provider.AUTHORIZATION_TYPE, c.AuthorizationType)
	d.Set(provider.ENABLED, c.Enabled)
	d.Set(provider.DMR_ENABLED, c.DmrEnabled)
	d.Set(provider.EVENT_LOG_TAG, c.EventLogTag)
	d.Set(provider.SERVICE_AMQP_MAX_CONNECTION_COUNT, c.ServiceAmqpMaxConnectionCount)
	d.Set(provider.SERVICE_AMQP_TLS_ENABLED, c.ServiceAmqpTLSEnabled)
	d.Set(provider.SERVICE_AMQP_TLS_LISTEN_PORT, c.ServiceAmqpTLSListenPort)
	d.Set(provider.SERVICE_MQTT_MAX_CONNECTION_COUNT, c.ServiceMqttMaxConnectionCount)
	d.Set(provider.SERVICE_MQTT_TLS_ENABLED, c.ServiceMqttTLSEnabled)
	d.Set(provider.SERVICE_MQTT_TLS_LISTEN_PORT, c.ServiceMqttTLSListenPort)
	return diags
}
