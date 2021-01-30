/****************

Author: yuansudong
Date: 2018.06.23
*****************/

package hkit

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/http2"
)

// New 用于新建立一个客户端
func _New(
	iOpt *Option, // iOpt 一个HttpClient需要启动的参数
) (*Engine, error) {
	iOpt._RetryTime = time.Duration(iOpt.RetryInterval) * time.Second
	tr := &http.Transport{
		MaxIdleConns:        iOpt.MaxIdleConns,
		IdleConnTimeout:     time.Duration(iOpt.IdleConnTimeout) * time.Second,
		MaxIdleConnsPerHost: iOpt.MaxIdleConnsPerHost,
		DisableCompression:  iOpt.DisableCompression,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: iOpt.SkipTLSVerify},
	}
	if iOpt.EnableHTTP2 {
		http2.ConfigureTransport(tr)
	}
	if iOpt.EnableProxy {
		tr.Proxy = func(_ *http.Request) (*url.URL, error) {
			return url.Parse(iOpt.ProxyAddr)
		}
	}
	return &Engine{
		_Client: &http.Client{Transport: tr},
		_Opt:    iOpt,
	}, nil
}

// BaseHeader 设置baseheader
func (E *Engine) BaseHeader(iBaseHeader []HTTPOption) {
	E._BaseHeader = iBaseHeader
}

// _BuildHTTPRequest 新建立一个http的请求
func (R *Request) _BuildHTTPRequest() (
	*http.Request,
	error,
) {
	var body io.Reader
	if R.Body != nil {
		body = bytes.NewBuffer(R.Body)
	}
	req, err := http.NewRequest(R.Method, R.URL, body)
	if err != nil {
		return nil, err
	}
	for _, opt := range R.Header {
		opt(req)
	}
	return req, nil
}

// _DO DO方法下的实际调用者
func (E *Engine) _DO(
	iReq *Request,
) (*Response, error) {
	request, err := iReq._BuildHTTPRequest()
	if err != nil {
		return nil, err
	}
	resp, err := E._Client.Do(request)

	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	return &Response{
		Status: resp.StatusCode,
		Header: resp.Header,
		Body:   body,
	}, nil
}

// DO 用于发起一个自定义请求的接口
func (E *Engine) DO(
	iCtx context.Context, // 浏览器接口
	iReq *Request, //  请求
) (*Response, error) {
	var (
		result *Response
		err    error
	)
	for retryTimes := 0; retryTimes < E._Opt.RetryCount; retryTimes++ {
		result, err = E._DO(iReq)
		if err == nil {
			break
		}
		if !E._PendingForRetry(iCtx) {
			break
		}
		// 指数退避
		time.Sleep(1 << retryTimes)
	}
	return result, err
}

// _PendingForRetry 重试接口的调用
func (E *Engine) _PendingForRetry(
	ctx context.Context,
) bool {
	if E._Opt.RetryInterval > 0 {
		select {
		case <-ctx.Done():
			return false
		case <-time.After(E._Opt._RetryTime):
			return true
		}
	}
	return false
}

// POST  HTTP POST类消息
func (E *Engine) POST(
	iCtx context.Context,
	iURL string, //  请求的URL
	iArgs map[string]string, //  请求参数
	iHeader ...HTTPOption, //  附加头部
) (
	[]byte,
	error,
) {

	mParam := new(url.Values)
	for key, val := range iArgs {
		mParam.Add(key, val)
	}
	hRsp, err := E.DO(iCtx, &Request{
		Method: http.MethodPost,
		URL:    iURL,
		Body:   []byte(mParam.Encode()),
		Header: _MergeHeader(iHeader, E._BaseHeader),
	})
	if err != nil {
		return nil, err
	}
	if hRsp.Status == http.StatusOK {
		return hRsp.Body, nil
	}
	return nil, fmt.Errorf(_LayoutError, hRsp.Status, string(hRsp.Body))
}

// JSON JSON类消息
func (E *Engine) JSON(
	iCtx context.Context,
	iURL string, //  请求的URL
	iObject interface{}, // 序列化的对象
	iHeader ...HTTPOption, //  附加头部
) (
	[]byte,
	error,
) {
	body, err := json.Marshal(iObject)
	if err != nil {
		return nil, err
	}
	hRsp, err := E.DO(iCtx, &Request{
		Method: http.MethodPost,
		URL:    iURL,
		Body:   body,
		Header: _MergeHeader(
			append(
				iHeader,
				SetHeader(
					_HeaderKeyContentType,
					_HeaderValueApplicationJSON,
				)),
			E._BaseHeader,
		),
	})
	if err != nil {
		return nil, err
	}
	if hRsp.Status == http.StatusOK {
		return hRsp.Body, nil
	}
	return nil, fmt.Errorf(_LayoutError, hRsp.Status, string(hRsp.Body))
}

// GET 用于发送get请求
func (E *Engine) GET(
	iCtx context.Context,
	iURL string, // 请求的URL
	iParams map[string]string, //  请求参数
	iHeader ...HTTPOption, // 需要添加的头部
) (
	[]byte,
	error,
) {
	mParam := new(url.Values)
	for key, val := range iParams {
		mParam.Add(key, val)
	}
	hRsp, err := E.DO(iCtx, &Request{
		Method: http.MethodGet,
		URL:    iURL + "?" + mParam.Encode(),
		Body:   nil,
		Header: _MergeHeader(iHeader, E._BaseHeader),
	})
	if err != nil {
		return nil, err
	}
	if hRsp.Status == http.StatusOK {
		return hRsp.Body, nil
	}
	return nil, fmt.Errorf(_LayoutError, hRsp.Status, string(hRsp.Body))
}
