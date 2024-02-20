// Copyright (c) 2020 tickstep.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/tickstep/cloudpan189-api/cloudpan"
	"github.com/tickstep/cloudpan189-api/cloudpan/apiutil"
	"github.com/tickstep/library-go/jsonhelper"
)

type (
	userpw struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
)

// SignatureOfHmacV2 HMAC签名
func SignatureOfHmacV2(secretKey, sessionKey, operate, url, dateOfGmt, param string) string {
	requestUri := strings.Split(url, "?")[0]
	requestUri = strings.ReplaceAll(requestUri, "https://", "")
	requestUri = strings.ReplaceAll(requestUri, "http://", "")
	idx := strings.Index(requestUri, "/")
	requestUri = requestUri[idx:]

	plainStr := &strings.Builder{}
	fmt.Fprintf(plainStr, "SessionKey=%s&Operate=%s&RequestURI=%s&Date=%s",
		sessionKey, operate, requestUri, dateOfGmt)
	if param != "" {
		plainStr.WriteString(fmt.Sprintf("&params=%s", param))
	}
	key := []byte(secretKey)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(plainStr.String()))
	return strings.ToUpper(hex.EncodeToString(mac.Sum(nil)))
}

// SignatureOfHmac HMAC签名
func SignatureOfHmac(secretKey, sessionKey, operate, url, dateOfGmt string) string {
	return SignatureOfHmacV2(secretKey, sessionKey, operate, url, dateOfGmt, "")
}

func main() {
	configFile, err := os.OpenFile("userpw.txt", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		fmt.Println("read user info error")
		return
	}
	defer configFile.Close()

	userpw := &userpw{}
	err = jsonhelper.UnmarshalData(configFile, userpw)
	if err != nil {
		fmt.Println("read user info error")
		return
	}

	// do login
	appToken, e := cloudpan.AppLogin(userpw.UserName, userpw.Password)
	if e != nil {
		fmt.Println(e)
		return
	}

	sessionKey := appToken.FamilySessionKey
	sessionSecret := appToken.FamilySessionSecret
	httpMethod := "POST"
	dateOfGmt := apiutil.DateOfGmtStr()
	signature := strings.ToLower(SignatureOfHmac(sessionSecret, sessionKey, httpMethod, "https://api.cloud.189.cn/batch/createBatchTask.action?rand=1706541310636&clientType=TELEANDROID&model=M2012K11AC&version=10.1.3", dateOfGmt))

	fmt.Println(appToken.SessionKey)
	fmt.Println(signature)
	fmt.Println(dateOfGmt)
	fmt.Println("login success")
	return
}
