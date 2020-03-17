/*
@Author : hrbc
@Time : 2020/3/13 4:57 PM
*/
package geecache

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"testing"
)

func TestGetterFunc_Get(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})
	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed")
	}
}

// 模拟耗时数据库 有一个scores表
var db = map[string]map[string]string{
	"scores": {
		"Tom":  "630",
		"Jack": "589",
		"Sam":  "567",
	},
}

func TestGroup_Get(t *testing.T) {
	loadCounts := make(map[string]int, len(db["scores"]))
	loadFromScores := GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key", key)
		if v, ok := db["scores"][key]; ok {
			if _, ok := loadCounts[key]; !ok {
				loadCounts[key] = 0
			}
			loadCounts[key] += 1
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	})
	gee := NewGroup("scores", 2<<10, loadFromScores)
	for k, v := range db["scores"] {
		if view, err := gee.Get(k); err != nil || view.String() != v {
			t.Fatal("failed to get value from db")
		}
		if _, err := gee.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}
	if view, err := gee.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}

func TestHTTPPool_ServeHTTP(t *testing.T) {
	t.Log("start")
	var db = map[string]string{
		"Tom":  "630",
		"Jack": "589",
		"Sam":  "567",
	}
	NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	peers := NewHTTPPool(addr)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
