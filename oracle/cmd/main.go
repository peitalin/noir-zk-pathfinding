package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"

	astar "noir.com/oracle/astar"
)

func main() {

	astar.Test1()

	port := 5555
	fmt.Printf("Starting server at %v...\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%v", port), &OracleRPC{})
	if err != nil {
		fmt.Println(err)
	}

}

type OracleRPC struct{}

type SingleForeignCallParam struct {
	Single string `json:"Single"`
}

type ArrayForeignCallParam struct {
	Array []string `json:"Array"`
}

// unmarshal as either Single or Array type in RPC call
type ForeignCallParam struct {
	SingleForeignCallParam `json:"Single,omitempty"`
	ArrayForeignCallParam  `json:"Array,omitempty"`
}

type IncomingRpc struct {
	JsonRpc string `json:"jsonrpc"`
	//   "jsonrpc": "2.0",

	Method string `json:"method"`
	//   "method": "GetSqrt", // method name is GetSqrt defined in Aztec main.nr

	Params []ForeignCallParam `json:"params"`
	// unmarshal as either Single or Array type in RPC call
	//   "params": [{ "Single": "hex" }] or [{"Array": ["hex"]}]

	Id uint64 `json:"id"`
	//   "id": 123
}

func (irpc *ForeignCallParam) UnmarshalJSON(data []byte) error {
	sdata := string(data)
	switch sdata[2] {
	case 'S':
		fmt.Println("\nUnmarshal data as Single:", sdata)
		if err := json.Unmarshal(data, &irpc.SingleForeignCallParam); err != nil {
			return err
		}
	case 'A':
		fmt.Println("\nUnmarshal data as Array:", sdata)
		irpc.Array = []string{}
		if err := json.Unmarshal(data, &irpc.ArrayForeignCallParam); err != nil {
			return err
		}
	}
	return nil
}

// response format: { jsonrpc: '2.0', id: 1, result: { values: [{ Single: 5 }] } }
type RpcResponse struct {
	JsonRpc string       `json:"jsonrpc"`
	Id      uint64       `json:"id"`
	Result  OutputValues `json:"result"`
}

type SingleForeignCallOutput struct {
	Single string `json:"Single"`
}

type ArrayForeignCallOutput struct {
	Array []string `json:"Array"`
}

type OutputValues struct {
	Values []SingleForeignCallOutput `json:"values"`
}

type RpcResponseArray struct {
	JsonRpc string            `json:"jsonrpc"`
	Id      uint64            `json:"id"`
	Result  OutputValuesArray `json:"result"`
}
type OutputValuesArray struct {
	Values []ArrayForeignCallOutput `json:"values"`
}

func (oracle *OracleRPC) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var payload IncomingRpc
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("\nIncomingRpc Payload: ", payload)

	//////////////////////////////////////////////////////////
	if payload.Method == "GetSqrt" {
		s := handleGetSqrt(payload)
		response_single := SingleForeignCallOutput{parseIntasHex(s)}

		sqrtRes, err2 := json.Marshal(&RpcResponse{
			JsonRpc: "2.0",
			Id:      payload.Id,
			Result:  OutputValues{Values: []SingleForeignCallOutput{response_single}},
		})
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("\nResponse with hex-encoded values:", string(sqrtRes))
		w.Write([]byte(sqrtRes))
	}

	//////////////////////////////////////////////////////////
	if payload.Method == "GetAstarPath" {

		fmt.Println("\nParams: ", payload.Params)
		// start
		x1, _ := parseHexToInt(payload.Params[0].Array[0])
		y1, _ := parseHexToInt(payload.Params[0].Array[1])

		// end
		x2, _ := parseHexToInt(payload.Params[1].Array[0])
		y2, _ := parseHexToInt(payload.Params[1].Array[1])

		fmt.Printf("\nstart:\t(%v, %v)\n", x1, y1)
		fmt.Printf("end:\t(%v, %v)\n", x2, y2)

		maze := astar.CreateMaze()
		A := astar.Node{Parent: nil, Position: astar.Position{Col: x1, Row: y1}}
		B := astar.Node{Parent: nil, Position: astar.Position{Col: x2, Row: y2}}

		astar_path := astar.AStar(maze, A, B)
		fmt.Println("\nA* Path to B: ", astar_path)
		astar.PrintPathOnMaze(&maze, A, astar_path)

		// hex value: 10000 => 65536 when decoded as i32. Use this value as "zero" in noir constraints
		NOIR_ZERO_VALUE_HEX := "0x10000"
		zero_point := ArrayForeignCallOutput{[]string{NOIR_ZERO_VALUE_HEX, NOIR_ZERO_VALUE_HEX}}

		// fixed 10 elements, oracle array response must have a fixed length
		astar_path_hex := []ArrayForeignCallOutput{
			zero_point,
			zero_point,
			zero_point,
			zero_point,
			zero_point,
			zero_point,
			zero_point,
			zero_point,
			zero_point,
			zero_point,
		}

		for index, p := range astar_path {
			var col string = parseIntasHex(p.Col)
			var row string = parseIntasHex(p.Row)
			path_point_hex := []string{col, row}
			astar_path_hex[index].Array = path_point_hex
		}

		sqrtRes, err2 := json.Marshal(&RpcResponseArray{
			JsonRpc: "2.0",
			Id:      payload.Id,
			Result:  OutputValuesArray{Values: astar_path_hex},
		})
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println("\nResponse with hex-encoded values:", string(sqrtRes))
		w.Write([]byte(sqrtRes))
	}

}

func extractSqrtParams(payload IncomingRpc) int32 {
	var numberToSquare string
	if len(payload.Params[0].Single) > 0 {
		numberToSquare = payload.Params[0].Single
	} else {
		numberToSquare = payload.Params[0].Array[0]
	}
	n, _ := parseHexToInt(numberToSquare)
	fmt.Println("\nIncoming number to sqrt: ", n)
	return n
}

func parseHexToInt(h string) (int32, error) {
	n, err := strconv.ParseUint(h, 16, 64)
	return int32(n), err
}

func handleGetSqrt(payload IncomingRpc) int32 {
	n := extractSqrtParams(payload)
	// parse hex to int
	squareRootedNumber := math.Sqrt(float64(n))
	// watchout for floats
	fmt.Println("sqrt(n): ", squareRootedNumber)
	return int32(squareRootedNumber)
}

func parseIntasHex(n int32) string {
	var s string = strconv.FormatUint(uint64(n), 16)
	return "0x" + s
}
