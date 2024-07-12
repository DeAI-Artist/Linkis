pragma solidity ^0.5.16;

interface IRegisterable {
    function register(uint _amount) external;
    function increaseDeposit(uint _additionalAmount) external;
    function unregister() external;
    function isRegistered(address _address) external view returns (bool);
    function setMinimumRegister(uint _minimumRegister) external;
}
