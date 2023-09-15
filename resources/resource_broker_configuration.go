package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/solacecommunity/solace-go-semp-client/client/all"
	"github.com/solacecommunity/solace-go-semp-client/models"
	"github.com/solacecommunity/solace-terraform-provider/provider"
)

// Main resource definition for Client Usernames
func ResourceBrokerConfiguration() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		// List of supported configuration fields for your resource
		Schema: schemaBrokerConfiguration(),
		// Provider CRUD functions
		CreateContext: createBrokerConfiguration,
		ReadContext:   readBrokerConfiguration,
		UpdateContext: updateBrokerConfiguration,
		DeleteContext: deleteBrokerConfiguration,
	}
}

func schemaBrokerConfiguration() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		TLS_SSH_CIPHER_SUITE_LIST: {
			Type:     schema.TypeString,
			Optional: true,
			Description: "The colon-separated list of cipher suites used for TLS secure shell connections (e.g. SSH, SFTP, SCP). " +
				"The value default implies all supported suites ordered from most secure to least secure.",
		},
		TLS_MSG_BACKBONE_CIPHER_SUITE_LIST: {
			Type:     schema.TypeString,
			Optional: true,
			Description: "The colon-separated list of cipher suites used for TLS data connections (e.g. client pub/sub). " +
				"The value default implies all supported suites ordered from most secure to least secure. The default value is default",
		},
		TLS_MANAGEMENT_CIPHER_SUITE_LIST: {
			Type:     schema.TypeString,
			Optional: true,
			Description: "The colon-separated list of cipher suites used for TLS management connections (e.g. SEMP, LDAP). " +
				"The value default implies all supported suites ordered from most secure to least secure. The default value is default.",
		},
		TLS_SERVER_CERTIFICATE: {
			Type:     schema.TypeString,
			Optional: true,
			Description: "The PEM formatted content for the server certificate used for TLS connections." +
				"It must consist of a private key and between one and three certificates comprising the certificate trust chain." +
				"This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4." +
				"Changing this attribute requires an HTTPS connection. The default value is \"\"",
		},
		TLS_SERVER_CERTIFICATE_PASSWORD: {
			Type:     schema.TypeString,
			Optional: true,
			Description: "The password for the server certificate used for TLS connections." +
				"This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4." +
				"Changing this attribute requires an HTTPS connection. The default value is \"\"",
		},
		GUARANTEED_MSGING_MAX_MSG_SPOOL_USAGE: {
			Type:     schema.TypeInt,
			Optional: true,
			Description: "The maximum total message spool usage allowed across all VPNs on this broker, in megabytes. " +
				"Recommendation: the maximum value should be less than 90 percent of the disk space allocated for the guaranteed message spool." +
				"The default value is `60000`",
			Default: 60000,
		},
		SERVICE_SEMP_SESSION_IDLE_TIMEOUT: {
			Type:     schema.TypeInt,
			Optional: true,
			Description: "The session idle timeout, in minutes. Sessions will be invalidated if there is no activity in this period of time." +
				"Changes to this attribute are synchronized to HA mates via config-sync.",
			Default: 480,
		},
		LOG_RETENTION_DURATION: {
			Type:     schema.TypeString,
			Optional: true,
			Description: "This configures the broker to retain log files for a maximum number of days specified by the max-num-days value " +
				"(the actual period that logs are retained for is subject to availability of disk space). The valid range of values is 2 to 90." +
				"Note that when there is heavy log output, multiple log files may be necessary to contain the logs for a single day.",
			Default: "30",
		},
	}
}

func getBrokerConfigurationModelFromResource(d *schema.ResourceData) *models.Broker {
	q := &models.Broker{}

	if v, s := d.GetOk(TLS_MANAGEMENT_CIPHER_SUITE_LIST); s {
		q.TLSCipherSuiteManagementList = v.(string)
	}
	if v, s := d.GetOk(TLS_MSG_BACKBONE_CIPHER_SUITE_LIST); s {
		q.TLSCipherSuiteMsgBackboneList = v.(string)
	}
	if v, s := d.GetOk(TLS_SSH_CIPHER_SUITE_LIST); s {
		q.TLSCipherSuiteSecureShellList = v.(string)
	}
	if v, s := d.GetOk(TLS_SERVER_CERTIFICATE); s {
		q.TLSServerCertContent = v.(string)
	}
	if v, s := d.GetOk(TLS_SERVER_CERTIFICATE_PASSWORD); s {
		q.TLSServerCertPassword = v.(string)
	}
	if v, s := d.GetOk(GUARANTEED_MSGING_MAX_MSG_SPOOL_USAGE); s {
		q.GuaranteedMsgingMaxMsgSpoolUsage = int64(v.(int))
	}
	return q
}

func createBrokerConfiguration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("broker_tls_configuration")
	diags = updateBrokerConfiguration(ctx, d, meta)
	return append(diags, readBrokerConfiguration(ctx, d, meta)...)
}

func readBrokerConfiguration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*provider.ProviderState).SempV1Client
	state := meta.(*provider.ProviderState)
	params := all.NewGetBrokerParamsWithContext(ctx)
	resp, err := state.Client.All.GetBroker(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	p := resp.Payload.Data
	resp2, err2 := client.GetLogRetention(ctx)
	if err2 != nil {
		return provider.AppendError(diags, err2)
	}
	retention_days := resp2.Rpc.Show.Logging.Config.Retention
	d.Set(TLS_MANAGEMENT_CIPHER_SUITE_LIST, p.TLSCipherSuiteManagementList)
	d.Set(TLS_MSG_BACKBONE_CIPHER_SUITE_LIST, p.TLSCipherSuiteMsgBackboneList)
	d.Set(TLS_SSH_CIPHER_SUITE_LIST, p.TLSCipherSuiteSecureShellList)
	d.Set(GUARANTEED_MSGING_MAX_MSG_SPOOL_USAGE, p.GuaranteedMsgingMaxMsgSpoolUsage)
	d.Set(SERVICE_SEMP_SESSION_IDLE_TIMEOUT, 480)
	d.Set(LOG_RETENTION_DURATION, retention_days)
	return diags
}

func updateBrokerConfiguration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	state := meta.(*provider.ProviderState)
	client := meta.(*provider.ProviderState).SempV1Client
	body := getBrokerConfigurationModelFromResource(d)

	params := all.NewUpdateBrokerParamsWithContext(ctx).WithBody(body)
	_, err := state.Client.All.UpdateBroker(params, state.Auth)
	if err != nil {
		return provider.AppendError(diags, err)
	}
	err2 := client.SetSEMPIdleTimeout(ctx, 480)
	if err2 != nil {
		return provider.AppendError(diags, err2)
	}
	err3 := client.UpdateLogRetention(ctx, d.Get(LOG_RETENTION_DURATION).(string))
	if err3 != nil {
		return provider.AppendError(diags, err3)
	}
	// read after update to make sure tf state is in sync with broker
	return readBrokerConfiguration(ctx, d, meta)
}

func deleteBrokerConfiguration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	// not implemented for obvious reasons - you can't delete the complete broker
	return diags
}
