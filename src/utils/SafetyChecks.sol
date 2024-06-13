//SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "@openzeppelin/contracts/access/Ownable.sol";

contract SafetyChecks is Ownable {

    constructor() Ownable(msg.sender) {
    }

    modifier nonZeroAddress(address _addr) {
        require(_addr != address(0), "0 address");
        _;
    }

    function isOwner() internal view returns(bool) {
        return owner() == msg.sender;
    }

}