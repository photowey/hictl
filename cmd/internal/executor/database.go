package executor

// Hictl the config of hictl cmd
type Hictl struct {
	DatabaseMap map[string]Database `json:"databases" yaml:"databases" toml:"databases"` // database map
}

type Database struct {
	Database string `json:"database" yaml:"database" toml:"database"` // database name
	Url      string `json:"url" yaml:"url" toml:"url"`                // database url
	Port     string `json:"port" yaml:"port" toml:"port"`             // database port
	UserName string `json:"username" yaml:"username" toml:"username"` // database user name
	Password string `json:"password" yaml:"password" toml:"password"` // database database password
}
