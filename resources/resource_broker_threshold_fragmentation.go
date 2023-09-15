package resources

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/solacecommunity/solace-terraform-provider/provider"
)

// Main resource definition for Client Usernames
func ResourceBrokerThresholdFragmentation() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaBrokerThresholdFragmentation(),
		// Provider CRUD functions
		CreateContext: createBrokerThresholdFragmentation,
		ReadContext:   readBrokerThresholdFragmentation,
		UpdateContext: updateBrokerThresholdFragmentation,
		DeleteContext: deleteBrokerThresholdFragmentation,
	}
}

func schemaBrokerThresholdFragmentation() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		FRAGMENTATION_THRESHOLD_ENABLED: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			Description: "Once enabled, threshold-based defragmentation runs are triggered by all of the following three criteria being true:" +
				"The amount of fragmentation has reached a certain percent (see fragmentatiton_threshold_fragmentation_percentage)." +
				"The amount of spool usage has reached a certain percent (see fragmentatiton_threshold_usage_percentage)." +
				"A minimum amount of time has passed since the previous defragmentation run regardless of its trigger (see fragmentatiton_threshold_min_interval)." +
				"The default value is false",
		},
		FRAGMENTATION_THRESHOLD_FRAGMENTATION_PERCENTAGE: {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "50",
			Description: "This value is the percentage of fragmentation between 30 and 100. The default value is '50'",
		},
		FRAGMENTATION_THRESHOLD_USAGE_PERCENTAGE: {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "50",
			Description: "This value is the percentage of spool usage between 30 and 100. The default value is '50'",
		},
		FRAGMENTATION_THRESHOLD_MIN_INTERVAL: {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "15",
			Description: "This value defines the minimum interval in minutes. The default value is '15'",
		},
	}
}

func createBrokerThresholdFragmentation(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*provider.ProviderState).SempV1Client
	fragmentation_threshold_enable := d.Get(FRAGMENTATION_THRESHOLD_ENABLED).(bool)

	if fragmentation_threshold_enable {
		fragmentation_percentage := d.Get(FRAGMENTATION_THRESHOLD_FRAGMENTATION_PERCENTAGE).(string)
		fragmentation_usage := d.Get(FRAGMENTATION_THRESHOLD_USAGE_PERCENTAGE).(string)
		fragmentation_interval := d.Get(FRAGMENTATION_THRESHOLD_MIN_INTERVAL).(string)
		err := client.EnableThresholdBrokerFragmentation(ctx)
		if err != nil {
			return provider.AppendError(diags, err)
		}
		err = client.SetThresholdBrokerFragmentation(ctx, "fragmentation-percentage", fragmentation_percentage)
		if err != nil {
			return provider.AppendError(diags, err)
		}
		err = client.SetThresholdBrokerFragmentation(ctx, "usage-percentage", fragmentation_usage)
		if err != nil {
			return provider.AppendError(diags, err)
		}
		err = client.SetThresholdBrokerFragmentation(ctx, "min-interval", fragmentation_interval)
		if err != nil {
			return provider.AppendError(diags, err)
		}
	}
	d.SetId("sol_broker_threshold_fragmentation")
	return append(diags, readBrokerThresholdFragmentation(ctx, d, meta)...)
	// return diags
}

func readBrokerThresholdFragmentation(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*provider.ProviderState).SempV1Client

	result, err := client.GetBrokerFragmentationSettings(ctx)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	fragmentation_threshold_enabled := result.Rpc.Show.MessageSpool.MessageSpoolInfo.DefragThresholdEnabled
	fragmentation_threshold_enabled_bool, err := strconv.ParseBool(fragmentation_threshold_enabled)
	if err != nil {
		return provider.AppendError(diags, err)
	}

	d.Set(FRAGMENTATION_THRESHOLD_ENABLED, fragmentation_threshold_enabled_bool)
	d.Set(FRAGMENTATION_THRESHOLD_FRAGMENTATION_PERCENTAGE, result.Rpc.Show.MessageSpool.MessageSpoolInfo.DefragThresholdSpoolFragmentationPercentage)
	d.Set(FRAGMENTATION_THRESHOLD_USAGE_PERCENTAGE, result.Rpc.Show.MessageSpool.MessageSpoolInfo.DefragThresholdSpoolUsagePercentage)
	d.Set(FRAGMENTATION_THRESHOLD_MIN_INTERVAL, result.Rpc.Show.MessageSpool.MessageSpoolInfo.DefragThresholdMinInterval)
	return diags
}

func updateBrokerThresholdFragmentation(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	fragmentation_threshold_enable := d.Get(FRAGMENTATION_THRESHOLD_ENABLED).(bool)
	client := meta.(*provider.ProviderState).SempV1Client
	if d.HasChanges(FRAGMENTATION_THRESHOLD_ENABLED) {
		if fragmentation_threshold_enable {
			createBrokerThresholdFragmentation(ctx, d, meta)
		} else {
			deleteBrokerThresholdFragmentation(ctx, d, meta)
		}
	}
	if fragmentation_threshold_enable {
		if d.HasChanges(FRAGMENTATION_THRESHOLD_FRAGMENTATION_PERCENTAGE, FRAGMENTATION_THRESHOLD_MIN_INTERVAL, FRAGMENTATION_THRESHOLD_USAGE_PERCENTAGE) {
			if d.HasChange(FRAGMENTATION_THRESHOLD_FRAGMENTATION_PERCENTAGE) {
				err := client.SetThresholdBrokerFragmentation(ctx, "fragmentation-percentage", d.Get(FRAGMENTATION_THRESHOLD_FRAGMENTATION_PERCENTAGE).(string))
				if err != nil {
					return provider.AppendError(diags, err)
				}
			}
			if d.HasChange(FRAGMENTATION_THRESHOLD_MIN_INTERVAL) {
				err := client.SetThresholdBrokerFragmentation(ctx, "min-interval", d.Get(FRAGMENTATION_THRESHOLD_MIN_INTERVAL).(string))
				if err != nil {
					return provider.AppendError(diags, err)
				}
			}
			if d.HasChange(FRAGMENTATION_THRESHOLD_USAGE_PERCENTAGE) {
				err := client.SetThresholdBrokerFragmentation(ctx, "usage-percentage", d.Get(FRAGMENTATION_THRESHOLD_USAGE_PERCENTAGE).(string))
				if err != nil {
					return provider.AppendError(diags, err)
				}
			}
		}
	}
	return diags
}

func deleteBrokerThresholdFragmentation(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*provider.ProviderState).SempV1Client
	err := client.DisableThresholdBrokerFragmentation(ctx)
	if err != nil {
		return provider.AppendError(diags, err)
	}

	return diags
}

// func validateBrokerThresholdFragmentationDailySettings(v interface{}, k string) (ws []string, errors []error) {
// 	value := v.(string)
// 	if value == "daily" {
// 		return
// 	}
// 	days := strings.Split(value, ",")
// 	if len(days) <= 0 {
// 		errors = append(errors, fmt.Errorf(
// 			"if using daily and providing numbers, those need to be sepereated by a comma => %q','", k))
// 	}
// 	for _, v2 := range days {
// 		if v_as_int, err := strconv.Atoi(v2); err != nil {
// 			errors = append(errors, fmt.Errorf(
// 				"if using daily you need to provide comma seperated, single digit numbers => %q','", v2))
// 		} else {
// 			if v_as_int >= 7 {
// 				errors = append(errors, fmt.Errorf(
// 					"if using daily you need to provide comma seperated, single digit numbers in the range of 0-6 (sunday-saturday) => %q','", v2))
// 			}
// 		}
// 	}
// 	return
// }
