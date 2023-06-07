package util

import (
	"bytes"
	"encoding/json"
)

//序列化
func ToJson(obj interface{}) string {
	json, err := json.Marshal(obj)
	if err != nil {
		LogError("序列化失败" + err.Error())
		return ""
	}
	return string(json)
}

//返回序列化为对象
func FromJson(val string, obj interface{}) error {
	err := json.Unmarshal([]byte(val), &obj)
	if err != nil {
		return err
	}
	return nil
}

//返回序列化对象
func FromByteJson(val []byte, obj interface{}) error {
	err := json.Unmarshal(val, obj)
	if err != nil {
		return err
	}
	return nil
}

//ToJsonSetEscapeHTML 序列化(特殊字符不转义)
func ToJsonSetEscapeHTML(obj interface{}) string {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(obj)
	if err != nil {
		LogError("序列化失败" + err.Error())
		return ""
	}
	return buf.String()
}
