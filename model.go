package qiniuyun_tools

type QiNiuClient struct {
	AccessKey       string `json:"access_key"`
	SecretKey       string `json:"secret_key"`
	QiNiuExpTimeKey string `json:"qiniu_exptime_key"`
	QiNiuExpTime    int    `json:"qiniu_exp_time"`
	Bucket          string `json:"bucket"`
	QiNiuUrlPrefix  string `json:"qiniu_url_prefix"`
}

type QiNiuUnixTimeReq struct {
	OrderParameter string `json:"order_parameter"`
	Path           string `json:"path"`
}
