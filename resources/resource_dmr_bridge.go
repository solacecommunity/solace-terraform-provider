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

// Main resource definition for DMR Cluster entities
func ResourceDmrBridge() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaDmrBridge(),
		// Provider CRUD functions
		CreateContext: createDmrBridgeFunc,
		ReadContext:   readDmrBridgeFunc,
		UpdateContext: updateDmrBridgeFunc,
		DeleteContext: deleteDmrBridgeFunc,
	}
}

func schemaDmrBridge() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		DMR_BRIDGE_MESSAGE_VPN: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the Message VPN.",
		},
		DMR_BRIDGE_REMOTE_MESSAGE_VPN: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The remote Message VPN of the DMR Bridge. Changes to this attribute are synchronized to HA mates and replication sites via config-sync.",
		},
		DMR_BRIDGE_REMOTE_NODE_NAME: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the node at the remote end of the DMR Bridge.",
		},
	}
}

// Creates a queue model based on the terraform resource state.
func getDmrBridgeModelFromResource(d *schema.ResourceData) *models.MsgVpnDmrBridge {
	q := &models.MsgVpnDmrBridge{
		MsgVpnName:       d.Get(DMR_BRIDGE_MESSAGE_VPN).(string),
		RemoteMsgVpnName: d.Get(DMR_BRIDGE_REMOTE_MESSAGE_VPN).(string),
		RemoteNodeName:   d.Get(DMR_BRIDGE_REMOTE_NODE_NAME).(string),
	}
	return q
}

func createDmrBridgeFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	body := getDmrBridgeModelFromResource(d)
	params := all.NewCreateMsgVpnDmrBridgeParamsWithContext(ctx).WithMsgVpnName(body.MsgVpnName).WithBody(body)
	_, err := state.Client.All.CreateMsgVpnDmrBridge(params, state.Auth)
	if err != nil {
		if semp, err2 := provider.GetSempError(err); err2 == nil {
			if *semp.Status == "ALREADY_EXISTS" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Dmr Bridge between local vpn %s and remote vpn %s already exists. Going to import state from Broker", body.MsgVpnName, body.RemoteMsgVpnName),
				})
			} else {
				return provider.AppendError(diags, err)
			}
		} else {
			return provider.AppendError(diags, err)
		}
	}

	d.SetId("sol_DmrBridge_" + body.RemoteNodeName + "_" + body.MsgVpnName + "_" + body.RemoteMsgVpnName)
	return append(diags, readDmrBridgeFunc(ctx, d, meta)...)
}

func updateDmrBridgeFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	RemoteNodeName := d.Get(DMR_BRIDGE_REMOTE_NODE_NAME).(string)
	body := getDmrBridgeModelFromResource(d)
	params := all.NewUpdateMsgVpnDmrBridgeParamsWithContext(ctx).WithMsgVpnName(body.MsgVpnName).WithRemoteNodeName(RemoteNodeName).WithBody(body)
	_, err := state.Client.All.UpdateMsgVpnDmrBridge(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}

	return readDmrBridgeFunc(ctx, d, meta)
}

func readDmrBridgeFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	RemoteNodeName := d.Get(DMR_BRIDGE_REMOTE_NODE_NAME).(string)
	MsgVpnName := d.Get(DMR_BRIDGE_MESSAGE_VPN).(string)
	params := all.NewGetMsgVpnDmrBridgeParamsWithContext(ctx).WithMsgVpnName(MsgVpnName).WithRemoteNodeName(RemoteNodeName)
	resp, err := state.Client.All.GetMsgVpnDmrBridge(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	c := resp.Payload.Data
	d.Set(DMR_BRIDGE_REMOTE_NODE_NAME, c.RemoteNodeName)
	d.Set(DMR_BRIDGE_REMOTE_MESSAGE_VPN, c.RemoteMsgVpnName)
	d.Set(DMR_BRIDGE_MESSAGE_VPN, c.MsgVpnName)
	return diags
}

func deleteDmrBridgeFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	RemoteNodeName := d.Get(DMR_BRIDGE_REMOTE_NODE_NAME).(string)
	MsgVpnName := d.Get(DMR_BRIDGE_MESSAGE_VPN).(string)
	params := all.NewDeleteMsgVpnDmrBridgeParamsWithContext(ctx).WithMsgVpnName(MsgVpnName).WithRemoteNodeName(RemoteNodeName)
	_, err := state.Client.All.DeleteMsgVpnDmrBridge(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	return diags
}
