package qiniuyun_tools

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func GetQiNiuFileURLPath(urlStr string, needVersion bool) (res string, err error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	res = strings.TrimPrefix(parsedURL.Path, "/")
	if needVersion {
		var version string
		queryParams := parsedURL.Query()
		if queryParams.Get("version") == "" {
			version = strconv.Itoa(int(time.Now().Unix()))
			res = fmt.Sprintf("%s?version=%s", res, version)
		} else {
			newQuery := url.Values{}
			newQuery.Set("version", queryParams.Get("version"))
			res = res + "?" + newQuery.Encode()
		}
	}
	return res, nil
}
