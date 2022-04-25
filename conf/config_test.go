package conf

import (
	"os"
	"testing"
)

func Test_ParseConfig( t *testing.T){
	data,err:=os.ReadFile("../conf.json")
	if err !=nil {
		t.Error(err)
	}
	conf,err:=ParseConfig(data)
	if err !=nil {
		t.Error(err)
	}
	t.Log(conf.ParseTopics[0].Abi.Events["Transfer"].ID.Hex())
	//返回值类型判断可以根据abi中参数类型进行解析
	for _, arg := range conf.ParseTopics[0].Abi.Events["Transfer"].Inputs {
		t.Logf("param name is %s,type is %v", arg.Name,arg.Type)
	}
}

