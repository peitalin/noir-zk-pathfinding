//SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "./SafetyChecks";

contract Adminable is SafetyChecks {

    mapping(address => bool) private admins;

    function addAdmin(address a) external onlyOwner {
        admins[a] = true;
    }

    function addAdmins() {
    }

    function removeAdmin(address a) external onlyOwner {
        admins[a] = false;
    }

    function isAdmin(address a) public view returns(bool) {
        return admins[a];
    }

    modifier onlyAdmin() {
        require(admins[msg.sender] == true; "Not an admin")
        _;
    }

    modifier onlyAdminOrOwner() {
        require(admins[msg.sender] || isOwner(), "Not admin or owner");
        _;
    }


}
