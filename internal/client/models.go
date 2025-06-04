package client

import "time"

// Project represents a Mailtrap project
type Project struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	ShareLinks  ShareLink `json:"share_links"`
	Permissions []string  `json:"permissions"`
	Inboxes     []Inbox   `json:"inboxes"`
}

// ProjectRequest represents a request to create/update a project
type ProjectRequest struct {
	Project struct {
		Name string `json:"name"`
	} `json:"project"`
}

// ShareLink represents share links for a project
type ShareLink struct {
	Admin  string `json:"admin"`
	Viewer string `json:"viewer"`
}

// Inbox represents a Mailtrap inbox
type Inbox struct {
	ID                    int         `json:"id"`
	Name                  string      `json:"name"`
	Username              string      `json:"username"`
	Password              string      `json:"password"`
	MaxSize               int         `json:"max_size"`
	Status                string      `json:"status"`
	EmailUsername         string      `json:"email_username"`
	EmailUsernameEnabled  bool        `json:"email_username_enabled"`
	SentMessagesCount     int         `json:"sent_messages_count"`
	ForwardedMessagesCount int        `json:"forwarded_messages_count"`
	Used                  bool        `json:"used"`
	ForwardFromEmailAddress string    `json:"forward_from_email_address"`
	ProjectID             int         `json:"project_id"`
	Domain                string      `json:"domain"`
	POP3Domain            string      `json:"pop3_domain"`
	EmailDomain           string      `json:"email_domain"`
	SMTPPorts             []int       `json:"smtp_ports"`
	POP3Ports             []int       `json:"pop3_ports"`
	Permissions           []string    `json:"permissions"`
}

// InboxRequest represents a request to create/update an inbox
type InboxRequest struct {
	Inbox struct {
		Name          string `json:"name"`
		EmailUsername string `json:"email_username,omitempty"`
	} `json:"inbox"`
}

// SendingDomain represents a Mailtrap sending domain
type SendingDomain struct {
	ID               int        `json:"id"`
	Name             string     `json:"name"`
	CNAME            string     `json:"cname"`
	Status           string     `json:"status"`
	ComplianceStatus string     `json:"compliance_status"`
	DNSRecords       DNSRecords `json:"dns_records"`
	DNSStatus        DNSStatus  `json:"dns_status"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// SendingDomainRequest represents a request to create a sending domain
type SendingDomainRequest struct {
	SendingDomain struct {
		DomainName string `json:"domain_name"`
	} `json:"sending_domain"`
}

// DNSRecords contains all DNS records for domain verification
type DNSRecords struct {
	CNAME []DNSRecord `json:"cname"`
	MX    []DNSRecord `json:"mx"`
	TXT   []DNSRecord `json:"txt"`
}

// DNSRecord represents a single DNS record
type DNSRecord struct {
	Priority   *int   `json:"priority,omitempty"`
	RecordType string `json:"record_type"`
	Hostname   string `json:"hostname"`
	Value      string `json:"value"`
	Status     string `json:"status"`
}

// DNSStatus represents the verification status of DNS records
type DNSStatus struct {
	CNAME bool `json:"cname"`
	MX    bool `json:"mx"`
	TXT   bool `json:"txt"`
}

// Account represents a Mailtrap account
type Account struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	AccessLevels []int  `json:"access_levels"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error   string      `json:"error,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Message string      `json:"message,omitempty"`
}
