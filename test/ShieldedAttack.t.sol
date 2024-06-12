// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Test, console} from "forge-std/Test.sol";
import {UltraVerifier} from "../circuits/contract/shielded_attack/plonk_vk.sol";
import {ShieldedAttack} from "../src/ShieldedAttack.sol";
import {IHasher} from "../src/MerkleTreeWithHistory.sol";


contract ShieldedAttackTest is Test {

	UltraVerifier public plonkVerifier;
	ShieldedAttack public shieldedAttack;

	// public inputs
	bytes32 _correctInput1 = 0x00000000000000000000000000000000000000000000000000000000000004d2;
	// missile tokenId
	bytes32 _returnedValue = 0x14c704ef63cef381567b698c7058e0e9e26b57bb92f255121b69ec58fe60186e;
	bytes32[] correctInputs = new bytes32[](2);

	function setUp() public {
		plonkVerifier = new UltraVerifier();

		uint amount = 1;
		uint32 merkleTreeHeight = 3;

		string memory hasherJson = vm.readLine("./src/Hasher.json");
		bytes memory data = vm.parseJson(hasherJson);
		IHasher hasher = abi.decode(data, (IHasher));

		shieldedAttack = new ShieldedAttack(
			plonkVerifier,
			hasher,
			amount,
			merkleTreeHeight
		);

	}

	function test_ShieldedAttack_CorrectInput() public {

		correctInputs[0] = _correctInput1;
		correctInputs[1] = _returnedValue;
		string memory proofStr = vm.readLine("./circuits/proofs/shielded_attack.proof");
		bytes memory proof = vm.parseBytes(proofStr);

		bool result = shieldedAttack.verifyShieldedAttack(proof, correctInputs);
		console.log("proof verified: ", result);
		assert(result);
	}

}