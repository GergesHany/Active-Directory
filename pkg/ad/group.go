package ad

import (
	"fmt"
	"log"

	"github.com/go-ldap/ldap/v3"
)

// GetGroupMembers retrieves members of a group
func (c *Client) GetGroupMembers(groupName string) ([]*ldap.Entry, error) {
	// First find the group DN
	groupFilter := fmt.Sprintf("(&(objectClass=group)(cn=%s))", groupName)

	groupSearch := ldap.NewSearchRequest(
		c.config.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		groupFilter,
		[]string{"member"},
		nil,
	)

	sr, err := c.conn.Search(groupSearch)
	if err != nil {
		return nil, fmt.Errorf("group search error: %w", err)
	}

	if len(sr.Entries) == 0 {
		return nil, fmt.Errorf("group not found")
	}

	members := sr.Entries[0].GetAttributeValues("member")
	var entries []*ldap.Entry

	// Get details for each member
	for _, memberDN := range members {
		memberSearch := ldap.NewSearchRequest(
			memberDN,
			ldap.ScopeBaseObject,
			ldap.NeverDerefAliases,
			0,
			0,
			false,
			"(objectClass=*)",
			[]string{"cn", "sAMAccountName", "mail", "displayName"},
			nil,
		)

		memberResult, err := c.conn.Search(memberSearch)
		if err != nil {
			log.Printf("Warning: could not fetch member %s: %v", memberDN, err)
			continue
		}

		if len(memberResult.Entries) > 0 {
			entries = append(entries, memberResult.Entries[0])
		}
	}

	return entries, nil
}

// AddUserToGroup adds a user to a group
func (c *Client) AddUserToGroup(username, groupName string) error {
	userDN, err := c.searchUserDN(username)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Find group DN
	groupFilter := fmt.Sprintf("(&(objectClass=group)(cn=%s))", groupName)
	groupSearch := ldap.NewSearchRequest(
		c.config.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		groupFilter,
		[]string{"dn"},
		nil,
	)

	sr, err := c.conn.Search(groupSearch)
	if err != nil {
		return fmt.Errorf("group search error: %w", err)
	}

	if len(sr.Entries) == 0 {
		return fmt.Errorf("group not found")
	}

	groupDN := sr.Entries[0].DN

	// Add user to group
	modifyRequest := ldap.NewModifyRequest(groupDN, nil)
	modifyRequest.Add("member", []string{userDN})

	err = c.conn.Modify(modifyRequest)
	if err != nil {
		return fmt.Errorf("failed to add user to group: %w", err)
	}

	return nil
}

// RemoveUserFromGroup removes a user from a group
func (c *Client) RemoveUserFromGroup(username, groupName string) error {
	userDN, err := c.searchUserDN(username)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Find group DN
	groupFilter := fmt.Sprintf("(&(objectClass=group)(cn=%s))", groupName)
	groupSearch := ldap.NewSearchRequest(
		c.config.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		groupFilter,
		[]string{"dn"},
		nil,
	)

	sr, err := c.conn.Search(groupSearch)
	if err != nil {
		return fmt.Errorf("group search error: %w", err)
	}

	if len(sr.Entries) == 0 {
		return fmt.Errorf("group not found")
	}

	groupDN := sr.Entries[0].DN

	// Remove user from group
	modifyRequest := ldap.NewModifyRequest(groupDN, nil)
	modifyRequest.Delete("member", []string{userDN})

	err = c.conn.Modify(modifyRequest)
	if err != nil {
		return fmt.Errorf("failed to remove user from group: %w", err)
	}

	return nil
}
