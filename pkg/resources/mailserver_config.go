// Copyright 2022 Linka Cloud  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resources

import (
	"fmt"
	"reflect"
	"strconv"
)

type LogLevel string

const (
	// LogLevelError is the error log level
	LogLevelError LogLevel = "error"
	// LogLevelWarn is the warning log level
	LogLevelWarn LogLevel = "warn"
	// LogLevelInfo is the info log level
	LogLevelInfo LogLevel = "info"
	// LogLevelDebug is the debug log level
	LogLevelDebug LogLevel = "debug"
)

type AccountProvisioner string

const (
	// AccountProvisionerLDAP is the LDAP account provisioner
	AccountProvisionerLDAP AccountProvisioner = "ldap"
	// AccountProvisionerOIDC is the OIDC account provisioner (not implemented yet)
	AccountProvisionerOIDC AccountProvisioner = "oidc"
	// AccountProvisionerFile is the file account provisioner (not implemented yet)
	AccountProvisionerFile AccountProvisioner = "file"
)

type MailServerConfig struct {
	// OverrideHostname
	// empty => uses the `hostname` command to get the mail server's canonical hostname
	// => Specify a fully-qualified domainname to serve mail for.  This is used for many of the config features so if you can't set your hostname (e.g. you're in a container platform that doesn't let you) specify it in this environment variable.
	OverrideHostname string `json:"overrideHostname,omitempty" env:"OVERRIDE_HOSTNAME"`

	// LogLevel
	// Set the log level for DMS.
	// This is mostly relevant for container startup scripts and change detection event feedback.
	//
	// Valid values (in order of increasing verbosity) are: `error`, `warn`, `info`, `debug` and `trace`.
	// The default log level is `info`.
	LogLevel LogLevel `json:"logLevel,omitempty" env:"LOG_LEVEL"`

	// SupervisorLoglevel
	// critical => Only show critical messages
	// error => Only show erroneous output
	// **warn** => Show warnings
	// info => Normal informational output
	// debug => Also show debug messages
	SupervisorLoglevel LogLevel `json:"supervisorLoglevel,omitempty" env:"SUPERVISOR_LOGLEVEL"`

	// OneDir
	// false => mail state in default directories
	// true => consolidate all states into a single directory (`/var/mail-state`) to allow persistence using docker volumes
	OneDir *bool `json:"oneDir,omitempty" env:"ONE_DIR"`

	// AccountProvisioner
	// **empty** => use FILE
	// LDAP => use LDAP authentication
	// OIDC => use OIDC authentication (not yet implemented)
	// FILE => use local files (this is used as the default)
	AccountProvisioner AccountProvisioner `json:"accountProvisioner,omitempty" env:"ACCOUNT_PROVISIONER"`
	// PostmasterAddress
	// empty => postmaster@domain.com
	// => Specify the postmaster address
	PostmasterAddress string `json:"postmasterAddress,omitempty" env:"POSTMASTER_ADDRESS"`

	// EnableUpdateCheck
	// Check for updates on container start and then once a day
	// If an update is available, a mail is sent to POSTMASTER_ADDRESS
	// false => Update check disabled
	// true => Update check enabled
	EnableUpdateCheck bool `json:"enableUpdateCheck,omitempty" env:"ENABLE_UPDATE_CHECK"`
	// UpdateCheckInterval
	// Customize the update check interval.
	// Number + Suffix. Suffix must be 's' for seconds, 'm' for minutes, 'h' for hours or 'd' for days.
	UpdateCheckInterval string `json:"updateCheckInterval,omitempty" env:"UPDATE_CHECK_INTERVAL"`
	// PermitDocker
	// Set different options for mynetworks option (can be overwrite in postfix-main.cf)
	// **WARNING**: Adding the docker network's gateway to the list of trusted hosts, e.g. using the `network` or
	// `connected-networks` option, can create an open relay
	// https://github.com/docker-mailserver/docker-mailserver/issues/1405//issuecomment-590106498
	// The same can happen for rootless podman. To prevent this, set the value to "none" or configure slirp4netns
	// https://github.com/docker-mailserver/docker-mailserver/issues/2377
	//
	// none => Explicitly force authentication
	// container => Container IP address only
	// host => Add docker container network (ipv4 only)
	// network => Add all docker container networks (ipv4 only)
	// connected-networks => Add all connected docker networks (ipv4 only)
	PermitDocker string `json:"permitDocker,omitempty" env:"PERMIT_DOCKER"`
	// TZ
	// Set the timezone. If this variable is unset, the container runtime will try to detect the time using
	// `/etc/localtime`, which you can alternatively mount into the container. The value of this variable
	// must follow the pattern `AREA/ZONE`, i.e. of you want to use Germany's time zone, use `Europe/Berlin`.
	// You can look up all available timezones here: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones//List
	TZ string `json:"tz,omitempty" env:"TZ"`
	// NetworkInterface
	// In case you network interface differs from 'eth0', e.g. when you are using HostNetworking in Kubernetes,
	// you can set NETWORK_INTERFACE to whatever interface you want. This interface will then be used.
	//  - **empty** => eth0
	NetworkInterface string `json:"networkInterface,omitempty" env:"NETWORK_INTERFACE"`
	// TLSLevel
	// empty => modern
	// modern => Enables TLSv1.2 and modern ciphers only. (default)
	// intermediate => Enables TLSv1, TLSv1.1 and TLSv1.2 and broad compatibility ciphers.
	TLSLevel string `json:"tlsLevel,omitempty" env:"TLS_LEVEL"`
	// SSLType
	// Please read [the SSL page in the documentation](https://docker-mailserver.github.io/docker-mailserver/edge/config/security/ssl) for more information.
	//
	// empty => SSL disabled
	// letsencrypt => Enables Let's Encrypt certificates
	// custom => Enables custom certificates
	// manual => Let you manually specify locations of your SSL certificates for non-standard cases
	// self-signed => Enables self-signed certificates
	SSLType string `json:"sslType,omitempty" env:"SSL_TYPE"`
	// SSLCertPath
	// These are only supported with `SSL_TYPE=manual`.
	// Provide the path to your cert and key files that you've mounted access to within the container.
	SSLCertPath string `json:"sslCertPath,omitempty" env:"SSL_CERT_PATH"`
	// SSLKeyPath
	SSLKeyPath string `json:"sslKeyPath,omitempty" env:"SSL_KEY_PATH"`
	// SSLAltCertPath
	// Optional: A 2nd certificate can be supported as fallback (dual cert support), eg ECDSA with an RSA fallback.
	// Useful for additional compatibility with older MTA and MUA (eg pre-2015).
	SSLAltCertPath string `json:"sslAltCertPath,omitempty" env:"SSL_ALT_CERT_PATH"`
	// SSLAltKeyPath
	SSLAltKeyPath string `json:"sslAltKeyPath,omitempty" env:"SSL_ALT_KEY_PATH"`

	// SMTPOnly only launch postfix smtp
	SMTPOnly bool `json:"smtpOnly,omitempty" env:"SMTP_ONLY"`

	// SpoofProtection configures the handling of creating mails with forged sender addresses.
	//
	// false => (not recommended, but default for backwards compatibility reasons)
	//           Mail address spoofing allowed. Any logged-in user may create email messages with a forged sender address.
	//           See also https://en.wikipedia.org/wiki/Email_spoofing
	// true => (recommended) Mail spoofing denied. Each user may only send with his own or his alias addresses.
	//       Addresses with extension delimiters(http://www.postfix.org/postconf.5.html//recipient_delimiter) are not able to send messages.
	SpoofProtection bool `json:"spoofProtection,omitempty" env:"SPOOF_PROTECTION"`

	// EnablePOP3 enables the POP3 service
	EnablePOP3 bool `json:"enablePop3,omitempty" env:"ENABLE_POP3"`
	// EnableClamav enables the ClamAV service
	EnableClamav bool `json:"enableClamav,omitempty" env:"ENABLE_CLAMAV"`
	// EnableAmavis enables the Amavis service
	EnableAmavis bool `json:"enableAmavis,omitempty" env:"ENABLE_AMAVIS"`
	// AmavisLoglevel
	// -1/-2/-3 => Only show errors
	// **0**    => Show warnings
	// 1/2      => Show default informational output
	// 3/4/5    => log debug information (very verbose)
	AmavisLoglevel string `json:"amavisLoglevel,omitempty" env:"AMAVIS_LOGLEVEL"`

	// EnableDNSBL enables DNS BlackList service
	// This enables the [zen.spamhaus.org](https://www.spamhaus.org/zen/) DNS block list in postfix
	// and various [lists](https://github.com/docker-mailserver/docker-mailserver/blob/f7465a50888eef909dbfc01aff4202b9c7d8bc00/target/postfix/main.cf//L58-L66) in postscreen.
	// Note: Emails will be rejected, if they don't pass the block list checks!
	EnableDNSBL bool `json:"enableDnsbl,omitempty" env:"ENABLE_DNSBL"`

	// EnableFail2ban enables Fail2ban service
	// If you enable Fail2Ban, don't forget to add the following lines to your `docker-compose.yml`:
	//    cap_add:
	//      - NET_ADMIN
	// Otherwise, `nftables` won't be able to ban IPs.
	EnableFail2ban bool `json:"enableFail2ban,omitempty" env:"ENABLE_FAIL2BAN"`
	// Fail2banBlockType
	// Fail2Ban blocktype
	// drop   => drop packet (send NO reply)
	// reject => reject packet (send ICMP unreachable)
	Fail2banBlockType string `json:"fail2banBlocktype,omitempty" env:"FAIL2BAN_BLOCKTYPE"`

	// EnableManageSieve
	// Enables Managesieve on port 4190
	EnableManageSieve bool `json:"enableManagesieve,omitempty" env:"ENABLE_MANAGESIEVE"`

	// PostscreenAction
	// **enforce** => Allow other tests to complete. Reject attempts to deliver mail with a 550 SMTP reply, and log the helo/sender/recipient information. Repeat this test the next time the client connects.
	// drop => Drop the connection immediately with a 521 SMTP reply. Repeat this test the next time the client connects.
	// ignore => Ignore the failure of this test. Allow other tests to complete. Repeat this test the next time the client connects. This option is useful for testing and collecting statistics without blocking mail.
	PostscreenAction string `json:"postscreenAction,omitempty" env:"POSTSCREEN_ACTION"`

	// ClamavMessageSizeLimit
	// Mails larger than this limit won't be scanned.
	// ClamAV must be enabled (ENABLE_CLAMAV=1) for this.
	//
	// empty => 25M (25 MB)
	ClamavMessageSizeLimit string `json:"clamavMessageSizeLimit,omitempty" env:"CLAMAV_MESSAGE_SIZE_LIMIT"`
	// VirusMailsDeleteDelay
	// Set how many days a virusmail will stay on the server before being deleted
	// empty => 7 days
	VirusMailsDeleteDelay string `json:"virusmailsDeleteDelay,omitempty" env:"VIRUSMAILS_DELETE_DELAY"`

	// EnablePostfixVirtualTransport
	// // This Option is activating the Usage of POSTFIX_DAGENT to specify a lmtp client different from default dovecot socket.
	EnablePostfixVirtualTransport bool `json:"enablePostfixVirtualTransport,omitempty" env:"ENABLE_POSTFIX_VIRTUAL_TRANSPORT"`
	// PostfixDagent
	// Enabled by ENABLE_POSTFIX_VIRTUAL_TRANSPORT. Specify the final delivery of postfix
	//
	// empty => fail
	// `lmtp:unix:private/dovecot-lmtp` (use socket)
	// `lmtps:inet:<host>:<port>` (secure lmtp with starttls, take a look at https://sys4.de/en/blog/2014/11/17/sicheres-lmtp-mit-starttls-in-dovecot/)
	// `lmtp:<kopano-host>:2003` (use kopano as mailstore)
	// etc.
	PostfixDagent string `json:"postfixDagent,omitempty" env:"POSTFIX_DAGENT"`
	// PostfixMailboxSizeLimit
	// Set the mailbox size limit for all users. If set to zero, the size will be unlimited (default).
	PostfixMailboxSizeLimit string `json:"postfixMailboxSizeLimit,omitempty" env:"POSTFIX_MAILBOX_SIZE_LIMIT"`
	// EnableQuotas
	// See https://docker-mailserver.github.io/docker-mailserver/edge/config/user-management/accounts///notes
	EnableQuotas bool `json:"enableQuotas,omitempty" env:"ENABLE_QUOTAS"`
	// PostfixMessageSizeLimit
	// Set the message size limit for all users. If set to zero, the size will be unlimited (not recommended!)
	//
	// empty => 10240000 (~10 MB)
	PostfixMessageSizeLimit string `json:"postfixMessageSizeLimit,omitempty" env:"POSTFIX_MESSAGE_SIZE_LIMIT"`

	// PflogsummTrigger Enables regular pflogsumm mail reports.
	// This is a new option. The old REPORT options are still supported for backwards compatibility. If this is not set and reports are enabled with the old options, logrotate will be used.
	//
	// not set => No report
	// daily_cron => Daily report for the previous day
	// logrotate => Full report based on the mail log when it is rotated
	PflogsummTrigger string `json:"pflogsummTrigger,omitempty" env:"PFLOGSUMM_TRIGGER"`
	// PflogsummRecipient: Recipient address for pflogsumm reports.
	//
	// not set => Use REPORT_RECIPIENT or POSTMASTER_ADDRESS
	// => Specify the recipient address(es)
	PflogsummRecipient string `json:"pflogsummRecipient,omitempty" env:"PFLOGSUMM_RECIPIENT"`
	// PflogsummSender: Sender address (`FROM`) for pflogsumm reports if pflogsumm reports are enabled.
	//
	// not set => Use REPORT_SENDER
	// => Specify the sender address
	PflogsummSender string `json:"pflogsummSender,omitempty" env:"PFLOGSUMM_SENDER"`
	// LogwatchInterval: Interval for logwatch report.
	//
	// none => No report is generated
	// daily => Send a daily report
	// weekly => Send a report every week
	LogwatchInterval string `json:"logwatchInterval,omitempty" env:"LOGWATCH_INTERVAL"`
	// LogwatchRecipient: Recipient address for logwatch reports if they are enabled.
	//
	// not set => Use REPORT_RECIPIENT or POSTMASTER_ADDRESS
	// => Specify the recipient address(es)
	LogwatchRecipient string `json:"logwatchRecipient,omitempty" env:"LOGWATCH_RECIPIENT"`
	// LogwatchSender: Sender address (`FROM`) for logwatch reports if logwatch reports are enabled.
	//
	// not set => Use REPORT_SENDER
	// => Specify the sender address
	LogwatchSender string `json:"logwatchSender,omitempty" env:"LOGWATCH_SENDER"`
	// ReportRecipient: Defines who receives reports if they are enabled.
	// **empty** => ${POSTMASTER_ADDRESS}
	// => Specify the recipient address
	ReportRecipient string `json:"reportRecipient,omitempty" env:"REPORT_RECIPIENT"`
	// ReportSender: Defines who sends reports if they are enabled.
	// **empty** => mailserver-report@${DOMAINNAME}
	// => Specify the sender address
	ReportSender string `json:"reportSender,omitempty" env:"REPORT_SENDER"`

	// LogrotateInterval: Changes the interval in which log files are rotated
	// **weekly** => Rotate log files weekly
	// daily => Rotate log files daily
	// monthly => Rotate log files monthly
	//
	// Note: This Variable actually controls logrotate inside the container
	// and rotates the log files depending on this setting. The main log output is
	// still available in its entirety via `docker logs mail` (Or your
	// respective container name). If you want to control logrotation for
	// the Docker-generated logfile see:
	// https://docs.docker.com/config/containers/logging/configure/
	//
	// Note: This variable can also determine the interval for Postfix's log summary reports, see [`PFLOGSUMM_TRIGGER`](//pflogsumm_trigger).
	LogrotateInterval string `json:"logrotateInterval,omitempty" env:"LOGROTATE_INTERVAL"`
	// PostfixInetProtocols: Choose TCP/IP protocols for postfix to use
	// **all** => All possible protocols.
	// ipv4 => Use only IPv4 traffic. Most likely you want this behind Docker.
	// ipv6 => Use only IPv6 traffic.
	//
	// Note: More details at http://www.postfix.org/postconf.5.html//inet_protocols
	PostfixInetProtocols string `json:"postfixInetProtocols,omitempty" env:"POSTFIX_INET_PROTOCOLS"`
	// DovecotInetProtocols: Choose TCP/IP protocols for dovecot to use
	// **all** => Listen on all interfaces
	// ipv4 => Listen only on IPv4 interfaces. Most likely you want this behind Docker.
	// ipv6 => Listen only on IPv6 interfaces.
	//
	// Note: More information at https://dovecot.org/doc/dovecot-example.conf
	DovecotInetProtocols string `json:"dovecotInetProtocols,omitempty" env:"DOVECOT_INET_PROTOCOLS"`
	// EnableSpamassassin Enables Spamassassin
	EnableSpamassassin bool `json:"enableSpamassassin,omitempty" env:"ENABLE_SPAMASSASSIN"`
	// SpamassassinSpamToInbox: deliver spam messages in the inbox (eventually tagged using SA_SPAM_SUBJECT)
	SpamassassinSpamToInbox string `json:"spamassassinSpamToInbox,omitempty" env:"SPAMASSASSIN_SPAM_TO_INBOX"`
	// EnableSpamassassinKam: KAM is a 3rd party SpamAssassin ruleset, provided by the McGrail Foundation.
	// If SpamAssassin is enabled, KAM can be used in addition to the default ruleset.
	// - **0** => KAM disabled
	// - 1 => KAM enabled
	//
	// Note: only has an effect if `ENABLE_SPAMASSASSIN=1`
	EnableSpamassassinKam bool `json:"enableSpamassassinKam,omitempty" env:"ENABLE_SPAMASSASSIN_KAM"`
	// MoveSpamToJunk: spam messages will be moved in the Junk folder (SPAMASSASSIN_SPAM_TO_INBOX=1 required)
	MoveSpamToJunk bool `json:"moveSpamToJunk,omitempty" env:"MOVE_SPAM_TO_JUNK"`
	// SATag: add spam info headers if at, or above that level:
	SATag string `json:"saTag,omitempty" env:"SA_TAG"`
	// SATag2 add 'spam detected' headers at that level
	SATag2 string `json:"saTag2,omitempty" env:"SA_TAG2"`
	// SAKill: triggers spam evasive actions
	SAKill string `json:"saKill,omitempty" env:"SA_KILL"`
	// SASpamSubject: add tag to subject if spam detected
	SASpamSubject string `json:"saSpamSubject,omitempty" env:"SA_SPAM_SUBJECT"`

	// EnableFetchmail enables fetchmail
	EnableFetchmail bool `json:"enableFetchmail,omitempty" env:"ENABLE_FETCHMAIL"`
	// FetchmailPoll: The interval to fetch mail in seconds
	FetchmailPoll string `json:"fetchmailPoll,omitempty" env:"FETCHMAIL_POLL"`

	// EnablePostgrey enables postgrey
	EnablePostgrey bool `json:"enablePostgrey,omitempty" env:"ENABLE_POSTGREY"`
	// PostgreyDelay: The delay for postgrey
	// greylist for N seconds
	PostgreyDelay string `json:"postgreyDelay,omitempty" env:"POSTGREY_DELAY"`
	// PostgreyMaxAge: The max age for postgrey
	// delete entries older than N days since the last time that they have been seen
	PostgreyMaxAge string `json:"postgreyMaxAge,omitempty" env:"POSTGREY_MAX_AGE"`
	// PostgreyText: response when a mail is greylisted
	PostgreyText string `json:"postgreyText,omitempty" env:"POSTGREY_TEXT"`
	// PostgreyAutoWhitelistClients: whitelist host after N successful deliveries (N=0 to disable whitelisting)
	PostgreyAutoWhitelistClients string `json:"postgreyAutoWhitelistClients,omitempty" env:"POSTGREY_AUTO_WHITELIST_CLIENTS"`
	// EnableSRS enables the Sender Rewriting Scheme. SRS is needed if your mail server acts as forwarder. See [postsrsd](https://github.com/roehling/postsrsd/blob/master/README.md//sender-rewriting-scheme-crash-course) for further explanation.
	EnableSRS bool `json:"enableSrs,omitempty" env:"ENABLE_SRS"`
	// SRSSenderClasses
	// envelope_sender => Rewrite only envelope sender address (default)
	// header_sender => Rewrite only header sender (not recommended)
	// envelope_sender,header_sender => Rewrite both senders
	// An email has an "envelope" sender (indicating the sending server) and a
	// "header" sender (indicating who sent it). More strict SPF policies may require
	// you to replace both instead of just the envelope sender.
	SRSSenderClasses string `json:"srsSenderClasses,omitempty" env:"SRS_SENDER_CLASSES"`
	// SRSExcludeDomains
	// empty => Envelope sender will be rewritten for all domains
	// provide comma separated list of domains to exclude from rewriting
	SRSExcludeDomains string `json:"srsExcludeDomains,omitempty" env:"SRS_EXCLUDE_DOMAINS"`
	// SRSSecret
	// empty => generated when the image is built
	// provide a secret to use in base64
	// you may specify multiple keys, comma separated. the first one is used for
	// signing and the remaining will be used for verification. this is how you
	// rotate and expire keys
	SRSSecret string `json:"srsSecret,omitempty" env:"SRS_SECRET"`

	// DefaultRelayHost: Setup relaying all mail through a default relay host
	//
	// empty => don't configure default relay host
	// default host and optional port to relay all mail through
	DefaultRelayHost string `json:"defaultRelayHost,omitempty" env:"DEFAULT_RELAY_HOST"`

	// RelayHost: Setup relaying for multiple domains based on the domain name of the sender
	// optionally uses usernames and passwords in postfix-sasl-password.cf and relay host mappings in postfix-relaymap.cf
	//
	// empty => don't configure relay host
	// default host to relay mail through
	RelayHost string `json:"relayHost,omitempty" env:"RELAY_HOST"`
	// RelayPort
	// empty => 25
	// default port to relay mail
	RelayPort int `json:"relayPort,omitempty" env:"RELAY_PORT"`
	// RelayUser
	// empty => no defaults
	// default relay username (if no specific entry exists in postfix-sasl-password.cf)
	RelayUser string `json:"relayUser,omitempty" env:"RELAY_USER"`
	// RelayPassword
	// empty => no default
	// password for default relay user
	RelayPassword string `json:"relayPassword,omitempty" env:"RELAY_PASSWORD"`

	// EnableLdap enables LDAP Authentication
	// A second container for the ldap service is necessary (i.e. https://github.com/osixia/docker-openldap)
	// For preparing the ldap server to use in combination with this container this article may be helpful: http://acidx.net/wordpress/2014/06/installing-a-mailserver-with-postfix-dovecot-sasl-ldap-roundcube/

	// with the :edge tag, use ACCOUNT_PROVISIONER=LDAP
	EnableLdap bool `json:"enableLdap,omitempty" env:"ENABLE_LDAP"`
	// LdapStartTLS enables STARTTLS for LDAP
	LdapStartTLS bool `json:"ldapStartTls,omitempty" env:"LDAP_START_TLS"`
	// LdapServerHost: The hostname of the LDAP server
	// empty => mail.domain.com
	// Specify the dns-name/ip-address where the ldap-server
	LdapServerHost string `json:"ldapServerHost,omitempty" env:"LDAP_SERVER_HOST"`
	// LdapSearchBase: The search base for LDAP
	// empty => ou=people,dc=domain,dc=com
	// => e.g. LDAP_SEARCH_BASE=dc=mydomain,dc=local
	LdapSearchBase string `json:"ldapSearchBase,omitempty" env:"LDAP_SEARCH_BASE"`
	// LdapBindDn: The bind dn for LDAP
	// empty => cn=admin,dc=domain,dc=com
	// => take a look at examples of SASL_LDAP_BIND_DN
	LdapBindDN string `json:"ldapBindDn,omitempty" env:"LDAP_BIND_DN"`
	// LdapBindPw: The bind password for LDAP
	// empty** => admin
	// => Specify the password to bind against ldap
	LdapBindPW string `json:"ldapBindPw,omitempty" env:"LDAP_BIND_PW"`
	// LdapQueryFilterUser: The query filter for users
	// e.g. `"(&(mail=%s)(mailEnabled=TRUE))"`
	// => Specify how ldap should be asked for users
	LdapQueryFilterUser string `json:"ldapQueryFilterUser,omitempty" env:"LDAP_QUERY_FILTER_USER"`
	// LdapQueryFilterGroup: The query filter for groups
	// e.g. `"(&(mailGroupMember=%s)(mailEnabled=TRUE))"`
	// => Specify how ldap should be asked for groups
	LdapQueryFilterGroup string `json:"ldapQueryFilterGroup,omitempty" env:"LDAP_QUERY_FILTER_GROUP"`
	// LdapQueryFilterAlias: The query filter for aliases
	// e.g. `"(&(mailAlias=%s)(mailEnabled=TRUE))"`
	// => Specify how ldap should be asked for aliases
	LdapQueryFilterAlias string `json:"ldapQueryFilterAlias,omitempty" env:"LDAP_QUERY_FILTER_ALIAS"`
	// LdapQueryFilterDomain: The query filter for domains
	// e.g. `"(&(|(mail=*@%s)(mailalias=*@%s)(mailGroupMember=*@%s))(mailEnabled=TRUE))"`
	// => Specify how ldap should be asked for domains
	LdapQueryFilterDomain string `json:"ldapQueryFilterDomain,omitempty" env:"LDAP_QUERY_FILTER_DOMAIN"`

	// DovecotTLS enables LDAP over TLS for Dovecot
	// LDAP over TLS enabled for Dovecot
	DovecotTLS         bool   `json:"dovecotTls,omitempty" env:"DOVECOT_TLS"`
	DovecotHosts       string `json:"dovecotHosts,omitempty" env:"DOVECOT_HOSTS"`
	DovecotLdapVersion string `json:"dovecotLdapVersion,omitempty" env:"DOVECOT_LDAP_VERSION"`

	// DovecotUserFilter: The user filter for Dovecot
	// e.g. `"(&(objectClass=PostfixBookMailAccount)(uniqueIdentifier=%n))"`
	DovecotUserFilter string `json:"dovecotUserFilter,omitempty" env:"DOVECOT_USER_FILTER"`
	// DovecotPassFilter: The password filter for Dovecot
	// e.g. `"(&(objectClass=PostfixBookMailAccount)(uniqueIdentifier=%n))"`
	DovecotPassFilter string `json:"dovecotPassFilter,omitempty" env:"DOVECOT_PASS_FILTER"`
	// DovecotUserAttrs: The user attributes for Dovecot
	// Define the mailbox format to be used
	// default is maildir, supported values are: sdbox, mdbox, maildir
	DovecotMailboxFormat string `json:"dovecotMailboxFormat,omitempty" env:"DOVECOT_MAILBOX_FORMAT"`
	// DovecotAuthBind: The authentication bind for Dovecot
	// Allow bind authentication for LDAP
	// https://wiki.dovecot.org/AuthDatabase/LDAP/AuthBinds
	DovecotAuthBind  bool   `json:"dovecotAuthBind,omitempty" env:"DOVECOT_AUTH_BIND"`
	DovecotScope     string `json:"dovecotScope,omitempty" env:"DOVECOT_SCOPE"`
	DovecotUserAttrs string `json:"dovecotUserAttrs,omitempty" env:"DOVECOT_USER_ATTRS"`

	// EnableSASLAuthd enables SASL DLAP authentication
	EnableSASLAuthd bool `json:"enableSaslauthd,omitempty" env:"ENABLE_SASLAUTHD"`
	// SASLAuthdMechanisms
	// empty => pam
	// `ldap` => authenticate against ldap server
	// `shadow` => authenticate against local user db
	// `mysql` => authenticate against mysql db
	// `rimap` => authenticate against imap server
	// Note: can be a list of mechanisms like pam ldap shadow
	SASLAuthdMechanisms string `json:"saslauthdMechanisms,omitempty" env:"SASLAUTHD_MECHANISMS"`
	// SASLAuthdMechOptions
	// empty => None
	// e.g. with SASLAUTHD_MECHANISMS rimap you need to specify the ip-address/servername of the imap server  ==> xxx.xxx.xxx.xxx
	SASLAuthdMechOptions string `json:"saslauthdMechOptions,omitempty" env:"SASLAUTHD_MECH_OPTIONS"`
	// SASLAuthdLdapServer
	// empty => Use value of LDAP_SERVER_HOST
	// Note: since version 10.0.0, you can specify a protocol here (like ldaps://); this deprecates SASLAUTHD_LDAP_SSL.
	SASLAuthdLdapServer string `json:"saslauthdLdapServer,omitempty" env:"SASLAUTHD_LDAP_SERVER"`
	// SASLAuthdLdapBindDn
	// empty => Use value of LDAP_BIND_DN
	// specify an object with priviliges to search the directory tree
	// e.g. active directory: SASLAUTHD_LDAP_BIND_DN=cn=Administrator,cn=Users,dc=mydomain,dc=net
	// e.g. openldap: SASLAUTHD_LDAP_BIND_DN=cn=admin,dc=mydomain,dc=net
	SASLAuthdLdapBindDn string `json:"saslauthdLdapBindDn,omitempty" env:"SASLAUTHD_LDAP_BIND_DN"`
	// SASLAuthdLdapPassword
	// empty => Use value of LDAP_BIND_PW
	SASLAuthdLdapPassword string `json:"saslauthdLdapPassword,omitempty" env:"SASLAUTHD_LDAP_PASSWORD"`
	// SASLAuthdLdapSearchBase
	// empty => Use value of LDAP_SEARCH_BASE
	// specify the search base
	SASLAuthdLdapSearchBase string `json:"saslauthdLdapSearchBase,omitempty" env:"SASLAUTHD_LDAP_SEARCH_BASE"`
	// SASLAuthdLdapFilter
	// empty => default filter `(&(uniqueIdentifier=%u)(mailEnabled=TRUE))`
	// e.g. for active directory: `(&(sAMAccountName=%U)(objectClass=person))`
	// e.g. for openldap: `(&(uid=%U)(objectClass=person))`
	SASLAuthdLdapFilter string `json:"saslauthdLdapFilter,omitempty" env:"SASLAUTHD_LDAP_FILTER"`
	// SASLAuthdLdapStartTls
	// empty => no
	// yes => LDAP over TLS enabled for SASL
	// If set to yes, the protocol in SASLAUTHD_LDAP_SERVER must be ldap:// or missing.
	SASLAuthdLdapStartTls bool `json:"saslauthdLdapStartTls,omitempty" env:"SASLAUTHD_LDAP_START_TLS"`
	// SASLAuthdLdapTlsCheckPeer
	// empty => no
	// yes => Require and verify server certificate
	// If yes you must/could specify SASLAUTHD_LDAP_TLS_CACERT_FILE or SASLAUTHD_LDAP_TLS_CACERT_DIR.
	SASLAuthdLdapTlsCheckPeer bool `json:"saslauthdLdapTlsCheckPeer,omitempty" env:"SASLAUTHD_LDAP_TLS_CHECK_PEER"`
	// SASLAuthdLdapTlsCacertFile: File containing CA (Certificate Authority) certificate(s).
	// empty => Nothing is added to the configuration
	// Any value => Fills the `ldap_tls_cacert_file` option
	SASLAuthdLdapTlsCacertFile string `json:"saslauthdLdapTlsCacertFile,omitempty" env:"SASLAUTHD_LDAP_TLS_CACERT_FILE"`
	// SASLAuthdLdapTlsCacertDir: Path to directory with CA (Certificate Authority) certificates.
	// empty => Nothing is added to the configuration
	// Any value => Fills the `ldap_tls_cacert_dir` option
	SASLAuthdLdapTlsCacertDir string `json:"saslauthdLdapTlsCacertDir,omitempty" env:"SASLAUTHD_LDAP_TLS_CACERT_DIR"`
	// SASLAuthdLdapPasswordAttr: Specify what password attribute to use for password verification.
	// empty => Nothing is added to the configuration but the documentation says it is `userPassword` by default.
	// Any value => Fills the `ldap_password_attr` option
	SASLAuthdLdapPasswordAttr string `json:"saslauthdLdapPasswordAttr,omitempty" env:"SASLAUTHD_LDAP_PASSWORD_ATTR"`
	// SASLPasswd
	// empty => No sasl_passwd will be created
	// string => `/etc/postfix/sasl_passwd` will be created with the string as password
	SASLPasswd string `json:"saslPasswd,omitempty" env:"SASL_PASSWD"`
	// SASLAuthdLdapAuthMethod
	// empty => `bind` will be used as a default value
	// `fastbind` => The fastbind method is used
	// `custom` => The custom method uses userPassword attribute to verify the password
	SASLAuthdLdapAuthMethod string `json:"saslauthdLdapAuthMethod,omitempty" env:"SASLAUTHD_LDAP_AUTH_METHOD"`
	// SASLAuthdLdapMech: Specify the authentication mechanism for SASL bind
	// empty => Nothing is added to the configuration
	// Any value => Fills the `ldap_mech` option
	SASLAuthdLdapMech string `json:"saslauthdLdapMech,omitempty" env:"SASLAUTHD_LDAP_MECH"`
}

func (c MailServerConfig) ToMap() map[string][]byte {
	// iterate over the struct and return a map of the values
	m := make(map[string][]byte)
	v := reflect.ValueOf(c)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		k := t.Field(i).Tag.Get("env")
		switch v := v.Field(i).Interface().(type) {
		case int:
			if v != 0 {
				m[k] = []byte(strconv.Itoa(v))
			} else {
				m[k] = []byte("")
			}
		case *bool:
			if v != nil && *v {
				m[k] = []byte("1")
			} else {
				m[k] = []byte("")
			}
		case bool:
			switch k {
			case "LDAP_START_TLS",
				"DOVECOT_TLS",
				"DOVECOT_AUTH_BIND",
				"SASLAUTHD_LDAP_START_TLS",
				"SASLAUTHD_LDAP_TLS_CHECK_PEER":
				if v {
					m[k] = []byte("yes")
				} else {
					m[k] = []byte("no")
				}
			default:
				if v {
					m[k] = []byte("1")
				} else {
					m[k] = []byte("")
				}
			}
		default:
			m[k] = []byte(fmt.Sprintf("%v", v))
		}
	}
	return m
}
