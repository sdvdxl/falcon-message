package util

import (
    "encoding/json"
    "fmt"
    "testing"
)

func TestHandleContent(t *testing.T) {
    c, err := HandleContent(`[P2][PROBLEM][10-13-33-153][][测试 all(#1) net.port.listen port=2][O3 2017-06-06 16:46:00]`)
    if err != nil {
        t.Error(err)
    }

    j, _ := json.Marshal(c)
    fmt.Println(string(j))
}
