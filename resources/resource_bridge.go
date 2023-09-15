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

// Main resource definition for Bridge entities
func ResourceBridge() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaBridge(),
		// Provider CRUD functions
		CreateContext: createBridgeFunc,
		ReadContext:   readBridgeFunc,
		UpdateContext: updateBridgeFunc,
		DeleteContext: deleteBridgeFunc,
	}
}

func schemaBridge() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		BRIDGE_NAME: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the Bridge.",
		},
		BRIDGE_MESSAGE_VPN: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the Message VPN.",
		},
		BRIDGE_ENABLED: {
			Type:     schema.TypeBool,
			Default:  false,
			Optional: true,
			Description: "Enable or disable the Bridge. Changes to this attribute are synchronized to HA mates" +
				"and replication sites via config-sync. The default value is false.",
		},
		BRIDGE_VIRTUAL_ROUTER: {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "auto",
			Description: "The virtual router of the Bridge. The allowed values and their meaning are:" +
				"'primary' - The Bridge is used for the primary virtual router." +
				"'backup' - The Bridge is used for the backup virtual router." +
				"'auto' - The Bridge is automatically assigned a virtual router at creation, depending on the broker's active-standby role.",
		},
		BRIDGE_MAX_TTL: {
			Type:     schema.TypeInt,
			Default:  8,
			Optional: true,
			Description: "The maximum time-to-live (TTL) in hops. Messages are discarded if their TTL exceeds this value. " +
				"Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is 8.",
		},
		BRIDGE_REMOTE_AUTHENTICATION_BASIC_CLIENT_USERNAME: {
			Type:     schema.TypeString,
			Required: true,
			Description: "The Client Username the Bridge uses to login to the remote Message VPN." +
				"Changes to this attribute are synchronized to HA mates and replication sites via config-sync.",
		},
		BRIDGE_REMOTE_AUTHENTICATION_BASIC_PASSWORD: {
			Type:     schema.TypeString,
			Required: true,
			Description: "The password for the Client Username. This attribute is absent from a GET and not updated when absent in a PUT," +
				"subject to the exceptions in note 4. Changes to this attribute are synchronized to HA mates and replication sites via config-sync.",
		},
		BRIDGE_REMOTE_AUTHENTICATION_SCHEME: {
			Type:     schema.TypeString,
			Default:  "basic",
			Optional: true,
			Description: "The authentication scheme for the remote Message VPN. Changes to this attribute are synchronized to HA mates and replication" +
				"sites via config-sync. The default value is 'basic'.",
		},
		BRIDGE_REMOTE_CONNECTION_RETRY_COUNT: {
			Type:     schema.TypeInt,
			Default:  0,
			Optional: true,
			Description: "The maximum number of retry attempts to establish a connection to the remote Message VPN." +
				"A value of 0 means to retry forever. Changes to this attribute are synchronized to HA mates and" +
				"replication sites via config-sync. The default value is 0.",
		},
		BRIDGE_REMOTE_CONNECTION_RETRY_DELAY: {
			Type:     schema.TypeInt,
			Default:  3,
			Optional: true,
			Description: "The number of seconds the broker waits for the bridge connection to be established before " +
				"attempting a new connection. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is 3.",
		},
		BRIDGE_REMOTE_DELIVER_TO_ONE_PRIORITY: {
			Type:     schema.TypeString,
			Default:  "p1",
			Optional: true,
			Description: "The priority for deliver-to-one (DTO) messages transmitted from the remote Message VPN. Changes to this attribute are synchronized" +
				"to HA mates and replication sites via config-sync. The default value is 'p1'. The allowed values and their meaning are:" +
				"'p1' - The 1st or highest priority." +
				"'p2' - The 2nd highest priority." +
				"'p3' - The 3rd highest priority." +
				"'p4' - The 4th highest priority." +
				"'da' - Ignore priority and deliver always.",
		},
		BRIDGE_TLS_CIPHER_SUITE_LIST: {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "default",
			Description: "The colon-separated list of cipher suites supported for TLS connections to the remote Message VPN. The value 'default'" +
				"implies all supported suites ordered from most secure to least secure." +
				"Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is 'default'.",
		},
	}
}

// Creates a bridge model based on the terraform resource state.
func getBridgeModelFromResource(d *schema.ResourceData) *models.MsgVpnBridge {
	q := &models.MsgVpnBridge{
		MsgVpnName:                              d.Get(BRIDGE_MESSAGE_VPN).(string),
		BridgeName:                              d.Get(BRIDGE_NAME).(string),
		RemoteAuthenticationBasicClientUsername: d.Get(BRIDGE_REMOTE_AUTHENTICATION_BASIC_CLIENT_USERNAME).(string),
		RemoteAuthenticationBasicPassword:       d.Get(BRIDGE_REMOTE_AUTHENTICATION_BASIC_PASSWORD).(string),
	}
	if v, s := d.GetOk(BRIDGE_ENABLED); s {
		q.Enabled = v.(bool)
	}
	if v, s := d.GetOk(BRIDGE_VIRTUAL_ROUTER); s {
		q.BridgeVirtualRouter = v.(string)
	}
	if v, s := d.GetOk(BRIDGE_MAX_TTL); s {
		q.MaxTTL = int64(v.(int))
	}
	if v, s := d.GetOk(BRIDGE_REMOTE_AUTHENTICATION_SCHEME); s {
		q.RemoteAuthenticationScheme = v.(string)
	}
	if v, s := d.GetOk(BRIDGE_REMOTE_CONNECTION_RETRY_COUNT); s {
		q.RemoteConnectionRetryCount = int64(v.(int))
	}
	if v, s := d.GetOk(BRIDGE_REMOTE_CONNECTION_RETRY_DELAY); s {
		q.RemoteConnectionRetryDelay = int64(v.(int))
	}
	if v, s := d.GetOk(BRIDGE_REMOTE_DELIVER_TO_ONE_PRIORITY); s {
		q.RemoteDeliverToOnePriority = v.(string)
	}
	if v, s := d.GetOk(BRIDGE_TLS_CIPHER_SUITE_LIST); s {
		q.TLSCipherSuiteList = v.(string)
	}
	return q
}

func createBridgeFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	body := getBridgeModelFromResource(d)
	params := all.NewCreateMsgVpnBridgeParamsWithContext(ctx).WithMsgVpnName(body.MsgVpnName).WithBody(body)
	_, err := state.Client.All.CreateMsgVpnBridge(params, state.Auth)
	if err != nil {
		if semp, err2 := provider.GetSempError(err); err2 == nil {
			if *semp.Status == "ALREADY_EXISTS" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Bridge %s already exists. Going to import state from Broker", body.BridgeName),
				})
			} else {
				return provider.AppendError(diags, err)
			}
		} else {
			return provider.AppendError(diags, err)
		}
	}

	d.SetId("sol_Bridge_" + body.BridgeName + "_" + body.MsgVpnName)
	return append(diags, readBridgeFunc(ctx, d, meta)...)
}

func updateBridgeFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	body := getBridgeModelFromResource(d)
	params := all.NewUpdateMsgVpnBridgeParamsWithContext(ctx).WithMsgVpnName(body.MsgVpnName).WithBridgeName(body.BridgeName).WithBridgeVirtualRouter(body.BridgeVirtualRouter).WithBody(body)
	_, err := state.Client.All.UpdateMsgVpnBridge(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}

	return readBridgeFunc(ctx, d, meta)
}

func readBridgeFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	MsgVpnName := d.Get(BRIDGE_MESSAGE_VPN).(string)
	BridgeName := d.Get(BRIDGE_NAME).(string)
	BridgeVirtualRouter := d.Get(BRIDGE_VIRTUAL_ROUTER).(string)
	params := all.NewGetMsgVpnBridgeParamsWithContext(ctx).WithMsgVpnName(MsgVpnName).WithBridgeName(BridgeName).WithBridgeVirtualRouter(BridgeVirtualRouter)
	resp, err := state.Client.All.GetMsgVpnBridge(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	c := resp.Payload.Data
	d.Set(BRIDGE_ENABLED, c.Enabled)
	d.Set(BRIDGE_MAX_TTL, c.MaxTTL)
	d.Set(BRIDGE_MESSAGE_VPN, c.MsgVpnName)
	d.Set(BRIDGE_NAME, c.BridgeName)
	d.Set(BRIDGE_REMOTE_AUTHENTICATION_BASIC_CLIENT_USERNAME, c.RemoteAuthenticationBasicClientUsername)
	d.Set(BRIDGE_REMOTE_AUTHENTICATION_SCHEME, c.RemoteAuthenticationScheme)
	d.Set(BRIDGE_REMOTE_CONNECTION_RETRY_COUNT, c.RemoteConnectionRetryCount)
	d.Set(BRIDGE_REMOTE_CONNECTION_RETRY_DELAY, c.RemoteConnectionRetryDelay)
	d.Set(BRIDGE_REMOTE_DELIVER_TO_ONE_PRIORITY, c.RemoteDeliverToOnePriority)
	d.Set(BRIDGE_VIRTUAL_ROUTER, c.BridgeVirtualRouter)
	d.Set(BRIDGE_TLS_CIPHER_SUITE_LIST, c.TLSCipherSuiteList)
	return diags
}

func deleteBridgeFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	MsgVpnName := d.Get(BRIDGE_MESSAGE_VPN).(string)
	BridgeName := d.Get(BRIDGE_NAME).(string)
	BridgeVirtualRouter := d.Get(BRIDGE_VIRTUAL_ROUTER).(string)
	params := all.NewDeleteMsgVpnBridgeParamsWithContext(ctx).WithMsgVpnName(MsgVpnName).WithBridgeName(BridgeName).WithBridgeVirtualRouter(BridgeVirtualRouter)
	_, err := state.Client.All.DeleteMsgVpnBridge(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	return diags
}
