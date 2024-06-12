// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {Script, console} from "forge-std/Script.sol";
import "../src/MockERC20.sol";

contract DeployMockTokenScript is Script {

    function run() external returns(address) {

        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        vm.startBroadcast(deployerPrivateKey);

        MockERC20 token = new TestERC20("EigenWifHat", "EWH");

        vm.stopBroadcast();

        return token;
    }

}