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

// Main resource definition for Queue entities
func ResourceQueue() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaQueue(),
		// Provider CRUD functions
		CreateContext: createQueueFunc,
		ReadContext:   readQueueFunc,
		UpdateContext: updateQueueFunc,
		DeleteContext: deleteQueueFunc,
	}
}

func schemaQueue() map[string]*schema.Schema {
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
		ACCESS_TYPE: {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     false,
			Default:      "non-exclusive", // Solace has a default of exclusive which is not what we want to be the default
			ValidateFunc: validation.StringInSlice([]string{"exclusive", "non-exclusive"}, false),
		},
		DEAD_MSG_QUEUE: {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the Dead Message Queue (DMQ) used by the Queue. The default value is \"#DEAD_MSG_QUEUE\".",
			Default:     "#DEAD_MSG_QUEUE",
		},
		EGRESS_ENABLED: {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Enable or disable the transmission of messages from the Queue. The default value is true.",
		},
		INGRESS_ENABLED: {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Enable or disable the reception of messages to the Queue. The default value is true.",
		},
		MAX_MSG_SIZE: {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     provider.DEFAULT_QUEUE_MAX_MSG_SIZE,
			Description: "The maximum message size allowed in the Queue, in bytes (B). The default value is " + fmt.Sprintf("%d", provider.DEFAULT_QUEUE_MAX_MSG_SIZE),
		},
		MAX_SPOOL_USAGE: {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     provider.DEFAULT_QUEUE_MAX_SPOOL_USAGE,
			Description: "The maximum message spool usage allowed by the Queue, in megabytes (MB). A value of 0 only allows spooling of the last message received and disables quota checking. The default value is 4000.",
		},
		MAX_SPOOL_USAGE_ALERT_SET_PCT: {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      provider.DEFAULT_QUEUE_MAX_SPOOL_USAGE_SET_PCT,
			Description:  "The thresholds for the message spool usage alert of the Queue, relative to Messages Queued Quota.",
			ValidateFunc: validation.IntBetween(0, 100),
		},
		MAX_SPOOL_USAGE_ALERT_CLEAR_PCT: {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      provider.DEFAULT_QUEUE_MAX_SPOOL_USAGE_CLEAR_PCT,
			Description:  "The thresholds for the message spool usage alert of the Queue, relative to Messages Queued Quota.",
			ValidateFunc: validation.IntBetween(0, 100),
		},
		MAX_REDELIVERY_COUNT: {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     provider.DEFAULT_MAX_REDELIVERY_COUNT,
			Description: "The maximum number of times the Queue will attempt redelivery of a message prior to it being discarded or moved to the DMQ. A value of 0 means to retry forever. The default value is 0.",
		},
		MAX_BIND_COUNT: {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      provider.DEFAULT_MAX_BIND_COUNT,
			Description:  "The maximum number of consumer flows that can bind to the Queue. The default value is 1000.",
			ValidateFunc: validation.IntAtLeast(1),
		},
		MAX_BIND_COUNT_ALERT_SET_PCT: {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      provider.DEFAULT_MAX_BIND_COUNT_ALERT_SET_PCT,
			ValidateFunc: validation.IntBetween(0, 100),
		},
		MAX_BIND_COUNT_ALERT_CLEAR_PCT: {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      provider.DEFAULT_MAX_BIND_COUNT_ALERT_CLEAR_PCT,
			ValidateFunc: validation.IntBetween(0, 100),
		},
		OWNER: {
			Type:        schema.TypeString,
			ForceNew:    false,
			Optional:    true,
			Description: "The Client Username that owns the Queue and has permission equivalent to \"delete\". The default value is \"\".",
		},
		PERMISSION: {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "no-access",
			ForceNew:     false,
			ValidateFunc: validation.StringInSlice([]string{"no-access", "read-only", "consume", "modify-topic", "delete"}, false),
			Description:  "The permission level for all consumers of the Queue, excluding the owner. The default value is \"no-access\"",
		},
		TOPIC_SUBSCRIPTION_LIST: {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type:             schema.TypeString,
				ValidateDiagFunc: validateTopicString,
			},
		},
		RESPECT_TTL_ENABLED: {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable or disable the respecting of the time-to-live (TTL) for messages in the Queue. When enabled, expired messages are discarded or moved to the DMQ. The default value is `false`.",
		},
		MAX_TTL: {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     false,
			Description: "The maximum time in seconds a message can stay in the Queue when respectTtlEnabled is true. A message expires when the lesser of the sender assigned time-to-live (TTL) in the message and the `maxTtl` configured for the Queue, is exceeded. A value of 0 disables expiry. The default value is `0`",
		},
	}
}

// Creates a queue model based on the terraform resource state.
func getQueueModelFromResource(d *schema.ResourceData) *models.MsgVpnQueue {
	q := &models.MsgVpnQueue{
		QueueName:  d.Get(QUEUE_NAME).(string),
		MsgVpnName: d.Get(MSG_VPN_NAME).(string),
	}

	if v, s := d.GetOk(ACCESS_TYPE); s {
		q.AccessType = v.(string)
	}
	if v, s := d.GetOk(DEAD_MSG_QUEUE); s {
		q.DeadMsgQueue = v.(string)
	}
	if v, s := d.GetOk(EGRESS_ENABLED); s {
		q.EgressEnabled = v.(bool)
	}
	if v, s := d.GetOk(INGRESS_ENABLED); s {
		q.IngressEnabled = v.(bool)
	}
	if v, s := d.GetOk(MAX_MSG_SIZE); s {
		q.MaxMsgSize = int32(v.(int)) // BUHU
	}
	if v, s := d.GetOk(MAX_REDELIVERY_COUNT); s {
		q.MaxRedeliveryCount = int64(v.(int))
	}
	if v, s := d.GetOk(OWNER); s {
		q.Owner = v.(string)
	}
	if v, s := d.GetOk(PERMISSION); s {
		q.Permission = v.(string)
	}

	if v, s := d.GetOk(MAX_SPOOL_USAGE); s {
		q.MaxMsgSpoolUsage = int64(v.(int))
	}
	if v, s := d.GetOk(MAX_SPOOL_USAGE_ALERT_SET_PCT); s {
		if q.EventMsgSpoolUsageThreshold == nil {
			q.EventMsgSpoolUsageThreshold = &models.EventThreshold{}
		}
		q.EventMsgSpoolUsageThreshold.SetPercent = int64(v.(int))
	}
	if v, s := d.GetOk(MAX_SPOOL_USAGE_ALERT_CLEAR_PCT); s {
		if q.EventMsgSpoolUsageThreshold == nil {
			q.EventMsgSpoolUsageThreshold = &models.EventThreshold{}
		}
		q.EventMsgSpoolUsageThreshold.ClearPercent = int64(v.(int))
	}

	if v, s := d.GetOk(MAX_BIND_COUNT); s {
		q.MaxBindCount = int64(v.(int))
	}
	if v, s := d.GetOk(MAX_BIND_COUNT_ALERT_SET_PCT); s {
		if q.EventBindCountThreshold == nil {
			q.EventBindCountThreshold = &models.EventThreshold{}
		}
		q.EventBindCountThreshold.SetPercent = int64(v.(int))
	}
	if v, s := d.GetOk(MAX_BIND_COUNT_ALERT_CLEAR_PCT); s {
		if q.EventBindCountThreshold == nil {
			q.EventBindCountThreshold = &models.EventThreshold{}
		}
		q.EventBindCountThreshold.ClearPercent = int64(v.(int))
	}
	if v, s := d.GetOk(RESPECT_TTL_ENABLED); s {
		q.RespectTTLEnabled = v.(bool)
	}
	if v, s := d.GetOk(MAX_TTL); s {
		q.MaxTTL = int64(v.(int))
	}

	return q
}

// Updates the terraform resource state according to the solace queue model returned from the SEMP API
func updateQueueResourceFromModel(d *schema.ResourceData, q *models.MsgVpnQueue) {
	d.Set(ACCESS_TYPE, q.AccessType)
	d.Set(DEAD_MSG_QUEUE, q.DeadMsgQueue)
	d.Set(EGRESS_ENABLED, q.EgressEnabled)
	d.Set(INGRESS_ENABLED, q.IngressEnabled)
	d.Set(MAX_MSG_SIZE, int(q.MaxMsgSize))
	d.Set(MAX_SPOOL_USAGE, int(q.MaxMsgSpoolUsage))
	d.Set(MAX_REDELIVERY_COUNT, int(q.MaxRedeliveryCount))
	d.Set(MAX_BIND_COUNT, int(q.MaxBindCount))
	d.Set(OWNER, q.Owner)
	d.Set(PERMISSION, q.Permission)
	if q.EventMsgSpoolUsageThreshold != nil {
		d.Set(MAX_SPOOL_USAGE_ALERT_SET_PCT, int(q.EventMsgSpoolUsageThreshold.SetPercent))
		d.Set(MAX_SPOOL_USAGE_ALERT_CLEAR_PCT, int(q.EventMsgSpoolUsageThreshold.ClearPercent))
	}
	if q.EventBindCountThreshold != nil {
		d.Set(MAX_BIND_COUNT_ALERT_SET_PCT, int(q.EventBindCountThreshold.SetPercent))
		d.Set(MAX_BIND_COUNT_ALERT_CLEAR_PCT, int(q.EventBindCountThreshold.ClearPercent))
	}
	d.Set(RESPECT_TTL_ENABLED, q.RespectTTLEnabled)
	d.Set(MAX_TTL, q.MaxTTL)
}

func createQueueFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	state := meta.(*provider.ProviderState)
	body := getQueueModelFromResource(d)

	params := all.NewCreateMsgVpnQueueParams().WithContext(ctx).WithMsgVpnName(msgvpn).WithBody(body)
	_, err := state.Client.All.CreateMsgVpnQueue(params, state.Auth)
	if err != nil {
		if semp, err2 := provider.GetSempError(err); err2 == nil {
			if *semp.Status == "ALREADY_EXISTS" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Queue %s already exists. Going to import state from Broker", body.QueueName),
				})
			} else {
				return provider.AppendError(diags, err)
			}
		} else {
			return provider.AppendError(diags, err)
		}
	}

	if d.HasChange(TOPIC_SUBSCRIPTION_LIST) {
		err := syncTopicSubscriptions(ctx, d, state, msgvpn, body.QueueName)
		if err != nil {
			return provider.AppendError(diags, err)
		}
	}

	d.SetId("sol_" + msgvpn + "_q_" + body.QueueName)
	return append(diags, readQueueFunc(ctx, d, meta)...)
}

func readQueueFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	queue := d.Get(QUEUE_NAME).(string)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	state := meta.(*provider.ProviderState)
	params := all.NewGetMsgVpnQueueParams().WithContext(ctx).WithQueueName(queue).WithMsgVpnName(msgvpn)
	resp, err := state.Client.All.GetMsgVpnQueue(params, state.Auth)
	if err != nil {
		if semp, err2 := provider.GetSempError(err); err2 == nil {
			if *semp.Status == provider.SEMP_STATUS_NOT_FOUND {
				d.SetId("")
				return append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Queue %s does not exist. Going to remove ressource from tfstate", queue),
				})
			} else {
				return provider.AppendError(diags, err)
			}
		} else {
			return provider.AppendError(diags, err)
		}
	}
	updateQueueResourceFromModel(d, resp.Payload.Data)

	// read topic subcriptions
	ts, err := readTopicSubscriptions(ctx, state, msgvpn, queue)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	d.Set(TOPIC_SUBSCRIPTION_LIST, sSliceTosSet(&ts))

	return diags
}

func updateQueueFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	queue := d.Get(QUEUE_NAME).(string)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	state := meta.(*provider.ProviderState)
	body := getQueueModelFromResource(d)

	_, ACCESS_TYPE_IS_SET := d.GetOk(ACCESS_TYPE)
	_, OWNER_IS_SET := d.GetOk(OWNER)
	_, PERMISSION_IS_SET := d.GetOk(PERMISSION)

	if ACCESS_TYPE_IS_SET || OWNER_IS_SET || PERMISSION_IS_SET {
		body.EgressEnabled = false
	}
	params := all.NewUpdateMsgVpnQueueParams().WithContext(ctx).WithQueueName(queue).WithMsgVpnName(msgvpn).WithBody(body)
	resp, err := state.Client.All.UpdateMsgVpnQueue(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	if ACCESS_TYPE_IS_SET || OWNER_IS_SET || PERMISSION_IS_SET {
		err := setEgressIngressQueueState(ctx, state, msgvpn, queue, true, true)
		if err != nil {
			return provider.AppendError(diags, err)
		}
	}
	updateQueueResourceFromModel(d, resp.Payload.Data)

	if d.HasChange(TOPIC_SUBSCRIPTION_LIST) {
		err := syncTopicSubscriptions(ctx, d, state, msgvpn, body.QueueName)
		if err != nil {
			return provider.AppendError(diags, err)
		}
	}

	return diags
}

func deleteQueueFunc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	queue := d.Get(QUEUE_NAME).(string)
	msgvpn := d.Get(MSG_VPN_NAME).(string)
	state := meta.(*provider.ProviderState)

	params := all.NewDeleteMsgVpnQueueParams().WithContext(ctx).WithMsgVpnName(msgvpn).WithQueueName(queue)
	_, err := state.Client.All.DeleteMsgVpnQueue(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	return diags
}

func setEgressIngressQueueState(ctx context.Context, state *provider.ProviderState, vpnName string, queueName string, ingress bool, egress bool) error {
	body := &models.MsgVpnQueue{
		MsgVpnName:     vpnName,
		QueueName:      queueName,
		EgressEnabled:  egress,
		IngressEnabled: ingress,
	}
	params := all.NewUpdateMsgVpnQueueParamsWithContext(ctx).WithMsgVpnName(vpnName).WithQueueName(queueName).WithBody(body)
	_, err := state.Client.All.UpdateMsgVpnQueue(params, state.Auth)
	return err
}
