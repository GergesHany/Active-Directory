package ad

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/go-ldap/ldap/v3"
)

// Client represents an Active Directory client
type Client struct {
	conn   *ldap.Conn
	config *Config
}

// NewClient creates a new Active Directory client
func NewClient(config *Config) (*Client, error) {
	var conn *ldap.Conn
	var err error

	// TLS configuration for secure connections
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Skip certificate verification for self-signed certs
	}

	// Connect to AD server with TLS
	if config.UseTLS {
		// Use LDAPS (LDAP over TLS) - typically port 636
		conn, err = ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", config.Server, config.Port), tlsConfig)
	} else {
		// Plain LDAP connection - not recommended for production
		conn, err = ldap.Dial("tcp", fmt.Sprintf("%s:%d", config.Server, config.Port))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	// Set connection timeout
	conn.SetTimeout(10 * time.Second)

	// Bind with credentials
	err = conn.Bind(config.Username, config.Password)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to bind: %w", err)
	}

	return &Client{
		conn:   conn,
		config: config,
	}, nil
}

// Close closes the AD connection
func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
