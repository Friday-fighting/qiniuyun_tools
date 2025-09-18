package qiniuyun_tools

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"errors"

	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
	"github.com/qiniu/go-sdk/v7/storagev2/uptoken"
)

func NewQiNiuClient(accessKey, secretKey, qiNiuUrlPrefix, bucket, qiNiuExpTimeKey string, qiNiuExpTime int) *QiNiuClient {
	if qiNiuExpTime <= 0 {
		qiNiuExpTime = 3600
	}
	return &QiNiuClient{
		AccessKey:       accessKey,
		SecretKey:       secretKey,
		QiNiuUrlPrefix:  qiNiuUrlPrefix,
		QiNiuExpTimeKey: qiNiuExpTimeKey,
		QiNiuExpTime:    qiNiuExpTime,
		Bucket:          bucket,
	}
}

func (client *QiNiuClient) SaveFileIntoQiNiu(filaName string, fileData []byte) (string, error) {
	mac := credentials.NewCredentials(client.AccessKey, client.SecretKey)
	putPolicy, err := uptoken.NewPutPolicyWithKey(client.Bucket, filaName, time.Now().Add(1*time.Hour))
	if err != nil {
		return "", err
	}
	fmt.Printf(">> upload %s\n", filaName)
	fileBaseName := filepath.Base(filaName)
	uploaderObjectOptions := &uploader.ObjectOptions{
		UpToken:    uptoken.NewSigner(putPolicy, mac),
		BucketName: client.Bucket,
		ObjectName: &filaName,
		FileName:   fileBaseName,
	}
	uploadManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: mac,
		},
	})
	if len(fileData) <= 1048576 { // 1MB
		fileReader := bytes.NewReader(fileData)
		err := uploadManager.UploadReader(context.Background(), fileReader, uploaderObjectOptions, nil)
		if err != nil {
			return "", err
		}
	} else {
		replaceFilaName := strings.ReplaceAll(filaName, "\\", "/")
		replaceFilaName = strings.ReplaceAll(replaceFilaName, "/", "_")
		// 如果文件大于1MB，保存到临时文件
		tmpFilename := filepath.Join(filepath.Dir(os.Args[0]), "tmp", "qiniu_file_tmp", replaceFilaName)
		if err := os.MkdirAll(filepath.Join(filepath.Dir(os.Args[0]), "tmp", "qiniu_file_tmp"), os.ModePerm); err != nil {
			return "", err
		}
		if err := os.WriteFile(tmpFilename, fileData, 0644); err != nil {
			return "", err
		}
		err := uploadManager.UploadFile(context.Background(), tmpFilename, uploaderObjectOptions, nil)
		if err != nil {
			return "", err
		}
		// 删除临时文件
		_ = os.Remove(tmpFilename)
	}

	return filaName, nil
}

func (client *QiNiuClient) QiNiuUrlUnixTime(req *QiNiuUnixTimeReq) string {
	path := req.Path
	orderParameter := req.OrderParameter
	if path == "" {
		return ""
	}
	if len(path) >= 4 && (path[:4] == "http") {
		return path
	}
	bakPath := path
	bakPath = strings.ReplaceAll(bakPath, "sign=", "sign_old=")
	if path[0] != '/' {
		path = "/" + path
	}
	if idx := strings.Index(path, "?"); idx != -1 {
		path = path[:idx]
	}
	expTimeHex := fmt.Sprintf("%x", time.Now().Unix()+int64(client.QiNiuExpTime))
	s := fmt.Sprintf("%s%s%s", client.QiNiuExpTimeKey, url.QueryEscape(path), expTimeHex)
	m := md5.New()
	m.Write([]byte(s))
	sign := hex.EncodeToString(m.Sum(nil))
	qiNiuUrl := client.QiNiuUrlPrefix
	if !strings.HasSuffix(qiNiuUrl, "/") {
		qiNiuUrl += "/"
	}
	if strings.Contains(bakPath, "?") {
		return fmt.Sprintf("%s%s&sign=%s&t=%s%s", qiNiuUrl, bakPath, sign, expTimeHex, orderParameter)
	} else {
		return fmt.Sprintf("%s%s?sign=%s&t=%s%s", qiNiuUrl, bakPath, sign, expTimeHex, orderParameter)
	}
}

func (client *QiNiuClient) GenQiNiuToken(key string) (token string, err error) {
	if key == "" {
		return "", errors.New("fileName must not be empty")
	}
	mac := credentials.NewCredentials(client.AccessKey, client.SecretKey)
	putPolicy, err := uptoken.NewPutPolicyWithKey(client.Bucket, key, time.Now().Add(1*time.Hour))
	if err != nil {
		return "", err
	}
	token, err = uptoken.NewSigner(putPolicy, mac).GetUpToken(context.Background())
	if err != nil {
		return "", err
	}
	return token, nil
}
