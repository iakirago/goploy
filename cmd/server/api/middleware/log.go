// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package middleware

import (
	"encoding/json"
	"github.com/zhenorzz/goploy/internal/log"
	"github.com/zhenorzz/goploy/internal/server"
	"github.com/zhenorzz/goploy/internal/server/response"
	"github.com/zhenorzz/goploy/model"
	"strconv"
	"time"
)

func AddLoginLog(gp *server.Goploy, resp server.Response) {
	respJson := resp.(response.JSON)
	account := ""
	if respJson.Code != response.IllegalParam {
		type ReqData struct {
			Account string `json:"account"`
		}
		var reqData ReqData
		_ = json.Unmarshal(gp.Body, &reqData)
		account = reqData.Account
	}

	err := model.LoginLog{
		Account:    account,
		RemoteAddr: gp.Request.RemoteAddr,
		UserAgent:  gp.Request.UserAgent(),
		Referer:    gp.Request.Referer(),
		Reason:     respJson.Message,
		LoginTime:  time.Now().Format("20060102150405"),
	}.AddRow()
	if err != nil {
		log.Error(err.Error())
	}
}

func AddOPLog(gp *server.Goploy, resp server.Response) {
	requestData := ""
	if gp.URLQuery != nil && len(gp.URLQuery) > 0 && len(gp.Body) > 2 {
		requestData = "{" + string(gp.Body[1:len(gp.Body)-1]) + ", \"_query\":\"" + gp.URLQuery.Encode() + "\"}"
	} else if gp.URLQuery != nil && len(gp.URLQuery) > 0 {
		requestData = "{\"_query\":\"" + gp.URLQuery.Encode() + "\"}"
	} else if len(gp.Body) > 2 {
		requestData = string(gp.Body)
	}
	responseData := ""
	switch resp.(type) {
	case response.JSON:
		jsonBytes, err := json.Marshal(resp)
		if err == nil {
			responseData = string(jsonBytes)
		}
	}
	err := model.OperationLog{
		NamespaceID:  gp.Namespace.ID,
		UserID:       gp.UserInfo.ID,
		Router:       gp.Request.Header.Get("Router"),
		API:          gp.Request.URL.Path,
		RequestTime:  gp.Request.Header.Get("_time"),
		RequestData:  requestData,
		ResponseTime: time.Now().Format("20060102150405"),
		ResponseData: responseData,
	}.AddRow()
	if err != nil {
		log.Error(err.Error())
	}
}

func AddUploadLog(gp *server.Goploy, resp server.Response) {
	var serverID int64 = 0
	var path = ""
	respJson := resp.(response.JSON)
	if respJson.Code != response.IllegalParam {
		serverID, _ = strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
		file, fileHandler, _ := gp.Request.FormFile("file")
		path = gp.URLQuery.Get("filePath") + "/" + fileHandler.Filename
		_ = file.Close()
	}

	err := model.SftpLog{
		NamespaceID: gp.Namespace.ID,
		UserID:      gp.UserInfo.ID,
		ServerID:    serverID,
		RemoteAddr:  gp.Request.RemoteAddr,
		UserAgent:   gp.Request.UserAgent(),
		Type:        model.SftpLogTypeUpload,
		Path:        path,
		Reason:      respJson.Message,
	}.AddRow()
	if err != nil {
		log.Error(err.Error())
	}
}

func AddDeleteLog(gp *server.Goploy, resp server.Response) {
	var serverID int64 = 0
	var path = ""
	respJson := resp.(response.JSON)
	if respJson.Code != response.IllegalParam {
		type ReqData struct {
			ServerID int64  `json:"account"`
			File     string `json:"file"`
		}
		var reqData ReqData
		_ = json.Unmarshal(gp.Body, &reqData)
		serverID = reqData.ServerID
		path = reqData.File
	}

	err := model.SftpLog{
		NamespaceID: gp.Namespace.ID,
		UserID:      gp.UserInfo.ID,
		ServerID:    serverID,
		RemoteAddr:  gp.Request.RemoteAddr,
		UserAgent:   gp.Request.UserAgent(),
		Type:        model.SftpLogTypeDelete,
		Path:        path,
		Reason:      respJson.Message,
	}.AddRow()
	if err != nil {
		log.Error(err.Error())
	}
}

func AddDownloadLog(gp *server.Goploy, resp server.Response) {
	msg := ""
	path := ""
	var serverID int64 = 0
	switch resp.(type) {
	case response.JSON:
		respJson := resp.(response.JSON)
		if respJson.Code != response.IllegalParam {
			msg = respJson.Message
			serverID, _ = strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
		}
	case response.SftpFile:
		path = resp.(response.SftpFile).Filename
	}

	err := model.SftpLog{
		NamespaceID: gp.Namespace.ID,
		UserID:      gp.UserInfo.ID,
		ServerID:    serverID,
		RemoteAddr:  gp.Request.RemoteAddr,
		UserAgent:   gp.Request.UserAgent(),
		Type:        model.SftpLogTypeDownload,
		Path:        path,
		Reason:      msg,
	}.AddRow()
	if err != nil {
		log.Error(err.Error())
	}
}

func AddPreviewLog(gp *server.Goploy, resp server.Response) {
	msg := ""
	path := ""
	var serverID int64 = 0
	switch resp.(type) {
	case response.JSON:
		respJson := resp.(response.JSON)
		if respJson.Code != response.IllegalParam {
			msg = respJson.Message
			serverID, _ = strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
		}
	case response.SftpFile:
		path = resp.(response.SftpFile).Filename
	}

	err := model.SftpLog{
		NamespaceID: gp.Namespace.ID,
		UserID:      gp.UserInfo.ID,
		ServerID:    serverID,
		RemoteAddr:  gp.Request.RemoteAddr,
		UserAgent:   gp.Request.UserAgent(),
		Type:        model.SftpLogTypePreview,
		Path:        path,
		Reason:      msg,
	}.AddRow()
	if err != nil {
		log.Error(err.Error())
	}
}
