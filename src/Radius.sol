pragma solidity ^0.8.17;

import "../circuits-radius/contract/radius/plonk_vk.sol";

contract Radius {

  UltraVerifier public verifier;

  constructor(UltraVerifier _verifier) {
    verifier = _verifier;
  }

  function verifyDistance(bytes calldata proof, bytes32[] calldata y) public view returns (bool) {
    bool proofResult = verifier.verify(proof, y);
    require(proofResult, "proof is invalid");
    return proofResult;
  }

}