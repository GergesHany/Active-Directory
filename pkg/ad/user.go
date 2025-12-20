package ad

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
)

// AuthenticateUser authenticates a user against Active Directory
func (c *Client) AuthenticateUser(username, password string) (bool, error) {
	// Search for the user first
	userDN, err := c.searchUserDN(username)
	if err != nil {
		return false, fmt.Errorf("user not found: %w", err)
	}

	// Try to bind with user credentials
	err = c.conn.Bind(userDN, password)
	if err != nil {
		return false, nil
	}

	// Re-bind with admin credentials for further operations
	err = c.conn.Bind(c.config.Username, c.config.Password)
	if err != nil {
		return false, fmt.Errorf("failed to re-bind: %w", err)
	}

	return true, nil
}

// searchUserDN finds the distinguished name for a user
func (c *Client) searchUserDN(username string) (string, error) {
	searchRequest := ldap.NewSearchRequest(
		c.config.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(sAMAccountName=%s)", username),
		[]string{"dn"},
		nil,
	)

	sr, err := c.conn.Search(searchRequest)
	if err != nil {
		return "", err
	}

	if len(sr.Entries) != 1 {
		return "", fmt.Errorf("user not found or too many entries returned")
	}

	return sr.Entries[0].DN, nil
}

// SearchUsers searches for users in Active Directory
func (c *Client) SearchUsers(filter string) ([]*ldap.Entry, error) {
	searchRequest := ldap.NewSearchRequest(
		c.config.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter,
		[]string{"dn", "cn", "sAMAccountName", "mail", "displayName", "memberOf"},
		nil,
	)

	sr, err := c.conn.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("search error: %w", err)
	}

	return sr.Entries, nil
}

// GetUserDetails retrieves detailed information about a user
func (c *Client) GetUserDetails(username string) (*ldap.Entry, error) {
	filter := fmt.Sprintf("(sAMAccountName=%s)", username)

	searchRequest := ldap.NewSearchRequest(
		c.config.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter,
		[]string{"*"}, // Request all attributes
		nil,
	)

	sr, err := c.conn.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("search error: %w", err)
	}

	if len(sr.Entries) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return sr.Entries[0], nil
}

// UpdateUserAttribute updates a user attribute
func (c *Client) UpdateUserAttribute(username, attribute, value string) error {
	userDN, err := c.searchUserDN(username)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	modifyRequest := ldap.NewModifyRequest(userDN, nil)
	modifyRequest.Replace(attribute, []string{value})

	err = c.conn.Modify(modifyRequest)
	if err != nil {
		return fmt.Errorf("failed to update user attribute: %w", err)
	}

	return nil
}
