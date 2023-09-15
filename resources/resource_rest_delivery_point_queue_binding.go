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

// Main resource definition for Client Usernames
func ResourceRestDeliveryPointQueueBinding() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaRestDeliveryPointQueueBinding(),
		// Provider CRUD functions
		CreateContext: createRestDeliveryPointQueueBinding,
		ReadContext:   readRestDeliveryPointQueueBinding,
		UpdateContext: updateRestDeliveryPointQueueBinding,
		DeleteContext: deleteRestDeliveryPointQueueBinding,
	}
}

func schemaRestDeliveryPointQueueBinding() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		REST_DELIVERY_POINT_QUEUE_BINDING_NAME: {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The name of a queue (to bind) in the Message VPN.",
		},
		REST_DELIVERY_POINT_POST_REQUEST_TARGET: {
			Type:     schema.TypeString,
			Required: true,
			Description: "The request-target string to use when sending requests." +
				"It identifies the target resource on the far-end REST Consumer upon which to apply the request." +
				"There are generally two common forms for the request-target." +
				"The origin-form is most often used in practice and contains the path and query components of the target URI." +
				"If the path component is empty then the client must generally send a / as the path.",
		},
		MSG_VPN_NAME: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the MSG VPN where the RDP gets created.",
		},
		REST_DELIVERY_POINT_NAME: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the Rest Delivery Point.",
		},
	}
}

func getRestDeliveryPointQueueBindingModelFromResource(d *schema.ResourceData) *models.MsgVpnRestDeliveryPointQueueBinding {
	q := &models.MsgVpnRestDeliveryPointQueueBinding{
		PostRequestTarget:     d.Get(REST_DELIVERY_POINT_POST_REQUEST_TARGET).(string),
		QueueBindingName:      d.Get(REST_DELIVERY_POINT_QUEUE_BINDING_NAME).(string),
		MsgVpnName:            d.Get(MSG_VPN_NAME).(string),
		RestDeliveryPointName: d.Get(REST_DELIVERY_POINT_NAME).(string),
	}
	return q
}

func createRestDeliveryPointQueueBinding(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	body := getRestDeliveryPointQueueBindingModelFromResource(d)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	rdp_name := d.Get(REST_DELIVERY_POINT_NAME).(string)

	params := all.NewCreateMsgVpnRestDeliveryPointQueueBindingParamsWithContext(ctx).WithMsgVpnName(msgvpn).WithRestDeliveryPointName(rdp_name).WithBody(body)
	_, err := state.Client.All.CreateMsgVpnRestDeliveryPointQueueBinding(params, state.Auth)
	if err != nil {
		if semp, err2 := provider.GetSempError(err); err2 == nil {
			if *semp.Status == "ALREADY_EXISTS" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Queue Binding %s  for RDP %s already exists. Going to import state from Broker", body.QueueBindingName, body.RestDeliveryPointName),
				})
			} else {
				return provider.AppendError(diags, err)
			}
		} else {
			return provider.AppendError(diags, err)
		}
	}
	d.SetId("sol_rdp_" + body.RestDeliveryPointName + "_qb_" + body.QueueBindingName)

	// read after create to make sure tf state is in sync with broker
	return append(diags, readRestDeliveryPointQueueBinding(ctx, d, meta)...)
}

func readRestDeliveryPointQueueBinding(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	rdp_name := d.Get(REST_DELIVERY_POINT_NAME).(string)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	qb_name := d.Get(REST_DELIVERY_POINT_QUEUE_BINDING_NAME).(string)
	params := all.NewGetMsgVpnRestDeliveryPointQueueBindingParamsWithContext(ctx).WithMsgVpnName(msgvpn).WithRestDeliveryPointName(rdp_name).WithQueueBindingName(qb_name)
	resp, err := state.Client.All.GetMsgVpnRestDeliveryPointQueueBinding(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	p := resp.Payload.Data
	d.Set(REST_DELIVERY_POINT_QUEUE_BINDING_NAME, p.QueueBindingName)
	d.Set(REST_DELIVERY_POINT_POST_REQUEST_TARGET, p.PostRequestTarget)
	d.Set(MSG_VPN_NAME, p.MsgVpnName)
	d.Set(REST_DELIVERY_POINT_NAME, p.RestDeliveryPointName)
	return diags
}

func updateRestDeliveryPointQueueBinding(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	rdp_name := d.Get(REST_DELIVERY_POINT_NAME).(string)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	qb_name := d.Get(REST_DELIVERY_POINT_QUEUE_BINDING_NAME).(string)
	body := getRestDeliveryPointQueueBindingModelFromResource(d)
	params := all.NewUpdateMsgVpnRestDeliveryPointQueueBindingParamsWithContext(ctx).WithMsgVpnName(msgvpn).WithRestDeliveryPointName(rdp_name).WithQueueBindingName(qb_name).WithBody(body)
	_, err := state.Client.All.UpdateMsgVpnRestDeliveryPointQueueBinding(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}

	// read after update to make sure tf state is in sync with broker
	return readRestDeliveryPointQueueBinding(ctx, d, meta)
}

func deleteRestDeliveryPointQueueBinding(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	rdp_name := d.Get(REST_DELIVERY_POINT_NAME).(string)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	qb_name := d.Get(REST_DELIVERY_POINT_QUEUE_BINDING_NAME).(string)

	params := all.NewDeleteMsgVpnRestDeliveryPointQueueBindingParamsWithContext(ctx).WithMsgVpnName(msgvpn).WithRestDeliveryPointName(rdp_name).WithQueueBindingName(qb_name)
	_, err := state.Client.All.DeleteMsgVpnRestDeliveryPointQueueBinding(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	return diags
}
