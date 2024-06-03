import { BarretenbergBackend } from "@noir-lang/backend_barretenberg";
import { Noir } from "@noir-lang/noir_js";
import radius from "../circuits-radius/target/radius.json" assert { type: "json" };
import { JSONRPCClient } from "json-rpc-2.0";

const client = new JSONRPCClient((jsonRPCRequest) => {
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

const oracleResolver = async (name, input) => {
	// oracleResolver automatically transforms public 'd' input to this format:
	// input = [ [ '0x0000000000000000000000000000000000000000000000000000000000000019' ] ]
	let inputD = input[0][0].toString(16).padStart(64, "0")
	const oracleReturn = await client.request(name, [
		{ Single: inputD }
	]);

	return [ oracleReturn.values[0].Single ];
};

async function main() {
	const backend = new BarretenbergBackend(radius);
	const noir = new Noir(radius, backend);

	const input = {
		x1: 1,
		x2: 4,
		y1: 1,
		y2: 5,
		d: 225
	};

	// const oracleReturn = await client.request("GetSqrt", [
	// 	{ Single: input.d.toString(16).padStart(64, '0') }
	// ]);
	// console.log('\n>>>> oracle.values.Single:', oracleReturn.values[0].Single);

	const proof = await noir.generateProof(input, oracleResolver);
	const proofStr = Buffer.from(proof.proof).toString('hex');
	console.log("proofString:", proofStr);
	console.log("generateProof:", proof);

	const verified = await noir.verifyProof({
		proof: proof.proof,
		publicInputs: ["225"]
	});
	console.log('verifyProof: ', verified);
	process.exit(1)
}


main();