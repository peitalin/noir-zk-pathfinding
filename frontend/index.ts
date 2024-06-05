import { BarretenbergBackend } from "@noir-lang/backend_barretenberg";
import { Noir } from "@noir-lang/noir_js";
import astar from "../circuits-astar/target/astar.json" assert { type: "json" };
import { JSONRPCClient } from "json-rpc-2.0";

const client = new JSONRPCClient((jsonRPCRequest) => {
	// console.log("jsonRPC request: ", jsonRPCRequest)
	// console.log("params: ", jsonRPCRequest.params)
	return fetch("http://localhost:5555", {
		method: "POST",
		headers: {
			"content-type": "application/json",
		},
		body: JSON.stringify(jsonRPCRequest),
	}).then((response) => {
		if (response.status === 200) {
			return response
				.json()
				.then((jsonRPCResponse) => client.receive(jsonRPCResponse));
		} else if (jsonRPCRequest.id !== undefined) {
			return Promise.reject(new Error(response.statusText));
		}
	});
});

const oracleResolverSqrt = async (name, input) => {
	// oracleResolver automatically transforms public 'd' input to this format:
	// input = [ [ '0x0000000000000000000000000000000000000000000000000000000000000019' ] ]
	let inputD = input[0][0].toString(16).padStart(64, "0")
	const oracleReturn = await client.request(name, [
		{ Single: inputD }
	]);

	// NOTE: must remove all println in main.nr
	return [ oracleReturn.values[0].Single ];
};

const oracleResolverAstar = async (name, input) => {
	// oracleResolver automatically transforms public 'd' input to this format:
	// input = [ [ 'hex', 'hex ], ['hex', 'hex'] ]
	let x1 = input[0][0].toString(16).padStart(64, "0").slice(2)
	let y1 = input[0][1].toString(16).padStart(64, "0").slice(2)
	let x2 = input[1][0].toString(16).padStart(64, "0").slice(2)
	let y2 = input[1][1].toString(16).padStart(64, "0").slice(2)

	const oracleReturn = await client.request(name, [
		{ Array: [ x1, y1 ] },
		{ Array: [ x2, y2 ] },
	]);
	// NOTE: must remove all println in main.nr
	let data = oracleReturn.values.map(e => e.Array)
	return data

};

async function main() {
	const backend = new BarretenbergBackend(astar);
	const noir = new Noir(astar, backend);

	const input = {
		// private inputs
		x1: 1,
		x2: 4,
		y1: 1,
		y2: 5,
		// public inputs
		max_steps: 10
	};

	// const oracleReturn = await client.request("GetSqrt", [
	// 	{ Single: input.d.toString(16).padStart(64, '0') }
	// ]);
	// console.log('\n>>> oracle.values.Single:', oracleReturn.values[0].Single);

	const proof = await noir.generateProof(input, oracleResolverAstar);
	const proofStr = Buffer.from(proof.proof).toString('hex');
	console.log("proofString:", proofStr);
	console.log("generateProof:", proof);

	const verified = await noir.verifyProof({
		proof: proof.proof,
		publicInputs: ["10"]
	});
	console.log('verifyProof: ', verified);
	process.exit(1)
}


main();