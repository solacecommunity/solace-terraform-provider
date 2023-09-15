package resources

// Unique MsgVpn fields
const JNDI_ENABLED string = "jndi_enabled"
const MAX_CONNECTION_COUNT string = "max_connection_count"
const MAX_MSG_SPOOL_USAGE string = "max_msg_spool_usage"
const MSG_VPN_ENABLED string = "msg_vpn_enabled"
const MSG_VPN_DMR_ENABLED string = "msg_vpn_dmr_enabled"
const AUTHENTICATION_BASIC_ENABLED string = "authentication_basic_enabled"
const AUTHENTICATION_BASIC_PROFILE_NAME string = "authentication_basic_profile_name"
const AUTHENTICATION_BASIC_TYPE string = "authentication_basic_type"
const SERVICE_MQTT_PLAIN_TEXT_ENABLED string = "service_mqtt_plain_text_enabled"
const SERVICE_AMQP_PLAIN_TEXT_ENABLED string = "service_amqp_plain_text_enabled"
const SERVICE_MQTT_WEB_SOCKET_ENABLED string = "service_mqtt_web_socket_enabled"
const SERVICE_REST_INCOMING_PLAIN_TEXT_ENABLED string = "service_rest_incoming_plain_text_enabled"
const SERVICE_SMF_PLAIN_TEXT_ENABLED string = "service_smf_plain_text_enabled"
const SERVICE_WEB_PLAIN_TEXT_ENABLED string = "service_web_plain_text_enabled"
const SERVICE_MQTT_TLS_LISTEN_PORT string = "service_mqtt_tls_listen_port"
const SERVICE_MQTT_TLS_ENABLED string = "service_mqtt_tls_enabled"
const SERVICE_AMQP_TLS_ENABLED string = "service_amqp_tls_enabled"
const SERVICE_AMQP_TLS_LISTEN_PORT string = "service_amqp_tls_listen_port"
const EVENT_LARGE_MSG_THRESHOLD string = "event_large_msg_threshold"
const SERVICE_REST_TLS_ENABLED string = "service_rest_tls_enabled"
const SERVICE_REST_TLS_LISTEN_PORT string = "service_rest_tls_listen_port"

// Unique Queue Fields
const QUEUE_NAME string = "queue_name"
const TOPIC_NAME = "topic_name"
const MSG_VPN_NAME string = "msg_vpn_name"
const SUBSCRIPTION_TOPIC string = "subscription_topic"
const JNDI_NAME string = "jndi_name"
const QUEUE_TEMPLATE_NAME string = "queue_template_name"

// CLI User
const CLI_USERNAME string = "username"
const CLI_PASSWORD string = "password"
const CLI_USERTYPE string = "user_type"
const CLI_GLOBAL_ACCESS_LEVEL string = "global_access_level"
const CLI_MSGVPN_DEFAULT_ACCESS_LEVEL string = "msgvpn_default_access_level"
const CLI_MSGVPN_ACCESS_LEVEL_EXCEPTION string = "msgvpn_access_level_exceptions"
const CLI_USER_AUTH_TYPE string = "cli_auth_type"
const CLI_USER_AUTH_TYPE_PROFILE string = "cli_auth_type_profile"
const CLI_USER_DEFAULT_GLOBAL_ACCESS_LEVEL string = "default_global_access_level"
const CLI_USER_DEFAULT_MESSAGEVPN_ACCESS_LEVEL string = "default_messagevpn_access_level"

// LDAP Profile
const LDAP_PROFILE_NAME string = "profile_name"
const LDAP_PROFILE_ENABLED string = "enabled"
const LDAP_PROFILE_HOST string = "host"
const LDAP_PROFILE_ADMIN_DN string = "admin_dn"
const LDAP_PROFILE_ADMIN_PWD string = "admin_password"
const LDAP_PROFILE_INDEX string = "index"
const LDAP_PROFILE_BASE_DN string = "base_dn"
const LDAP_PROFILE_SEARCH_FILER string = "search_filter"

const ACCESS_TYPE string = "access_type"
const PERMISSION string = "permission"
const CONSUMER_ACK_PROPAGATION_ENABLED string = "consumer_ack_propagation_enabled"
const DEAD_MSG_QUEUE string = "dead_msg_queue"
const EGRESS_ENABLED string = "egress_enabled"
const INGRESS_ENABLED string = "ingress_enabled"
const MAX_BIND_COUNT string = "max_bind_count"
const MAX_BIND_COUNT_ALERT_SET_PCT string = "max_bind_count_alert_set_pct"
const MAX_BIND_COUNT_ALERT_CLEAR_PCT string = "max_bind_count_alert_clear_pct"
const MAX_DELIVERED_UNACKED_MSGS_PER_FLOW string = "max_delivered_unacked_msgs_per_flow"
const MAX_MSG_SIZE string = "max_msg_size"
const MAX_SPOOL_USAGE string = "max_spool_usage"
const MAX_SPOOL_USAGE_ALERT_SET_PCT string = "max_spool_usage_alert_set_pct"
const MAX_SPOOL_USAGE_ALERT_CLEAR_PCT string = "max_spool_usage_alert_clear_pct"
const MAX_REDELIVERY_COUNT string = "max_redelivery_count"
const MAX_TTL string = "max_ttl"
const OWNER string = "owner"
const REJECT_LOW_PRIORITY_MSG_ENABLED string = "reject_low_priority_msg_enabled"
const REJECT_LOW_PRIORITY_MSG_LIMIT string = "reject_low_priority_msg_limit"
const REJECT_MSG_TO_SENDER_ON_DISCARD_BEHAVIOR string = "reject_msg_to_sender_on_discard_behavior"
const RESPECT_TTL_ENABLED string = "respect_ttl_enabled"
const TOPIC_SUBSCRIPTION_LIST string = "topic_subscription_list"

const CLIENT_USERNAME string = "client_username"
const ENABLED string = "enabled"
const ACL_PROFILE_NAME string = "acl_profile_name"
const PASSWORD string = "password"
const SUBSCRIPTION_MANAGER_ENABLED string = "subscription_manager_enabled"

// Client Profiles
const CLIENT_PROFILE_NAME string = "client_profile_name"
const ALLOW_GUARANTEED_MSG_RECEIVE string = "allow_guaranteed_msg_receive"
const ALLOW_GUARANTEED_MSG_SEND string = "allow_guaranteed_msg_send"
const ALLOW_DOWNGRADE_TO_PLAIN_TEXT string = "allow_downgrade_to_plain_text"
const ALLOW_TRANSACTED_SESSIONS_ENABLED string = "allow_transacted_sessions_enabled"
const ALLOW_BRIDGE_CONNECTIONS string = "allow_bridge_connections"
const ALLOW_CUT_THROUGH_FORWARDING string = "allow_cut_through_forwarding"
const ALLOW_GUARANTEED_ENDPOINT_CREATE string = "allow_guaranteed_endpoint_create"
const ALLOW_GUARANTEED_ENDPOINT_CREATE_DURABILITY string = "allow_guaranteed_endpoint_create_durability"
const ALLOW_SHARED_SUBSCRIPTIONS string = "allow_shared_subscriptions"
const COMPRESSION_ENABLED string = "compression_enabled"
const SERVICE_SMF_MAX_CONNECTION_COUNT_PER_CLIENT_USERNAME string = "service_smf_max_connection_count_per_client_username"
const MAX_TRANSACTED_SESSION_COUNT string = "max_transacted_session_count"

// ACL Profiles
const CLIENT_CONNECT_DEFAULT_ACTION string = "client_connect_default_action"
const PUBLISH_TOPIC_DEFAULT_ACTION string = "publish_topic_default_action"
const SUBSCRIBE_TOPIC_DEFAULT_ACTION string = "subscribe_topic_default_action"
const SUBSCRIBE_SHARE_NAME_DEFAULT_ACTION string = "subscribe_share_name_default_action"
const ACL_CONNECT_EXCEPTION_LIST string = "connect_exception_list"
const ACL_PUBLISH_EXCEPTION_LIST string = "publish_exception_list"
const ACL_SUBSCRIBE_EXCEPTION_LIST string = "subscribe_exception_list"
const ACL_SUBSCRIBE_SHARE_NAME_EXCEPTION_LIST string = "subscribe_share_name_exception_list"

// JNDI Connection Factories
const CONNECTION_FACTORY_NAME string = "connection_factory_name"
const TRANSPORT_CONNECT_TIMEOUT string = "transport_connect_timeout"
const TRANSPORT_REPLY_TIMEOUT string = "transport_reply_timeout"
const TRANSPORT_CONNECT_RETRY_COUNT string = "transport_connect_retry_count"
const TRANSPORT_RECONNECT_RETRY_COUNT string = "transport_reconnect_retry_count"
const TRANSPORT_RECONNECT_RETRY_WAIT string = "transport_reconnect_retry_wait"
const TRANSPORT_CONNECT_RETRY_PER_HOST_COUNT string = "transport_connect_retry_per_host_count"
const DYNAMIC_CREATE_DURABLE_ENDPOINT string = "dynamic_create_durable_endpoint"
const XA_ENABLED string = "xa_enabled"
const ALLOW_DUPLICATE_CLIENT_ID_ENABLED string = "allow_duplicate_client_id_enabled"
const DIRECT_TRANSPORT_ENABLED string = "direct_transport_enabled"
const DYNAMIC_ENDPOINT_RESPECT_TTL_ENABLED string = "dynamic_endpoint_respect_ttl_enabled"
const TRANSPORT_KEEPALIVE_INTERVAL string = "transport_keepalive_interval"
const DMQ_ELIGIBLE_ENABLED string = "dmq_eligible_enabled"

// Client Cert Authority
const CERT_AUTHORITY_NAME string = "cert_authority_name"
const CERT_CONTENT string = "cert_content"

// LDAP Cli Groups
const LDAP_CLI_GROUPNAME string = "group_name"
const LDAP_CLI_GLOBAL_ACCESS_LEVEL string = "global_access_level"
const LDAP_CLI_MSGVPN_DEFAULT_ACCESS_LEVEL string = "msgvpn_default_access_level"
const LDAP_CLI_MSGVPN_ACCESS_LEVEL_EXCEPTION string = "msgvpn_access_level_exceptions"

// Broker Settings
const TLS_SSH_CIPHER_SUITE_LIST string = "tls_ssh_cipher_suite_list"
const TLS_MSG_BACKBONE_CIPHER_SUITE_LIST string = "tls_msg_backbone_cipher_suite_list"
const TLS_MANAGEMENT_CIPHER_SUITE_LIST string = "tls_management_cipher_suite_list"
const TLS_SERVER_CERTIFICATE string = "tls_server_certificate"
const TLS_SERVER_CERTIFICATE_PASSWORD string = "tls_server_certificate_password"
const GUARANTEED_MSGING_MAX_MSG_SPOOL_USAGE string = "guaranteed_msging_max_msg_spool_usage"
const SERVICE_SEMP_SESSION_IDLE_TIMEOUT string = "service_semp_session_idle_timeout"
const LOG_RETENTION_DURATION string = "log_retention_duration"

// Syslog Settings
const SYSLOG_NAME string = "syslog_name"
const SYSLOG_HOST string = "syslog_host"
const SYSLOG_TRANSPORT string = "syslog_transport"
const SYSLOG_FACILITIES string = "syslog_facilties"

// Fragemntation Settings
const FRAGMENTATION_SCHEDULED_ENABLED string = "fragmentation_scheduled_enabled"
const FRAGMENTATION_SCHEDULED_DAYS string = "fragmentation_scheduled_days"
const FRAGMENTATION_SCHEDULED_TIMES string = "fragmentation_scheduled_times"
const FRAGMENTATION_THRESHOLD_ENABLED string = "fragmentation_threshold_enabled"
const FRAGMENTATION_THRESHOLD_FRAGMENTATION_PERCENTAGE string = "fragmentation_threshold_fragmentation_percentage"
const FRAGMENTATION_THRESHOLD_USAGE_PERCENTAGE string = "fragmentation_threshold_usage_percentage"
const FRAGMENTATION_THRESHOLD_MIN_INTERVAL string = "fragmentation_threshold_min_interval"

// Broker Authentication
const LDAP_ACCESS_LEVEL_GROUP_MEMBERSHIP_ATTRIBUTE_NAME string = "ldap_access_level_group_membership_attribute_name"

// DMR Cluster
const DMR_CLUSTER_AUTHENTICATION_BASIC_PASSWORD string = "dmr_cluster_authentication_basic_password"
const DMR_CLUSTER_NAME string = "dmr_cluster_name"
const DMR_CLUSTER_ENABLED string = "dmr_cluster_enabled"
const DMR_NODE_NAME string = "dmr_node_name"

// Backups
const BACKUP_DAYS_OF_WEEK string = "backup_days_of_week"
const BACKUP_TIMES_OF_DAY string = "backup_times_of_day"
const BACKUP_MAXIMUM_BACKUPS string = "backup_maximum_backups"

// External Link
const EXTERNAL_LINK_DMR_CLUSTER_NAME string = "external_link_dmr_cluster_name"
const EXTERNAL_LINK_AUTHENTICATION_BASIC_PASSWORD string = "external_link_authentication_basic_password"
const EXTERNAL_LINK_REMOTE_NODE_NAME string = "external_link_remote_node_name"
const EXTERNAL_LINK_INITIATOR string = "external_link_initiator"
const EXTERNAL_LINK_ENABLED string = "external_link_enabled"
const EXTERNAL_LINK_SPAN string = "external_link_span"
const EXTERNAL_LINK_TRANSPORT_TLS_ENABLED string = "external_link_transport_tls_enabled"

// Remote Address
const REMOTE_ADDRESS_DMR_CLUSTER_NAME string = "remote_address_dmr_cluster_name"
const REMOTE_ADDRESS_REMOTE_NODE_NAME string = "remote_address_remote_node_name"
const REMOTE_ADDRESS string = "remote_address"

// Bridge
const BRIDGE_NAME string = "bridge_name"
const BRIDGE_MESSAGE_VPN string = "bridge_message_vpn"
const BRIDGE_ENABLED string = "bridge_enabled"
const BRIDGE_VIRTUAL_ROUTER string = "bridge_virtual_router"
const BRIDGE_MAX_TTL string = "bridge_max_ttl"
const BRIDGE_REMOTE_AUTHENTICATION_BASIC_CLIENT_USERNAME string = "bridge_remote_authentication_basic_client_username"
const BRIDGE_REMOTE_AUTHENTICATION_BASIC_PASSWORD string = "bridge_remote_authentication_basic_password"
const BRIDGE_REMOTE_AUTHENTICATION_SCHEME string = "bridge_remote_authentication_scheme"
const BRIDGE_REMOTE_CONNECTION_RETRY_COUNT string = "bridge_remote_connection_retry_count"
const BRIDGE_REMOTE_CONNECTION_RETRY_DELAY string = "bridge_remote_connection_retry_delay"
const BRIDGE_REMOTE_DELIVER_TO_ONE_PRIORITY string = "bridge_remote_deliver_to_one_priority"
const BRIDGE_TLS_CIPHER_SUITE_LIST string = "bridge_tls_cipher_suite_list"

// Bridge Remote MsgVPN
const BRIDGE_REMOTE_MSGVPN_BRIDGE_NAME string = "bridge_remote_msgvpn_bridge_name"
const BRIDGE_REMOTE_MSGVPN_MESSAGE_VPN string = "bridge_remote_msgvpn_message_vpn"
const BRIDGE_REMOTE_MSGVPN_CLIENT_USERNAME string = "bridge_remote_msgvpn_client_username"
const BRIDGE_REMOTE_MSGVPN_BRIDGE_VIRTUAL_ROUTER string = "bridge_remote_msgvpn_bridge_virtual_router"
const BRIDGE_REMOTE_MSGVPN_COMPRESSED_DATA_ENABLED string = "bridge_remote_msgvpn_compressed_data_enabled"
const BRIDGE_REMOTE_MSGVPN_CONNECTOR_ORDER string = "bridge_remote_msgvpn_connector_order"
const BRIDGE_REMOTE_MSGVPN_EGRESS_FLOW_WINDOW_SIZE string = "bridge_remote_msgvpn_egress_flow_window_size"
const BRIDGE_REMOTE_MSGVPN_ENABLED string = "bridge_remote_msgvpn_enabled"
const BRIDGE_REMOTE_MSGVPN_PASSWORD string = "bridge_remote_msgvpn_password"
const BRIDGE_REMOTE_MSGVPN_QUEUE_BINDING string = "bridge_remote_msgvpn_queue_binding"
const BRIDGE_REMOTE_MSGVPN_REMOTE_MSGVPN_INTERFACE string = "bridge_remote_msgvpn_remote_msgvpn_interface"
const BRIDGE_REMOTE_MSGVPN_REMOTE_MSGVPN_LOCATION string = "bridge_remote_msgvpn_remote_msgvpn_location"
const BRIDGE_REMOTE_MSGVPN_REMOTE_MSGVPN_NAME string = "bridge_remote_msgvpn_remote_msgvpn_name"
const BRIDGE_REMOTE_MSGVPN_TLS_ENABLED string = "bridge_remote_msgvpn_tls_enabled"
const BRIDGE_REMOTE_MSGVPN_UNIDIRECTIONAL_CLIENT_PROFILE string = "bridge_remote_msgvpn_unidirectional_client_profile"

// Bridge Remote Subscription
const BRIDGE_REMOTE_SUBSCRIPTION_BRIDGE_VIRTUAL_ROUTER string = "bridge_remote_subscription_bridge_virtual_router"
const BRIDGE_REMOTE_SUBSCRIPTION_BRIDGE_NAME string = "bridge_remote_subscription_bridge_name"
const BRIDGE_REMOTE_SUBSCRIPTION_MESSAGE_VPN string = "bridge_remote_subscription_message_vpn"
const BRIDGE_REMOTE_SUBSCRIPTION_DELIVERY_ALWAYS_ENABLED string = "bridge_remote_subscription_delivery_always_enabled"
const BRIDGE_REMOTE_SUBSCRIPTION_REMOTE_SUBSCRIPTION_TOPIC string = "bridge_remote_subscription_remote_subscription_topic"

// DMR Bridge
const DMR_BRIDGE_MESSAGE_VPN string = "dmr_bridge_message_vpn"
const DMR_BRIDGE_REMOTE_MESSAGE_VPN string = "dmr_bridge_remote_message_vpn"
const DMR_BRIDGE_REMOTE_NODE_NAME string = "dmr_bridge_remote_node_name"

// RDP
const REST_DELIVERY_POINT_NAME string = "rest_delivery_point_name"
const REST_DELIVERY_QUEUE_BINDING string = "rest_delivery_queue_binding"
const REST_DELIVERY_QUEUE_BINDING_REQUEST_HEADER string = "rest_delivery_queue_binding_request_header"

// RDP QUEUE BINDING
const REST_DELIVERY_POINT_QUEUE_BINDING_NAME string = "rest_delivery_point_queue_binding_name"
const REST_DELIVERY_POINT_POST_REQUEST_TARGET string = "rest_delivery_point_post_request_target"

// RDP REST CONSUMER
const REST_DELIVERY_POINT_REST_CONSUMER_HTTP_METHOD string = "rest_delivery_point_rest_consumer_http_method"
const REST_DELIVERY_POINT_REST_CONSUMER_MAX_POST_WAIT_TIME string = "rest_delivery_point_rest_consumer_max_post_wait_time"
const REST_DELIVERY_POINT_REST_CONSUMER_OUTGOING_CONNECTION_COUNT string = "rest_delivery_point_rest_consumer_outgoing_connection_count"
const REST_DELIVERY_POINT_REST_CONSUMER_NAME string = "rest_delivery_point_rest_consumer_name"
const REST_DELIVERY_POINT_REST_CONSUMER_RETRY_DELAY string = "rest_delivery_point_rest_consumer_retry_delay"
const REST_DELIVERY_POINT_REST_CONSUMER_REMOTE_HOST string = "rest_delivery_point_rest_consumer_remote_host"
const REST_DELIVERY_POINT_REST_CONSUMER_REMOTE_PORT string = "rest_delivery_point_rest_consumer_remote_port"
const REST_DELIVERY_POINT_REST_CONSUMER_AUTHENTICATION_SCHEME string = "rest_delivery_point_rest_consumer_authentication_scheme"
const REST_DELIVERY_POINT_REST_CONSUMER_HTTP_BASIC_USERNAME string = "rest_delivery_point_rest_consumer_http_basic_username"
const REST_DELIVERY_POINT_REST_CONSUMER_HTTP_BASIC_PASSWORD string = "rest_delivery_point_rest_consumer_http_basic_password"
const REST_DELIVERY_POINT_REST_CONSUMER_TLS_ENABLED string = "rest_delivery_point_rest_consumer_tls_enabled"
