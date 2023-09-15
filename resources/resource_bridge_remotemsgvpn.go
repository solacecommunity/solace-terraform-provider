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

// Main resource definition for BridgeRemoteMsgVpn entities
func ResourceBridgeRemoteMsgVpn() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaBridgeRemoteMsgVpn(),
		// Provider CRUD functions
		CreateContext: createBridgeRemoteMsgVpnFunc,
		ReadContext:   readBridgeRemoteMsgVpnFunc,
		UpdateContext: updateBridgeRemoteMsgVpnFunc,
		DeleteContext: deleteBridgeRemoteMsgVpnFunc,
	}
}

func schemaBridgeRemoteMsgVpn() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		BRIDGE_REMOTE_MSGVPN_BRIDGE_NAME: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the Bridge.",
		},
		BRIDGE_REMOTE_MSGVPN_MESSAGE_VPN: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the Message VPN.",
		},
		BRIDGE_REMOTE_MSGVPN_CLIENT_USERNAME: {
			Type:     schema.TypeString,
			Required: true,
			Description: "The Client Username the Bridge uses to login to the remote Message VPN. This per remote Message VPN value overrides " +
				"the value provided for the Bridge overall. Changes to this attribute are synchronized to HA mates and replication sites via config-sync.",
		},
		BRIDGE_REMOTE_MSGVPN_BRIDGE_VIRTUAL_ROUTER: {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "auto",
			Description: "The virtual router of the Bridge. The allowed values and their meaning are:" +
				"'primary' - The Bridge is used for the primary virtual router." +
				"'backup' - The Bridge is used for the backup virtual router." +
				"'auto' - The Bridge is automatically assigned a virtual router at creation, depending on the broker's active-standby role.",
		},
		BRIDGE_REMOTE_MSGVPN_COMPRESSED_DATA_ENABLED: {
			Type:     schema.TypeBool,
			Default:  false,
			Optional: true,
			Description: "Enable or disable data compression for the remote Message VPN connection. " +
				"Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is false.",
		},
		BRIDGE_REMOTE_MSGVPN_CONNECTOR_ORDER: {
			Type:     schema.TypeInt,
			Default:  4,
			Optional: true,
			Description: "The preference given to incoming connections from remote Message VPN hosts, from 1 (highest priority) to 4 (lowest priority)." +
				"Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is 4.",
		},
		BRIDGE_REMOTE_MSGVPN_EGRESS_FLOW_WINDOW_SIZE: {
			Type:     schema.TypeInt,
			Default:  255,
			Optional: true,
			Description: "The number of outstanding guaranteed messages that can be transmitted over the remote Message VPN connection before an acknowledgement " +
				"is received. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is 255.",
		},
		BRIDGE_REMOTE_MSGVPN_ENABLED: {
			Type:        schema.TypeBool,
			Default:     false,
			Optional:    true,
			Description: "Enable or disable the remote Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is false.",
		},
		BRIDGE_REMOTE_MSGVPN_PASSWORD: {
			Type:     schema.TypeString,
			Required: true,
			Description: "The password for the Client Username. This attribute is absent from a GET and not updated when absent in a PUT. " +
				"Changes to this attribute are synchronized to HA mates and replication sites via config-sync.",
		},
		BRIDGE_REMOTE_MSGVPN_QUEUE_BINDING: {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The queue binding of the Bridge in the remote Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync.",
		},
		BRIDGE_REMOTE_MSGVPN_REMOTE_MSGVPN_INTERFACE: {
			Type:     schema.TypeString,
			Required: true,
			Description: "The physical interface on the local Message VPN host for connecting to the remote Message VPN. By default, an interface is chosen automatically (recommended)," +
				"but if specified, remoteMsgVpnLocation must not be a virtual router name.",
		},
		BRIDGE_REMOTE_MSGVPN_REMOTE_MSGVPN_LOCATION: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The location of the remote Message VPN as either an FQDN with port, IP address with port, or virtual router name (starting with 'v:').",
		},
		BRIDGE_REMOTE_MSGVPN_REMOTE_MSGVPN_NAME: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the remote Message VPN.",
		},
		BRIDGE_REMOTE_MSGVPN_TLS_ENABLED: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			Description: "Enable or disable encryption (TLS) for the remote Message VPN connection. Changes to this attribute are synchronized to HA mates and replication sites via config-sync." +
				"The default value is false.",
		},
		BRIDGE_REMOTE_MSGVPN_UNIDIRECTIONAL_CLIENT_PROFILE: {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "#client-profile",
			Description: "The Client Profile for the unidirectional Bridge of the remote Message VPN. The Client Profile must exist in the local Message VPN," +
				"and it is used only for the TCP parameters. Note that the default client profile has a TCP maximum window size of 2MB. " +
				"	Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is '#client-profile'.",
		},
	}
}

// Creates a BridgeRemoteMsgVpn model based on the terraform resource state.
func getBridgeRemoteMsgVpnModelFromResource(d *schema.ResourceData) *models.MsgVpnBridgeRemoteMsgVpn {
	q := &models.MsgVpnBridgeRemoteMsgVpn{
		MsgVpnName:            d.Get(BRIDGE_REMOTE_MSGVPN_MESSAGE_VPN).(string),
		BridgeName:            d.Get(BRIDGE_REMOTE_MSGVPN_BRIDGE_NAME).(string),
		BridgeVirtualRouter:   d.Get(BRIDGE_REMOTE_MSGVPN_BRIDGE_VIRTUAL_ROUTER).(string),
		RemoteMsgVpnInterface: d.Get(BRIDGE_REMOTE_MSGVPN_REMOTE_MSGVPN_INTERFACE).(string),
		RemoteMsgVpnLocation:  d.Get(BRIDGE_REMOTE_MSGVPN_REMOTE_MSGVPN_LOCATION).(string),
		RemoteMsgVpnName:      d.Get(BRIDGE_REMOTE_MSGVPN_REMOTE_MSGVPN_NAME).(string),
	}
	if v, s := d.GetOk(BRIDGE_REMOTE_MSGVPN_CLIENT_USERNAME); s {
		q.ClientUsername = v.(string)
	}
	if v, s := d.GetOk(BRIDGE_REMOTE_MSGVPN_COMPRESSED_DATA_ENABLED); s {
		q.CompressedDataEnabled = v.(bool)
	}
	if v, s := d.GetOk(BRIDGE_REMOTE_MSGVPN_CONNECTOR_ORDER); s {
		q.ConnectOrder = int32(v.(int))
	}
	if v, s := d.GetOk(BRIDGE_REMOTE_MSGVPN_EGRESS_FLOW_WINDOW_SIZE); s {
		q.EgressFlowWindowSize = int64(v.(int))
	}
	if v, s := d.GetOk(BRIDGE_REMOTE_MSGVPN_ENABLED); s {
		q.Enabled = v.(bool)
	}
	if v, s := d.GetOk(BRIDGE_REMOTE_MSGVPN_PASSWORD); s {
		q.Password = v.(string)
	}
	if v, s := d.GetOk(BRIDGE_REMOTE_MSGVPN_QUEUE_BINDING); s {
		q.QueueBinding = v.(string)
	}
	if v, s := d.GetOk(BRIDGE_REMOTE_MSGVPN_TLS_ENABLED); s {
		q.TLSEnabled = v.(bool)
	}
	if v, s := d.GetOk(BRIDGE_REMOTE_MSGVPN_UNIDIRECTIONAL_CLIENT_PROFILE); s {
		q.UnidirectionalClientProfile = v.(string)
	}
	return q
}

func createBridgeRemoteMsgVpnFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	body := getBridgeRemoteMsgVpnModelFromResource(d)
	params := all.NewCreateMsgVpnBridgeRemoteMsgVpnParamsWithContext(ctx).WithMsgVpnName(body.MsgVpnName).WithBridgeName(body.BridgeName).WithBridgeVirtualRouter(body.BridgeVirtualRouter).WithBody(body)
	_, err := state.Client.All.CreateMsgVpnBridgeRemoteMsgVpn(params, state.Auth)
	if err != nil {
		if semp, err2 := provider.GetSempError(err); err2 == nil {
			if *semp.Status == "ALREADY_EXISTS" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("BridgeRemoteMsgVpn %s already exists in Bridge %s. Going to import state from Broker", body.RemoteMsgVpnName, body.BridgeName),
				})
			} else {
				return provider.AppendError(diags, err)
			}
		} else {
			return provider.AppendError(diags, err)
		}
	}
	d.SetId("sol_BridgeRemoteMsgVpn_" + body.BridgeName + "_" + body.MsgVpnName + "_to_" + body.RemoteMsgVpnName)
	return append(diags, readBridgeRemoteMsgVpnFunc(ctx, d, meta)...)
}

func updateBridgeRemoteMsgVpnFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	body := getBridgeRemoteMsgVpnModelFromResource(d)
	params := all.NewUpdateMsgVpnBridgeRemoteMsgVpnParamsWithContext(ctx).WithMsgVpnName(body.MsgVpnName).WithBridgeName(body.BridgeName).WithBridgeVirtualRouter(body.BridgeVirtualRouter).WithRemoteMsgVpnName(body.RemoteMsgVpnName).WithRemoteMsgVpnLocation(body.RemoteMsgVpnLocation).WithRemoteMsgVpnInterface(body.RemoteMsgVpnInterface).WithBody(body)
	_, err := state.Client.All.UpdateMsgVpnBridgeRemoteMsgVpn(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	return readBridgeRemoteMsgVpnFunc(ctx, d, meta)
}

func readBridgeRemoteMsgVpnFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	BridgeRemoteMsgVpnMsgVpnName := d.Get(BRIDGE_REMOTE_MSGVPN_MESSAGE_VPN).(string)
	BridgeRemoteMsgVpnBridgeName := d.Get(BRIDGE_REMOTE_MSGVPN_BRIDGE_NAME).(string)
	BridgeRemoteMsgVpnVirtualRouter := d.Get(BRIDGE_REMOTE_MSGVPN_BRIDGE_VIRTUAL_ROUTER).(string)
	BridgeRemoteMsgVpnRemoteMsgVpnName := d.Get(BRIDGE_REMOTE_MSGVPN_REMOTE_MSGVPN_NAME).(string)
	BridgeRemoteMsgVpnRemoteMsgVpnInterface := d.Get(BRIDGE_REMOTE_MSGVPN_REMOTE_MSGVPN_INTERFACE).(string)
	BridgeRemoteMsgVpnRemoteMsgVpnLocation := d.Get(BRIDGE_REMOTE_MSGVPN_REMOTE_MSGVPN_LOCATION).(string)
	params := all.NewGetMsgVpnBridgeRemoteMsgVpnParamsWithContext(ctx).WithMsgVpnName(BridgeRemoteMsgVpnMsgVpnName).WithBridgeName(BridgeRemoteMsgVpnBridgeName).WithBridgeVirtualRouter(BridgeRemoteMsgVpnVirtualRouter).WithRemoteMsgVpnName(BridgeRemoteMsgVpnRemoteMsgVpnName).WithRemoteMsgVpnLocation(BridgeRemoteMsgVpnRemoteMsgVpnLocation).WithRemoteMsgVpnInterface(BridgeRemoteMsgVpnRemoteMsgVpnInterface)
	resp, err := state.Client.All.GetMsgVpnBridgeRemoteMsgVpn(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	c := resp.Payload.Data
	d.Set(BRIDGE_REMOTE_MSGVPN_BRIDGE_NAME, c.BridgeName)
	d.Set(BRIDGE_REMOTE_MSGVPN_BRIDGE_VIRTUAL_ROUTER, c.BridgeVirtualRouter)
	d.Set(BRIDGE_REMOTE_MSGVPN_CLIENT_USERNAME, c.ClientUsername)
	d.Set(BRIDGE_REMOTE_MSGVPN_COMPRESSED_DATA_ENABLED, c.CompressedDataEnabled)
	d.Set(BRIDGE_REMOTE_MSGVPN_CONNECTOR_ORDER, c.ConnectOrder)
	d.Set(BRIDGE_REMOTE_MSGVPN_EGRESS_FLOW_WINDOW_SIZE, c.EgressFlowWindowSize)
	d.Set(BRIDGE_REMOTE_MSGVPN_ENABLED, c.Enabled)
	d.Set(BRIDGE_REMOTE_MSGVPN_MESSAGE_VPN, c.MsgVpnName)
	d.Set(BRIDGE_REMOTE_MSGVPN_QUEUE_BINDING, c.QueueBinding)
	d.Set(BRIDGE_REMOTE_MSGVPN_REMOTE_MSGVPN_INTERFACE, c.RemoteMsgVpnInterface)
	d.Set(BRIDGE_REMOTE_MSGVPN_REMOTE_MSGVPN_LOCATION, c.RemoteMsgVpnLocation)
	d.Set(BRIDGE_REMOTE_MSGVPN_REMOTE_MSGVPN_NAME, c.RemoteMsgVpnName)
	d.Set(BRIDGE_REMOTE_MSGVPN_TLS_ENABLED, c.TLSEnabled)
	d.Set(BRIDGE_REMOTE_MSGVPN_UNIDIRECTIONAL_CLIENT_PROFILE, c.UnidirectionalClientProfile)
	return diags
}

func deleteBridgeRemoteMsgVpnFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	BridgeRemoteMsgVpnMsgVpnName := d.Get(BRIDGE_REMOTE_MSGVPN_MESSAGE_VPN).(string)
	BridgeRemoteMsgVpnBridgeName := d.Get(BRIDGE_REMOTE_MSGVPN_BRIDGE_NAME).(string)
	BridgeRemoteMsgVpnVirtualRouter := d.Get(BRIDGE_REMOTE_MSGVPN_BRIDGE_VIRTUAL_ROUTER).(string)
	BridgeRemoteMsgVpnRemoteMsgVpnName := d.Get(BRIDGE_REMOTE_MSGVPN_REMOTE_MSGVPN_NAME).(string)
	BridgeRemoteMsgVpnRemoteMsgVpnInterface := d.Get(BRIDGE_REMOTE_MSGVPN_REMOTE_MSGVPN_INTERFACE).(string)
	BridgeRemoteMsgVpnRemoteMsgVpnLocation := d.Get(BRIDGE_REMOTE_MSGVPN_REMOTE_MSGVPN_LOCATION).(string)
	params := all.NewDeleteMsgVpnBridgeRemoteMsgVpnParamsWithContext(ctx).WithMsgVpnName(BridgeRemoteMsgVpnMsgVpnName).WithBridgeName(BridgeRemoteMsgVpnBridgeName).WithBridgeVirtualRouter(BridgeRemoteMsgVpnVirtualRouter).WithRemoteMsgVpnName(BridgeRemoteMsgVpnRemoteMsgVpnName).WithRemoteMsgVpnLocation(BridgeRemoteMsgVpnRemoteMsgVpnLocation).WithRemoteMsgVpnInterface(BridgeRemoteMsgVpnRemoteMsgVpnInterface)
	_, err := state.Client.All.DeleteMsgVpnBridgeRemoteMsgVpn(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	return diags
}
