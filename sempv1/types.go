package sempv1

import (
	"encoding/xml"
	"fmt"
)

type ExecuteResult struct {
	ExecuteResult struct {
		Result     string `xml:"code,attr"`
		Reason     string `xml:"reason,attr"`
		ReasonCode int    `xml:"reasonCode,attr"`
	} `xml:"execute-result"`
	ParseError string `xml:"parse-error"`
}

func (r *ExecuteResult) checkResult() error {
	if r.ExecuteResult.Result == "ok" {
		return nil
	}
	return r
}

func (r *ExecuteResult) Error() string {
	if r.ExecuteResult.Result == "fail" {
		return fmt.Sprintf("command failed: %s (%d)", r.ExecuteResult.Reason, r.ExecuteResult.ReasonCode)
	}
	if r.ParseError != "" {
		return fmt.Sprintf("parseError: %s", r.ParseError)
	}
	return fmt.Sprintf("code %s: %s (%d)", r.ExecuteResult.Result, r.ExecuteResult.Reason, r.ExecuteResult.ReasonCode)
}

/*
No worries, the following structs have not been coded by hand - they have been generated with the help of
- cli2semp
- xmltogo: https://www.onlinetool.io/xmltogo/

Use cli2semp to get a valid rpc structure and call it via Postman or VSCode Rest Client.
Paste the result into xmltogo and you're done.
*/
type BrokerVersionReply struct {
	RPC struct {
		Show struct {
			Version struct {
				Description string `xml:"description"`
				CurrentLoad string `xml:"current-load"`
				Uptime      struct {
					Days      float64 `xml:"days"`
					Hours     float64 `xml:"hours"`
					Mins      float64 `xml:"mins"`
					Secs      float64 `xml:"secs"`
					TotalSecs float64 `xml:"total-secs"`
				} `xml:"uptime"`
			} `xml:"version"`
		} `xml:"show"`
	} `xml:"rpc"`
	ExecuteResult
}

type CliUserReply struct {
	XMLName     xml.Name `xml:"rpc-reply"`
	Text        string   `xml:",chardata"`
	SempVersion string   `xml:"semp-version,attr"`
	Rpc         struct {
		Text string `xml:",chardata"`
		Show struct {
			Text     string `xml:",chardata"`
			Username struct {
				Text  string `xml:",chardata"`
				Users struct {
					Text string `xml:",chardata"`
					User []struct {
						Text                    string `xml:",chardata"`
						Name                    string `xml:"name"`
						UserType                string `xml:"user-type"`
						GlobalAccessLevel       string `xml:"global-access-level"`
						DefaultVpnAccessLevel   string `xml:"default-vpn-access-level"`
						VpnAccessLevelException []struct {
							Text        string `xml:",chardata"`
							VpnName     string `xml:"vpn-name"`
							AccessLevel string `xml:"access-level"`
						} `xml:"vpn-access-level-exception"`
					} `xml:"user"`
				} `xml:"users"`
			} `xml:"username"`
		} `xml:"show"`
	} `xml:"rpc"`
	ExecuteResult
}

// Exceptions returns a list of all vpn exceptions and the access level.
// This will only return the exceptions for the first user in the collection.
func (r *CliUserReply) Exceptions() map[string]string {
	exceptions := make(map[string]string)

	if len(r.Rpc.Show.Username.Users.User) == 1 {
		for _, ex := range r.Rpc.Show.Username.Users.User[0].VpnAccessLevelException {
			exceptions[ex.VpnName] = ex.AccessLevel
		}
	}

	return exceptions
}

type CliAuthenticationReply struct {
	XMLName     xml.Name `xml:"rpc-reply"`
	Text        string   `xml:",chardata"`
	SempVersion string   `xml:"semp-version,attr"`
	Rpc         struct {
		Text string `xml:",chardata"`
		Show struct {
			Text           string `xml:",chardata"`
			Authentication struct {
				Text            string `xml:",chardata"`
				Authentications struct {
					Text           string `xml:",chardata"`
					Authentication struct {
						Text                     string `xml:",chardata"`
						UserClass                string `xml:"user-class"`
						AuthDomain               string `xml:"auth-domain"`
						RadiusDomain             string `xml:"radius-domain"`
						AuthType                 string `xml:"auth-type"`
						ProfileName              string `xml:"profile-name"`
						AccessLevelConfiguration struct {
							Text    string `xml:",chardata"`
							Default struct {
								Text                  string `xml:",chardata"`
								GlobalAccessLevel     string `xml:"global-access-level"`
								DefaultVpnAccessLevel string `xml:"default-vpn-access-level"`
							} `xml:"default"`
							Ldap struct {
								Text                         string `xml:",chardata"`
								GroupMembershipAttributeName string `xml:"group-membership-attribute-name"`
							} `xml:"ldap"`
						} `xml:"access-level-configuration"`
					} `xml:"authentication"`
				} `xml:"authentications"`
			} `xml:"authentication"`
		} `xml:"show"`
	} `xml:"rpc"`
	ExecuteResult
}

type LdapProfileReply struct {
	XMLName     xml.Name `xml:"rpc-reply"`
	Text        string   `xml:",chardata"`
	SempVersion string   `xml:"semp-version,attr"`
	Rpc         struct {
		Text string `xml:",chardata"`
		Show struct {
			Text        string `xml:",chardata"`
			LdapProfile struct {
				Text        string `xml:",chardata"`
				LdapProfile struct {
					Text                          string `xml:",chardata"`
					ProfileName                   string `xml:"profile-name"`
					Shutdown                      string `xml:"shutdown"`
					Tls                           string `xml:"tls"`
					Starttls                      string `xml:"starttls"`
					UnauthenticatedAuthentication string `xml:"unauthenticated-authentication"`
					AdminDn                       string `xml:"admin-dn"`
					Search                        struct {
						Text                         string `xml:",chardata"`
						BaseDn                       string `xml:"base-dn"`
						Filter                       string `xml:"filter"`
						Scope                        string `xml:"scope"`
						Deref                        string `xml:"deref"`
						Timeout                      string `xml:"timeout"`
						FollowContinuationReferences string `xml:"follow-continuation-references"`
					} `xml:"search"`
					GroupMembershipSecondarySearch struct {
						Text                             string `xml:",chardata"`
						Shutdown                         string `xml:"shutdown"`
						BaseDn                           string `xml:"base-dn"`
						FilterAttributeFromPrimarySearch string `xml:"filter-attribute-from-primary-search"`
						Filter                           string `xml:"filter"`
						Scope                            string `xml:"scope"`
						Deref                            string `xml:"deref"`
						Timeout                          int    `xml:"timeout"`
						FollowContinuationReferences     string `xml:"follow-continuation-references"`
					} `xml:"group-membership-secondary-search"`
					LdapServers struct {
						Text       string `xml:",chardata"`
						LdapServer []struct {
							Text                string `xml:",chardata"`
							Index               int    `xml:"index"`
							LdapURI             string `xml:"ldap-uri"`
							NumAdminConnections int    `xml:"num-admin-connections"`
							AdminLastError      string `xml:"admin-last-error"`
							AdminLastErrorTime  string `xml:"admin-last-error-time"`
							NumBindConnections  int    `xml:"num-bind-connections"`
							BindLastError       string `xml:"bind-last-error"`
							BindLastErrorTime   string `xml:"bind-last-error-time"`
						} `xml:"ldap-server"`
					} `xml:"ldap-servers"`
					LdapServersV2 struct {
						Text       string `xml:",chardata"`
						LdapServer []struct {
							Text    string `xml:",chardata"`
							Index   string `xml:"index"`
							LdapURI string `xml:"ldap-uri"`
						} `xml:"ldap-server"`
					} `xml:"ldap-servers-v2"`
					ReferralSession struct {
						Text            string `xml:",chardata"`
						ReferralHostURI string `xml:"referral-host-uri"`
						LastError       string `xml:"last-error"`
						LastErrorTime   string `xml:"last-error-time"`
					} `xml:"referral-session"`
					Users struct {
						Text string `xml:",chardata"`
						User string `xml:"user"`
					} `xml:"users"`
				} `xml:"ldap-profile"`
			} `xml:"ldap-profile"`
		} `xml:"show"`
	} `xml:"rpc"`
	ExecuteResult
}

type LdapCliGroupReply struct {
	XMLName     xml.Name `xml:"rpc-reply"`
	Text        string   `xml:",chardata"`
	SempVersion string   `xml:"semp-version,attr"`
	Rpc         struct {
		Text string `xml:",chardata"`
		Show struct {
			Text           string `xml:",chardata"`
			Authentication struct {
				Text            string `xml:",chardata"`
				Authentications struct {
					Text           string `xml:",chardata"`
					Authentication struct {
						Text                     string `xml:",chardata"`
						UserClass                string `xml:"user-class"`
						AuthDomain               string `xml:"auth-domain"`
						RadiusDomain             string `xml:"radius-domain"`
						AuthType                 string `xml:"auth-type"`
						ProfileName              string `xml:"profile-name"`
						AccessLevelConfiguration struct {
							Text string `xml:",chardata"`
							Ldap struct {
								Text                         string `xml:",chardata"`
								GroupMembershipAttributeName string `xml:"group-membership-attribute-name"`
								Group                        []struct {
									Text                    string `xml:",chardata"`
									Name                    string `xml:"name"`
									GlobalAccessLevel       string `xml:"global-access-level"`
									DefaultVpnAccessLevel   string `xml:"default-vpn-access-level"`
									VpnAccessLevelException []struct {
										Text        string `xml:",chardata"`
										VpnName     string `xml:"vpn-name"`
										AccessLevel string `xml:"access-level"`
									} `xml:"vpn-access-level-exception"`
								} `xml:"group"`
							} `xml:"ldap"`
						} `xml:"access-level-configuration"`
					} `xml:"authentication"`
				} `xml:"authentications"`
			} `xml:"authentication"`
		} `xml:"show"`
	} `xml:"rpc"`
	ExecuteResult
}

// Exceptions returns a list of all vpn exceptions and the access level.
// This will only return the exceptions for the first user in the collection.
func (r *LdapCliGroupReply) LdapCliGroupExceptions() map[string]string {
	exceptions := make(map[string]string)

	if len(r.Rpc.Show.Authentication.Authentications.Authentication.AccessLevelConfiguration.Ldap.Group) == 1 {
		for _, ex := range r.Rpc.Show.Authentication.Authentications.Authentication.AccessLevelConfiguration.Ldap.Group[0].VpnAccessLevelException {
			exceptions[ex.VpnName] = ex.AccessLevel
		}
	}

	return exceptions
}

type BrokerAuthenticationReply struct {
	XMLName     xml.Name `xml:"rpc-reply"`
	Text        string   `xml:",chardata"`
	SempVersion string   `xml:"semp-version,attr"`
	Rpc         struct {
		Text string `xml:",chardata"`
		Show struct {
			Text           string `xml:",chardata"`
			Authentication struct {
				Text            string `xml:",chardata"`
				Authentications struct {
					Text           string `xml:",chardata"`
					Authentication struct {
						Text                     string `xml:",chardata"`
						UserClass                string `xml:"user-class"`
						AuthDomain               string `xml:"auth-domain"`
						RadiusDomain             string `xml:"radius-domain"`
						AuthType                 string `xml:"auth-type"`
						ProfileName              string `xml:"profile-name"`
						AccessLevelConfiguration struct {
							Text    string `xml:",chardata"`
							Default struct {
								Text                  string `xml:",chardata"`
								GlobalAccessLevel     string `xml:"global-access-level"`
								DefaultVpnAccessLevel string `xml:"default-vpn-access-level"`
							} `xml:"default"`
							Ldap struct {
								Text                         string `xml:",chardata"`
								GroupMembershipAttributeName string `xml:"group-membership-attribute-name"`
								Group                        []struct {
									Text                    string `xml:",chardata"`
									Name                    string `xml:"name"`
									GlobalAccessLevel       string `xml:"global-access-level"`
									DefaultVpnAccessLevel   string `xml:"default-vpn-access-level"`
									VpnAccessLevelException []struct {
										Text        string `xml:",chardata"`
										VpnName     string `xml:"vpn-name"`
										AccessLevel string `xml:"access-level"`
									} `xml:"vpn-access-level-exception"`
								} `xml:"group"`
							} `xml:"ldap"`
						} `xml:"access-level-configuration"`
					} `xml:"authentication"`
				} `xml:"authentications"`
			} `xml:"authentication"`
		} `xml:"show"`
	} `xml:"rpc"`
	ExecuteResult
}

type MessageVpnReply struct {
	XMLName     xml.Name `xml:"rpc-reply"`
	Text        string   `xml:",chardata"`
	SempVersion string   `xml:"semp-version,attr"`
	Rpc         struct {
		Text string `xml:",chardata"`
		Show struct {
			Text       string `xml:",chardata"`
			MessageVpn struct {
				Text string `xml:",chardata"`
				Vpn  struct {
					Text                                   string `xml:",chardata"`
					Name                                   string `xml:"name"`
					Alias                                  string `xml:"alias"`
					IsManagementMessageVpn                 string `xml:"is-management-message-vpn"`
					Enabled                                string `xml:"enabled"`
					Operational                            string `xml:"operational"`
					LocallyConfigured                      string `xml:"locally-configured"`
					LocalStatus                            string `xml:"local-status"`
					DistributedCacheManagementEnabled      string `xml:"distributed-cache-management-enabled"`
					SslPlainTextDowngradeAllowed           string `xml:"ssl-plain-text-downgrade-allowed"`
					RestMode                               string `xml:"rest-mode"`
					ServiceRestAuthorizationHeaderHandling string `xml:"service-rest-authorization-header-handling"`
					UniqueSubscriptions                    string `xml:"unique-subscriptions"`
					TotalLocalUniqueSubscriptions          string `xml:"total-local-unique-subscriptions"`
					TotalRemoteUniqueSubscriptions         string `xml:"total-remote-unique-subscriptions"`
					TotalUniqueSubscriptions               string `xml:"total-unique-subscriptions"`
					MaxSubscriptionsEffective              string `xml:"max-subscriptions-effective"`
					MaxSubscriptions                       string `xml:"max-subscriptions"`
					ExportSubscriptions                    string `xml:"export-subscriptions"`
					ExportSubscriptionsPercentComplete     string `xml:"export-subscriptions-percent-complete"`
					PreferIpVersion                        string `xml:"prefer-ip-version"`
					Connections                            string `xml:"connections"`
					ConnectionsServiceSmf                  string `xml:"connections-service-smf"`
					ConnectionsServiceWeb                  string `xml:"connections-service-web"`
					ConnectionsServiceRestIncoming         string `xml:"connections-service-rest-incoming"`
					ConnectionsServiceMqtt                 string `xml:"connections-service-mqtt"`
					ConnectionsServiceAmqp                 string `xml:"connections-service-amqp"`
					ConnectionsServiceRestOutgoing         string `xml:"connections-service-rest-outgoing"`
					MaxConnections                         string `xml:"max-connections"`
					MaxConnectionsServiceSmf               string `xml:"max-connections-service-smf"`
					MaxConnectionsServiceWeb               string `xml:"max-connections-service-web"`
					MaxConnectionsServiceRestIncoming      string `xml:"max-connections-service-rest-incoming"`
					MaxConnectionsServiceMqtt              string `xml:"max-connections-service-mqtt"`
					MaxConnectionsServiceAmqp              string `xml:"max-connections-service-amqp"`
					MaxConnectionsServiceRestOutgoing      string `xml:"max-connections-service-rest-outgoing"`
					AuthType                               string `xml:"auth-type"`
					AuthProfile                            string `xml:"auth-profile"`
					RadiusDomain                           string `xml:"radius-domain"`
					Authentication                         struct {
						Text      string `xml:",chardata"`
						BasicAuth struct {
							Text         string `xml:",chardata"`
							Enabled      string `xml:"enabled"`
							AuthType     string `xml:"auth-type"`
							AuthProfile  string `xml:"auth-profile"`
							RadiusDomain string `xml:"radius-domain"`
						} `xml:"basic-auth"`
						ClientCertAuth struct {
							Text                     string `xml:",chardata"`
							Enabled                  string `xml:"enabled"`
							MaxChainDepth            string `xml:"max-chain-depth"`
							ValidateCertDate         string `xml:"validate-cert-date"`
							AllowApiProvidedUsername string `xml:"allow-api-provided-username"`
							UsernameSource           string `xml:"username-source"`
							RevocationCheckMode      string `xml:"revocation-check-mode"`
						} `xml:"client-cert-auth"`
						KerberosAuth struct {
							Text                     string `xml:",chardata"`
							Enabled                  string `xml:"enabled"`
							AllowApiProvidedUsername string `xml:"allow-api-provided-username"`
						} `xml:"kerberos-auth"`
						Oauth struct {
							Text            string `xml:",chardata"`
							Enabled         string `xml:"enabled"`
							DefaultProvider string `xml:"default-provider"`
						} `xml:"oauth"`
					} `xml:"authentication"`
					MaximumSpoolUsageMb             string `xml:"maximum-spool-usage-mb"`
					MaximumTransactedSessions       string `xml:"maximum-transacted-sessions"`
					MaximumTransactions             string `xml:"maximum-transactions"`
					SempOverMessageBusConfiguration struct {
						Text                      string `xml:",chardata"`
						SempOverMessageBusAllowed string `xml:"semp-over-message-bus-allowed"`
						AdminConfiguration        struct {
							Text                            string `xml:",chardata"`
							AdminCommandsAllowed            string `xml:"admin-commands-allowed"`
							ClientCommandsAllowed           string `xml:"client-commands-allowed"`
							DistributedCacheCommandsAllowed string `xml:"distributed-cache-commands-allowed"`
						} `xml:"admin-configuration"`
						ShowConfiguration struct {
							Text                string `xml:",chardata"`
							ShowCommandsAllowed string `xml:"show-commands-allowed"`
						} `xml:"show-configuration"`
						LegacyConfiguration struct {
							Text                           string `xml:",chardata"`
							LegacyShowClearCommandsAllowed string `xml:"legacy-show-clear-commands-allowed"`
						} `xml:"legacy-configuration"`
					} `xml:"semp-over-message-bus-configuration"`
					EventConfiguration struct {
						Text                  string `xml:",chardata"`
						LargeMessageThreshold string `xml:"large-message-threshold"`
						EventLogTag           string `xml:"event-log-tag"`
						PublishTopicFormat    struct {
							Text string `xml:",chardata"`
							Smf  string `xml:"smf"`
							Mqtt string `xml:"mqtt"`
						} `xml:"publish-topic-format"`
						PublishClientEventMessagesEnabled                                      string `xml:"publish-client-event-messages-enabled"`
						PublishMessageVpnEventMessagesEnabled                                  string `xml:"publish-message-vpn-event-messages-enabled"`
						PublishSubscriptionEventMessagesEnabled                                string `xml:"publish-subscription-event-messages-enabled"`
						PublishSubscriptionEventMessagesEnabledNoUnsubscribeEventsOnDisconnect string `xml:"publish-subscription-event-messages-enabled-no-unsubscribe-events-on-disconnect"`
						EventThresholds                                                        []struct {
							Text            string `xml:",chardata"`
							Name            string `xml:"name"`
							SetPercentage   string `xml:"set-percentage"`
							ClearPercentage string `xml:"clear-percentage"`
							SetValue        string `xml:"set-value"`
							ClearValue      string `xml:"clear-value"`
						} `xml:"event-thresholds"`
					} `xml:"event-configuration"`
				} `xml:"vpn"`
			} `xml:"message-vpn"`
		} `xml:"show"`
	} `xml:"rpc"`
	ExecuteResult
}

type SyslogReply struct {
	XMLName     xml.Name `xml:"rpc-reply"`
	Text        string   `xml:",chardata"`
	SempVersion string   `xml:"semp-version,attr"`
	Rpc         struct {
		Text string `xml:",chardata"`
		Show struct {
			Text   string `xml:",chardata"`
			Syslog struct {
				Text          string `xml:",chardata"`
				SyslogStatus  string `xml:"syslog-status"`
				SyslogElement struct {
					Text       string `xml:",chardata"`
					Name       string `xml:"name"`
					Facilities struct {
						Text     string   `xml:",chardata"`
						Facility []string `xml:"facility"`
					} `xml:"facilities"`
					Files string `xml:"files"`
					Hosts struct {
						Text string `xml:",chardata"`
						Host struct {
							Text      string `xml:",chardata"`
							Address   string `xml:"address"`
							Transport string `xml:"transport"`
						} `xml:"host"`
					} `xml:"hosts"`
				} `xml:"syslog-element"`
			} `xml:"syslog"`
		} `xml:"show"`
	} `xml:"rpc"`
	ExecuteResult
}

type BrokerBackupReply struct {
	XMLName     xml.Name `xml:"rpc-reply"`
	Text        string   `xml:",chardata"`
	SempVersion string   `xml:"semp-version,attr"`
	Rpc         struct {
		Text string `xml:",chardata"`
		Show struct {
			Text   string `xml:",chardata"`
			Backup struct {
				Text                string `xml:",chardata"`
				Schedule            string `xml:"schedule"`
				DaysOfWeek          string `xml:"days-of-week"`
				TimesOfDay          string `xml:"times-of-day"`
				MaxBackups          int    `xml:"max-backups"`
				BackupStatusChanged string `xml:"backup-status-changed"`
			} `xml:"backup"`
		} `xml:"show"`
	} `xml:"rpc"`
	ExecuteResult
}

type LogRetentionReply struct {
	XMLName     xml.Name `xml:"rpc-reply"`
	Text        string   `xml:",chardata"`
	SempVersion string   `xml:"semp-version,attr"`
	Rpc         struct {
		Text string `xml:",chardata"`
		Show struct {
			Text    string `xml:",chardata"`
			Logging struct {
				Text   string `xml:",chardata"`
				Config struct {
					Text                        string `xml:",chardata"`
					MillisecondTimestampEnabled string `xml:"millisecond-timestamp-enabled"`
					Retention                   string `xml:"retention"`
				} `xml:"config"`
			} `xml:"logging"`
		} `xml:"show"`
	} `xml:"rpc"`
	ExecuteResult
}

type MessageSpoolReply struct {
	XMLName     xml.Name `xml:"rpc-reply"`
	Text        string   `xml:",chardata"`
	SempVersion string   `xml:"semp-version,attr"`
	Rpc         struct {
		Text string `xml:",chardata"`
		Show struct {
			Text         string `xml:",chardata"`
			MessageSpool struct {
				Text             string `xml:",chardata"`
				MessageSpoolInfo struct {
					Text                                           string `xml:",chardata"`
					ConfigStatus                                   string `xml:"config-status"`
					CapabilityGuaranteedMsgingSupport              string `xml:"capability-guaranteed-msging-support"`
					MaxDiskUsage                                   string `xml:"max-disk-usage"`
					SpoolWhileCharging                             string `xml:"spool-while-charging"`
					SpoolWithoutFlash                              string `xml:"spool-without-flash"`
					UsingInternalDisk                              string `xml:"using-internal-disk"`
					DiskArrayWwn                                   string `xml:"disk-array-wwn"`
					DiskArrayWwnV2                                 string `xml:"disk-array-wwn-v2"`
					OperationalStatus                              string `xml:"operational-status"`
					DatapathUp                                     string `xml:"datapath-up"`
					SynchronizationStatus                          string `xml:"synchronization-status"`
					SpoolSyncStatus                                string `xml:"spool-sync-status"`
					SpoolSyncLastFailureReason                     string `xml:"spool-sync-last-failure-reason"`
					SpoolSyncLastFailureTime                       string `xml:"spool-sync-last-failure-time"`
					MaxMessageCount                                string `xml:"max-message-count"`
					MaxQueueMessages                               string `xml:"max-queue-messages"`
					MessageCountUtilizationPercentage              string `xml:"message-count-utilization-percentage"`
					TransactionResourceUtilizationPercentage       string `xml:"transaction-resource-utilization-percentage"`
					TransactedSessionResourceUtilizationPercentage string `xml:"transacted-session-resource-utilization-percentage"`
					TransactedSessionCountUtilizationPercentage    string `xml:"transacted-session-count-utilization-percentage"`
					DeliveredUnackedMsgsUtilizationPercentage      string `xml:"delivered-unacked-msgs-utilization-percentage"`
					SpoolFilesUtilizationPercentage                string `xml:"spool-files-utilization-percentage"`
					ActiveDiskPartitionUsage                       string `xml:"active-disk-partition-usage"`
					MateDiskPartitionUsage                         string `xml:"mate-disk-partition-usage"`
					NextMessageID                                  string `xml:"next-message-id"`
					DefragStatusState                              string `xml:"defrag-status-state"`
					DefragStatus                                   string `xml:"defrag-status"`
					DefragScheduleEnabled                          string `xml:"defrag-schedule-enabled"`
					DefragScheduleDays                             string `xml:"defrag-schedule-days"`
					DefragScheduleTimes                            string `xml:"defrag-schedule-times"`
					DefragThresholdEnabled                         string `xml:"defrag-threshold-enabled"`
					DefragThresholdSpoolFragmentationPercentage    string `xml:"defrag-threshold-spool-fragmentation-percentage"`
					DefragThresholdSpoolUsagePercentage            string `xml:"defrag-threshold-spool-usage-percentage"`
					DefragThresholdMinInterval                     string `xml:"defrag-threshold-min-interval"`
					DefragEstFragmentationPercentage               string `xml:"defrag-est-fragmentation-percentage"`
					DefragEstRecoverableSpace                      string `xml:"defrag-est-recoverable-space"`
					NumDeleteInProgress                            string `xml:"num-delete-in-progress"`
					MaxMessageSpoolEntities                        string `xml:"max-message-spool-entities"`
					MessageSpoolEntitiesAllowedByQendpt            string `xml:"message-spool-entities-allowed-by-qendpt"`
					MessageSpoolEntitiesUsedByQueue                string `xml:"message-spool-entities-used-by-queue"`
					MessageSpoolEntitiesUsedByDte                  string `xml:"message-spool-entities-used-by-dte"`
					TransactedSessionsUsed                         string `xml:"transacted-sessions-used"`
					MaxTransactedSessions                          string `xml:"max-transacted-sessions"`
					TransactedSessionsLocalUsed                    string `xml:"transacted-sessions-local-used"`
					TransactedSessionsXaUsed                       string `xml:"transacted-sessions-xa-used"`
					TransactionsUsed                               string `xml:"transactions-used"`
					MaxTransactions                                string `xml:"max-transactions"`
					TransactionsLocalUsed                          string `xml:"transactions-local-used"`
					TransactionsXaUsed                             string `xml:"transactions-xa-used"`
					SequencedTopicsUsed                            string `xml:"sequenced-topics-used"`
					MaxSequencedTopics                             string `xml:"max-sequenced-topics"`
					QueueTopicSubscriptionsUsed                    string `xml:"queue-topic-subscriptions-used"`
					MaxQueueTopicSubscriptions                     string `xml:"max-queue-topic-subscriptions"`
					IngressFlowCount                               string `xml:"ingress-flow-count"`
					IngressFlowsAllowed                            string `xml:"ingress-flows-allowed"`
					FlowsAllowed                                   string `xml:"flows-allowed"`
					ActiveFlowCount                                string `xml:"active-flow-count"`
					InactiveFlowCount                              string `xml:"inactive-flow-count"`
					BrowserFlowCount                               string `xml:"browser-flow-count"`
					MessageReplaysInitializing                     string `xml:"message-replays-initializing"`
					MessageReplaysActive                           string `xml:"message-replays-active"`
					MessageReplaysPendingComplete                  string `xml:"message-replays-pending-complete"`
					MessageReplaysFailed                           string `xml:"message-replays-failed"`
					CvridConfigReady                               string `xml:"cvrid-config-ready"`
					CardReady                                      string `xml:"card-ready"`
					FlashModuleReady                               string `xml:"flash-module-ready"`
					PowerModuleReady                               string `xml:"power-module-ready"`
					CardContentsReady                              string `xml:"card-contents-ready"`
					LocalContentsKey                               string `xml:"local-contents-key"`
					MateContentsKey                                string `xml:"mate-contents-key"`
					RouterSchemaMatch                              string `xml:"router-schema-match"`
					DiskReady                                      string `xml:"disk-ready"`
					DiskContentsStatus                             string `xml:"disk-contents-status"`
					DiskKeyPrimary                                 string `xml:"disk-key-primary"`
					DiskKeyBackup                                  string `xml:"disk-key-backup"`
					LastFailureReason                              string `xml:"last-failure-reason"`
					LastFailureTime                                string `xml:"last-failure-time"`
					CurrentRfadUsage                               string `xml:"current-rfad-usage"`
					CurrentDiskUsage                               string `xml:"current-disk-usage"`
					CurrentPersistUsage                            string `xml:"current-persist-usage"`
					RfadMessagesCurrentlySpooled                   string `xml:"rfad-messages-currently-spooled"`
					DiskMessagesCurrentlySpooled                   string `xml:"disk-messages-currently-spooled"`
					TotalMessagesCurrentlySpooled                  string `xml:"total-messages-currently-spooled"`
					DiskInfos                                      struct {
						Text     string `xml:",chardata"`
						DiskInfo struct {
							Text      string `xml:",chardata"`
							Partition string `xml:"partition"`
							Blocks    string `xml:"blocks"`
							Used      string `xml:"used"`
							Available string `xml:"available"`
							Use       string `xml:"use"`
							MountedOn string `xml:"mounted-on"`
						} `xml:"disk-info"`
					} `xml:"disk-infos"`
					SpoolFiles struct {
						Text               string `xml:",chardata"`
						TotalMaximum       string `xml:"total-maximum"`
						TotalUsed          string `xml:"total-used"`
						TotalAvailable     string `xml:"total-available"`
						TotalUsedPercent   string `xml:"total-used-percent"`
						TotalPendingDelete string `xml:"total-pending-delete"`
					} `xml:"spool-files"`
					SpoolSync struct {
						Text                   string `xml:",chardata"`
						Mode                   string `xml:"mode"`
						AvgMessageAckLatency   string `xml:"avg-message-ack-latency"`
						MaxMessageAckLatency   string `xml:"max-message-ack-latency"`
						MessageAckTimeout      string `xml:"message-ack-timeout"`
						AvgSpoolFileAckLatency string `xml:"avg-spool-file-ack-latency"`
						MaxSpoolFileAckLatency string `xml:"max-spool-file-ack-latency"`
						SpoolFileAckTimeout    string `xml:"spool-file-ack-timeout"`
					} `xml:"spool-sync"`
					TransactionReplicationCompatibilityMode string `xml:"transaction-replication-compatibility-mode"`
					CacheStatus                             string `xml:"cache-status"`
					MaxCacheUsage                           string `xml:"max-cache-usage"`
					CurrentCacheUsage                       string `xml:"current-cache-usage"`
					CacheHighWaterMark                      string `xml:"cache-high-water-mark"`
					EventConfiguration                      struct {
						Text            string `xml:",chardata"`
						EventThresholds []struct {
							Text            string `xml:",chardata"`
							Name            string `xml:"name"`
							SetPercentage   string `xml:"set-percentage"`
							ClearPercentage string `xml:"clear-percentage"`
						} `xml:"event-thresholds"`
					} `xml:"event-configuration"`
				} `xml:"message-spool-info"`
			} `xml:"message-spool"`
		} `xml:"show"`
	} `xml:"rpc"`
	ExecuteResult
}
