package sempv1

import (
	"context"
	"fmt"
	"strconv"
)

func (c *SempV1Client) GetBrokerVersion(ctx context.Context) (*BrokerVersionReply, error) {
	command := "<rpc><show><version/></show></rpc>"
	var result = BrokerVersionReply{}

	err := c.executeCommand(ctx, command, &result)
	if err != nil {
		return nil, err
	}
	err = result.checkResult()
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *SempV1Client) SetCliUserAuthType(ctx context.Context, authType string, profileName string) error {
	command_s := "<rpc><authentication><auth-type>"
	command_e := "</auth-type></authentication></rpc>"
	command := ""
	if authType == CLI_AUTH_TYPE_INTERNAL {
		command = command_s + "<internal />" + command_e
	} else if authType == CLI_AUTH_TYPE_LDAP {
		command = command_s + fmt.Sprintf("<ldap></ldap><ldap-profile>%s</ldap-profile>", profileName) + command_e
	} else if authType == CLI_AUTH_TYPE_RADIUS {
		command = command_s + fmt.Sprintf("<radius></radius><radius-profile>%s</radius-profile>", profileName) + command_e
	} else {
		return fmt.Errorf("%s is not a valid authentication type", authType)
	}

	return c.executeEmptyCommand(ctx, command)
}

// Create a local CLI User with Username and Password.
func (c *SempV1Client) CreateCliUser(ctx context.Context, username string, password string) error {
	command := fmt.Sprintf("<rpc><create><username><name>%s</name><password>%s</password></username></create></rpc>", username, password)
	err := c.executeEmptyCommand(ctx, command)
	if err != nil {
		if result, ok := err.(*ExecuteResult); ok {
			if result.ExecuteResult.ReasonCode == SEMP_ERR_ALREADY_EXISTS {
				return nil
			}
		}
	}
	return err
}

func (c *SempV1Client) GetCliUser(ctx context.Context, username string) (*CliUserReply, error) {
	command := fmt.Sprintf("<rpc><show><username><username-pattern>%s</username-pattern><detail></detail></username></show></rpc>", username)
	var result = CliUserReply{}

	err := c.executeCommand(ctx, command, &result)
	if err != nil {
		return nil, err
	}
	err = result.checkResult()
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete a local CLI user.
func (c *SempV1Client) DeleteCliUser(ctx context.Context, username string) error {
	command := fmt.Sprintf("<rpc><no><username><name>%s</name></username></no></rpc>", username)
	return c.executeEmptyCommand(ctx, command)
}

// Change the password for an existing cli-user
func (c *SempV1Client) SetCliUserPassword(ctx context.Context, username string, password string) error {
	command := fmt.Sprintf("<rpc><username><name>%s</name><change-password><password>%s</password></change-password></username></rpc>", username, password)
	return c.executeEmptyCommand(ctx, command)
}

// Set the global access level for a Broker for CLI Users.
// For access, you can use one of the predefined contants, like
// sempv1.SEMPv1_ACCESS_LEVEL_NONE or SEMPv1_ACCESS_LEVEL_READ_WRITE
func (c *SempV1Client) SetCliUserGlobalAccessLevel(ctx context.Context, username string, access string) error {
	command := fmt.Sprintf("<rpc><username><name>%s</name><global-access-level><access-level>%s</access-level></global-access-level></username></rpc>", username, access)
	return c.executeEmptyCommand(ctx, command)
}

// Set the Default Access Level for all Message VPNs for CLI Users.
// For access, you can use one of the predefined contants, like
// sempv1.SEMPv1_ACCESS_LEVEL_NONE or SEMPv1_ACCESS_LEVEL_READ_WRITE
func (c *SempV1Client) SetCliUserMsgVpnDefaultAccessLevel(ctx context.Context, username string, access string) error {
	command := fmt.Sprintf("<rpc><username><name>%s</name><message-vpn><default-access-level><access-level>%s</access-level></default-access-level></message-vpn></username></rpc>", username, access)
	return c.executeEmptyCommand(ctx, command)
}

func (c *SempV1Client) CreateCliUserVpnAccessLevelException(ctx context.Context, username string, msgVpn string, accessLevel string) error {
	command := fmt.Sprintf("<rpc><username><name>%s</name><message-vpn><create><access-level-exception><vpn-name>%s</vpn-name></access-level-exception></create></message-vpn></username></rpc>", username, msgVpn)
	err := c.executeEmptyCommand(ctx, command)
	if err != nil {
		if e, ok := err.(*ExecuteResult); ok {
			if e.ExecuteResult.ReasonCode != SEMP_ERR_ALREADY_EXISTS {
				return err
			}
		} else {
			return err
		}
	}
	// conveniently call the set function with the expected value
	if accessLevel != SEMPv1_ACCESS_LEVEL_NONE {
		return c.SetCliUserVpnAccessLevelException(ctx, username, msgVpn, accessLevel)
	}
	return nil
}

func (c *SempV1Client) DeleteCliUserVpnAccessLevelException(ctx context.Context, username string, msgVpn string) error {
	command := fmt.Sprintf("<rpc><username><name>%s</name><message-vpn><no><access-level-exception><vpn-name>%s</vpn-name></access-level-exception></no></message-vpn></username></rpc>", username, msgVpn)
	return c.executeEmptyCommand(ctx, command)
}

func (c *SempV1Client) SetCliUserVpnAccessLevelException(ctx context.Context, username string, msgVpn string, accessLevel string) error {
	command := fmt.Sprintf("<rpc><username><name>%s</name><message-vpn><access-level-exception><vpn-name>%s</vpn-name>"+
		"<access-level><access-level>%s</access-level></access-level></access-level-exception></message-vpn>"+
		"</username></rpc>", username, msgVpn, accessLevel)
	return c.executeEmptyCommand(ctx, command)
}

func (c *SempV1Client) CreateLdapProfile(ctx context.Context, profileName string) error {
	command := fmt.Sprintf("<rpc><authentication><create><ldap-profile><profile-name>%s</profile-name></ldap-profile></create></authentication></rpc>", profileName)
	return c.executeEmptyCommand(ctx, command)
}

func (c *SempV1Client) DeleteLdapProfile(ctx context.Context, profileName string) error {
	command := fmt.Sprintf("<rpc><authentication><no><ldap-profile><profile-name>%s</profile-name></ldap-profile></no></authentication></rpc>", profileName)
	return c.executeEmptyCommand(ctx, command)
}

func (c *SempV1Client) SetLdapProfileAdminDN(ctx context.Context, profileName string, adminDN string, adminPassword string) error {
	command := fmt.Sprintf("<rpc><authentication><ldap-profile><profile-name>%s</profile-name>"+
		"<admin><admin-dn>%s</admin-dn><admin-password>%s</admin-password></admin>"+
		"</ldap-profile></authentication></rpc>", profileName, adminDN, adminPassword)
	return c.executeEmptyCommand(ctx, command)
}

func (c *SempV1Client) SetLdapProfileServer(ctx context.Context, profileName string, ldapHost string, ldapIndex int) error {
	command := fmt.Sprintf("<rpc><authentication><ldap-profile><profile-name>%s</profile-name><ldap-server>"+
		"<ldap-host>%s</ldap-host><server-index>%d</server-index></ldap-server></ldap-profile></authentication></rpc>", profileName, ldapHost, ldapIndex)
	return c.executeEmptyCommand(ctx, command)
}

func (c *SempV1Client) SetLdapProfileSeachBaseDN(ctx context.Context, profileName string, baseDN string) error {
	command := fmt.Sprintf("<rpc><authentication><ldap-profile><profile-name>%s</profile-name><search><base-dn>"+
		"<distinguished-name>%s</distinguished-name></base-dn></search></ldap-profile></authentication></rpc>", profileName, baseDN)
	return c.executeEmptyCommand(ctx, command)
}

func (c *SempV1Client) SetLdapProfileSeachFilter(ctx context.Context, profileName string, searchFilter string) error {
	command := fmt.Sprintf("<rpc><authentication><ldap-profile><profile-name>%s</profile-name><search><filter>"+
		"<filter>%s</filter></filter></search></ldap-profile></authentication></rpc>", profileName, searchFilter)
	return c.executeEmptyCommand(ctx, command)
}

func (c *SempV1Client) EnableLdapProfile(ctx context.Context, profileName string, enable bool) error {
	command := fmt.Sprintf("<rpc><authentication><ldap-profile><profile-name>%s</profile-name>", profileName)
	if enable {
		command += "<no><shutdown></shutdown></no>"
	} else {
		command += "<shutdown></shutdown>"
	}
	command += "</ldap-profile></authentication></rpc>"
	return c.executeEmptyCommand(ctx, command)
}

func (c *SempV1Client) GetLdapProfile(ctx context.Context, profileName string) (*LdapProfileReply, error) {
	command := fmt.Sprintf("<rpc><show><ldap-profile><profile-name>%s</profile-name><detail></detail></ldap-profile></show></rpc>", profileName)
	var result = LdapProfileReply{}

	err := c.executeCommand(ctx, command, &result)
	if err != nil {
		return nil, err
	}
	err = result.checkResult()
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *SempV1Client) GetLdapCliGroup(ctx context.Context, groupName string) (*LdapCliGroupReply, error) {
	command := fmt.Sprintf("<rpc><show><authentication><access-level></access-level><ldap></ldap><group></group>"+
		"<group-name-pattern>%s</group-name-pattern><detail></detail></authentication></show></rpc>", groupName)
	var result = LdapCliGroupReply{}

	err := c.executeCommand(ctx, command, &result)
	if err != nil {
		return nil, err
	}
	err = result.checkResult()
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *SempV1Client) CreateLdapCliGroup(ctx context.Context, groupName string) error {
	command := fmt.Sprintf("<rpc><authentication><access-level><ldap><create><group><group-name>%s</group-name></group></create>"+
		"</ldap></access-level></authentication></rpc>", groupName)
	err := c.executeEmptyCommand(ctx, command)
	if err != nil {
		if e, ok := err.(*ExecuteResult); ok {
			if e.ExecuteResult.ReasonCode != SEMP_ERR_ALREADY_EXISTS {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func (c *SempV1Client) DeleteLdapCliGroup(ctx context.Context, groupName string) error {
	command := fmt.Sprintf("<rpc><authentication><access-level><ldap><no><group><group-name>%s</group-name>"+
		"</group></no></ldap></access-level></authentication></rpc>", groupName)
	return c.executeEmptyCommand(ctx, command)
}

func (c *SempV1Client) SetLdapCliGroupGlobalAccessLevel(ctx context.Context, groupName string, globalAccessLevel string) error {
	command := fmt.Sprintf("<rpc><authentication><access-level><ldap><group><group-name>%s</group-name>"+
		"<global-access-level><access-level>%s</access-level></global-access-level></group></ldap></access-level></authentication></rpc>", groupName, globalAccessLevel)
	err := c.executeEmptyCommand(ctx, command)
	if err != nil {
		return err
	}
	return nil
}

func (c *SempV1Client) SetLdapCliGroupMessageVpnDefaultAccessLevel(ctx context.Context, groupName string, msgVpnDefaultAccessLevel string) error {
	command := fmt.Sprintf("<rpc><authentication><access-level><ldap><group><group-name>%s</group-name>"+
		"<message-vpn><default-access-level><access-level>%s</access-level></default-access-level></message-vpn></group></ldap></access-level></authentication></rpc>", groupName, msgVpnDefaultAccessLevel)
	err := c.executeEmptyCommand(ctx, command)
	if err != nil {
		return err
	}
	return nil
}

func (c *SempV1Client) CreateLdapCliGroupMessageVpnAccessLevelException(ctx context.Context, groupName string, msgVpn string) error {
	command := fmt.Sprintf("<rpc><authentication><access-level><ldap><group><group-name>%s</group-name>"+
		"<message-vpn><create><access-level-exception><vpn-name>%s</vpn-name></access-level-exception></create>"+
		"</message-vpn></group></ldap></access-level></authentication></rpc>", groupName, msgVpn)
	err := c.executeEmptyCommand(ctx, command)
	if err != nil {
		if e, ok := err.(*ExecuteResult); ok {
			if e.ExecuteResult.ReasonCode != SEMP_ERR_ALREADY_EXISTS {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func (c *SempV1Client) DeleteLdapCliGroupMessageVpnAccessLevelException(ctx context.Context, groupName string, msgVpn string) error {
	command := fmt.Sprintf("<rpc><authentication><access-level><ldap><group><group-name>%s</group-name>"+
		"<message-vpn><no><access-level-exception><vpn-name>%s</vpn-name></access-level-exception></no>"+
		"</message-vpn></group></ldap></access-level></authentication></rpc>", groupName, msgVpn)
	return c.executeEmptyCommand(ctx, command)
}

func (c *SempV1Client) SetLdapCliGroupMessageVpnAccesslevelException(ctx context.Context, groupName string, msgVpn string, accessLevel string) error {
	command := fmt.Sprintf("<rpc><authentication><access-level><ldap><group><group-name>%s</group-name>"+
		"<message-vpn><access-level-exception><vpn-name>%s</vpn-name><access-level><access-level>%s</access-level></access-level>"+
		"</access-level-exception></message-vpn></group></ldap></access-level></authentication></rpc>", groupName, msgVpn, accessLevel)
	err := c.executeEmptyCommand(ctx, command)
	if err != nil {
		return err
	}
	return nil

}

func (c *SempV1Client) SetLdapCliGroupMembershipAttributeName(ctx context.Context, attribute string) error {
	command := fmt.Sprintf("<rpc><authentication><access-level><ldap>"+
		"<group-membership-attribute-name><attribute-name>%s</attribute-name></group-membership-attribute-name>"+
		"</ldap></access-level></authentication></rpc>", attribute)
	err := c.executeEmptyCommand(ctx, command)
	if err != nil {
		return err
	}
	return nil
}

func (c *SempV1Client) DeleteLdapCliGroupMembershipAttributeName(ctx context.Context) error {
	command := fmt.Sprintf("<rpc><authentication><access-level><ldap>" +
		"<no><group-membership-attribute-name></group-membership-attribute-name></no>" +
		"</ldap></access-level></authentication></rpc>")
	err := c.executeEmptyCommand(ctx, command)
	if err != nil {
		return err
	}
	return nil
}

func (c *SempV1Client) GetBrokerAuthentication(ctx context.Context) (*BrokerAuthenticationReply, error) {
	command := "<rpc><show><authentication><access-level></access-level><detail></detail></authentication></show></rpc>"
	var result = BrokerAuthenticationReply{}

	err := c.executeCommand(ctx, command, &result)
	if err != nil {
		return nil, err
	}
	err = result.checkResult()
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *SempV1Client) SetBrokerAuthentication(ctx context.Context, auth_type string, access_level string) error {
	command_s := "<rpc><authentication><access-level><default>"
	command_e := "</default></access-level></authentication></rpc>"
	command := ""
	if auth_type == CLI_AUTH_ACCESS_LEVEL_DEFAULT_TYPE_GLOBAL {
		command = "<global-access-level><access-level>" + access_level + "</access-level></global-access-level>"
		command = command_s + command + command_e
	} else if auth_type == CLI_AUTH_ACCESS_LEVEL_DEFAULT_TYPE_MSGVPN {
		command = "<message-vpn><default-access-level><access-level>" + access_level + "</access-level></default-access-level></message-vpn>"
		command = command_s + command + command_e
	} else {
		return fmt.Errorf("%s is not a valid authentication type", auth_type)
	}
	err := c.executeEmptyCommand(ctx, command)
	if err != nil {
		return err
	}
	return nil
}

func (c *SempV1Client) DisableDefaultMessageVPN(ctx context.Context) error {
	command := "<rpc><message-vpn><vpn-name>default</vpn-name><shutdown></shutdown></message-vpn></rpc>"
	err := c.executeEmptyCommand(ctx, command)
	if err != nil {
		return err
	}
	return nil
}

func (c *SempV1Client) EnableDefaultMessageVPN(ctx context.Context) error {
	command := "<rpc><message-vpn><vpn-name>default</vpn-name><no><shutdown></shutdown></no></message-vpn></rpc>"
	err := c.executeEmptyCommand(ctx, command)
	if err != nil {
		return err
	}
	return nil
}

func (c *SempV1Client) SetSEMPIdleTimeout(ctx context.Context, idleTimeout int) error {
	command := fmt.Sprintf("<rpc><service><semp><session-idle-timeout><value>%d</value></session-idle-timeout></semp></service></rpc>", idleTimeout)
	err := c.executeEmptyCommand(ctx, command)
	if err != nil {
		return err
	}
	return nil
}

func (c *SempV1Client) CreateSyslog(ctx context.Context, syslog_name string) error {
	command := "<rpc><create><syslog><name>" + syslog_name + "</name></syslog></create></rpc>"
	err := c.executeEmptyCommand(ctx, command)
	if err != nil {
		if result, ok := err.(*ExecuteResult); ok {
			if result.ExecuteResult.ReasonCode == SEMP_ERR_ALREADY_EXISTS {
				return nil
			}
		}
	}
	return nil
}

func (c *SempV1Client) ConfigureSyslogHostAndTransport(ctx context.Context, syslog_name string, host string, transport string) error {
	command_s := "<rpc><syslog><name>" + syslog_name + "</name><host><hostname-or-address>" + host + "</hostname-or-address><transport></transport>"
	command_m := "<tcp></tcp>"
	if transport == "udp" {
		command_m = "<udp></udp>"
	}
	command_e := "</host></syslog></rpc>"
	err := c.executeEmptyCommand(ctx, command_s+command_m+command_e)
	if err != nil {
		return err
	}
	return nil
}

func (c *SempV1Client) ConfigureSyslogFacilities(ctx context.Context, syslog_name string, facility string) error {
	command_s := "<rpc><syslog><name>" + syslog_name + "</name><facility>"
	command_m := "<event></event>"
	if facility == "command" {
		command_m = "<command></command>"
	} else if facility == "system" {
		command_m = "<system></system>"
	}
	command_e := "</facility></syslog></rpc>"
	err := c.executeEmptyCommand(ctx, command_s+command_m+command_e)
	if err != nil {
		if result, ok := err.(*ExecuteResult); ok {
			if result.ExecuteResult.ReasonCode == SEMP_ERR_ALREADY_EXISTS {
				return nil
			}
		}
	}
	return nil
}

func (c *SempV1Client) RemoveSyslogFacilities(ctx context.Context, syslog_name string) error {
	command := "<rpc><syslog><name>" + syslog_name + "</name><no><facility><event></event></facility></no></syslog></rpc>"
	err := c.executeEmptyCommand(ctx, command)
	if err != nil {
		if result, ok := err.(*ExecuteResult); ok {
			if result.ExecuteResult.ReasonCode == 6 {
				// do nothing
			} else {
				return err
			}
		}
	}

	command = "<rpc><syslog><name>" + syslog_name + "</name><no><facility><system></system></facility></no></syslog></rpc>"
	err = c.executeEmptyCommand(ctx, command)
	if err != nil {
		if result, ok := err.(*ExecuteResult); ok {
			if result.ExecuteResult.ReasonCode == 6 {
				// do nothing
			} else {
				return err
			}
		}
	}

	command = "<rpc><syslog><name>" + syslog_name + "</name><no><facility><command></command></facility></no></syslog></rpc>"
	err = c.executeEmptyCommand(ctx, command)
	if err != nil {
		if result, ok := err.(*ExecuteResult); ok {
			if result.ExecuteResult.ReasonCode == 6 {
				// do nothing
			} else {
				return err
			}
		}
	}
	return nil
}

func (c *SempV1Client) GetSyslog(ctx context.Context, syslog_name string) (*SyslogReply, error) {
	var result = SyslogReply{}
	command := "<rpc><show><syslog><name>" + syslog_name + "</name></syslog></show></rpc>"
	err := c.executeCommand(ctx, command, &result)
	if err != nil {
		return nil, err
	}
	err = result.checkResult()
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *SempV1Client) SetSyslogHostAndTransport(ctx context.Context, syslog_name string, syslog_host string, syslog_transport string) error {
	err := c.ConfigureSyslogHostAndTransport(ctx, syslog_name, syslog_host, syslog_transport)
	if err != nil {
		return err
	}
	return nil
}

func (c *SempV1Client) SetSyslogFacilities(ctx context.Context, syslog_name string, syslog_facility string) error {
	err := c.ConfigureSyslogFacilities(ctx, syslog_name, syslog_facility)
	if err != nil {
		return err
	}
	return nil
}

func (c *SempV1Client) DeleteSyslog(ctx context.Context, syslog_name string) error {
	command := "<rpc><no><syslog><name>" + syslog_name + "</name></syslog></no></rpc>"
	return c.executeEmptyCommand(ctx, command)
}

func (c *SempV1Client) CreateBrokerBackup(ctx context.Context, days_of_week string, times_of_day string, max_backups int) error {
	command := "<rpc><schedule><backup><days-of-week>" + days_of_week + "</days-of-week><times-of-day>" + times_of_day + "</times-of-day>" +
		"<max-backups>" + strconv.Itoa(max_backups) + "</max-backups></backup></schedule></rpc>"
	err := c.executeEmptyCommand(ctx, command)
	if err != nil {
		if result, ok := err.(*ExecuteResult); ok {
			if result.ExecuteResult.ReasonCode == SEMP_ERR_ALREADY_EXISTS {
				return nil
			}
		}
	}
	return nil
}

func (c *SempV1Client) UpdateBrokerBackup(ctx context.Context, days_of_week string, times_of_day string, max_backups int) error {
	command := "<rpc><schedule><backup><days-of-week>" + days_of_week + "</days-of-week><times-of-day>" + times_of_day + "</times-of-day>" +
		"<max-backups>" + strconv.Itoa(max_backups) + "</max-backups></backup></schedule></rpc>"
	err := c.executeEmptyCommand(ctx, command)
	if err != nil {
		if result, ok := err.(*ExecuteResult); ok {
			if result.ExecuteResult.ReasonCode == SEMP_ERR_ALREADY_EXISTS {
				return nil
			}
		}
	}
	return nil
}

func (c *SempV1Client) GetBrokerBackup(ctx context.Context) (*BrokerBackupReply, error) {
	var result = BrokerBackupReply{}
	command := "<rpc><show><backup></backup></show></rpc>"
	err := c.executeCommand(ctx, command, &result)
	if err != nil {
		return nil, err
	}
	err = result.checkResult()
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *SempV1Client) DeleteBrokerBackup(ctx context.Context) error {
	command := "<rpc><schedule><no><backup></backup></no></schedule></rpc>"
	return c.executeEmptyCommand(ctx, command)
}

func (c *SempV1Client) GetLogRetention(ctx context.Context) (*LogRetentionReply, error) {
	var result = LogRetentionReply{}
	command := "<rpc><show><logging><config></config></logging></show></rpc>"
	err := c.executeCommand(ctx, command, &result)
	if err != nil {
		return nil, err
	}
	err = result.checkResult()
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *SempV1Client) UpdateLogRetention(ctx context.Context, retention string) error {
	command := ""
	if retention == "max-size" {
		command = "<rpc><logging><retention><max-size></max-size></retention></logging></rpc>"
	} else {
		command = "<rpc><logging><retention><days></days><max-num-days>" + retention + "</max-num-days></retention></logging></rpc>"
	}
	err := c.executeEmptyCommand(ctx, command)
	if err != nil {
		if result, ok := err.(*ExecuteResult); ok {
			if result.ExecuteResult.ReasonCode == SEMP_ERR_ALREADY_EXISTS {
				return nil
			}
		}
	}
	return nil
}

func (c *SempV1Client) EnableThresholdBrokerFragmentation(ctx context.Context) error {
	command_start := "<rpc><hardware><message-spool><defragment-spool-files><threshold>"
	command_mid := "<no><shutdown></shutdown></no>"
	command_end := "</threshold></defragment-spool-files></message-spool></hardware></rpc>"
	command := command_start + command_mid + command_end
	err := c.executeEmptyCommand(ctx, command)
	if err != nil {
		return err
	}
	return nil
}

func (c *SempV1Client) SetThresholdBrokerFragmentation(ctx context.Context, usageType, value string) error {
	command_start := "<rpc><hardware><message-spool><defragment-spool-files><threshold>"
	command_mid := ""
	command_end := "</threshold></defragment-spool-files></message-spool></hardware></rpc>"
	if usageType == "fragmentation-percentage" {
		command_mid = "<fragmentation-percentage><percentage>" + value + "</percentage></fragmentation-percentage>"
	} else if usageType == "usage-percentage" {
		command_mid = "<usage-percentage><percentage>" + value + "</percentage></usage-percentage>"
	} else if usageType == "min-interval" {
		command_mid = "<min-interval><interval>" + value + "</interval></min-interval>"
	} else {
		return fmt.Errorf("invalid usage type. Only fragmentation-percentage, usage-percentage and min-interval are possible")
	}
	command := command_start + command_mid + command_end
	err := c.executeEmptyCommand(ctx, command)
	if err != nil {
		if result, ok := err.(*ExecuteResult); ok {
			if result.ExecuteResult.ReasonCode == SEMP_ERR_ALREADY_EXISTS {
				return nil
			}
		}
	}
	return nil
}

func (c *SempV1Client) DisableThresholdBrokerFragmentation(ctx context.Context) error {
	command_start := "<rpc><hardware><message-spool><defragment-spool-files><threshold>"
	command_mid := "<shutdown></shutdown>"
	command_end := "</threshold></defragment-spool-files></message-spool></hardware></rpc>"
	command := command_start + command_mid + command_end
	return c.executeEmptyCommand(ctx, command)
}

func (c *SempV1Client) EnableScheduledBrokerFragmentation(ctx context.Context) error {
	command_start := "<rpc><hardware><message-spool><defragment-spool-files><schedule>"
	command_mid := "<no><shutdown></shutdown></no>"
	command_end := "</schedule></defragment-spool-files></message-spool></hardware></rpc>"
	command := command_start + command_mid + command_end
	err := c.executeEmptyCommand(ctx, command)
	if err != nil {
		return err
	}
	return nil
}

func (c *SempV1Client) SetScheduledBrokerFragmentation(ctx context.Context, usageType, value string) error {
	command_start := "<rpc><hardware><message-spool><defragment-spool-files><schedule>"
	command_mid := ""
	command_end := "</schedule></defragment-spool-files></message-spool></hardware></rpc>"
	if usageType == "days" {
		command_mid = "<days><days-of-week>" + value + "</days-of-week></days>"
	} else if usageType == "times" {
		command_mid = "<times><times-of-day>" + value + "</times-of-day></times>"
	} else {
		return fmt.Errorf("invalid usage type. Only days and times")
	}
	command := command_start + command_mid + command_end
	err := c.executeEmptyCommand(ctx, command)
	if err != nil {
		if result, ok := err.(*ExecuteResult); ok {
			if result.ExecuteResult.ReasonCode == SEMP_ERR_ALREADY_EXISTS {
				return nil
			}
		}
	}
	return nil
}

func (c *SempV1Client) DisableScheduledBrokerFragmentation(ctx context.Context) error {
	command_start := "<rpc><hardware><message-spool><defragment-spool-files><schedule>"
	command_mid := "<shutdown></shutdown>"
	command_end := "</schedule></defragment-spool-files></message-spool></hardware></rpc>"
	command := command_start + command_mid + command_end
	return c.executeEmptyCommand(ctx, command)
}

func (c *SempV1Client) GetBrokerFragmentationSettings(ctx context.Context) (*MessageSpoolReply, error) {
	var result = MessageSpoolReply{}
	command := "<rpc><show><message-spool><detail></detail></message-spool></show></rpc>"
	err := c.executeCommand(ctx, command, &result)
	if err != nil {
		return nil, err
	}
	err = result.checkResult()
	if err != nil {
		return nil, err
	}
	return &result, nil
}
