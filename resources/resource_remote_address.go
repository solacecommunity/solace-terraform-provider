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
func ResourceRemoteAddress() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaRemoteAddress(),
		// Provider CRUD functions
		CreateContext: createRemoteAddressFunc,
		ReadContext:   readRemoteAddressFunc,
		UpdateContext: updateRemoteAddressFunc,
		DeleteContext: deleteRemoteAddressFunc,
	}
}

func schemaRemoteAddress() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		REMOTE_ADDRESS_DMR_CLUSTER_NAME: {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    false,
			Description: "The name of the Cluster.",
		},
		REMOTE_ADDRESS: {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			ForceNew:    true,
			Description: "The password used to authenticate with the remote node when using basic internal authentication. If this per-Link password is not configured, the Cluster's password is used instead",
		},
		REMOTE_ADDRESS_REMOTE_NODE_NAME: {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    false,
			Description: "The name of the node at the remote end of the Link.",
		},
	}
}

// Creates a queue model based on the terraform resource state.
func getRemoteAddressModelFromResource(d *schema.ResourceData) *models.DmrClusterLinkRemoteAddress {
	q := &models.DmrClusterLinkRemoteAddress{
		DmrClusterName: d.Get(REMOTE_ADDRESS_DMR_CLUSTER_NAME).(string),
		RemoteNodeName: d.Get(REMOTE_ADDRESS_REMOTE_NODE_NAME).(string),
		RemoteAddress:  d.Get(REMOTE_ADDRESS).(string),
	}

	return q
}

func createRemoteAddressFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	body := getRemoteAddressModelFromResource(d)
	err := disableAndEnableLinkWhileCreatingRemoteAddress(ctx, d, meta)
	if err != nil {
		if semp, err2 := provider.GetSempError(err); err2 == nil {
			if *semp.Status == "ALREADY_EXISTS" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Remote address %s to remote node %s already exists. Going to import state from Broker", body.RemoteAddress, body.RemoteNodeName),
				})
			} else {
				return provider.AppendError(diags, err)
			}
		} else {
			return provider.AppendError(diags, err)
		}
	}

	d.SetId("sol_RemoteAddress_" + body.RemoteNodeName)
	return append(diags, readRemoteAddressFunc(ctx, d, meta)...)
}

func updateRemoteAddressFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	RemoteNodeName := d.Get(REMOTE_ADDRESS_REMOTE_NODE_NAME).(string)
	ClusterName := d.Get(REMOTE_ADDRESS_DMR_CLUSTER_NAME).(string)
	RemoteAddress := d.Get(REMOTE_ADDRESS).(string)
	params := all.NewDeleteDmrClusterLinkRemoteAddressParamsWithContext(ctx).WithDmrClusterName(ClusterName).WithRemoteNodeName(RemoteNodeName).WithRemoteAddress(RemoteAddress)
	_, err := state.Client.All.DeleteDmrClusterLinkRemoteAddress(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	err_2 := disableAndEnableLinkWhileCreatingRemoteAddress(ctx, d, meta)
	if err_2 != nil {
		return provider.AppendError(diags, err)
	}
	return readRemoteAddressFunc(ctx, d, meta)
}

func readRemoteAddressFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	RemoteNodeName := d.Get(REMOTE_ADDRESS_REMOTE_NODE_NAME).(string)
	ClusterName := d.Get(REMOTE_ADDRESS_DMR_CLUSTER_NAME).(string)
	RemoteAddress := d.Get(REMOTE_ADDRESS).(string)
	params := all.NewGetDmrClusterLinkRemoteAddressParamsWithContext(ctx).WithDmrClusterName(ClusterName).WithRemoteNodeName(RemoteNodeName).WithRemoteAddress(RemoteAddress)
	resp, err := state.Client.All.GetDmrClusterLinkRemoteAddress(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	c := resp.Payload.Data
	d.Set(REMOTE_ADDRESS_DMR_CLUSTER_NAME, c.DmrClusterName)
	d.Set(REMOTE_ADDRESS_REMOTE_NODE_NAME, c.RemoteNodeName)
	d.Set(REMOTE_ADDRESS, c.RemoteAddress)
	return diags
}

func deleteRemoteAddressFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	RemoteNodeName := d.Get(REMOTE_ADDRESS_REMOTE_NODE_NAME).(string)
	ClusterName := d.Get(REMOTE_ADDRESS_DMR_CLUSTER_NAME).(string)
	RemoteAddress := d.Get(REMOTE_ADDRESS).(string)
	params := all.NewDeleteDmrClusterLinkRemoteAddressParamsWithContext(ctx).WithDmrClusterName(ClusterName).WithRemoteNodeName(RemoteNodeName).WithRemoteAddress(RemoteAddress)
	_, err := state.Client.All.DeleteDmrClusterLinkRemoteAddress(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	return diags
}

func disableAndEnableLinkWhileCreatingRemoteAddress(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	// to create a remote address the underlying link has to be disabled.
	state := meta.(*provider.ProviderState)
	RemoteNodeName := d.Get(REMOTE_ADDRESS_REMOTE_NODE_NAME).(string)
	ClusterName := d.Get(REMOTE_ADDRESS_DMR_CLUSTER_NAME).(string)
	q_disable := &models.DmrClusterLink{
		Enabled: false,
	}
	q_enable := &models.DmrClusterLink{
		Enabled:             true,
		TransportTLSEnabled: true,
	}
	// disable the link
	params_disable_link := all.NewUpdateDmrClusterLinkParamsWithContext(ctx).WithDmrClusterName(ClusterName).WithRemoteNodeName(RemoteNodeName).WithBody(q_disable)
	_, err := state.Client.All.UpdateDmrClusterLink(params_disable_link, state.Auth)
	if err != nil {
		return err
	}
	// add the remote address
	body := getRemoteAddressModelFromResource(d)
	params := all.NewCreateDmrClusterLinkRemoteAddressParamsWithContext(ctx).WithDmrClusterName(ClusterName).WithRemoteNodeName(RemoteNodeName).WithBody(body)
	_, err_2 := state.Client.All.CreateDmrClusterLinkRemoteAddress(params, state.Auth)
	if err_2 != nil {
		return err_2
	}
	// enable the link again
	params_enable_link := all.NewUpdateDmrClusterLinkParamsWithContext(ctx).WithDmrClusterName(ClusterName).WithRemoteNodeName(RemoteNodeName).WithBody(q_enable)
	_, err_3 := state.Client.All.UpdateDmrClusterLink(params_enable_link, state.Auth)
	if err_3 != nil {
		return err_3
	}

	return nil
}
