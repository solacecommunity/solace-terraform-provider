package provider

const DEFAULT_QUEUE_MAX_MSG_SIZE int = 10 * 1000 * 1000 // Maximum size of messages in Bytes

const DEFAULT_QUEUE_MAX_SPOOL_USAGE int = 4000         // Maximum size of spooled messages to disk in MB
const DEFAULT_QUEUE_MAX_SPOOL_USAGE_SET_PCT int = 25   // The thresholds for the message spool usage alert of the Queue, relative to Messages Queued Quota.
const DEFAULT_QUEUE_MAX_SPOOL_USAGE_CLEAR_PCT int = 18 // The thresholds for the message spool usage alert of the Queue, relative to Messages Queued Quota.

const DEFAULT_MAX_BIND_COUNT int = 1000 // The maximum number of consumer flows that can bind to the Queue. The default value is 1000.
const DEFAULT_MAX_BIND_COUNT_ALERT_SET_PCT int = 25
const DEFAULT_MAX_BIND_COUNT_ALERT_CLEAR_PCT int = 18
const DEFAULT_MAX_REDELIVERY_COUNT int = 0
