pragma solidity ^0.8.17;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "./utils/Adminable.sol";

contract MockERC20 is Adminable, ERC20 {

    uint8 public decimals = 18;
    mapping(address => uint256) balances;

    event Transfer(address indexed _from, address indexed _to, uint256 _value);

    constructor(
        string memory name,
        string memory symbol
    ) ERC20(name, symbol, decimals) {
    }

    function mint(address to, uint256 amount) public onlyAdminOrOwner {
        _mint(to, amount);
    }

    function burn(address from, uint256 amount) public onlyOwner {
        _burn(from, amount);
    }

}