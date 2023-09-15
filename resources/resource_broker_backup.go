package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/solacecommunity/solace-terraform-provider/provider"
)

// Main resource definition for Broker
func ResourceBrokerBackup() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaBrokerBackup(),
		// Provider CRUD functions
		CreateContext: createBrokerBackup,
		ReadContext:   readBrokerBackup,
		UpdateContext: updateBrokerBackup,
		DeleteContext: deleteBrokerBackup,
	}
}

func schemaBrokerBackup() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		BACKUP_DAYS_OF_WEEK: {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "daily",
			Description: "This field is either the entry “daily”, or a list of named days from Sunday to Saturday separated by commas with no spaces" +
				"or a list of numbers from 0 to 6 representing the named days separated by commas with no spaces," +
				"where 0 is Sunday, 1 is Monday, on through to 6 for Saturday. Default is “daily”.",
		},
		BACKUP_TIMES_OF_DAY: {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "03:30",
			Description: "This field is either the entry “hourly”, or a list of up to four times of day in the format hh:mm separated by commas without spaces" +
				"where hh is 0 to 23 representing hours, and mm is 0 to 59 representing minutes.",
		},
		BACKUP_MAXIMUM_BACKUPS: {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  5,
			Description: "This field is the maximum number of scheduled backups to keep from 1 to 25. When a new scheduled backup causes the number of backups to" +
				"exceed the set maximum, the oldest backup file is deleted. The default value is 5 backups if this parameter is not provided.",
		},
	}
}
func createBrokerBackup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*provider.ProviderState).SempV1Client
	days_of_week := d.Get(BACKUP_DAYS_OF_WEEK).(string)
	times_of_day := d.Get(BACKUP_TIMES_OF_DAY).(string)
	maximum_backups := d.Get(BACKUP_MAXIMUM_BACKUPS).(int)
	err := client.CreateBrokerBackup(ctx, days_of_week, times_of_day, maximum_backups)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	d.SetId("broker_backup")
	return append(diags, readBrokerBackup(ctx, d, meta)...)
}

func readBrokerBackup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*provider.ProviderState).SempV1Client

	result, err := client.GetBrokerBackup(ctx)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	d.Set(BACKUP_DAYS_OF_WEEK, result.Rpc.Show.Backup.DaysOfWeek)
	d.Set(CLI_USER_DEFAULT_GLOBAL_ACCESS_LEVEL, result.Rpc.Show.Backup.TimesOfDay)
	d.Set(CLI_USER_DEFAULT_MESSAGEVPN_ACCESS_LEVEL, result.Rpc.Show.Backup.MaxBackups)
	return diags
}

func updateBrokerBackup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*provider.ProviderState).SempV1Client
	if d.HasChanges(BACKUP_DAYS_OF_WEEK, BACKUP_TIMES_OF_DAY, BACKUP_MAXIMUM_BACKUPS) {
		days_of_week := d.Get(BACKUP_DAYS_OF_WEEK).(string)
		times_of_day := d.Get(BACKUP_TIMES_OF_DAY).(string)
		maximum_backups := d.Get(BACKUP_MAXIMUM_BACKUPS).(int)
		err := client.UpdateBrokerBackup(ctx, days_of_week, times_of_day, maximum_backups)
		if err != nil {
			return provider.AppendError(diags, err)
		}
	}
	return append(diags, readBrokerBackup(ctx, d, meta)...)
}

func deleteBrokerBackup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*provider.ProviderState).SempV1Client
	err := client.DeleteBrokerBackup(ctx)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	return diags
}
