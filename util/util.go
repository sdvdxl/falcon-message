package util

import (
	"bytes"
	"encoding/json"
	"log"
	"strconv"
	"strings"
)

// EncodeJSON json序列化(禁止 html 符号转义)
func EncodeJSON(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

//StringToInt string 类型转 int
func StringToInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("agent 类型转换失败, 请检查配置文件中 agentid 配置是否为纯数字(%v)", err)
		return 0
	}
	return n
}

// HandleContent [P2][PROBLEM][10-13-33-153][][测试 all(#1) net.port.listen port=2 0==0][O3 2017-06-06 16:46:00]
func HandleContent(content string) string {
	content = strings.Replace(content, "][", "\n", -1)
	if content[0] == '[' {
		content = content[1:]
	}

	if content[len(content)-1] == ']' {
		content = content[:len(content)-1]
	}

	return content
}
