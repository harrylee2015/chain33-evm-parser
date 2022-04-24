package main

import (
	"chain33-evm-parser/common"
	"chain33-evm-parser/common/abi"
	"chain33-evm-parser/conf"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/33cn/chain33/types"
	. "github.com/33cn/plugin/plugin/dapp/evm/executor/vm/common"
)

// 主动关闭服务器
var server *http.Server

//需要解析的parseMap
var parseMap *common.ParseMap

func main() {

	data, err := os.ReadFile("conf.json")
	if err != nil {
		log.Panic(err)
	}
	cfg, err := conf.ParseConfig(data)
	if err != nil {
		log.Panic(err)
	}
	//初始化并赋值
	parseMap = &common.ParseMap{
		TopicsContractMap: make(map[string]map[Hash]abi.Event),
		TopicsEventMap:    make(map[Hash]abi.Event),
	}
	for _, parseTopic := range cfg.ParseTopics {
		eventMap := make(map[Hash]abi.Event)
		for _, event := range parseTopic.EventNames {
			parseMap.TopicsEventMap[parseTopic.Abi.Events[event].ID] = parseTopic.Abi.Events[event]
			eventMap[parseTopic.Abi.Events[event].ID] = parseTopic.Abi.Events[event]
		}
		if parseTopic.ContractAddr != "" {
			parseMap.TopicsContractMap[parseTopic.ContractAddr] = eventMap
		}
	}

	// 一个通知退出的chan
	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt)

	mux := http.NewServeMux()
	mux.Handle("/", &Handler{cfg: cfg})
	//mux.HandleFunc("/v1/heath", types.HealthCheck{})

	server = &http.Server{
		Addr:         cfg.ListenServer.ListenAddr,
		WriteTimeout: time.Second * 4,
		Handler:      mux,
	}

	go func() {
		// 接收退出信号
		<-exit
		if err := server.Close(); err != nil {
			log.Fatal("Close server:", err)
		}
	}()

	log.Println("Starting v3 httpserver")
	err = server.ListenAndServe()
	if err != nil {
		// 正常退出
		if err == http.ErrServerClosed {
			log.Fatal("Server closed under request")
		} else {
			log.Fatal("Server closed unexpected", err)
		}
	}
	log.Fatal("Server exited")

}

type Handler struct {
	cfg *conf.Config
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)
	if len(r.Header["Content-Encoding"]) >= 1 && r.Header["Content-Encoding"][0] == "gzip" {
		gr, err := gzip.NewReader(r.Body)
		body, err := ioutil.ReadAll(gr)
		if err != nil {
			log.Fatal("Error while serving JSON request: %v", err)
			return
		}

		err = handlerReq(body, h.cfg)
		if err == nil {
			w.Write([]byte("OK"))
		} else {
			w.Write([]byte(err.Error()))
		}
	}
}

//解析evm订阅
func handlerReq(body []byte, cfg *conf.Config) error {
	//TODO 这里暂时只支持区块订阅类型事件解析
	if cfg.Topic.Type == 0 {
		var reqs types.BlockSeqs
		if cfg.Topic.Encode == "jrpc" {
			err := types.JSONToPB(body, &reqs)
			if err != nil {
				log.Fatal("Decoding JSON body have err: %v", err)
				return err
			}
		} else {
			err := types.Decode(body, &reqs)
			if err != nil {
				log.Fatal("Decoding proto body have err: %v", err)
				return err
			}
		}
		results := common.ParseBlockReceipts(&reqs, parseMap)
		//TODO 后续处理
		log.Println(results)
		return nil
	}

	if cfg.Topic.Type == 4 {
		var reqs types.EVMTxLogsInBlks
		if cfg.Topic.Encode == "jrpc" {
			err := types.JSONToPB(body, &reqs)
			if err != nil {
				log.Fatal("Decoding JSON body have err: %v", err)
				return err
			}
		} else {
			err := types.Decode(body, &reqs)
			if err != nil {
				log.Fatal("Decoding proto body have err: %v", err)
				return err
			}
		}
		results := common.ParseEVMTxLogs(&reqs, parseMap)
		//TODO 后续处理
		log.Println(results)
		return nil
	}
	return fmt.Errorf("unknown type")
}
