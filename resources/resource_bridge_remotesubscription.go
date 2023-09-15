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

// Main resource definition for BridgeRemoteSubscription entities
func ResourceBridgeRemoteSubscription() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaBridgeRemoteSubscription(),
		// Provider CRUD functions
		CreateContext: createBridgeRemoteSubscriptionFunc,
		ReadContext:   readBridgeRemoteSubscriptionFunc,
		//UpdateContext: updateBridgeRemoteSubscriptionFunc,
		DeleteContext: deleteBridgeRemoteSubscriptionFunc,
	}
}

func schemaBridgeRemoteSubscription() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		BRIDGE_REMOTE_SUBSCRIPTION_BRIDGE_NAME: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the Bridge.",
			ForceNew:    true,
		},
		BRIDGE_REMOTE_SUBSCRIPTION_MESSAGE_VPN: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the Message VPN.",
			ForceNew:    true,
		},
		BRIDGE_REMOTE_SUBSCRIPTION_REMOTE_SUBSCRIPTION_TOPIC: {
			Type:     schema.TypeString,
			Required: true,
			Description: "The Client Username the Bridge uses to login to the remote Message VPN. This per remote Message VPN value overrides " +
				"the value provided for the Bridge overall. Changes to this attribute are synchronized to HA mates and replication sites via config-sync.",
			ForceNew: true,
		},
		BRIDGE_REMOTE_SUBSCRIPTION_BRIDGE_VIRTUAL_ROUTER: {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "auto",
			Description: "The virtual router of the Bridge. The allowed values and their meaning are:" +
				"'primary' - The Bridge is used for the primary virtual router." +
				"'backup' - The Bridge is used for the backup virtual router." +
				"'auto' - The Bridge is automatically assigned a virtual router at creation, depending on the broker's active-standby role.",
			ForceNew: true,
		},
		BRIDGE_REMOTE_SUBSCRIPTION_DELIVERY_ALWAYS_ENABLED: {
			Type:     schema.TypeBool,
			Required: true,
			Description: "Enable or disable deliver-always for the Bridge remote subscription topic instead of a deliver-to-one remote priority." +
				"A given topic for the Bridge may be deliver-to-one or deliver-always but not both.",
			ForceNew: true,
		},
	}
}

// Creates a BridgeRemoteSubscription model based on the terraform resource state.
func getBridgeRemoteSubscriptionModelFromResource(d *schema.ResourceData) *models.MsgVpnBridgeRemoteSubscription {
	q := &models.MsgVpnBridgeRemoteSubscription{
		MsgVpnName:              d.Get(BRIDGE_REMOTE_SUBSCRIPTION_MESSAGE_VPN).(string),
		BridgeName:              d.Get(BRIDGE_REMOTE_SUBSCRIPTION_BRIDGE_NAME).(string),
		RemoteSubscriptionTopic: d.Get(BRIDGE_REMOTE_SUBSCRIPTION_REMOTE_SUBSCRIPTION_TOPIC).(string),
		DeliverAlwaysEnabled:    d.Get(BRIDGE_REMOTE_SUBSCRIPTION_DELIVERY_ALWAYS_ENABLED).(bool),
	}
	if v, s := d.GetOk(BRIDGE_REMOTE_SUBSCRIPTION_BRIDGE_VIRTUAL_ROUTER); s {
		q.BridgeVirtualRouter = v.(string)
	}
	return q
}

func createBridgeRemoteSubscriptionFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	body := getBridgeRemoteSubscriptionModelFromResource(d)
	params := all.NewCreateMsgVpnBridgeRemoteSubscriptionParamsWithContext(ctx).WithMsgVpnName(body.MsgVpnName).WithBridgeName(body.BridgeName).WithBridgeVirtualRouter(body.BridgeVirtualRouter).WithBody(body)
	_, err := state.Client.All.CreateMsgVpnBridgeRemoteSubscription(params, state.Auth)
	if err != nil {
		if semp, err2 := provider.GetSempError(err); err2 == nil {
			if *semp.Status == "ALREADY_EXISTS" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("BridgeRemoteSubscription %s already exists in Bridge %s. Going to import state from Broker", body.RemoteSubscriptionTopic, body.BridgeName),
				})
			} else {
				return provider.AppendError(diags, err)
			}
		} else {
			return provider.AppendError(diags, err)
		}
	}
	d.SetId("sol_BridgeRemoteSubscription_" + body.BridgeName + "_" + body.MsgVpnName + "_" + body.RemoteSubscriptionTopic)
	return append(diags, readBridgeRemoteSubscriptionFunc(ctx, d, meta)...)
}

func readBridgeRemoteSubscriptionFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	BridgeRemoteSubscriptionMsgVpnName := d.Get(BRIDGE_REMOTE_SUBSCRIPTION_MESSAGE_VPN).(string)
	BridgeRemoteSubscriptionBridgeName := d.Get(BRIDGE_REMOTE_SUBSCRIPTION_BRIDGE_NAME).(string)
	BridgeRemoteSubscriptionVirtualRouter := d.Get(BRIDGE_REMOTE_SUBSCRIPTION_BRIDGE_VIRTUAL_ROUTER).(string)
	BridgeRemoteSubscriptionTopic := d.Get(BRIDGE_REMOTE_SUBSCRIPTION_REMOTE_SUBSCRIPTION_TOPIC).(string)
	params := all.NewGetMsgVpnBridgeRemoteSubscriptionParamsWithContext(ctx).WithMsgVpnName(BridgeRemoteSubscriptionMsgVpnName).WithBridgeName(BridgeRemoteSubscriptionBridgeName).WithBridgeVirtualRouter(BridgeRemoteSubscriptionVirtualRouter).WithRemoteSubscriptionTopic(BridgeRemoteSubscriptionTopic)
	resp, err := state.Client.All.GetMsgVpnBridgeRemoteSubscription(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	c := resp.Payload.Data
	d.Set(BRIDGE_REMOTE_SUBSCRIPTION_BRIDGE_NAME, c.BridgeName)
	d.Set(BRIDGE_REMOTE_SUBSCRIPTION_BRIDGE_VIRTUAL_ROUTER, c.BridgeVirtualRouter)
	d.Set(BRIDGE_REMOTE_SUBSCRIPTION_REMOTE_SUBSCRIPTION_TOPIC, c.RemoteSubscriptionTopic)
	d.Set(BRIDGE_REMOTE_SUBSCRIPTION_DELIVERY_ALWAYS_ENABLED, c.DeliverAlwaysEnabled)
	d.Set(BRIDGE_REMOTE_SUBSCRIPTION_MESSAGE_VPN, c.MsgVpnName)
	return diags
}

func deleteBridgeRemoteSubscriptionFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	BridgeRemoteSubscriptionMsgVpnName := d.Get(BRIDGE_REMOTE_SUBSCRIPTION_MESSAGE_VPN).(string)
	BridgeRemoteSubscriptionBridgeName := d.Get(BRIDGE_REMOTE_SUBSCRIPTION_BRIDGE_NAME).(string)
	BridgeRemoteSubscriptionVirtualRouter := d.Get(BRIDGE_REMOTE_SUBSCRIPTION_BRIDGE_VIRTUAL_ROUTER).(string)
	BridgeRemoteSubscriptionTopic := d.Get(BRIDGE_REMOTE_SUBSCRIPTION_REMOTE_SUBSCRIPTION_TOPIC).(string)
	params := all.NewDeleteMsgVpnBridgeRemoteSubscriptionParamsWithContext(ctx).WithMsgVpnName(BridgeRemoteSubscriptionMsgVpnName).WithBridgeName(BridgeRemoteSubscriptionBridgeName).WithBridgeVirtualRouter(BridgeRemoteSubscriptionVirtualRouter).WithRemoteSubscriptionTopic(BridgeRemoteSubscriptionTopic)
	_, err := state.Client.All.DeleteMsgVpnBridgeRemoteSubscription(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	return diags
}
