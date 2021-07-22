package nacos_service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chujieyang/commonops/ops/opslog"
	"io/ioutil"
	"net/http"
	"strings"
)

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 7/20/21 10:53 AM
 * @Desc:
 */

type nacosClient struct {
	endPoint string
	username string
	password string
	accessToken string
}

func (r *nacosClient) requestPost(url string, body string, headers map[string]string) (respData string, statusCode int, err error) {
	url = fmt.Sprintf("http://%s/%s?accessToken=%s", r.endPoint, url, r.accessToken)
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		opslog.Error().Println(err.Error())
		return
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		opslog.Error().Println(err.Error())
		return
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			opslog.Error().Println(err.Error())
		}
	}()
	statusCode = resp.StatusCode
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		opslog.Error().Println(err.Error())
		return
	}
	respData = string(respBytes)
	return
}

func (r *nacosClient) requestDelete(url string, params string, body string, headers map[string]string) (respData string, statusCode int, err error) {
	url = fmt.Sprintf("http://%s/%s?accessToken=%s&%s", r.endPoint, url, r.accessToken, params)
	fmt.Println(url)
	req, err := http.NewRequest("DELETE", url, strings.NewReader(body))
	if err != nil {
		opslog.Error().Println(err.Error())
		return
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		opslog.Error().Println(err.Error())
		return
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			opslog.Error().Println(err.Error())
		}
	}()
	statusCode = resp.StatusCode
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		opslog.Error().Println(err.Error())
		return
	}
	respData = string(respBytes)
	return
}

func (r *nacosClient) requestGet(url string, params string, headers map[string]string) (respData string, respHeader map[string][]string, statusCode int, err error) {
	url = fmt.Sprintf("http://%s/%s?accessToken=%s&%s", r.endPoint, url, r.accessToken, params)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		opslog.Error().Println(err.Error())
		return
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		opslog.Error().Println(err.Error())
		return
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			opslog.Error().Println(err.Error())
		}
	}()
	respHeader = resp.Header.Clone()
	statusCode = resp.StatusCode
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		opslog.Error().Println(err.Error())
		return
	}
	respData = string(respBytes)
	return
}

type LoginResp struct {
	AccessToken string `json:"accessToken"`
}

func NewNacosClient(endpoint, username, password string) (nacos *nacosClient, err error) {
	nacos = &nacosClient{
		endPoint: endpoint,
		username: username,
		password: password,
		accessToken: "",
	}
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	respData, statusCode, err := nacos.requestPost("nacos/v1/auth/login", fmt.Sprintf("username=%s&password=%s", username, password), headers)
	if err != nil {
		return
	}
	if statusCode != 200 {
		err = errors.New("操作失败")
		return
	}
	var loginResult LoginResp
	if err = json.Unmarshal([]byte(respData), &loginResult); err != nil {
		return
	}
	nacos.accessToken = loginResult.AccessToken
	return
}

type namespaceItem struct {
	Namespace string
	NamespaceShowName string
	Quota int
	ConfigCount int
	Type int
}

type namespaceResp struct {
	Data []namespaceItem `json:"data"`
}

func (r *nacosClient) GetNamespace() (namespaceList namespaceResp, err error) {
	data, _, statusCode, err := r.requestGet("nacos/v1/console/namespaces", "", nil)
	if statusCode != 200 {
		err = errors.New(fmt.Sprintf("%s \n %s", err, data))
		return
	}
	if err = json.Unmarshal([]byte(data), &namespaceList); err != nil {
		return
	}
	return
}

func (r *nacosClient) GetConfig(namespace, dataId, group string) (data, configType string, err error) {
	data, respHeader, statusCode, err := r.requestGet("nacos/v1/cs/configs",
		fmt.Sprintf("tenant=%s&dataId=%s&group=%s", namespace, dataId, group), nil)
	if statusCode != 200 {
		err = errors.New(fmt.Sprintf("%s \n %s", err, data))
		return
	}
	configType = respHeader["Config-Type"][0]
	return
}

func (r *nacosClient) PublishConfig(namespace, dataId, group, content, configType string) (err error) {
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	body := fmt.Sprintf("tenant=%s&dataId=%s&group=%s&content=%s&type=%s", namespace, dataId, group, content, configType)
	data, statusCode, err := r.requestPost("nacos/v1/cs/configs", body, headers)
	if statusCode != 200 {
		err = errors.New(fmt.Sprintf("%s \n %s", err, data))
		return
	}
	return
}

func (r *nacosClient) CopyConfig(srcNamespace, srcDataId, srcGroup, dstNamespace, dstDataId, dstGroup string) (err error) {
	srcConfig, srcConfigType, err := r.GetConfig(srcNamespace, srcDataId, srcGroup)
	if err != nil {
		return
	}
	if err = r.PublishConfig(dstNamespace, dstDataId, dstGroup, srcConfig, srcConfigType); err != nil {
		return
	}
	return
}

func (r *nacosClient) GetNsConfigs(namespace string, page, size int) (data string, err error) {
	data, _, statusCode, err := r.requestGet("nacos/v1/cs/configs",
		fmt.Sprintf("pageNo=%d&pageSize=%d&search=accurate&dataId=&group=&tenant=%s", page, size, namespace), nil)
	if statusCode != 200 {
		err = errors.New(fmt.Sprintf("%s \n %s", err, data))
		return
	}
	return
}

func (r *nacosClient) DeleteConfig(namespace, dataId, group string) (err error) {
	params := fmt.Sprintf("tenant=%s&dataId=%s&group=%s", namespace, dataId, group)
	data, statusCode, err := r.requestDelete("nacos/v1/cs/configs", params, "", nil)
	if statusCode != 200 {
		err = errors.New(fmt.Sprintf("%s \n %s", err, data))
		return
	}
	return
}
