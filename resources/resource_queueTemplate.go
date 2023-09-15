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

// Main resource definition for Queue entities
func ResourceQueueTemplate() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaQueueTemplate(),
		// Provider CRUD functions
		CreateContext: createQueueTemplateFunc,
		ReadContext:   readQueueTemplateFunc,
		//UpdateContext: updateQueueTemplateFunc,
		DeleteContext: deleteQueueTemplateFunc,
	}
}

func schemaQueueTemplate() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		MSG_VPN_NAME: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		QUEUE_TEMPLATE_NAME: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
	}
}

// Updates the terraform resource state according to the solace model returned from the SEMP API
func updateQueueTemplateResourceFromModel(d *schema.ResourceData, q *models.MsgVpnQueueTemplate) {
}

func getQueueTemplateModelFromResource(d *schema.ResourceData) *models.MsgVpnQueueTemplate {
	q := &models.MsgVpnQueueTemplate{
		QueueTemplateName: d.Get(QUEUE_TEMPLATE_NAME).(string),
		MsgVpnName:        d.Get(MSG_VPN_NAME).(string),
	}
	return q
}

func createQueueTemplateFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	state := meta.(*provider.ProviderState)
	body := getQueueTemplateModelFromResource(d)

	params := all.NewCreateMsgVpnQueueTemplateParamsWithContext(ctx).WithMsgVpnName(msgvpn).WithBody(body)
	_, err := state.Client.All.CreateMsgVpnQueueTemplate(params, state.Auth)
	if err != nil {
		if semp, err2 := provider.GetSempError(err); err2 == nil {
			if *semp.Status == provider.SEMP_STATUS_ALREADY_EXISTS {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("QueueTemplate %s already exists. Going to import state from Broker", body.QueueTemplateName),
				})
			} else {
				return provider.AppendError(diags, err)
			}
		} else {
			return provider.AppendError(diags, err)
		}
	}
	d.SetId("sol_" + msgvpn + "_qtpl_" + body.QueueTemplateName)
	return append(diags, readQueueTemplateFunc(ctx, d, meta)...)
}

func readQueueTemplateFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	queue := d.Get(QUEUE_TEMPLATE_NAME).(string)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	state := meta.(*provider.ProviderState)
	params := all.NewGetMsgVpnQueueTemplateParamsWithContext(ctx).WithQueueTemplateName(queue).WithMsgVpnName(msgvpn)
	resp, err := state.Client.All.GetMsgVpnQueueTemplate(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	updateQueueTemplateResourceFromModel(d, resp.Payload.Data)
	return diags
}

// func updateQueueTemplateFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	var diags diag.Diagnostics
// 	queue := d.Get(QUEUE_NAME).(string)
// 	msgvpn := d.Get(MSG_VPN_NAME).(string)
// 	state := meta.(*provider.ProviderState)
// 	body := getQueueTemplateModelFromResource(d)

// 	params := all.NewUpdateMsgVpnQueueTemplateParamsWithContext(ctx).WithQueueTemplateName(queue).WithMsgVpnName(msgvpn).WithBody(body)
// 	resp, err := state.Client.All.UpdateMsgVpnQueueTemplate(params, state.Auth)
// 	if err != nil {
// 		return provider.AppendError(diags, err)
// 	}
// 	updateQueueTemplateResourceFromModel(d, resp.Payload.Data)
// 	return diags
// }

func deleteQueueTemplateFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	queue := d.Get(QUEUE_TEMPLATE_NAME).(string)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	state := meta.(*provider.ProviderState)

	params := all.NewDeleteMsgVpnQueueTemplateParamsWithContext(ctx).WithMsgVpnName(msgvpn).WithQueueTemplateName(queue)
	_, err := state.Client.All.DeleteMsgVpnQueueTemplate(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}

	return diags
}
