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
func ResourceRestDeliveryPointRestConsumer() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaRestDeliveryPointRestConsumer(),
		// Provider CRUD functions
		CreateContext: createRestDeliveryPointRestConsumer,
		ReadContext:   readRestDeliveryPointRestConsumer,
		UpdateContext: updateRestDeliveryPointRestConsumer,
		DeleteContext: deleteRestDeliveryPointRestConsumer,
	}
}

func schemaRestDeliveryPointRestConsumer() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		ENABLED: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
			Description: "Enable or disable the REST Consumer. When disabled, no connections are initiated " +
				"or messages delivered to this particular REST Consumer. The default value is true",
		},
		REST_DELIVERY_POINT_REST_CONSUMER_HTTP_METHOD: {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "post",
			Description: "The HTTP method to use (POST or PUT).The default value is post.",
		},
		REST_DELIVERY_POINT_REST_CONSUMER_MAX_POST_WAIT_TIME: {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  30,
			Description: "The maximum amount of time (in seconds) to wait for an HTTP POST response from the REST Consumer." +
				"Once this time is exceeded, the TCP connection is reset. The default value is 30",
		},
		REST_DELIVERY_POINT_REST_CONSUMER_OUTGOING_CONNECTION_COUNT: {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     3,
			Description: "The number of concurrent TCP connections open to the REST Consumer. The default value is 3.",
		},
		REST_DELIVERY_POINT_REST_CONSUMER_NAME: {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The name of the REST Consumer.",
		},
		REST_DELIVERY_POINT_REST_CONSUMER_RETRY_DELAY: {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     3,
			Description: "The number of seconds that must pass before retrying the remote REST Consumer connection. The default value is 3",
		},
		REST_DELIVERY_POINT_REST_CONSUMER_REMOTE_HOST: {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
			Description: "The IP address or DNS name to which the broker is to connect to deliver messages for the REST Consumer." +
				"A host value must be configured for the REST Consumer to be operationally up. The default value is ''",
		},
		REST_DELIVERY_POINT_REST_CONSUMER_REMOTE_PORT: {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     8080,
			Description: "The port associated with the host of the REST Consumer. The default value is 8080",
		},
		REST_DELIVERY_POINT_REST_CONSUMER_AUTHENTICATION_SCHEME: {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "http-basic",
			Description: "The authentication scheme used by the REST Consumer to login to the REST host." +
				"The default value is 'http-basic'. The allowed values and their meaning are:" +
				"'none' - Login with no authentication. This may be useful for anonymous connections or when a REST Consumer does not require authentication." +
				"'http-basic' - Login with a username and optional password according to HTTP Basic authentication as per RFC2616.",
		},
		REST_DELIVERY_POINT_REST_CONSUMER_HTTP_BASIC_USERNAME: {
			Type:     schema.TypeString,
			Default:  "",
			Optional: true,
			Description: "The username that the REST Consumer will use to login to the REST host. " +
				"Normally a username is only configured when basic authentication is selected for the REST Consumer" +
				"The default value is '' ",
		},
		REST_DELIVERY_POINT_REST_CONSUMER_HTTP_BASIC_PASSWORD: {
			Type:     schema.TypeString,
			Default:  "",
			Optional: true,
			Description: "The password for the username. This attribute is absent from a GET and not updated when absent in a PUT." +
				"The default value is '' ",
		},
		MSG_VPN_NAME: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the MSG VPN where the RDP gets created.",
		},
		REST_DELIVERY_POINT_NAME: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the Rest Delivery Point.",
		},
		REST_DELIVERY_POINT_REST_CONSUMER_TLS_ENABLED: {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Enable or disable encryption (TLS) for the REST Consumer. The default value is true",
		},
	}
}

func getRestDeliveryPointRestConsumerModelFromResource(d *schema.ResourceData) *models.MsgVpnRestDeliveryPointRestConsumer {
	q := &models.MsgVpnRestDeliveryPointRestConsumer{
		RestDeliveryPointName: d.Get(REST_DELIVERY_POINT_NAME).(string),
		MsgVpnName:            d.Get(MSG_VPN_NAME).(string),
		RestConsumerName:      d.Get(REST_DELIVERY_POINT_REST_CONSUMER_NAME).(string),
	}
	if v, s := d.GetOk(REST_DELIVERY_POINT_REST_CONSUMER_HTTP_METHOD); s {
		q.HTTPMethod = v.(string)
	}
	if v, s := d.GetOk(ENABLED); s {
		q.Enabled = v.(bool)
	}
	if v, s := d.GetOk(REST_DELIVERY_POINT_REST_CONSUMER_TLS_ENABLED); s {
		q.TLSEnabled = v.(bool)
	}
	if v, s := d.GetOk(REST_DELIVERY_POINT_REST_CONSUMER_MAX_POST_WAIT_TIME); s {
		q.MaxPostWaitTime = int32(v.(int))
	}
	if v, s := d.GetOk(REST_DELIVERY_POINT_REST_CONSUMER_OUTGOING_CONNECTION_COUNT); s {
		q.OutgoingConnectionCount = int32(v.(int))
	}
	if v, s := d.GetOk(REST_DELIVERY_POINT_REST_CONSUMER_RETRY_DELAY); s {
		q.RetryDelay = int32(v.(int))
	}
	if v, s := d.GetOk(REST_DELIVERY_POINT_REST_CONSUMER_REMOTE_HOST); s {
		q.RemoteHost = v.(string)
	}
	if v, s := d.GetOk(REST_DELIVERY_POINT_REST_CONSUMER_REMOTE_PORT); s {
		q.RemotePort = int64(v.(int))
	}
	if v, s := d.GetOk(REST_DELIVERY_POINT_REST_CONSUMER_AUTHENTICATION_SCHEME); s {
		q.AuthenticationScheme = v.(string)
	}
	if v, s := d.GetOk(REST_DELIVERY_POINT_REST_CONSUMER_HTTP_BASIC_USERNAME); s {
		q.AuthenticationHTTPBasicUsername = v.(string)
	}
	if v, s := d.GetOk(REST_DELIVERY_POINT_REST_CONSUMER_HTTP_BASIC_PASSWORD); s {
		q.AuthenticationHTTPBasicPassword = v.(string)
	}
	return q
}

func createRestDeliveryPointRestConsumer(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	body := getRestDeliveryPointRestConsumerModelFromResource(d)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	rdp_name := d.Get(REST_DELIVERY_POINT_NAME).(string)

	params := all.NewCreateMsgVpnRestDeliveryPointRestConsumerParamsWithContext(ctx).WithMsgVpnName(msgvpn).WithRestDeliveryPointName(rdp_name).WithBody(body)
	_, err := state.Client.All.CreateMsgVpnRestDeliveryPointRestConsumer(params, state.Auth)
	if err != nil {
		if semp, err2 := provider.GetSempError(err); err2 == nil {
			if *semp.Status == "ALREADY_EXISTS" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Rest Consumer %s  for RDP %s already exists. Going to import state from Broker", body.RestConsumerName, body.RestDeliveryPointName),
				})
			} else {
				return provider.AppendError(diags, err)
			}
		} else {
			return provider.AppendError(diags, err)
		}
	}
	d.SetId("sol_rdp_" + body.RestDeliveryPointName + "_rc_" + body.RestConsumerName)

	// read after create to make sure tf state is in sync with broker
	return append(diags, readRestDeliveryPointRestConsumer(ctx, d, meta)...)
}

func readRestDeliveryPointRestConsumer(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	rdp_name := d.Get(REST_DELIVERY_POINT_NAME).(string)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	rc_name := d.Get(REST_DELIVERY_POINT_REST_CONSUMER_NAME).(string)
	params := all.NewGetMsgVpnRestDeliveryPointRestConsumerParamsWithContext(ctx).WithMsgVpnName(msgvpn).WithRestDeliveryPointName(rdp_name).WithRestConsumerName(rc_name)
	resp, err := state.Client.All.GetMsgVpnRestDeliveryPointRestConsumer(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	p := resp.Payload.Data
	d.Set(MSG_VPN_NAME, p.MsgVpnName)
	d.Set(REST_DELIVERY_POINT_NAME, p.RestDeliveryPointName)
	d.Set(REST_DELIVERY_POINT_REST_CONSUMER_NAME, p.RestConsumerName)
	d.Set(REST_DELIVERY_POINT_REST_CONSUMER_HTTP_METHOD, p.HTTPMethod)
	d.Set(ENABLED, p.Enabled)
	d.Set(REST_DELIVERY_POINT_REST_CONSUMER_MAX_POST_WAIT_TIME, p.MaxPostWaitTime)
	d.Set(REST_DELIVERY_POINT_REST_CONSUMER_OUTGOING_CONNECTION_COUNT, p.OutgoingConnectionCount)
	d.Set(REST_DELIVERY_POINT_REST_CONSUMER_RETRY_DELAY, p.RetryDelay)
	d.Set(REST_DELIVERY_POINT_REST_CONSUMER_REMOTE_HOST, p.RemoteHost)
	d.Set(REST_DELIVERY_POINT_REST_CONSUMER_REMOTE_PORT, p.RemotePort)
	d.Set(REST_DELIVERY_POINT_REST_CONSUMER_AUTHENTICATION_SCHEME, p.AuthenticationScheme)
	d.Set(REST_DELIVERY_POINT_REST_CONSUMER_HTTP_BASIC_USERNAME, p.AuthenticationHTTPBasicUsername)
	return diags
}

func updateRestDeliveryPointRestConsumer(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	rdp_name := d.Get(REST_DELIVERY_POINT_NAME).(string)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	rc_name := d.Get(REST_DELIVERY_POINT_REST_CONSUMER_NAME).(string)
	//current_state := d.Get(ENABLED).(bool)
	body := getRestDeliveryPointRestConsumerModelFromResource(d)
	// disable rest consumer before applying changes. This has to be done for almost all changes thats why we do it per default always.
	err_disable := setEnableRestConsumerStatus(ctx, state, msgvpn, rdp_name, rc_name, false)
	if err_disable != nil {
		return provider.AppendError(diags, err_disable)
	}
	params := all.NewUpdateMsgVpnRestDeliveryPointRestConsumerParamsWithContext(ctx).WithMsgVpnName(msgvpn).WithRestDeliveryPointName(rdp_name).WithRestConsumerName(rc_name).WithBody(body)
	_, err := state.Client.All.UpdateMsgVpnRestDeliveryPointRestConsumer(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	// read after update to make sure tf state is in sync with broker
	return readRestDeliveryPointRestConsumer(ctx, d, meta)
}

func deleteRestDeliveryPointRestConsumer(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	rdp_name := d.Get(REST_DELIVERY_POINT_NAME).(string)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	rc_name := d.Get(REST_DELIVERY_POINT_REST_CONSUMER_NAME).(string)

	params := all.NewDeleteMsgVpnRestDeliveryPointRestConsumerParamsWithContext(ctx).WithMsgVpnName(msgvpn).WithRestDeliveryPointName(rdp_name).WithRestConsumerName(rc_name)
	_, err := state.Client.All.DeleteMsgVpnRestDeliveryPointRestConsumer(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	return diags
}

func setEnableRestConsumerStatus(ctx context.Context, state *provider.ProviderState, vpnName string, rdp_name string, rc_name string, enabled bool) error {
	body := &models.MsgVpnRestDeliveryPointRestConsumer{
		MsgVpnName:            vpnName,
		RestDeliveryPointName: rdp_name,
		Enabled:               enabled,
		RestConsumerName:      rc_name,
	}
	params := all.NewUpdateMsgVpnRestDeliveryPointRestConsumerParamsWithContext(ctx).WithMsgVpnName(vpnName).WithRestDeliveryPointName(rdp_name).WithRestConsumerName(rc_name).WithBody(body)
	_, err := state.Client.All.UpdateMsgVpnRestDeliveryPointRestConsumer(params, state.Auth)
	return err
}
