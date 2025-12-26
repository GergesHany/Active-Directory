package main

import (
	"ad-example/pkg/ad"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

var once sync.Once

func getEnv(key, defaultValue string) string {
	once.Do(func() {
		// load the .env file here only once
		godotenv.Load()
	})

	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	port := 636
	if portStr := getEnv("AD_PORT", ""); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	config := &ad.Config{
		Server:   getEnv("AD_SERVER", "localhost"),
		Port:     port,
		BaseDN:   getEnv("AD_BASE_DN", "DC=EXAMPLE,DC=COM"),
		Username: getEnv("AD_USERNAME", "CN=Administrator,CN=Users,DC=EXAMPLE,DC=COM"),
		Password: getEnv("AD_PASSWORD", "P@ssw0rd123!"),
	}

	client, err := ad.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create AD client: %v", err)
	}
	defer client.Close()

	// Interactive CLI loop
	for {
		fmt.Println("\nAvailable commands: get-user, get-group, auth-user, update-user, add-user-to-group, remove-user-from-group, search-users, exit")
		fmt.Print("Enter command: ")
		var input string
		fmt.Scanln(&input)

		switch input {
		case "get-user":
			var username string
			fmt.Print("Username: ")
			fmt.Scanln(&username)
			if username == "" {
				fmt.Println("Username is required")
				continue
			}
			user, err := client.GetUserDetails(username)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			fmt.Printf("User: %s\nEmail: %s\nUsername: %s\nGroups: %v\n",
				user.GetAttributeValue("displayName"),
				user.GetAttributeValue("mail"),
				user.GetAttributeValue("sAMAccountName"),
				user.GetAttributeValues("memberOf"))
		case "get-group":
			var group string
			fmt.Print("Group name: ")
			fmt.Scanln(&group)
			if group == "" {
				fmt.Println("Group name is required")
				continue
			}
			members, err := client.GetGroupMembers(group)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			fmt.Printf("Group '%s' has %d members:\n", group, len(members))
			for _, member := range members {
				fmt.Printf("  - %s (%s)\n",
					member.GetAttributeValue("displayName"),
					member.GetAttributeValue("sAMAccountName"))
			}
		case "auth-user":
			var username, password string
			fmt.Print("Username: ")
			fmt.Scanln(&username)
			fmt.Print("Password: ")
			fmt.Scanln(&password)
			if username == "" || password == "" {
				fmt.Println("Username and password are required")
				continue
			}
			ok, err := client.AuthenticateUser(username, password)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			if ok {
				fmt.Println("Authenticated successfully!")
			} else {
				fmt.Println("Authentication failed!")
			}
		case "update-user":
			var username, attr, value string
			fmt.Print("Username: ")
			fmt.Scanln(&username)
			fmt.Print("Attribute: ")
			fmt.Scanln(&attr)
			fmt.Print("New value: ")
			fmt.Scanln(&value)
			if username == "" || attr == "" || value == "" {
				fmt.Println("Username, attribute, and value are required")
				continue
			}
			err := client.UpdateUserAttribute(username, attr, value)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			fmt.Println("User attribute updated successfully!")
		case "add-user-to-group":
			var username, group string
			fmt.Print("Username: ")
			fmt.Scanln(&username)
			fmt.Print("Group name: ")
			fmt.Scanln(&group)
			if username == "" || group == "" {
				fmt.Println("Username and group are required")
				continue
			}
			err := client.AddUserToGroup(username, group)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			fmt.Println("User added to group successfully!")
		case "remove-user-from-group":
			var username, group string
			fmt.Print("Username: ")
			fmt.Scanln(&username)
			fmt.Print("Group name: ")
			fmt.Scanln(&group)
			if username == "" || group == "" {
				fmt.Println("Username and group are required")
				continue
			}
			err := client.RemoveUserFromGroup(username, group)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			fmt.Println("User removed from group successfully!")
		case "search-users":
			var filter string
			fmt.Print("LDAP filter (default: (&(objectClass=user)(objectCategory=person))): ")
			fmt.Scanln(&filter)
			if filter == "" {
				filter = "(&(objectClass=user)(objectCategory=person))"
			}
			users, err := client.SearchUsers(filter)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			fmt.Printf("Found %d users:\n", len(users))
			for i, user := range users {
				if i >= 5 {
					fmt.Println("... (showing first 5)")
					break
				}
				fmt.Printf("  - %s (%s)\n",
					user.GetAttributeValue("displayName"),
					user.GetAttributeValue("sAMAccountName"))
			}
		case "exit":
			fmt.Println("Exiting CLI.")
			return
		default:
			fmt.Println("Unknown command.")
		}
	}
}
