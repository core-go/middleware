package middleware

type LogConfig struct {
	Separate       bool              `mapstructure:"separate"`
	Build          bool              `mapstructure:"build"`
	Log            bool              `mapstructure:"log"`
	Skips          string            `mapstructure:"skips"`
	Ip             string            `mapstructure:"ip"`
	Duration       string            `mapstructure:"duration"`
	Uri            string            `mapstructure:"uri"`
	Body           string            `mapstructure:"body"`
	Size           string            `mapstructure:"size"`
	ReqId          string            `mapstructure:"req_id"`
	Scheme         string            `mapstructure:"scheme"`
	Proto          string            `mapstructure:"proto"`
	Method         string            `mapstructure:"method"`
	RemoteAddr     string            `mapstructure:"remote_addr"`
	RemoteIp       string            `mapstructure:"remote_ip"`
	UserAgent      string            `mapstructure:"user_agent"`
	ResponseStatus string            `mapstructure:"status"`
	Request        string            `mapstructure:"request"`
	Response       string            `mapstructure:"response"`
	Fields         string            `mapstructure:"fields"`
	Masks          string            `mapstructure:"masks"`
	Map            map[string]string `mapstructure:"map"`
	Constants      map[string]string `mapstructure:"constants"`
}

type FieldConfig struct {
	Log       bool              `mapstructure:"log"`
	Ip        string            `mapstructure:"ip"`
	Map       map[string]string `mapstructure:"map"`
	Constants map[string]string `mapstructure:"constants"`
	Duration  string            `mapstructure:"duration"`
	Fields    []string          `mapstructure:"fields"`
	Masks     []string          `mapstructure:"masks"`
	Skips     []string          `mapstructure:"skips"`
}
