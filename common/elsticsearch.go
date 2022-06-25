package common

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"net"
	"net/http"
	"time"
)

type DB interface {
	//存储
	Save(index string,ty string,info interface{})error
	//查询
    Query(index string,ty string,param map[string]interface{})(interface{},error)
}

type ESClient struct {
	*elasticsearch.Client
}

func NewESClient(url string)(*ESClient,error){
	cfg := elasticsearch.Config{
		Addresses: []string{url},
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
	if err !=nil {
		return nil,err
	}
	return &ESClient{es},nil
}

func(es *ESClient)Save(index string,ty string,info interface{})error{


	data,err:=json.Marshal(info)
	if err !=nil {
		return err
	}
	// Set up the request object.
	req := esapi.IndexRequest{
		//相当于数据库名称
		Index:      index,
		//相当于关系型数据库中的table,严格意义来讲，不是这样一层关系
		DocumentType: ty,
		//DocumentID: strconv.Itoa(i + 1),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), es)
	if err != nil {
		return fmt.Errorf("Error getting response: %s", res.Body)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("respon [%s] Error", res.Status())
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			return fmt.Errorf("Error parsing the response body: %s", err)
		} else {

		}
		return nil
	}
}

func(es *ESClient)Query(index string,ty string,param map[string]interface{})(interface{},error){
	var buf bytes.Buffer
	var r  map[string]interface{}
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": param,
		},
		"size":10,
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil,fmt.Errorf("Error encoding query: %s", err)
	}

	// Perform the search request.
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(index),
		es.Search.WithDocumentType(ty),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		return nil,fmt.Errorf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil,fmt.Errorf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			return nil,fmt.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil,fmt.Errorf("Error parsing the response body: %s", err)
	}
	return r,nil
}