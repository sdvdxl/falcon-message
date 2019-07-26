## 告警

- 告警等级: {{.Level}}
- 告警类型: {{.Type}}
- 告警指标: {{.Counter}} {{.Tags}}
- 表达式：  {{.Expression}}
- 告警主机: {{.Endpoint}}
- 通知时间: {{timeFormat .Time "2006-01-02 15:04:05"}}
- 告警时间: {{elapse .Count 60 .TriggerCount 300 | timeDiffFormat .Time "2006-01-02 15:04:05"}}
- 当前次数: {{.Count}}
- 告警说明: {{.Desc}}，已持续 {{with elapse .Count 60 .TriggerCount 300}}{{divide . 60}}{{end}}分钟
