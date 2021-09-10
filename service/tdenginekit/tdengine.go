package tdenginekit

import (
	"encoding/base64"
	"fmt"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/httpkit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"net/http"
)

func Query(sql string) []map[string]interface{} {
	return queryResult(sql)
}

func url() string {
	url := configkit.GetStringD(ConfigKeyTDEngineURL)
	if url == "" {
		panic(exception.New("tdengine url is null"))
	}
	queryUrl := fmt.Sprintf("%s/rest/sql", url)
	return queryUrl
}

func auth() string {
	user := configkit.GetStringD(ConfigKeyTDEngineUser)
	if user == "" {
		panic(exception.New("tdengine user is null"))
	}
	pwd := configkit.GetStringD(ConfigKeyTDEnginePwd)
	if pwd == "" {
		panic(exception.New("tdengine pwd is null"))
	}
	return encodeBase64([]byte(fmt.Sprintf("%s:%s", user, pwd)))
}

func encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func queryResult(sql string) []map[string]interface{} {
	ret, code := httpkit.Request(httpkit.Req{
		Method: http.MethodPost,
		Url:    url(),
		Header: map[string]string{
			"Authorization": fmt.Sprintf("Basic %s", auth()),
		},
		BinaryData: []byte(sql),
	})
	status := gjson.Get(ret, "status").String()
	if status != "succ" {
		panic(exception.New("tdengine query error: " + gjson.Get(ret, "desc").String()))
	}
	if code > 300 {
		panic(exception.New("tdengine query error: " + cast.ToString(code)))
	}
	rows := gjson.Get(ret, "rows").Int()
	if rows > 0 {
		columns := gjson.Get(ret, "column_meta").Array()
		values := gjson.Get(ret, "data").Array()
		if len(columns) > 0 && len(values) > 0 {
			list := make([]map[string]interface{}, len(values))
			for i, v := range values {
				e := map[string]interface{}{}
				for ii, vv := range columns {
					e[vv.Array()[0].String()] = v.Array()[ii].Value()
				}
				list[i] = e
			}
			return list
		}
	}
	return []map[string]interface{}{}
}
