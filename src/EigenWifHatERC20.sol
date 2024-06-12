pragma solidity ^0.8.17;

import "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import "./utils/Adminable.sol";

contract EigenWifHatERC20 is Adminable, ERC20Upgradeable {

    using SafeERC20 for IERC20;

    IERC20 public token;

    string public name = "EigenWifHat";
    string public symbol = "EWH";
    uint8 public decimals = 18;

    mapping(address => uint256) balances;
    mapping(address => mapping(address => uint256)) allowed;


    function initialize() external initializer {
        ERC20Upgradeable.__ERC20_init("EigenWifHat", "$EWH");
    }

    function mint(address to, uint256 amount) public onlyAdminOrOwner {
        _mint(to, amount);
    }

    function burn(address from, uint256 amount) public onlyOwner {
        _burn(from, amount);
    }

}