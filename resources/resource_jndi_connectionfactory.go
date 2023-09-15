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

// Main resource definition for Jndi Queue entities
func ResourceJndiConnectionFactory() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaJndiConnectionFactory(),
		// Provider CRUD functions
		CreateContext: createJndiCF,
		ReadContext:   readJndiCF,
		UpdateContext: updateJndiCF,
		DeleteContext: deleteJndiCF,
	}
}

func schemaJndiConnectionFactory() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		MSG_VPN_NAME: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		CONNECTION_FACTORY_NAME: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		TRANSPORT_CONNECT_TIMEOUT: {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     30000,
			Description: "The timeout for establishing an initial connection to the broker, in milliseconds. The default value is 30000.",
		},
		TRANSPORT_REPLY_TIMEOUT: {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     10000,
			Description: "The timeout for reading a reply from the broker, in milliseconds. The default value is 10000.",
		},
		TRANSPORT_CONNECT_RETRY_COUNT: {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
			Description: "The maximum number of retry attempts to establish an initial connection to the host or list of hosts. " +
				"The value \"0\" means a single attempt (no retries), and the value \"-1\" means to retry forever. The default value is 0.",
		},
		TRANSPORT_RECONNECT_RETRY_COUNT: {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  3,
			Description: "The maximum number of attempts to reconnect to the host or list of hosts after the connection has been lost. " +
				"The value \"-1\" means to retry forever. The default value is 3.",
		},
		TRANSPORT_RECONNECT_RETRY_WAIT: {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  3000,
			Description: "The amount of time before making another attempt to connect or reconnect to the host after the connection has " +
				"been lost, in milliseconds. The default value is 3000.",
		},
		TRANSPORT_CONNECT_RETRY_PER_HOST_COUNT: {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
			Description: "The maximum number of retry attempts to establish an initial connection to each host on the list of hosts. " +
				"The value \"0\" means a single attempt (no retries), and the value \"-1\" means to retry forever. The default value is 0.",
		},
		DYNAMIC_CREATE_DURABLE_ENDPOINT: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			Description: "Enable or disable whether a durable endpoint will be dynamically created on the broker when the client calls " +
				"\"Session.createDurableSubscriber()\" or \"Session.createQueue()\". The default value is false.",
		},
		XA_ENABLED: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			Description: "Enable or disable this as an XA Connection Factory. When enabled, the Connection Factory can be cast to " +
				"'XAConnectionFactory', 'XAQueueConnectionFactory' or 'XATopicConnectionFactory'. The default value is `false`.",
		},
		DMQ_ELIGIBLE_ENABLED: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			Description: "Enable or disable whether messages sent by the Publisher (Producer) are Dead Message Queue (DMQ) eligible by default." +
				"Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
		},
		ALLOW_DUPLICATE_CLIENT_ID_ENABLED: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			Description: "Enable or disable whether new JMS connections can use the same Client identifier (ID) as an existing connection." +
				"Changes to this attribute are synchronized to HA mates and replication sites via config-sync.",
		},
		DIRECT_TRANSPORT_ENABLED: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
			Description: "Enable or disable usage of the Direct Transport mode for sending non-persistent messages." +
				"When disabled, the Guaranteed Transport mode is used. Changes to this attribute are synchronized to HA mates and replication sites via config-sync.",
		},
		DYNAMIC_ENDPOINT_RESPECT_TTL_ENABLED: {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Enable or disable whether dynamically created durable and non-durable endpoints respect the message time-to-live (TTL) property. The default value is `true`..",
		},
		TRANSPORT_KEEPALIVE_INTERVAL: {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     3000,
			Description: "The interval between application-level keepalive messages, in milliseconds. The default value is `3000`",
		},
	}
}

func getJndiCFModelFromResource(d *schema.ResourceData) *models.MsgVpnJndiConnectionFactory {
	q := &models.MsgVpnJndiConnectionFactory{
		MsgVpnName:            d.Get(MSG_VPN_NAME).(string),
		ConnectionFactoryName: d.Get(CONNECTION_FACTORY_NAME).(string),
	}
	if v, s := d.GetOk(TRANSPORT_CONNECT_TIMEOUT); s {
		q.TransportConnectTimeout = int32(v.(int))
	}
	if v, s := d.GetOk(TRANSPORT_REPLY_TIMEOUT); s {
		q.TransportReadTimeout = int32(v.(int))
	}
	if v, s := d.GetOk(TRANSPORT_CONNECT_RETRY_COUNT); s {
		q.TransportConnectRetryCount = int32(v.(int))
	}
	if v, s := d.GetOk(TRANSPORT_RECONNECT_RETRY_COUNT); s {
		q.TransportReconnectRetryCount = int32(v.(int))
	}
	if v, s := d.GetOk(TRANSPORT_RECONNECT_RETRY_WAIT); s {
		q.TransportReconnectRetryWait = int32(v.(int))
	}
	if v, s := d.GetOk(TRANSPORT_CONNECT_RETRY_PER_HOST_COUNT); s {
		q.TransportConnectRetryPerHostCount = int32(v.(int))
	}
	if v, s := d.GetOk(XA_ENABLED); s {
		q.XaEnabled = v.(bool)
	}
	if v, s := d.GetOk(DMQ_ELIGIBLE_ENABLED); s {
		q.MessagingDefaultDmqEligibleEnabled = v.(bool)
	}
	if v, s := d.GetOk(DYNAMIC_CREATE_DURABLE_ENDPOINT); s {
		q.DynamicEndpointCreateDurableEnabled = v.(bool)
	}
	if v, s := d.GetOk(ALLOW_DUPLICATE_CLIENT_ID_ENABLED); s {
		q.AllowDuplicateClientIDEnabled = v.(bool)
	}
	if v, s := d.GetOk(DIRECT_TRANSPORT_ENABLED); s {
		q.TransportDirectTransportEnabled = v.(bool)
	}
	if v, s := d.GetOk(DYNAMIC_ENDPOINT_RESPECT_TTL_ENABLED); s {
		q.DynamicEndpointRespectTTLEnabled = v.(bool)
	}
	if v, s := d.GetOk(TRANSPORT_KEEPALIVE_INTERVAL); s {
		q.TransportKeepaliveInterval = int32(v.(int))
	}

	return q
}

func createJndiCF(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	body := getJndiCFModelFromResource(d)

	params := all.NewCreateMsgVpnJndiConnectionFactoryParams().WithContext(ctx).WithMsgVpnName(msgvpn).WithBody(body)
	_, err := state.Client.All.CreateMsgVpnJndiConnectionFactory(params, state.Auth)
	if err != nil {
		if semp, err2 := provider.GetSempError(err); err2 == nil {
			if *semp.Status == "ALREADY_EXISTS" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Connection Factory %s already exists. Going to import state from Broker", body.ConnectionFactoryName),
				})
			} else {
				return provider.AppendError(diags, err)
			}
		} else {
			return provider.AppendError(diags, err)
		}
	}
	d.SetId("sol_" + msgvpn + "_jcf_" + body.ConnectionFactoryName)
	return append(diags, readJndiCF(ctx, d, meta)...)
}

func readJndiCF(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	cfname := d.Get(CONNECTION_FACTORY_NAME).(string)

	params := all.NewGetMsgVpnJndiConnectionFactoryParams().WithContext(ctx).WithMsgVpnName(msgvpn).WithConnectionFactoryName(cfname)
	resp, err := state.Client.All.GetMsgVpnJndiConnectionFactory(params, state.Auth)
	if err != nil {
		if semp, err2 := provider.GetSempError(err); err2 == nil {
			if *semp.Status == provider.SEMP_STATUS_NOT_FOUND {
				d.SetId("")
				return append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("JNDI ConnectionFactory %s does not exist. Going to remove ressource from tfstate", cfname),
				})
			} else {
				return provider.AppendError(diags, err)
			}
		} else {
			return provider.AppendError(diags, err)
		}
	}
	obj := resp.Payload.Data

	d.Set(TRANSPORT_CONNECT_TIMEOUT, obj.TransportConnectTimeout)
	d.Set(TRANSPORT_REPLY_TIMEOUT, obj.TransportReadTimeout)
	d.Set(TRANSPORT_CONNECT_RETRY_COUNT, obj.TransportConnectRetryCount)
	d.Set(TRANSPORT_RECONNECT_RETRY_COUNT, obj.TransportReconnectRetryCount)
	d.Set(TRANSPORT_RECONNECT_RETRY_WAIT, obj.TransportReconnectRetryWait)
	d.Set(TRANSPORT_CONNECT_RETRY_PER_HOST_COUNT, obj.TransportConnectRetryPerHostCount)
	d.Set(DYNAMIC_CREATE_DURABLE_ENDPOINT, obj.DynamicEndpointCreateDurableEnabled)
	d.Set(XA_ENABLED, obj.XaEnabled)
	d.Set(DMQ_ELIGIBLE_ENABLED, obj.MessagingDefaultDmqEligibleEnabled)
	d.Set(ALLOW_DUPLICATE_CLIENT_ID_ENABLED, obj.AllowDuplicateClientIDEnabled)
	d.Set(DIRECT_TRANSPORT_ENABLED, obj.TransportDirectTransportEnabled)
	d.Set(DYNAMIC_ENDPOINT_RESPECT_TTL_ENABLED, obj.DynamicEndpointRespectTTLEnabled)
	d.Set(TRANSPORT_KEEPALIVE_INTERVAL, obj.TransportKeepaliveInterval)
	return diags
}

func updateJndiCF(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	cfname := d.Get(CONNECTION_FACTORY_NAME).(string)
	body := getJndiCFModelFromResource(d)

	params := all.NewUpdateMsgVpnJndiConnectionFactoryParams().WithContext(ctx).WithMsgVpnName(msgvpn).WithConnectionFactoryName(cfname).WithBody(body)
	_, err := state.Client.All.UpdateMsgVpnJndiConnectionFactory(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	return readJndiCF(ctx, d, meta)
}

func deleteJndiCF(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	cfname := d.Get(CONNECTION_FACTORY_NAME).(string)

	params := all.NewDeleteMsgVpnJndiConnectionFactoryParams().WithContext(ctx).WithMsgVpnName(msgvpn).WithConnectionFactoryName(cfname)
	_, err := state.Client.All.DeleteMsgVpnJndiConnectionFactory(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	return diags
}
