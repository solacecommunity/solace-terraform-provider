[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.0%20adopted-ff69b4.svg)](CODE_OF_CONDUCT.md)

# Solace Terraform Provider

**This repo contains a terraform provider for Solace Event Brokers (with a focus on self-managed software brokers). We're always improving it, but if there are features you want please feel free to open an enhancement request, or even better a PR!**

## Disclaimer
This provider aims to support the creators so far to manage Solace Broker instances for company internal use. 
Therefore not all properties of a managed resource are exposed by the provider. That is to keep the 
management as simple as possible and hide unused features unless they become relevant. 

That also means, that some default values set by the provider do not match with the default values 
defined by the Solace Broker. Instead, this provider overwrites those values which are supposed to 
work best for most production requirements. See the respective documentation of the resources for informations 
what has been changed. You can still re-overwrite these values in the actual terraform configuration as you prefer. 

This provider is independend from the official Solace Terraform Provider: https://github.com/SolaceProducts/terraform-provider-solacebroker/blob/main/README.md
It`s suggested to use the official Terraform Provider for all general requirements (e.g. Solace Cloud, Appliance or managing assets within a software broker). 
This provider here is mainly used for self-managed software brokers, where additional Terraform resources might be required. 

## Getting started quickly
### Build local
To install the provider locally, you can either clone this repo and run the Makefile script 
```
make build
```
followed by deploying the binary into a specific directy which depends on your local operating system.
For Mac OS and Linux (and other Unix-like systems) it is usually
```
$HOME/.terraform.d/plugins
```
Followed by the directory layout `HOSTNAME/NAMESPACE/TYPE/VERSION/TARGET`
A full directory path for the Solace provider could look like
`$HOME/.terraform.d/plugins/localhost/local/solace/1.0.0/linux_amd64`
in which you can paste your binary from the build operation.
It might be needed that you rename your binary while moving it into the directory path to follow terraform conventions.
The complete path (including the binary) would look like this.
`$HOME/.terraform.d/plugins/localhost/local/solace/1.0.0/linux_amd64/terraform-provider-solace_v1.0.0`
Afterwards you can refer the provider in your terraform script like so:
```terraform
terraform {
  required_providers {
    solace = {
      source  = "localhost/local/solace"
      version = "1.0.0"
    }
  }
}
```
If this does not work because Terraform can not find the specified provider please follow those steps below:
For Terraform to find it, you must tell it where to search local for providers. Create a .terraformrc file inside your home directory
and fill it with the following content:

```terraform
provider_installation {
  filesystem_mirror {
    path    = "$HOME/.terraform.d/plugins/"
    include = ["localhost/*/*"]
  }
  direct {
    exclude = ["localhost/*/*"]
  }
}
```
This will tell Terraform that you want every provider which starts with localhost to be searched for on the local filesystem (filesystem_mirror) while all other providers will be searched with the direct path (meaning the official terraform registry).

## Configuration
The provider supports a set of configuation values. There are 3 parameters that must be set, that is `host`, `admin_user` and `admin_password`. It is recommended to keep those values externally and use variables for that.

```terraform
provider "solace" {
  admin_user     = <username!>
  admin_password = <password!>
  host           = <host!>
  port           = <port default=443>
  schema         = <schema default="https">
}
```

#### **_Warning : Some resources that are listed below can only be changed/created/updated with the Creation-Admin-User of the broker instance!_**
Following are the (already noted) "super-admin" resources:
- solace_broker_authentication
- solace_msgvpn
- ... more to follow...

## Documentation
### Datasources

The provides supports queriing the following datasources:

| Datasource Name | Description              |
| --------------- | ------------------------ |
| `solace_broker` | Query data of the Broker |
| `solace_msgvpn` | Query data of a MsgVPN   |

#### Broker [`solace_broker`]

To query global data from a Broker, there is no parameter required.

```tf
data "solace_broker" "default" {
}
```

You can now refer to some globel variables, e.g. in output statements:

```tf
output "Broker_Version" {
  value = data.solace_broker.default.version
}
```
This will print something similar like this, after a `terraform apply`:
```
Broker_Version = "Solace PubSub+ Standard Version 9.8.0.12"
```

The following properties are available:
| Property  | Description                               | Example Value                              |
| --------- | ----------------------------------------- | ------------------------------------------ |
| `version` | Verbose version and edition of the broker | `Solace PubSub+ Standard Version 9.8.0.12` |
| `build`   | Version of the broker                     | `9.8.0.12`                                 |

#### MsgVpn [`solace_msgvpn`]

To query data from a MsgVPN you simply need to define it's name:

```tf
data "solace_msgvpn" "default" {
  msg_vpn_name = "default"
}
```

### Resources

The provider supports management of the following resources:

| Resource Name                   | Description                                      |
| ------------------------------- | ------------------------------------------------ |
| `solace_msgvpn`                 | Manage MsgVPN                                    |
| `solace_ldap_profile`           | Manage LDAP Profiles                             |
| `solace_ldap_cli_group`         | Manage LDAP CLI Groups                           |
| `solace_cli_user`               | Manage local CLI Users                           |
| `solace_queue`                  | Manage Queues for persistent messaging           |
| `solace_queue_subscription`     | Manage topic subscriptions for queues            |
| `solace_jndi_queue`             | Manage JNDI Aliases for queues                   |
| `solace_jndi_connectionfactory` | Manage JNDI Connection Factories for JMS Clients |
| `solace_client_username`        | Manage client usernames                          |
| `solace_client_profile`         | Manage client profiles                           |
| `solace_acl_profile`            | Manage ACL profiles                              |
| `solace_client_cert_authority`  | Manage Client Certificate Authorities            |
| `solace_domain_cert_authority`  | Manage Domain Certificate Authorities            |
| `solace_broker_authentication`  | Manage Broker Authentication                     |
| `solace_broker_configuration`   | Manage Broker Configuration                      |
| `solace_broker_backup`          | Manage Broker Backups                            |
| `solace_dmrcluster`             | Manage Broker Clustering                         |
| `solace_syslog_configuration`   | Manage Broker Syslog Configuration               |

#### MsgVpn [`solace_msgvpn`]
For a detailed explanation of all parameters, see: [Manage MsgVpn Spec](https://docs.solace.com/API-Developer-Online-Ref-Documentation/swagger-ui/config/index.html#/msgVpn)


| Parameter      | Provider Default | Solace's Default | Comment                            |
| -------------- | ---------------- | ---------------- | ---------------------------------- |
| `jndi_enabled` | `true`           | `false`          | JNDI should be enabled per default |


#### LDAP Profiles [`solace_ldap_profiles`]
Manage LDAP Profiles for cli and client user authentication. 


| Parameter        | Comment                                                                      |
| ---------------- | ---------------------------------------------------------------------------- |
| `profile_name`   | Name of the profile                                                          |
| `index`          | Priority of the profile. Must be a value between 1 and 3                     |
| `enabled`        | Enable or disable the profile                                                |
| `host`           | Hostname / URL to the ldap server like `ldaps://hostname:636`     |
| `admin_dn`       | LDAP Proxy User used for binding with the ldap server                        |
| `admin_password` | Password for LDAP Proxy user                                                 |
| `base_dn`        | Root DN for searching users, e.g. `ou=usr,o=employee`                        |
| `search_filter`  | Search criteria to match LDAP User with Solace User `(uid=$CLIENT_USERNAME)` |

#### LDAP Cli Groups [`solace_ldap_cli_group`]
Manage LDAP groups for cli and client user authentication.


| Parameter                        | Comment                                                                                                            |
| -------------------------------- | ------------------------------------------------------------------------------------------------------------------ |
| `group_name`                     | FQN of the LDAP Group (e.g. cn=solace-admin,ou=solace,ou=apps,o=global)                                            |
| `global_access_level`            | Global access-level of the user. Can be any of `none`, `admin`, `read-only` or `read-write`. Default is `none`     |
| `msgvpn_default_access_level`    | Default MsgVpn access-level of the user. Can be any of `none`, `read-only` or `read-write`. Default is `none`      |
| `msgvpn_access_level_exceptions` | Map with specific access-level exceptions for MsgVpn. `<msgvpn> = <access-level>`                                  |

#### CLI Users [`solace_cli_user`]
Manage CLI users that have administrative rights on the broker or specific MsgVpn's. 
This resource uses SEMPv1 cli commands. See [CLI Reference](https://docs.solace.com/Solace-CLI/CLI-Reference/VMR_CLI_Commands.html#Root_enable_configure_username)

| Parameter                        | Comment                                                                                                            |
| -------------------------------- | ------------------------------------------------------------------------------------------------------------------ |
| `username`                       | Username of the cli user                                                                                           |
| `password`                       | Password of the cli user                                                                                           |
| `global_access_level`            | Global access-level of the cli user. Can be any of `none`, `admin`, `read-only` or `read-write`. Default is `none` |
| `msgvpn_default_access_level`    | Default MsgVpn access-level of the cli user. Can be any of `none`, `read-only` or `read-write`. Default is `none`  |
| `msgvpn_access_level_exceptions` | Map with specific access-level exceptions for MsgVpn. `<msgvpn> = <access-level>`                                  |


#### Queues [`solace_queue`]
For a detailed explanation of all parameters, see: [Manage Queue Spec](https://docs.solace.com/API-Developer-Online-Ref-Documentation/swagger-ui/config/index.html#/queue)

##### Changed Defaults

| Parameter         | Provider Default | Solace's Default | Comment                                                                                                     |
| ----------------- | ---------------- | ---------------- | ----------------------------------------------------------------------------------------------------------- |
| `access_type`     | `non-exclusive`  | `exclusive`      | For scalability and HA, most clients expect to have a round-robin / loadbalancing when connecting to queues |
| `egress_enabled`  | `true`           | `false`          | Queues will be created with egress enabled.                                                                 |
| `ingress_enabled` | `true`           | `false`          | Queues will be created with ingress enabled.                                                                |

##### Known problems / Todo's
- When terraform recreates the queue, for example because you've changed the owner, the dependent resources like Topic Subscriptions will still exist for terraform, even 
  if they are deleted on the Broker. Currently the only way to fix that is to manually remove the resources from terraform's state and run another `terraform apply`. To remove 
  the state of a resource use the following command:
  
```sh
$ terrraform state rm <resource_type>.<resource_name>
```

#### Topic Subscriptions [`solace_queue_subscription`]
For a detailed explanation of all parameters, see: [Manage Queue Spec](https://docs.solace.com/API-Developer-Online-Ref-Documentation/swagger-ui/config/index.html#/queue)

Be aware that queue subscriptions cannot be modified. Instead if you want to change a subscribed topic, this will always result in a DELETE of the old, and CREATE of the new subscription.

#### JNDI Connection Factories [`solace_jndi_connectionfactory`]
For a detailed explanation of all parameters, see: [Manage Jndi Spec](https://docs.solace.com/API-Developer-Online-Ref-Documentation/swagger-ui/config/index.html#/jndi)

##### Changed Defaults
- ConnectionFactories have a configuration option that is called `Transport Reply Timeout (ms)` in the UI, however in the SEMP API it is called `transportReadTimeout`. To be consistent with the UI, the terraform property to set this value is `transport_reply_timeout`.

#### JNDI Queues Aliases [`solace_jndi_queue`]
For a detailed explanation of all parameters, see: [Manage Jndi Spec](https://docs.solace.com/API-Developer-Online-Ref-Documentation/swagger-ui/config/index.html#/jndi)

##### :warning: Changed Property Names
Be aware that this provider uses a different naming convention for properties, compared to the Solace SEMP API.
That has been done to use the `queue_name` property consistently throughout the scripts.

| Provider Property | Solace Parameter | Comment                             |
| ----------------- | ---------------- | ----------------------------------- |
| `msg_vpn_name`    | `msg_vpn_name`   | The name of the Message VPN.        |
| `queue_name`      | `physicalName`   | The physical name of the JMS Queue. |
| `jndi_name`       | `queueName`      | The JNDI name of the JMS Queue.     |

#### Client Username [`solace_client_username`]

##### Changed Defaults

| Parameter             | Provider Default  | Solace's Default | Comment                                                                                                                        |
| --------------------- | ----------------- | ---------------- | ------------------------------------------------------------------------------------------------------------------------------ |
| `enabled`             | `true`            | `false`          | Enable or disable the Client Username. When disabled, all clients currently connected as the Client Username are disconnected. |
| `acl_profile_name`    | `nil`, *required* | `default`        | The profile `default` is discouraged and therfore you must set a valid option                                                  |
| `client_profile_name` | `nil`, *required* | `default`        | The profile `default` is discouraged and therfore you must set a valid option                                                  |

#### Client Profiles [`solace_client_profile`]
For a detailed explanation of all parameters, see: [Client Profile Spec](https://docs.solace.com/API-Developer-Online-Ref-Documentation/swagger-ui/config/index.html#/clientProfile)

##### Changed Defaults

| Parameter                              | Provider Default | Solace's Default | Comment                                                |
| -------------------------------------- | ---------------- | ---------------- | ------------------------------------------------------ |
| `allow_guaranteed_msg_receive_enabled` | `true`           | `false`          | Receiving of guaranteed messages is allowed by default |
| `allow_guaranteed_msg_send_enabled`    | `true`           | `false`          | Sending of guaranteed messages is allowed by default   |

#### ACL Profiles [`solace_acl_profile`]

##### Publish / Subscribe exception list

To manage ACL pub/sub exception lists, you can use the following snippet:

```terraform
resource "solace_acl_profile" "demo_acl_profile" {
  msg_vpn_name                  = var.msgvpn
  acl_profile_name              = "tfdemo_profile"
  client_connect_default_action = "allow"

  publish_exception_list = [
    "b/topic/to/publish@smf",
    "#P2P/QUE/terraform/queue2@smf", # the prefix #P2P/QUE/ is used to create publish permission on queues
    "a/topic/to/publish@smf",
    "b/topic/to/publish/>@mqtt",
  ]
```

To create publish permissions to queues, use `#P2P/QUE/` as prefix before the name of the queue.

Make sure to use prepend all entries with the so called **syntax definition** which clarifies if
you're topic expression uses `smf` Syntax (Wildcards are `*` and `>`) or `mqtt` Syntax (Wildcards
are `+` and `#`). 

* `my/topic/exception@smf`
* `my/topic/exception/with/wildcard/>@smf`
* `my/topic/+/with/wildcard/#@mqtt`

For a detailed explanation of all parameters, see: [ACL Profile Spec](https://docs.solace.com/API-Developer-Online-Ref-Documentation/swagger-ui/config/index.html#/aclProfile)

#### Changed Defaults

| Parameter                       | Provider Default | Solace's Default | Comment                                                     |
| ------------------------------- | ---------------- | ---------------- | ----------------------------------------------------------- |
| `client_connect_default_action` | `allow`          | `disallow`       | By default, all new profiles will allow clients to connect. |

#### Client Certificate Authorities [`solace_client_cert_authority`]
For a detailed explanation of all parameters, see: [Manage Client Certificate Authorities](https://docs.solace.com/API-Developer-Online-Ref-Documentation/swagger-ui/config/index.html#/clientCertAuthority)

| Parameter             | Comment                                                                                                                                                                                              |
| --------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `cert_authority_name` | The name of the certificate authority                                                                                                                                                                |
| `cert_content`        | The path to the PEM formatted content for the trusted root certificate of a client Certificate  Authority wrapped in the file() function of terraform (e.g. file("certificates/my-certificate.pem")) |


#### Domain Certificate Authorities [`solace_domain_cert_authority`]
For a detailed explanation of all parameters, see: [Manage Domain Certificate Authorities](https://docs.solace.com/API-Developer-Online-Ref-Documentation/swagger-ui/config/index.html#/domainCertAuthority)

| Parameter             | Comment                                                                                                                                                                                              |
| --------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `cert_authority_name` | The name of the certificate authority |
| `cert_content`        | The path to the PEM formatted content for the trusted root certificate of a domain Certificate  Authority wrapped in the file() function of terraform (e.g. file("certificates/my-certificate.pem")) |


#### Solace Broker Authentication [`solace_broker_authentication`]
Changes the way how people authenticate with the Solace Broker. Also sets default access levels for global and messagevpn if no other rules can be applied to a user.

| Parameter                        |  Comment                                                                                                               |
| ---------------------------------| ---------------------------------------------------------------------------------------------------------------------- |
| `default_global_access_level`    | Default Global access-level of the user. Can be any of `none`, `admin`, `read-only` or `read-write`. Default is `none` |
| `default_msgvpn_access_level`    | Default MsgVpn access-level of the user. Can be any of `none`, `read-only` or `read-write`. Default is `none`          |
| `cli_auth_type`                  | Defines which authentication type should be used for users. Can be any of `internal`, `ldap` or `radius`               |
| `cli_auth_type_profile`          | Defines which [`solace_ldap_profiles`] gets used for either the `ldap` or `radius` cli_auth_type                       |

#### Solace Broker Configuration [`solace_broker_configuration`]
Change some configuration settings for the broker including TLS/SSL settings and the maximum spool size.

| Parameter                               |  Comment                                                                                                               |
| --------------------------------------- | ---------------------------------------------------------------------------------------------------------------------- |
| `tls_ssh_cipher_suite_list`             | The colon-separated list of cipher suites used for TLS secure shell connections (e.g. SSH, SFTP, SCP). The value default implies all supported suites ordered from most secure to least secure.                                                                                             |
| `tls_msg_backbone_cipher_suite_list`    | The colon-separated list of cipher suites used for TLS data connections (e.g. client pub/sub). The value default implies all supported suites ordered from most secure to least secure. The default value is default                                                                            | 
| `tls_management_cipher_suite_list`      | The colon-separated list of cipher suites used for TLS management connections (e.g. SEMP, LDAP). The value default implies all supported suites ordered from most secure to least secure. The default value is default.                                                                       |
| `tls_server_certificate`                | The PEM formatted content for the server certificate used for TLS connections. It must consist of a private key and between one and three certificates comprising the certificate trust chain Changing this attribute requires an HTTPS connection. The default value is "",                   |
| `tls_server_certificate_password`       | The password for the server certificate used for TLS connections. Changing this attribute requires an HTTPS connection. The default value is ""                                                                                                                                            |
| `guaranteed_msging_max_msg_spool_usage` | The maximum total message spool usage allowed across all VPNs on this broker, in megabytes. Recommendation: the maximum value should be less than 90 percent of the disk space allocated for the guaranteed message spool.                                                                 |
| `service_semp_session_idle_timeout`     | The session idle timeout, in minutes. Sessions will be invalidated if there is no activity in this period of time. Changes to this attribute are synchronized to HA mates via config-sync.                                                                                                    |

#### Solace Broker Backup [`solace_broker_backup`]
Defines they way in which the Solace Broker creates automatic backups.

| Parameter                               |  Comment                                                                                                               |
| --------------------------------------- | ---------------------------------------------------------------------------------------------------------------------- |
| `backup_days_of_week`       | This field is either the entry “daily”, or a list of named days from Sunday to Saturday separated by commas with no spaces or a list of numbers from 0 to 6 representing the named days separated by commas with no spaces, where 0 is Sunday, 1 is Monday, on through to 6 for Saturday. Default is  daily.                                                                                                                                                             |
| `backup_times_of_day`       | This field is either the entry “hourly”, or a list of up to four times of day in the format hh:mm separated by commas without spaces where hh is 0 to 23 representing hours, and mm is 0 to 59 representing minutes.                                                                                    |
| `backup_maximum_backups`    | This field is the maximum number of scheduled backups to keep from 1 to 25. When a new scheduled backup causes the number of backups to exceed the set maximum, the oldest backup file is deleted. The default value is 5 backups if this parameter is not provided.                                       |

#### Solace Broker DMR Cluster [`solace_dmrcluster`]
Defines a local Solace Broker DMR Cluster (this is a needed step to create or join a DMR Mesh).

| Parameter                                   |  Comment                                                                                                               |
| ---------------------------------------     | ---------------------------------------------------------------------------------------------------------------------- |
| `dmr_cluster_authentication_basic_password` | Defines a password for the cluster. This password is needed if other Brokers want to connect to this Broker for a Mesh |
| `dmr_cluster_name`                          | Defines a name for the Cluster.                                                                                        |
| `dmr_cluster_enabled`                       | Enables or disables the cluster through a boolean `true` or `false`                                                    |

#### Solace Syslog Configuration [`solace_syslog_configuration`]
Configures syslog to send events to a defines instance (e.g. logstash).

| Parameter                                   |  Comment                                                                                                               |
| ---------------------------------------     | ---------------------------------------------------------------------------------------------------------------------- |
| `syslog_name`       | Defines the name of the syslog asset. This can be any value but must be unique per Broker.                                                     |
| `syslog_host`       | Defines the remote host where to send the syslog events                                                                                        |
| `syslog_transport`  | Defines the transport to the remote where to send the syslog events - can be either `UDP` or `TCP`                                             |
| `syslog_facilties`  | Defines the facilities for which we want syslog to send files inside a list. List can contain `event`, `command` and `system` or any combination of those values                                                                                                                                                        |


## Resources
This is not an officially supported Solace product.

For more information try these resources:
- Ask the [Solace Community](https://solace.community)
- The Solace Developer Portal website at: https://solace.dev


## Contributing
Contributions are encouraged! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Authors
See the list of [contributors](https://github.com/solacecommunity/<github-repo>/graphs/contributors) who participated in this project.

## License
See the [LICENSE](LICENSE) file for details.
