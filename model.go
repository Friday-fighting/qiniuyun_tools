package qiniuyun_tools

type Client struct {
	AccessKey  string `json:"access_key"`
	SecretKey  string `json:"secret_key"`
	ExpTimeKey string `json:"exp_time_key"`
	ExpTime    int    `json:"exp_time"`
	Bucket     string `json:"bucket"`
	UrlPrefix  string `json:"url_prefix"`
}

type UnixTimeReq struct {
	OrderParameter string `json:"order_parameter"`
	Path           string `json:"path"`
}
