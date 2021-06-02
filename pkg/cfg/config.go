package cfg

type MongoConfig struct {
	Host        string     `yaml:"host" envconfig:"MONGO_HOST"`
	User        string     `yaml:"user" envconfig:"MONGO_USER"`
	Password    SafeString `yaml:"password" envconfig:"MONGO_PASSWORD"`
	Database    string     `yaml:"database" envconfig:"MONGO_DATABASE"`
	Collections struct {
		Users string `yaml:"users" envconfig:"MONGO_COLLECTIONS_USERS"`
	} `yaml:"collections"`
}

type SessionConfig struct {
	Cookie  string `yaml:"cookie" envconfig:"SESSION_COOKIE"`
	Expires int    `yaml:"expires" envconfig:"SESSION_EXPIRES"`
}

type MaponConfig struct {
	URL       string     `yaml:"url" envconfig:"MAPON_URL"`
	Key       SafeString `yaml:"key" envconfig:"MAPON_KEY"`
	Endpoints struct {
		Unit  string `yaml:"unit" envconfig:"MAPON_ENDPOINTS_UNIT"`
		Route string `yaml:"route" envconfig:"MAPON_ENDPOINTS_ROUTE"`
	} `yaml:"endpoints"`
}

// SafeString wraps string to implement fmt.Stringer and fmt.GoStringer funcs
// that return a constant string instead of the one stored. It is intended to
// avoid accidental value leak through fmt.
type SafeString string

// String returns a constant string representing the "native" format of a
// string.
func (SafeString) String() string {
	return "[REDACTED]"
}

// GoString returns a constant string representing the "Go-syntax"
// representation of a string.
func (SafeString) GoString() string {
	return `"[REDACTED]"`
}
