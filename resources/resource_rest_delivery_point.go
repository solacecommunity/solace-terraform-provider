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
func ResourceRestDeliveryPoint() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaRestDeliveryPoint(),
		// Provider CRUD functions
		CreateContext: createRestDeliveryPoint,
		ReadContext:   readRestDeliveryPoint,
		UpdateContext: updateRestDeliveryPoint,
		DeleteContext: deleteRestDeliveryPoint,
	}
}

func schemaRestDeliveryPoint() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		CLIENT_PROFILE_NAME: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the Client Profile",
		},
		ENABLED: {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Enable or disable the RDP",
		},
		MSG_VPN_NAME: {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The name of the MSG VPN where the RDP gets created.",
		},
		REST_DELIVERY_POINT_NAME: {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The name of the Rest Delivery Point.",
		},
	}
}

func getRestDeliveryPointModelFromResource(d *schema.ResourceData) *models.MsgVpnRestDeliveryPoint {
	q := &models.MsgVpnRestDeliveryPoint{
		ClientProfileName:     d.Get(CLIENT_PROFILE_NAME).(string),
		Enabled:               d.Get(ENABLED).(bool),
		MsgVpnName:            d.Get(MSG_VPN_NAME).(string),
		RestDeliveryPointName: d.Get(REST_DELIVERY_POINT_NAME).(string),
	}
	return q
}

func createRestDeliveryPoint(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	body := getRestDeliveryPointModelFromResource(d)
	msgvpn := d.Get(MSG_VPN_NAME).(string)

	params := all.NewCreateMsgVpnRestDeliveryPointParamsWithContext(ctx).WithMsgVpnName(msgvpn).WithBody(body)
	_, err := state.Client.All.CreateMsgVpnRestDeliveryPoint(params, state.Auth)
	if err != nil {
		if semp, err2 := provider.GetSempError(err); err2 == nil {
			if *semp.Status == "ALREADY_EXISTS" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Rest Delivery Point %s already exists. Going to import state from Broker", body.RestDeliveryPointName),
				})
			} else {
				return provider.AppendError(diags, err)
			}
		} else {
			return provider.AppendError(diags, err)
		}
	}
	d.SetId("sol_rdp_" + body.RestDeliveryPointName)

	// read after create to make sure tf state is in sync with broker
	return append(diags, readRestDeliveryPoint(ctx, d, meta)...)
}

func readRestDeliveryPoint(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	rdp_name := d.Get(REST_DELIVERY_POINT_NAME).(string)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	params := all.NewGetMsgVpnRestDeliveryPointParamsWithContext(ctx).WithMsgVpnName(msgvpn).WithRestDeliveryPointName(rdp_name)
	resp, err := state.Client.All.GetMsgVpnRestDeliveryPoint(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	p := resp.Payload.Data
	d.Set(CLIENT_PROFILE_NAME, p.ClientProfileName)
	d.Set(ENABLED, p.Enabled)
	d.Set(MSG_VPN_NAME, p.MsgVpnName)
	d.Set(REST_DELIVERY_POINT_NAME, p.RestDeliveryPointName)
	return diags
}

func updateRestDeliveryPoint(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	rdp_name := d.Get(REST_DELIVERY_POINT_NAME).(string)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	body := getRestDeliveryPointModelFromResource(d)
	params := all.NewUpdateMsgVpnRestDeliveryPointParamsWithContext(ctx).WithMsgVpnName(msgvpn).WithRestDeliveryPointName(rdp_name).WithBody(body)
	_, err := state.Client.All.UpdateMsgVpnRestDeliveryPoint(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}

	// read after update to make sure tf state is in sync with broker
	return readRestDeliveryPoint(ctx, d, meta)
}

func deleteRestDeliveryPoint(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	rdp_name := d.Get(REST_DELIVERY_POINT_NAME).(string)
	msgvpn := d.Get(MSG_VPN_NAME).(string)

	params := all.NewDeleteMsgVpnRestDeliveryPointParamsWithContext(ctx).WithMsgVpnName(msgvpn).WithRestDeliveryPointName(rdp_name)
	_, err := state.Client.All.DeleteMsgVpnRestDeliveryPoint(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	return diags
}
