package util

import (
    "bytes"
    "encoding/json"
    "errors"
    "github.com/sdvdxl/falcon-message/config"
    "log"
    "regexp"
    "strconv"
    "strings"
    "time"
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

// HandleContent [P2][PROBLEM][10-13-33-153][][测试 all(#1) net.port.listen port=2][O3 2017-06-06 16:46:00]
func HandleContent(content string) (*config.AlarmMessage, error) {
    args := strings.Split(content, "][")
    if len(args) < 6 {
        return nil, errors.New("告警消息格式不匹配，可能是版本不一致导致")
    }

    args[0] = string([]rune(args[0])[1:])
    arg := args[5]
    args[5] = string([]rune(arg)[:len(arg)-1])

    // 描述和条件
    argStr := args[4]
    p := regexp.MustCompile(`\w+\(#(\d+)\)`)

    // 表达式
    expression := p.FindString(argStr)
    pos := strings.Index(argStr, expression)
    // pos 前面的是自定义描述
    desc := argStr[:pos-1] // 去掉最后一个空格

    // 提取表达式内的时间间隔
    triggerCount, err := strconv.Atoi(p.FindStringSubmatch(expression)[1])
    if err != nil {
        return nil, err
    }

    // metric tags
    argStr = strings.TrimSpace(argStr[pos+len(expression):])
    pos = strings.Index(argStr, " ")
    counter := argStr[:pos]
    tags := strings.TrimSpace(argStr[pos:])

    // 次数和时间
    argStr = args[5]
    subArgs := strings.Split(argStr, " ")
    count, err := strconv.Atoi(subArgs[0][1:])
    if err != nil {
        return nil, err
    }
    timeStr := strings.Join(subArgs[1:], " ")
    t, err := time.Parse("2006-01-02 15:04:05", timeStr)
    if err != nil {
        return nil, err
    }
    return &config.AlarmMessage{Level: args[0], Type: args[1], Endpoint: args[2], Expression: expression,
        Desc: desc, Counter: counter, Tags: tags, TriggerCount: triggerCount, Count: count, Time: t}, nil
}
