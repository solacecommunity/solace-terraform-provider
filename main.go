package main

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/go-openapi/strfmt"

	httptransport "github.com/go-openapi/runtime/client"
	apiclient "github.com/solacecommunity/solace-go-semp-client/client"
	"github.com/solacecommunity/solace-terraform-provider/datasources"
	"github.com/solacecommunity/solace-terraform-provider/provider"
	"github.com/solacecommunity/solace-terraform-provider/resources"
	"github.com/solacecommunity/solace-terraform-provider/sempv1"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"solace_dmr_cluster":                       resources.ResourceDmrCluster(),
			"solace_dmr_bridge":                        resources.ResourceDmrBridge(),
			"solace_bridge":                            resources.ResourceBridge(),
			"solace_bridge_remote_subscription":        resources.ResourceBridgeRemoteSubscription(),
			"solace_bridge_remote_msg_vpn":             resources.ResourceBridgeRemoteMsgVpn(),
			"solace_external_link":                     resources.ResourceExternalLink(),
			"solace_remote_address":                    resources.ResourceRemoteAddress(),
			"solace_broker_backup":                     resources.ResourceBrokerBackup(),
			"solace_broker_authentication":             resources.ResourceBrokerAuthentication(),
			"solace_msgvpn":                            resources.ResourceMsgVpn(),
			"solace_cli_user":                          resources.ResourceCliUser(),
			"solace_ldap_profile":                      resources.ResourceLdapProfile(),
			"solace_queue":                             resources.ResourceQueue(),
			"solace_queue_template":                    resources.ResourceQueueTemplate(),
			"solace_jndi_queue":                        resources.ResourceJndiQueue(),
			"solace_jndi_topic":                        resources.ResourceJndiTopic(),
			"solace_jndi_connectionfactory":            resources.ResourceJndiConnectionFactory(),
			"solace_client_username":                   resources.ResourceClientUsername(),
			"solace_client_profile":                    resources.ResourceClientProfile(),
			"solace_acl_profile":                       resources.ResourceAclProfile(),
			"solace_client_cert_authority":             resources.ResourceClientCertAuthoritoy(),
			"solace_domain_cert_authority":             resources.ResourceDomainCertAuthoritoy(),
			"solace_ldap_cli_group":                    resources.ResourceLdapCliGroup(),
			"solace_broker_configuration":              resources.ResourceBrokerConfiguration(),
			"solace_syslog_configuration":              resources.ResourceSyslog(),
			"solace_rest_delivery_point":               resources.ResourceRestDeliveryPoint(),
			"solace_rest_delivery_point_queue_binding": resources.ResourceRestDeliveryPointQueueBinding(),
			"solace_rest_delivery_point_rest_consumer": resources.ResourceRestDeliveryPointRestConsumer(),
			"solace_broker_threshold_fragmentation":    resources.ResourceBrokerThresholdFragmentation(),
			"solace_broker_scheduled_fragmentation":    resources.ResourceBrokerScheduledFragmentation(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"solace_broker": datasources.ResourceBroker(),
			"solace_msgvpn": datasources.ResourceMsgVpn(),
		},
		Schema:        provider.Schema(),
		ConfigureFunc: configure,
	}
}

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return Provider()
		},
	})
}

func configure(d *schema.ResourceData) (interface{}, error) {
	host := d.Get(provider.HOST).(string)
	port := d.Get(provider.PORT).(int)
	transport := httptransport.New(host+":"+strconv.Itoa(port), "/SEMP/v2/config", []string{d.Get(provider.SCHEMA).(string)})

	state := &provider.ProviderState{
		Auth:   httptransport.BasicAuth(d.Get(provider.ADMIN_USER).(string), d.Get(provider.ADMIN_PASSWD).(string)),
		Client: apiclient.New(transport, strfmt.Default),
		SempV1Client: &sempv1.SempV1Client{
			Username: d.Get(provider.ADMIN_USER).(string),
			Passwort: d.Get(provider.ADMIN_PASSWD).(string),
			Url:      d.Get(provider.SCHEMA).(string) + "://" + host + ":" + strconv.Itoa(port) + "/SEMP",
			Timeout:  time.Second * 60,
		},
	}
	return state, nil
}
