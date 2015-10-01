package main

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/jasonmoo/oc"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/vmihailenco/msgpack.v2"
)

func main() {

	fmt.Println("running...")
	defer fmt.Println("done!")

	nums := make([]int64, 1<<20)

	for i, _ := range nums {
		nums[i] = int64(i)
	}

	bnums := []byte{}
	bn := (*reflect.SliceHeader)(unsafe.Pointer(&bnums))
	bn.Data = (uintptr)(unsafe.Pointer(&nums[0]))
	bn.Cap = cap(nums) * 8
	bn.Len = len(nums) * 8

	set := oc.NewOc()

	set.Increment("[]int64", len(nums)*8)
	set.Increment("[]byte", len(bnums))

	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	w.Write(bnums)
	w.Close()
	set.Increment("gzip []byte", buf.Len())

	buf.Reset()
	json.NewEncoder(&buf).Encode(bnums)
	set.Increment("json []int64", buf.Len())

	buf.Reset()
	w.Reset(&buf)
	json.NewEncoder(w).Encode(bnums)
	w.Close()
	set.Increment("gzip json []int64", buf.Len())

	buf.Reset()
	gob.NewEncoder(&buf).Encode(bnums)
	set.Increment("gob []int64", buf.Len())

	buf.Reset()
	w.Reset(&buf)
	gob.NewEncoder(w).Encode(bnums)
	w.Close()
	set.Increment("gzip gob []int64", buf.Len())

	buf.Reset()
	msgpack.NewEncoder(&buf).Encode(bnums)
	set.Increment("msgpack []int64", buf.Len())

	buf.Reset()
	w.Reset(&buf)
	msgpack.NewEncoder(w).Encode(bnums)
	w.Close()
	set.Increment("gzip msgpack []int64", buf.Len())

	data, _ := bson.Marshal(bnums)
	set.Increment("bson []int64", len(data))

	buf.Reset()
	w.Reset(&buf)
	w.Write(data)
	w.Close()
	set.Increment("gzip bson []int64", buf.Len())

	set.SortByCt(oc.ASC)

	for set.Next() {
		fmt.Println(set.KeyValue())
	}

}

// Output:
//
// running...
// gzip []byte 1587209
// gzip msgpack []int64 1587217
// gzip gob []int64 1587221
// gzip json []int64 1946868
// []int64 8388608
// []byte 8388608
// msgpack []int64 8388613
// gob []int64 8388618
// json []int64 11184815
// gzip bson []int64 21968449
// bson []int64 107940799
// done!
