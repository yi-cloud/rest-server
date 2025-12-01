package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yi-cloud/rest-server/pkg/config"
	"io"
	"net/http"
	"strconv"
)

func EC2Tokens(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		GinResponseData(c, nil, errors.New("Error reading request body"), http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest(c.Request.Method, config.KeystoneOpt.EndPoint+"/ec2tokens", bytes.NewReader(body))
	if err != nil {
		GinResponseData(c, nil, err)
		return
	}

	req.ContentLength, err = strconv.ParseInt(c.Request.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		GinResponseData(c, nil, err)
		return
	}

	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	req.Header.Set("X-Forwarded-For", c.Request.RemoteAddr)
	req.Header.Set("X-Proxy-Server", "Go-Proxy")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		GinResponseData(c, nil, err, resp.StatusCode)
		return
	}

	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	switch {
	case (resp.StatusCode >= http.StatusOK) && (resp.StatusCode <= http.StatusNoContent):
		var respBody any
		if resp.StatusCode != http.StatusNoContent {
			if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
				respBody = resp.Body
			} else {
				loginResp := MakeLoginResponse(respBody)
				loginResp.Token = resp.Header.Get("X-Subject-Token")
				if config.Regenerate {
					loginResp.Token, err = RegenerateToken(loginResp.Token, loginResp.Roles)
					if err != nil {
						resp.StatusCode = http.StatusInternalServerError
					}
				}
				respBody = loginResp
			}
		} else {
			respBody = resp.Body
		}
		GinResponseData(c, respBody, nil, resp.StatusCode)
	default:
		var respBody any
		var data []byte
		var ret any
		err := json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			data, err = io.ReadAll(resp.Body)
			if err == nil {
				err = fmt.Errorf(string(data))
			} else {
				err = fmt.Errorf(resp.Status)
			}
		} else {
			ret = respBody
		}
		GinResponseData(c, ret, err, resp.StatusCode)
	}
}
