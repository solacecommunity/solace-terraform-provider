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
func ResourceJndiQueue() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaJndiQueue(),
		// Provider CRUD functions
		CreateContext: createJndiQueue,
		ReadContext:   readJndiQueue,
		DeleteContext: deleteJndiQueue,
	}
}

func schemaJndiQueue() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		QUEUE_NAME: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		MSG_VPN_NAME: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		JNDI_NAME: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
	}
}

func createJndiQueue(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	msgvpn := d.Get(MSG_VPN_NAME).(string)

	// don't get confused here with queueName vs. jndiName
	// Solace API uses QueueName as the name of the JNDI Alias and PhysicalName for the real queue
	// For better understanding we keep QueueName as property for the physical queue and introduce JndiName
	body := &models.MsgVpnJndiQueue{
		QueueName:    d.Get(JNDI_NAME).(string),
		MsgVpnName:   msgvpn,
		PhysicalName: d.Get(QUEUE_NAME).(string),
	}

	params := all.NewCreateMsgVpnJndiQueueParams().WithContext(ctx).WithMsgVpnName(msgvpn).WithBody(body)
	_, err := state.Client.All.CreateMsgVpnJndiQueue(params, state.Auth)
	if err != nil {
		if semp, err2 := provider.GetSempError(err); err2 == nil {
			if *semp.Status == "ALREADY_EXISTS" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("JNDI Queue %s already exists. Going to import state from Broker", body.QueueName),
				})
			} else {
				return provider.AppendError(diags, err)
			}
		} else {
			return provider.AppendError(diags, err)
		}
	}
	d.SetId("sol_" + msgvpn + "_jq_" + body.QueueName)
	return append(diags, readJndiQueue(ctx, d, meta)...)
}

func readJndiQueue(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	jndi := d.Get(JNDI_NAME).(string)
	msgvpn := d.Get(MSG_VPN_NAME).(string)

	params := all.NewGetMsgVpnJndiQueueParams().WithContext(ctx).WithQueueName(jndi).WithMsgVpnName(msgvpn)
	resp, err := state.Client.All.GetMsgVpnJndiQueue(params, state.Auth)
	if err != nil {
		if semp, err2 := provider.GetSempError(err); err2 == nil {
			if *semp.Status == provider.SEMP_STATUS_NOT_FOUND {
				d.SetId("")
				return append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("JNDI Queue %s does not exist. Going to remove ressource from tfstate", jndi),
				})
			} else {
				return provider.AppendError(diags, err)
			}
		} else {
			return provider.AppendError(diags, err)
		}
	}
	d.Set(QUEUE_NAME, resp.Payload.Data.PhysicalName)
	return diags
}

func deleteJndiQueue(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	jndi := d.Get(JNDI_NAME).(string)
	msgvpn := d.Get(MSG_VPN_NAME).(string)

	params := all.NewDeleteMsgVpnJndiQueueParams().WithContext(ctx).WithMsgVpnName(msgvpn).WithQueueName(jndi)
	_, err := state.Client.All.DeleteMsgVpnJndiQueue(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	return diags
}
