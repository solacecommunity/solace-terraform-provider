package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solacecommunity/solace-terraform-provider/provider"
)

// Main resource definition for Broker
func ResourceSyslog() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaSyslog(),
		// Provider CRUD functions
		CreateContext: createSyslog,
		ReadContext:   readSyslog,
		UpdateContext: updateSyslog,
		DeleteContext: deleteSyslog,
	}
}

func schemaSyslog() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		SYSLOG_NAME: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Defines the name of the syslog asset",
		},
		SYSLOG_HOST: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Defines the remote host where to send the syslog events",
		},
		SYSLOG_TRANSPORT: {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Defines the transport to the remote where to send the syslog events - either UDP or TCP",
			ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP"}, false),
		},
		SYSLOG_FACILITIES: {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional:    true,
			Description: "Defines the facilities for which we want syslog to send files",
		},
	}
}
func createSyslog(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*provider.ProviderState).SempV1Client
	syslog_name := d.Get(SYSLOG_NAME).(string)
	syslog_host := d.Get(SYSLOG_HOST).(string)
	syslog_transport := d.Get(SYSLOG_TRANSPORT).(string)

	err := client.CreateSyslog(ctx, syslog_name)
	if err != nil {
		return provider.AppendError(diags, err)
	}

	err = client.ConfigureSyslogHostAndTransport(ctx, syslog_name, syslog_host, syslog_transport)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	syslog_facilities_interface_array := d.Get(SYSLOG_FACILITIES).([]interface{})
	syslog_facilties_string_array := iSliceTosSlice(&syslog_facilities_interface_array)
	for _, facility := range syslog_facilties_string_array {
		err = client.ConfigureSyslogFacilities(ctx, syslog_name, facility)
		if err != nil {
			return provider.AppendError(diags, fmt.Errorf("cannot create syslog facility for syslog %s and facility %s: %s", syslog_name, facility, err))
		}
	}

	d.SetId("sol_syslog_" + syslog_name)
	return append(diags, readSyslog(ctx, d, meta)...)
	// return diags
}

func readSyslog(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*provider.ProviderState).SempV1Client
	syslog_name := d.Get(SYSLOG_NAME).(string)

	result, err := client.GetSyslog(ctx, syslog_name)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	d.Set(SYSLOG_NAME, result.Rpc.Show.Syslog.SyslogElement.Name)
	d.Set(SYSLOG_TRANSPORT, result.Rpc.Show.Syslog.SyslogElement.Hosts.Host.Transport)
	d.Set(SYSLOG_HOST, result.Rpc.Show.Syslog.SyslogElement.Hosts.Host.Address)
	if result.Rpc.Show.Syslog.SyslogElement.Facilities.Facility == nil {
		empty_string_array := make([]string, 0)
		empty_interface_array := sSliceToiSlice(&empty_string_array)
		d.Set(SYSLOG_FACILITIES, empty_interface_array)
	} else {
		string_array := result.Rpc.Show.Syslog.SyslogElement.Facilities.Facility
		interface_array := sSliceToiSlice(&string_array)
		d.Set(SYSLOG_FACILITIES, interface_array)
	}
	return diags
}

func updateSyslog(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*provider.ProviderState).SempV1Client
	syslog_name := d.Get(SYSLOG_NAME).(string)
	syslog_host := d.Get(SYSLOG_HOST).(string)
	syslog_transport := d.Get(SYSLOG_TRANSPORT).(string)
	syslog_facilities_interface_array := d.Get(SYSLOG_FACILITIES).([]interface{})
	syslog_facilties_string_array := iSliceTosSlice(&syslog_facilities_interface_array)
	// Change CLI user auth type
	if d.HasChanges(SYSLOG_HOST, SYSLOG_TRANSPORT) {
		err := client.SetSyslogHostAndTransport(ctx, syslog_name, syslog_host, syslog_transport)
		if err != nil {
			return provider.AppendError(diags, fmt.Errorf("cannot set syslog host %s for syslog %s: %s", syslog_host, syslog_name, err))
		}
	}
	if d.HasChanges(SYSLOG_FACILITIES) {
		err := client.RemoveSyslogFacilities(ctx, syslog_name)
		if err != nil {
			return provider.AppendError(diags, fmt.Errorf("cannot remove syslog facilities for syslog %s: %s", syslog_name, err))
		}
		for _, facility := range syslog_facilties_string_array {
			err := client.SetSyslogFacilities(ctx, syslog_name, facility)
			if err != nil {
				return provider.AppendError(diags, fmt.Errorf("cannot set syslog transport %s for syslog %s: %s", syslog_transport, syslog_name, err))
			}
		}

	}
	return diags
}

func deleteSyslog(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*provider.ProviderState).SempV1Client
	syslog_name := d.Get(SYSLOG_NAME).(string)

	err := client.DeleteSyslog(ctx, syslog_name)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	d.SetId("")
	return diags
}
