package sempv1

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

/* see: https://docs.solace.com/SEMP/SEMP-Error-Codes.htm */
const SEMPv1_EXECUTE_RESULT_OK = "ok"
const (
	SEMP_OK                                      = 0
	SEMP_ERR_FAIL                                = 2
	SEMP_ERR_NOT_INITIALIZED                     = 3
	SEMP_ERR_IN_PROGRESS_TRANSACTION_NOT_ALLOWED = 4
	SEMP_ERR_NOT_ENOUGH_SPACE                    = 5
	SEMP_ERR_NOT_FOUND                           = 6
	SEMP_ERR_COULD_NOT_INITIALIZE_MODES          = 7
	SEMP_ERR_ALREADY_EXISTS                      = 10
	SEMP_ERR_INVALID_PARAMETER                   = 11
	SEMP_ERR_NOT_SUPPORTED                       = 14
	SEMP_ERR_PARSE_ERROR                         = 27
	SEMP_ERR_UNAUTHORIZED                        = 72
	SEMP_ERR_NOT_ALLOWED                         = 89
)

const (
	CLI_AUTH_TYPE_INTERNAL     = "internal"
	CLI_AUTH_TYPE_INTERNAL_VAL = "Internal database authentication" // this is the response from the SEMP api
	CLI_AUTH_TYPE_RADIUS       = "radius"
	CLI_AUTH_TYPE_LDAP         = "ldap"
	CLI_AUTH_TYPE_LDAP_VAL     = "LDAP authentication" // this is the response from the SEMP api
)

const (
	CLI_AUTH_ACCESS_LEVEL_DEFAULT_TYPE_GLOBAL = "global"
	CLI_AUTH_ACCESS_LEVEL_DEFAULT_TYPE_MSGVPN = "message-vpn"
)
const (
	SEMPv1_ACCESS_LEVEL_READ_WRITE = "read-write"
	SEMPv1_ACCESS_LEVEL_READ_ONLY  = "read-only"
	SEMPv1_ACCESS_LEVEL_ADMIN      = "admin"
	SEMPv1_ACCESS_LEVEL_NONE       = "none"
)

type SempV1Client struct {
	Username string
	Passwort string
	Url      string
	Timeout  time.Duration
}

/*
credits to
https://github.com/solacecommunity/solace-prometheus-exporter/blob/master/solace_prometheus_exporter.go
*/

// Call HTTP POST with the given command as body at the SEMPv1 legacy endpoint and return the result as io stream.
func (c *SempV1Client) callSempRpc(ctx context.Context, body string) (io.ReadCloser, error) {
	transport := &http.Transport{}

	req, err := http.NewRequestWithContext(ctx, "POST", c.Url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	// ProxyFromEnvironment returns the URL of the proxy to use for a given request, as indicated by the environment
	// variables HTTP_PROXY, HTTPS_PROXY and NO_PROXY
	proxyURL, err := http.ProxyFromEnvironment(req)
	if err != nil {
		return nil, err
	}
	if proxyURL != nil {
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	client := http.Client{
		Timeout:   c.Timeout,
		Transport: transport,
	}

	req.SetBasicAuth(c.Username, c.Passwort)
	req.Header.Set("Content-Type", "application/xml")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("HTTP status %d (%s)", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	return resp.Body, nil
}

func (c *SempV1Client) executeEmptyCommand(ctx context.Context, command string) error {
	var result = ExecuteResult{}
	err := c.executeCommand(ctx, command, &result)
	if err != nil {
		return err
	}
	return result.checkResult()
}

// Execute a SEMPv1 legacy command and parse it into the given result struct.
// Make sure to put a reference (&struct) into `result`
func (c *SempV1Client) executeCommand(ctx context.Context, command string, result interface{}) error {
	body, err := c.callSempRpc(ctx, command)
	if err != nil {
		return err
	}
	defer body.Close()
	decoder := xml.NewDecoder(body)
	err = decoder.Decode(result)
	if err != nil {
		return err
	}

	return nil
}
