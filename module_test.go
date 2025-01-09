package smart

import (
	"github.com/ywengineer/smart/utility"
	"go.uber.org/zap"
	"regexp"
	"testing"
)

type Req struct {
	Ping  int    `json:"ping"`
	Extra string `json:"extra"`
}

func (r *Req) Reset() {
	r.Ping = -1
	r.Extra = ""
}

type Res struct {
	Pong int `json:"pong"`
}

type TestModule struct {
}

func (m *TestModule) Name() string {
	return "TestModule"
}

func (m *TestModule) RegisterAccount1001(channel *SocketChannel, req *Req) {
	utility.DefaultLogger().Info("RegisterAccount1001 invoked", zap.Any("req", *req))
}

func (m *TestModule) FindFriend1002(channel *SocketChannel, req *Req) {
	utility.DefaultLogger().Info("FindFriend1002 invoked", zap.Any("req", *req))
}

func (m *TestModule) UseItem1003(channel *SocketChannel, req *Req) {
	utility.DefaultLogger().Info("UseItem1003 invoked", zap.Any("req", *req))
}

func (m *TestModule) StartFight1004(channel *SocketChannel, req *Req) *Res {
	utility.DefaultLogger().Info("StartFight1004 invoked", zap.Any("req", *req))
	return &Res{
		Pong: req.Ping,
	}
}

func TestRegisterModule(t *testing.T) {
	err := RegisterModule(&TestModule{})
	if err != nil {
		t.Errorf("%v", err)
		t.FailNow()
	}
}

func TestRegex(t *testing.T) {
	r := regexp.MustCompile("^\\D+([1-9][0-9]*)$")
	t.Log(r.String())
	t.Logf("%v", r.FindAllString("Test1001", -1))
	t.Logf("%v", r.FindAllStringSubmatch("Test1001", -1))
	t.Logf("%v", r.FindAllStringSubmatch("aTest1001", -1))
	t.Logf("%v", r.FindAllStringSubmatch("fTest10021", -1))
	t.Logf("%v", r.FindStringSubmatch("Test100113"))

	r1 := regexp.MustCompile("\\D+")
	t.Logf("%s", r1.ReplaceAllString("Test100113a", ""))
}
