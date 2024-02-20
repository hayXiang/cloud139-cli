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
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/tickstep/cloudpan189-api/cloudpan"
	"github.com/tickstep/cloudpan189-api/cloudpan/apiutil"
)

func PostHeader(url string, msg []byte, headers map[string]string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(msg)))
	if err != nil {
		return "", err
	}
	for key, header := range headers {
		req.Header.Set(key, header)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func main() {
	var command = ""
	flag.StringVar(&command, "c", "", "command type")
	var strInfo = ""
	flag.StringVar(&strInfo, "d", "", "sessionKey and sessionSecret")
	var familyId = ""
	flag.StringVar(&familyId, "f", "", "famaily id")
	flag.Parse()

	if command == "login" {
		// do login
		appToken, e := cloudpan.AppLogin(os.Args[3], os.Args[4])
		if e != nil {
			fmt.Println(e)
			return
		}
		fmt.Printf("{\"sessionKey\":\"%s\", \"sessionSecret\" :\"%s\"}", appToken.FamilySessionKey, appToken.SessionSecret)
	} else if command == "clear" {
		var jsonMap map[string]string
		err := json.Unmarshal([]byte(strInfo), &jsonMap)
		if err != nil {
			fmt.Println(err)
			return
		}
		sessionKey := jsonMap["sessionKey"]
		sessionSecret := jsonMap["sessionSecret"]
		httpMethod := "POST"
		dateOfGmt := apiutil.DateOfGmtStr()
		url := "https://api.cloud.189.cn/batch/createBatchTask.action?rand=" + apiutil.Rand() + "&clientType=TELEANDROID&model=M2012K11AC&version=10.1.3"
		signature := strings.ToLower(apiutil.SignatureOfHmac(sessionSecret, sessionKey, httpMethod, url, dateOfGmt))

		headers := map[string]string{
			"accept":       "application/json;charset=UTF-8",
			"user-agent":   "Ecloud/10.1.3 (M2012K11AC; ; xiaomi) Android/33",
			"x-request-id": apiutil.XRequestId(),
			"content-type": "multipart/form-data; boundary=4c1ed063-2af1-426a-947a-802d807a4700",
			"date":         dateOfGmt,
			"signature":    signature,
			"sessionkey":   sessionKey,
		}

		payload := fmt.Sprintf("--4c1ed063-2af1-426a-947a-802d807a4700\r\nContent-Disposition: form-data; name=\"familyId\"\r\nContent-Transfer-Encoding: binary\r\nContent-Type: multipart/form-data; charset=utf-8\r\nContent-Length: %d\r\n\r\n%s\r\n--4c1ed063-2af1-426a-947a-802d807a4700\r\nContent-Disposition: form-data; name=\"groupId\"\r\nContent-Transfer-Encoding: binary\r\nContent-Type: multipart/form-data; charset=utf-8\r\nContent-Length: 4\r\n\r\nnull\r\n--4c1ed063-2af1-426a-947a-802d807a4700\r\nContent-Disposition: form-data; name=\"targetFolderId\"\r\nContent-Transfer-Encoding: binary\r\nContent-Type: multipart/form-data; charset=utf-8\r\nContent-Length: 4\r\n\r\nnull\r\n--4c1ed063-2af1-426a-947a-802d807a4700\r\nContent-Disposition: form-data; name=\"shareId\"\r\nContent-Transfer-Encoding: binary\r\nContent-Type: multipart/form-data; charset=utf-8\r\nContent-Length: 4\r\n\r\nnull\r\n--4c1ed063-2af1-426a-947a-802d807a4700\r\nContent-Disposition: form-data; name=\"type\"\r\nContent-Transfer-Encoding: binary\r\nContent-Type: multipart/form-data; charset=utf-8\r\nContent-Length: 13\r\n\r\nEMPTY_RECYCLE\r\n--4c1ed063-2af1-426a-947a-802d807a4700--", len(familyId), familyId)
		body, err := PostHeader(url, []byte(payload), headers)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(body)
	}
}
