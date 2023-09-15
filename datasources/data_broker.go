package datasources

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/solacecommunity/solace-terraform-provider/provider"
)

func ResourceBroker() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of the message broker.",
			},
			"build": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of the message broker.",
			},
		},
		ReadContext: readBroker,
		Description: "Datasource for Solace PubSub+ Broker",
	}
}

func readBroker(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	client := state.SempV1Client

	version, err := client.GetBrokerVersion(ctx)
	if err != nil {
		return provider.AppendError(diags, err)
	}

	d.Set("version", version.RPC.Show.Version.Description)
	d.Set("build", strings.TrimPrefix(version.RPC.Show.Version.CurrentLoad, "soltr_"))
	d.SetId("sol_broker")

	return diags
}
