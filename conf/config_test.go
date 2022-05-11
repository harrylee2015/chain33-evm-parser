package conf

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
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

func Test_ElasticSearch( t *testing.T){
	cfg := elasticsearch.Config{
		Addresses: []string{"http://106.54.178.39:9200"},
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				InsecureSkipVerify: true,
			},
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		t.Errorf("Error creating the client: %s", err)
	} else {
		info,err:=es.Info()
		if err !=nil {
			t.Error(err)
		}
		t.Logf("info %v",info)
	}


	var (
		r  map[string]interface{}
		wg sync.WaitGroup
	)
	//2. 插入数据
	for i, value := range []string{"Test One", "Test Two"} {
		wg.Add(1)

		go func(i int, value string) {
			defer wg.Done()

			// Build the request body.
			var b strings.Builder
			b.WriteString(`{"value" : "`)
			b.WriteString(value)
			b.WriteString(`"}`)

			// Set up the request object.
			req := esapi.IndexRequest{
				//相当于数据库名称
				Index:      "database",
				//相当于关系型数据库中的table,严格意义来讲，不是这样一层关系
				DocumentType: "table",
				//DocumentID: strconv.Itoa(i + 1),
				Body:       strings.NewReader(b.String()),
				Refresh:    "true",
			}

			// Perform the request with the client.
			res, err := req.Do(context.Background(), es)
			if err != nil {
				t.Errorf("Error getting response: %s", res.Body)
			}
			defer res.Body.Close()

			if res.IsError() {
				t.Errorf("[%s] Error %v", res.Status(), i+1)
			} else {
				// Deserialize the response into a map.
				var r map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
					t.Errorf("Error parsing the response body: %s", err)
				} else {
					t.Logf("respone:%v",r)
				}
			}
		}(i, value)
	}
	wg.Wait()

	// 3. 请求体查询，match
	// 3. Search for the indexed documents
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			//"match": map[string]interface{}{
			//	"sex":  "woman",
			//},
			"range": map[string]map[string]interface{}{
				//小于等于20的妹子
                "age": { "lte":20 },
			},
		},
		"size":10,
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		t.Errorf("Error encoding query: %s", err)
	}

	// Perform the search request.
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("test"),
		es.Search.WithDocumentType("student"),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		t.Errorf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			t.Errorf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			t.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		t.Errorf("Error parsing the response body: %s", err)
	}
	// Print the response status, number of results, and request duration.
	//t.Logf(
	//	"[%s] %d hits; took: %dms",
	//	res.Status(),
	//	int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
	//	int(r["took"].(float64)),
	//)
	// Print the ID and document source for each hit.
	t.Logf("value:%v",r)
}