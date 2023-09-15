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

// Main resource definition for Jndi Topic entities
func ResourceJndiTopic() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaJndiTopic(),
		// Provider CRUD functions
		CreateContext: createJndiTopic,
		ReadContext:   readJndiTopic,
		DeleteContext: deleteJndiTopic,
	}
}

func schemaJndiTopic() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		TOPIC_NAME: {
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

func createJndiTopic(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	msgvpn := d.Get(MSG_VPN_NAME).(string)

	// don't get confused here with topicname vs. jndiName
	// Solace API uses Topicname as the name of the JNDI Alias and PhysicalName for the real topic
	// For better understanding we keep Topicname as property for the physical Topic and introduce JndiName
	body := &models.MsgVpnJndiTopic{
		TopicName:    d.Get(JNDI_NAME).(string),
		MsgVpnName:   msgvpn,
		PhysicalName: d.Get(TOPIC_NAME).(string),
	}

	params := all.NewCreateMsgVpnJndiTopicParams().WithContext(ctx).WithMsgVpnName(msgvpn).WithBody(body)
	_, err := state.Client.All.CreateMsgVpnJndiTopic(params, state.Auth)
	if err != nil {
		if semp, err2 := provider.GetSempError(err); err2 == nil {
			if *semp.Status == "ALREADY_EXISTS" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("JNDI Topic %s already exists. Going to import state from Broker", body.TopicName),
				})
			} else {
				return provider.AppendError(diags, err)
			}
		} else {
			return provider.AppendError(diags, err)
		}
	}
	d.SetId("sol_" + msgvpn + "_jq_" + body.TopicName)
	return append(diags, readJndiTopic(ctx, d, meta)...)
}

func readJndiTopic(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	jndi := d.Get(JNDI_NAME).(string)
	msgvpn := d.Get(MSG_VPN_NAME).(string)

	params := all.NewGetMsgVpnJndiTopicParams().WithContext(ctx).WithTopicName(jndi).WithMsgVpnName(msgvpn)
	resp, err := state.Client.All.GetMsgVpnJndiTopic(params, state.Auth)
	if err != nil {
		if semp, err2 := provider.GetSempError(err); err2 == nil {
			if *semp.Status == provider.SEMP_STATUS_NOT_FOUND {
				d.SetId("")
				return append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("JNDI Topic %s does not exist. Going to remove ressource from tfstate", jndi),
				})
			} else {
				return provider.AppendError(diags, err)
			}
		} else {
			return provider.AppendError(diags, err)
		}
	}
	d.Set(TOPIC_NAME, resp.Payload.Data.PhysicalName)
	return diags
}

func deleteJndiTopic(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	jndi := d.Get(JNDI_NAME).(string)
	msgvpn := d.Get(MSG_VPN_NAME).(string)

	params := all.NewDeleteMsgVpnJndiTopicParams().WithContext(ctx).WithMsgVpnName(msgvpn).WithTopicName(jndi)
	_, err := state.Client.All.DeleteMsgVpnJndiTopic(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	return diags
}
