package log

type ChiLogConfig struct {
	Single         bool               `mapstructure:"single"`
	Build          bool               `mapstructure:"build"`
	Duration       string             `mapstructure:"duration"`
	Uri            string             `mapstructure:"uri"`
	Body           string             `mapstructure:"body"`
	Size           string             `mapstructure:"size"`
	ReqId          string             `mapstructure:"req_id"`
	Scheme         string             `mapstructure:"scheme"`
	Proto          string             `mapstructure:"proto"`
	Method         string             `mapstructure:"method"`
	RemoteAddr     string             `mapstructure:"remote_addr"`
	RemoteIp       string             `mapstructure:"remote_ip"`
	UserAgent      string             `mapstructure:"user_agent"`
	ResponseStatus string             `mapstructure:"status"`
	Request        string             `mapstructure:"request"`
	Response       string             `mapstructure:"response"`
	FieldMap       string             `mapstructure:"field_map"`
	Fields         string             `mapstructure:"fields"`
	Map            *map[string]string `mapstructure:"map"`
}

type FieldConfig struct {
	Map      *map[string]string `mapstructure:"map"`
	FieldMap string             `mapstructure:"field_map"`
	Duration string             `mapstructure:"duration"`
	Fields   *[]string          `mapstructure:"fields"`
}
