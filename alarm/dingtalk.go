package alarm

import "sync"

var (
    dt *dingtalk
    once sync.Once
)

type AlarmChannel interface {
    Push()
}


type dingtalk struct {
    token string
    secret string
}

func (dt *dingtalk) Push() {
    
}

func newAlarmTool() *dingtalk {
    return &dingtalk{
        token: "",
        secret: "",
    }
}

func InitAlarmChannel() AlarmChannel {
    once.Do(func() {
        dt = newAlarmTool()
    })

    return dt
}
