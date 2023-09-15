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

// Main resource definition for DMR Cluster entities
func ResourceExternalLink() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaExternalLink(),
		// Provider CRUD functions
		CreateContext: createExternalLinkFunc,
		ReadContext:   readExternalLinkFunc,
		UpdateContext: updateExternalLinkFunc,
		DeleteContext: deleteExternalLinkFunc,
	}
}

func schemaExternalLink() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		EXTERNAL_LINK_DMR_CLUSTER_NAME: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the Cluster.",
		},
		EXTERNAL_LINK_AUTHENTICATION_BASIC_PASSWORD: {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "The password used to authenticate with the remote node when using basic internal authentication. If this per-Link password is not configured, the Cluster's password is used instead",
		},
		EXTERNAL_LINK_REMOTE_NODE_NAME: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the node at the remote end of the Link.",
		},
		EXTERNAL_LINK_INITIATOR: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"local", "remote"}, false),
			Description:  "The initiator of the Link's TCP connections. Changes to this attribute are synchronized to HA mates via config-sync.",
		},
		EXTERNAL_LINK_ENABLED: {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Enable or disable the Link. When disabled, subscription sets of this and the remote node are not kept up-to-date, and messages are not exchanged with the remote node. Published guaranteed messages will be queued up for future delivery based on current subscription sets.",
		},
		EXTERNAL_LINK_SPAN: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"internal", "external"}, false),
			Description:  "The span of the Link, either internal or external. Internal Links connect nodes within the same Cluster. External Links connect nodes within different Clusters",
		},
		EXTERNAL_LINK_TRANSPORT_TLS_ENABLED: {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Enable or disable encryption (TLS) on the Link. Changes to this attribute are synchronized to HA mates via config-sync",
		},
	}
}

// Creates a queue model based on the terraform resource state.
func getExternalLinkModelFromResource(d *schema.ResourceData) *models.DmrClusterLink {
	q := &models.DmrClusterLink{
		AuthenticationBasicPassword: d.Get(EXTERNAL_LINK_AUTHENTICATION_BASIC_PASSWORD).(string),
		DmrClusterName:              d.Get(EXTERNAL_LINK_DMR_CLUSTER_NAME).(string),
		RemoteNodeName:              d.Get(EXTERNAL_LINK_REMOTE_NODE_NAME).(string),
		Initiator:                   d.Get(EXTERNAL_LINK_INITIATOR).(string),
		Span:                        d.Get(EXTERNAL_LINK_SPAN).(string),
	}

	if v, s := d.GetOk(EXTERNAL_LINK_ENABLED); s {
		q.Enabled = v.(bool)
	}
	if v, s := d.GetOk(EXTERNAL_LINK_TRANSPORT_TLS_ENABLED); s {
		q.TransportTLSEnabled = v.(bool)
	}
	return q
}

func createExternalLinkFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	ClusterName := d.Get(EXTERNAL_LINK_DMR_CLUSTER_NAME).(string)
	body := getExternalLinkModelFromResource(d)
	params := all.NewCreateDmrClusterLinkParamsWithContext(ctx).WithDmrClusterName(ClusterName).WithBody(body)
	_, err := state.Client.All.CreateDmrClusterLink(params, state.Auth)
	if err != nil {
		if semp, err2 := provider.GetSempError(err); err2 == nil {
			if *semp.Status == "ALREADY_EXISTS" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("External Link to remote node %s already exists. Going to import state from Broker", body.RemoteNodeName),
				})
			} else {
				return provider.AppendError(diags, err)
			}
		} else {
			return provider.AppendError(diags, err)
		}
	}

	d.SetId("sol_ExternalLink_" + body.RemoteNodeName)
	return append(diags, readExternalLinkFunc(ctx, d, meta)...)
}

func updateExternalLinkFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	RemoteNodeName := d.Get(EXTERNAL_LINK_REMOTE_NODE_NAME).(string)
	ClusterName := d.Get(EXTERNAL_LINK_DMR_CLUSTER_NAME).(string)
	body := getExternalLinkModelFromResource(d)
	params := all.NewUpdateDmrClusterLinkParamsWithContext(ctx).WithRemoteNodeName(RemoteNodeName).WithDmrClusterName(ClusterName).WithBody(body)

	_, err := state.Client.All.UpdateDmrClusterLink(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}

	return readExternalLinkFunc(ctx, d, meta)
}

func readExternalLinkFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	RemoteNodeName := d.Get(EXTERNAL_LINK_REMOTE_NODE_NAME).(string)
	ClusterName := d.Get(EXTERNAL_LINK_DMR_CLUSTER_NAME).(string)
	params := all.NewGetDmrClusterLinkParamsWithContext(ctx).WithRemoteNodeName(RemoteNodeName).WithDmrClusterName(ClusterName)
	resp, err := state.Client.All.GetDmrClusterLink(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	c := resp.Payload.Data
	d.Set(EXTERNAL_LINK_DMR_CLUSTER_NAME, c.DmrClusterName)
	d.Set(EXTERNAL_LINK_ENABLED, c.Enabled)
	d.Set(EXTERNAL_LINK_INITIATOR, c.Initiator)
	d.Set(EXTERNAL_LINK_REMOTE_NODE_NAME, c.RemoteNodeName)
	d.Set(EXTERNAL_LINK_SPAN, c.Span)
	d.Set(EXTERNAL_LINK_TRANSPORT_TLS_ENABLED, c.TransportTLSEnabled)
	return diags
}

func deleteExternalLinkFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	RemoteNodeName := d.Get(EXTERNAL_LINK_REMOTE_NODE_NAME).(string)
	ClusterName := d.Get(EXTERNAL_LINK_DMR_CLUSTER_NAME).(string)
	params := all.NewDeleteDmrClusterLinkParamsWithContext(ctx).WithRemoteNodeName(RemoteNodeName).WithDmrClusterName(ClusterName)
	_, err := state.Client.All.DeleteDmrClusterLink(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	return diags
}
