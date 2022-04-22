package common

import (
	"chain33-evm-parser/common/abi"
	"fmt"
	"strings"

	"github.com/33cn/chain33/types"

	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/common"
)


//多个合约订阅事件解析,全局变量   contractAddr-->eventID--->event
var TopicsContractsMap = make(map[string]map[common.Hash]abi.Event)

func ParseTopics(event abi.Event, topics []string) (map[string]interface{}, error) {
	var hashs []common.Hash
	for _, topic := range topics {
		hashs = append(hashs, common.BytesToHash(common.FromHex(topic)))
	}
	//判断eventID 是否相等,如果不等说明不是该事件
	if len(hashs) == 0 || hashs[0] != event.ID {
		return nil, fmt.Errorf("It's not a listen event!")
	}
	outMap := make(map[string]interface{})
	err := abi.ParseTopicsIntoMap(outMap, event.Inputs, hashs[1:])
	if err != nil {
		return outMap, fmt.Errorf("ParseTopics have a err: %s /n", err.Error())
	}
	return outMap, nil
}

//直接解析evm订阅事件
func ParseEVMTxLogs(event abi.Event, blks *types.EVMTxLogsInBlks) []map[string]interface{} {
	var replys []map[string]interface{}
	for _, blk := range blks.GetLogs4EVMPerBlk() {
		for _, txLog := range blk.GetTxAndLogs() {
			for _, log := range txLog.GetLogsPerTx().GetLogs() {
				//判断eventID 是否相等,如果不等说明不是该事件
				if len(log.GetTopic()) == 0 || common.BytesToHash(log.GetTopic()[0]) != event.ID {
					continue
				}
				var hashs []common.Hash
				for _, topic := range log.GetTopic() {
					hashs = append(hashs, common.BytesToHash(topic))
				}
				outMap := make(map[string]interface{})
				err := abi.ParseTopicsIntoMap(outMap, event.Inputs, hashs[1:])
				if err != nil {
					continue
				}
				replys = append(replys, outMap)
			}
		}
	}
	return replys
}

// 直接解析block订阅日志,返回数据存储： txhash--->contractAddr--->eventID--->paramName--->value
func ParseBlockReceipts(reqs *types.BlockSeqs) map[string]map[string]map[common.Hash]map[string]interface{} {

	//txhash--->contractAddr--->eventID--->paramName--->value
	var results = make(map[string]map[string]map[common.Hash]map[string]interface{})
	for _, req := range reqs.GetSeqs() {
		for txIndex, tx := range req.GetDetail().Block.Txs {
			//确认是订阅的交易类型
			if !strings.Contains(string(tx.Execer), "evm") {
				continue
			}
			var evmAction types.EVMContractAction4Chain33
			err := types.Decode(tx.Payload, &evmAction)
			if nil != err {
				continue
			}
			//如果全局变量中存该合约
			topicsEvent, ok := TopicsContractsMap[evmAction.ContractAddr]
			if ok {
				//因为只有交易执行成功时，才会存证log信息，所以需要事先判断
				if types.ExecOk != req.GetDetail().Receipts[txIndex].Ty {
					continue
				}
				results[common.Bytes2Hex(tx.Hash())] = make(map[string]map[common.Hash]map[string]interface{})
				for _, log := range req.GetDetail().Receipts[txIndex].Logs {
					//TyLogEVMEventData = 605 这个log类型定义在evm合约内部
					if 605 != log.Ty {
						continue
					}
					var evmLog types.EVMLog
					err := types.Decode(log.Log, &evmLog)
					if nil != err {
						continue
					}
					//从topicsEvent中匹配相关事件
					event, ok := topicsEvent[common.BytesToHash(evmLog.GetTopic()[0])]
					if ok {
						results[common.Bytes2Hex(tx.Hash())][evmAction.ContractAddr] = make(map[common.Hash]map[string]interface{})
						var hashs []common.Hash
						for _, topic := range evmLog.GetTopic() {
							hashs = append(hashs, common.BytesToHash(topic))
						}
						outMap := make(map[string]interface{})
						err := abi.ParseTopicsIntoMap(outMap, event.Inputs, hashs[1:])
						if err != nil {
							continue
						}
						results[common.Bytes2Hex(tx.Hash())][evmAction.ContractAddr][event.ID] = outMap
					}
				}
			}
		}
	}
	return results
}
