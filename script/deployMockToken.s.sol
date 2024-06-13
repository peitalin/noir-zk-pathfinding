// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "forge-std/console.sol";
import "forge-std/Script.sol";
import "../src/MockERC20.sol";

contract DeployMockTokenScript is Script {

    function run() external {

        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        vm.startBroadcast(deployerPrivateKey);

        MockERC20 token = new MockERC20("EigenWifHat", "EWH");
        console.log("Deployed token at address: ", address(token));

        vm.stopBroadcast();
    }

}