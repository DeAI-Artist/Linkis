pragma solidity ^0.5.16;

import "./interfaces/IERC20.sol";

interface ILock {
    // Events
    event DepositMade(address indexed depositor, uint amount, uint unlockTime);
    event Withdrawal(address indexed withdrawer, uint amount);

    // Constructor is not included in the interface

    // Functions
    function setUnlockTime(uint _durationInSeconds) external;
    function deposit(uint _amount) external;
    function withdraw() external;
    function getStakedBalance(address _user) external view returns (uint);

    // Public variables
    function token() external view returns (IERC20);
    function unlockTime() external view returns (uint);
    function owner() external view returns (address);
}
