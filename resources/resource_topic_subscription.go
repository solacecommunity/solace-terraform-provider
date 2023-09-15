package resources

import (
	"context"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/solacecommunity/solace-go-semp-client/client/all"
	"github.com/solacecommunity/solace-go-semp-client/models"
	"github.com/solacecommunity/solace-terraform-provider/provider"
)

func syncTopicSubscriptions(ctx context.Context, d *schema.ResourceData, state *provider.ProviderState, msgvpn string, queue string) error {
	// get a list of existing exceptions from broker
	ref, err := readTopicSubscriptions(ctx, state, msgvpn, queue)
	if err != nil {
		return err
	}
	// get a list of estimated state
	t := d.Get(TOPIC_SUBSCRIPTION_LIST).(*schema.Set).List()
	target := iSliceTosSlice(&t)

	// compare both lists
	news, olds := SliceDelta(&ref, &target)
	// add new entries
	for _, topic := range news {
		err := createTopicSubscription(ctx, state, msgvpn, queue, topic)
		if err != nil {
			return err
		}
	}
	// delete old ones
	for _, topic := range olds {
		err := deleteTopicSubscription(ctx, state, msgvpn, queue, topic)
		if err != nil {
			return err
		}
	}
	return nil
}

// validateTopicString validated a string and returns an error if it is not valid
func validateTopicString(v interface{}, p cty.Path) diag.Diagnostics {
	topic := v.(string)
	var diags diag.Diagnostics

	if strings.HasPrefix(topic, "/") {
		diag := diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "wrong topic subscription",
			Detail:   "topic string MUST NOT start with /",
		}
		return append(diags, diag)
	}
	if strings.HasSuffix(topic, "/") {
		diag := diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "wrong topic subscription",
			Detail:   "topic string MUST NOT end with /",
		}
		return append(diags, diag)
	}
	if strings.Contains(topic, "//") {
		diag := diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "wrong topic subscription",
			Detail:   "topic string MUST NOT contain empty levels, that is //",
		}
		return append(diags, diag)
	}
	return diags
}

func createTopicSubscription(ctx context.Context, state *provider.ProviderState, msgvpn string, queue string, topic string) error {
	body := &models.MsgVpnQueueSubscription{
		QueueName:         queue,
		MsgVpnName:        msgvpn,
		SubscriptionTopic: topic,
	}

	params := all.NewCreateMsgVpnQueueSubscriptionParamsWithContext(ctx).WithMsgVpnName(msgvpn).WithQueueName(queue).WithBody(body)
	_, err := state.Client.All.CreateMsgVpnQueueSubscription(params, state.Auth)
	return err
}

func readTopicSubscriptions(ctx context.Context, state *provider.ProviderState, msgvpn string, queue string) ([]string, error) {
	count := int64(100000)
	params := all.NewGetMsgVpnQueueSubscriptionsParamsWithContext(ctx).WithQueueName(queue).WithMsgVpnName(msgvpn).WithCount(&count)
	resp, err := state.Client.All.GetMsgVpnQueueSubscriptions(params, state.Auth)
	if err != nil {
		return nil, err
	}

	result := make([]string, len(resp.Payload.Data))
	for i, v := range resp.Payload.Data {
		result[i] = v.SubscriptionTopic
	}
	return result, nil
}

func deleteTopicSubscription(ctx context.Context, state *provider.ProviderState, msgvpn string, queue string, topic string) error {
	params := all.NewDeleteMsgVpnQueueSubscriptionParamsWithContext(ctx).WithMsgVpnName(msgvpn).
		WithQueueName(queue).WithSubscriptionTopic(topic)
	_, err := state.Client.All.DeleteMsgVpnQueueSubscription(params, state.Auth)
	return err
}
