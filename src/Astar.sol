pragma solidity ^0.8.17;

import "../circuits/contract/astar/plonk_vk.sol";

contract Astar {

  UltraVerifier public verifier;

  constructor(UltraVerifier _verifier) {
    verifier = _verifier;
  }

  function verifyDistance(bytes calldata proof, bytes32[] calldata max_steps) public view returns (bool) {
    bool proofResult = verifier.verify(proof, max_steps);
    require(proofResult, "proof is invalid");
    return proofResult;
  }

}