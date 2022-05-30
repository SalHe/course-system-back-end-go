package test

import (
	"github.com/se2022-qiaqia/course-system/services"
	"testing"
)

func Test_InTestDev(t *testing.T) {
	InitTest()
	if services.Services.User.GetUserCount() > 0 {
		t.Error("怀疑您不在测试环境")
	} else {
		t.Log("您的环境应该是正常的")
	}
}
