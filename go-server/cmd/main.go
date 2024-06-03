package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
)

func main() {

	port := 5555
	fmt.Printf("Starting server at %v...\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%v", port), &GetSqrt{})
	if err != nil {
		fmt.Println(err)
	}

}

type GetSqrt struct{}

type SingleForeignCallParam struct {
	Single string `json:"Single"`
}

type IncomingRpc struct {
	JsonRpc string `json:"jsonrpc"`
	//   "jsonrpc": "2.0",

	Method string `json:"method"`
	//   "method": "Calculator.Add",

	Params []SingleForeignCallParam `json:"params"`
	//   "params": { "Single": 13 },

	Id uint64 `json:"id"`
	//   "id": 123
}

// response format: { jsonrpc: '2.0', id: 1, result: { values: [{ Single: 5 }] } }
type RpcResponse struct {
	JsonRpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  Values `json:"result"`
}

type Values struct {
	Values []SingleForeignCallParam `json:"values"`
}

func (sqrt *GetSqrt) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var payload IncomingRpc
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("incoming request payload", payload)

	numberToSquare := payload.Params[0].Single
	n, _ := strconv.ParseUint(numberToSquare, 16, 64)
	fmt.Println("incoming number to square: ", n)

	// parse hex to int
	squaredNumber := math.Sqrt(float64(n))
	// watchout for floats
	fmt.Println("sqrt n: ", squaredNumber)
	var s string = strconv.FormatUint(uint64(squaredNumber), 10)
	fmt.Println("toString s: ", s)
	xsingle := SingleForeignCallParam{s}

	sqrtRes, err2 := json.Marshal(&RpcResponse{
		JsonRpc: "2.0",
		Id:      1,
		Result:  Values{Values: []SingleForeignCallParam{xsingle}},
	})
	if err2 != nil {
		fmt.Println(err2)
		http.Error(w, err2.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(string(sqrtRes))
	w.Write([]byte(sqrtRes))
}
