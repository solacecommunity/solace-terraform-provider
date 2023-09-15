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
func ResourceDmrCluster() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaDmrCluster(),
		// Provider CRUD functions
		CreateContext: createDmrClusterFunc,
		ReadContext:   readDmrClusterFunc,
		UpdateContext: updateDmrClusterFunc,
		DeleteContext: deleteDmrClusterFunc,
	}
}

func schemaDmrCluster() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		DMR_CLUSTER_AUTHENTICATION_BASIC_PASSWORD: {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		DMR_CLUSTER_NAME: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		DMR_CLUSTER_ENABLED: {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Enable or disable the dmr cluster.",
		},
	}
}

// Creates a queue model based on the terraform resource state.
func getDmrClusterModelFromResource(d *schema.ResourceData) *models.DmrCluster {
	q := &models.DmrCluster{}
	q.DmrClusterName = d.Get(DMR_CLUSTER_NAME).(string)
	q.AuthenticationBasicPassword = d.Get(DMR_CLUSTER_AUTHENTICATION_BASIC_PASSWORD).(string)
	q.Enabled = d.Get(DMR_CLUSTER_ENABLED).(bool)
	return q
}

func createDmrClusterFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)

	body := getDmrClusterModelFromResource(d)
	params := all.NewCreateDmrClusterParamsWithContext(ctx).WithBody(body)
	_, err := state.Client.All.CreateDmrCluster(params, state.Auth)
	if err != nil {
		if semp, err2 := provider.GetSempError(err); err2 == nil {
			if *semp.Status == "ALREADY_EXISTS" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Cluster %s already exists. Going to import state from Broker", body.DmrClusterName),
				})
			} else {
				return provider.AppendError(diags, err)
			}
		} else {
			return provider.AppendError(diags, err)
		}
	}

	d.SetId("sol_dmrcluster_" + body.DmrClusterName)
	return append(diags, readDmrClusterFunc(ctx, d, meta)...)
}

func updateDmrClusterFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	dmrClusterName := d.Get(DMR_CLUSTER_NAME).(string)
	body := getDmrClusterModelFromResource(d)
	params := all.NewUpdateDmrClusterParamsWithContext(ctx).WithDmrClusterName(dmrClusterName).WithBody(body)

	_, err := state.Client.All.UpdateDmrCluster(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}

	return readDmrClusterFunc(ctx, d, meta)
}

func readDmrClusterFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	dmrCluster := d.Get(DMR_CLUSTER_NAME).(string)

	params := all.NewGetDmrClusterParamsWithContext(ctx).WithDmrClusterName(dmrCluster)
	resp, err := state.Client.All.GetDmrCluster(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	c := resp.Payload.Data
	d.Set(DMR_CLUSTER_NAME, c.DmrClusterName)
	d.Set(DMR_CLUSTER_ENABLED, c.Enabled)
	return diags
}

func deleteDmrClusterFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	dmrCluster := d.Get(DMR_CLUSTER_NAME).(string)

	params := all.NewDeleteDmrClusterParamsWithContext(ctx).WithDmrClusterName(dmrCluster)
	_, err := state.Client.All.DeleteDmrCluster(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	return diags
}
