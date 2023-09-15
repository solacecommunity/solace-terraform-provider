package resources

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/solacecommunity/solace-terraform-provider/provider"
)

// Main resource definition for Client Usernames
func ResourceBrokerScheduledFragmentation() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaBrokerScheduledFragmentation(),
		// Provider CRUD functions
		CreateContext: createBrokerFragmentation,
		ReadContext:   readBrokerFragmentation,
		UpdateContext: updateBrokerFragmentation,
		DeleteContext: deleteBrokerFragmentation,
	}
}

func schemaBrokerScheduledFragmentation() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		FRAGMENTATION_SCHEDULED_ENABLED: {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "To schedule days of the week to trigger a defragmentation run enable this option. The default value is false",
		},
		FRAGMENTATION_SCHEDULED_DAYS: {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "daily",
			Description: "This value is either 'daily' or a comma-separated list of days or numbers where 0 is Sunday, 1 is Monday, etc. The default value is 'daily'",
		},
		FRAGMENTATION_SCHEDULED_TIMES: {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "hourly",
			Description: "This value is either hourly or a comma-separated list of up to four times of the form hh:mm where hh is between 0 and 23, and mm is between 0 and 59. The default value is '00:00'",
		},
	}
}

func createBrokerFragmentation(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*provider.ProviderState).SempV1Client
	fragmentation_schedule_enable := d.Get(FRAGMENTATION_SCHEDULED_ENABLED).(bool)

	if fragmentation_schedule_enable {
		fragmentation_days := d.Get(FRAGMENTATION_SCHEDULED_DAYS).(string)
		fragmentation_times := d.Get(FRAGMENTATION_SCHEDULED_TIMES).(string)
		err := client.EnableScheduledBrokerFragmentation(ctx)
		if err != nil {
			return provider.AppendError(diags, err)
		}
		err = client.SetScheduledBrokerFragmentation(ctx, "days", fragmentation_days)
		if err != nil {
			return provider.AppendError(diags, err)
		}
		err = client.SetScheduledBrokerFragmentation(ctx, "times", fragmentation_times)
		if err != nil {
			return provider.AppendError(diags, err)
		}
	}

	d.SetId("sol_broker_scheduled_fragmentation")
	return append(diags, readBrokerFragmentation(ctx, d, meta)...)
	// return diags
}

func readBrokerFragmentation(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*provider.ProviderState).SempV1Client

	result, err := client.GetBrokerFragmentationSettings(ctx)
	if err != nil {
		return provider.AppendError(diags, err)
	}

	fragmentation_scheduled_enabled := result.Rpc.Show.MessageSpool.MessageSpoolInfo.DefragScheduleEnabled
	fragmentation_scheduled_enabled_bool, err := strconv.ParseBool(fragmentation_scheduled_enabled)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	d.Set(FRAGMENTATION_SCHEDULED_ENABLED, fragmentation_scheduled_enabled_bool)
	d.Set(FRAGMENTATION_SCHEDULED_DAYS, result.Rpc.Show.MessageSpool.MessageSpoolInfo.DefragScheduleDays)
	d.Set(FRAGMENTATION_SCHEDULED_TIMES, result.Rpc.Show.MessageSpool.MessageSpoolInfo.DefragScheduleTimes)
	return diags
}

func updateBrokerFragmentation(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	fragmentation_schedule_enable := d.Get(FRAGMENTATION_SCHEDULED_ENABLED).(bool)
	client := meta.(*provider.ProviderState).SempV1Client
	if d.HasChanges(FRAGMENTATION_SCHEDULED_ENABLED) {
		if fragmentation_schedule_enable {
			createBrokerFragmentation(ctx, d, meta)
		} else {
			deleteBrokerFragmentation(ctx, d, meta)
		}
	}
	if fragmentation_schedule_enable {
		if d.HasChanges(FRAGMENTATION_SCHEDULED_DAYS, FRAGMENTATION_SCHEDULED_TIMES) {
			if d.HasChange(FRAGMENTATION_SCHEDULED_DAYS) {
				err := client.SetScheduledBrokerFragmentation(ctx, "days", d.Get(FRAGMENTATION_SCHEDULED_DAYS).(string))
				if err != nil {
					return provider.AppendError(diags, err)
				}
			}
			if d.HasChange(FRAGMENTATION_SCHEDULED_TIMES) {
				err := client.SetScheduledBrokerFragmentation(ctx, "times", d.Get(FRAGMENTATION_SCHEDULED_TIMES).(string))
				if err != nil {
					return provider.AppendError(diags, err)
				}
			}
		}
	}
	return diags
}

func deleteBrokerFragmentation(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*provider.ProviderState).SempV1Client
	err := client.DisableScheduledBrokerFragmentation(ctx)
	if err != nil {
		return provider.AppendError(diags, err)
	}

	return diags
}

// func validateBrokerFragmentationDailySettings(v interface{}, k string) (ws []string, errors []error) {
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
