package oss

import (
	"bytes"
	"fmt"
	"path"
	"runtu666/common/shared/utils"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/globalsign/mgo/bson"
	"github.com/tal-tech/go-zero/core/logx"
)

type Client struct {
	AccessKeyId     string
	AccessKeySecret string
	RoleArn         string
	RoleSessionName string
	Endpoint        string
	BucketName      string
	Domain          string
	Mode            string
}

func NewOssClient(c *Client) *Client {
	return c
}

func (c *Client) GetClient() (*oss.Client, error) {
	client, err := oss.New(c.Endpoint, c.AccessKeyId, c.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	return client, err
}

func (c *Client) GetBucket() (*oss.Bucket, error) {
	client, err := c.GetClient()
	if err != nil {
		return nil, err
	}
	bucket, err := client.Bucket(c.BucketName)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

func (c *Client) Upload(b []byte, filename string) (string, error) {
	var url string
	bucket, err := c.GetBucket()
	if err != nil {
		logx.Error(err)
		return "", err
	}
	ext := path.Ext(path.Base(filename))
	fileName := fmt.Sprintf("brush/%s/%s/%s%s", c.Mode, time.Now().Format(utils.DateRawYMD), bson.NewObjectId().Hex(), ext)
	err = bucket.PutObject(fileName, bytes.NewBuffer(b))
	if err != nil {
		logx.Error(err)
		return "", err
	}
	url = fmt.Sprintf("https://%s/%s", c.Domain, fileName)
	return url, nil
}
