package common

import (
	"chain33-evm-parser/common/abi"
	"fmt"

	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/common"
)

//监听事件
//type ListenEvent struct {
//	//合约地址,可以为空
//	ContractAddr  string
//	//监听合约事件,可能是多组事件
//	Events       []*abi.Event
//}

func ParseTopics(event abi.Event, topics []string)(map[string]interface{},error) {
	var hashs []common.Hash
	for _, topic := range topics {
		hashs = append(hashs, common.BytesToHash(common.FromHex(topic)))
	}
	//判断eventID 是否相等,如果不等说明不是该事件
	if len(hashs) == 0 || hashs[0] != event.ID {
		return nil,fmt.Errorf("It's not a listen event!")
	}
	outMap := make(map[string]interface{})
	err := abi.ParseTopicsIntoMap(outMap, event.Inputs, hashs[1:])
	if err != nil {
		return outMap,fmt.Errorf("ParseTopics have a err: %s /n", err.Error())
	}
	return outMap,nil
}

