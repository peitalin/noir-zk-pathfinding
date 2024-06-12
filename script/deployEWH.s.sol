// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "forge-std/Script.sol";
import "../src/EigenWifHatERC20.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";

contract DeployEigenWifHatScript is Script {

    function run() external returns (address, address) {

        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        vm.startBroadcast(deployerPrivateKey);

        // deploy proxy
        address _proxyAddress = Upgrades.deployTransparentProxy(
            "EigenWifHatERC20.sol",
            msg.sender,
            abi.encodeCall(EigenWifHatERC20.initialize, (msg.sender))
        );

        // get implementation address
        address implementationAddress = Upgrades.getImplementationAddress(_proxyAddress);

        vm.stopBroadcast();

        return (implementationAddress, _proxyAddress);
    }

}