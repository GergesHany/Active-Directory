package ad

// Config holds Active Directory connection configuration
type Config struct {
	Server   string
	Port     int
	BaseDN   string // The base distinguished name for searches
	Username string
	Password string
}
