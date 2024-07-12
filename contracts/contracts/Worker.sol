pragma solidity ^0.5.16;

import "./RegisterLock.sol";
import "./Owned.sol";

contract Worker is RegisterLock, Owned {

    struct WorkerStruct {
        address addr;
        uint balance;
    }

    WorkerStruct[] public workers;

    uint public minimumRegister = 500 * 1e18;

    // Constructor
    constructor(address _tokenAddress, address _owner, uint _unlockTime)
    RegisterLock(_tokenAddress, _owner, _unlockTime)
    Owned(_owner)
    public {}

    function register(uint _amount) public {
        require(_amount >= minimumRegister, "Deposit is below the minimum required amount");
        deposit(_amount);

        // Add new worker to the list
        workers.push(WorkerStruct({
            addr: msg.sender,
            balance: deposits[msg.sender].amount // Assuming deposit updates this balance
        }));
    }

    function increaseDeposit(uint _additionalAmount) public {
        require(_additionalAmount > 0, "Additional amount must be greater than 0");
        deposit(_additionalAmount);

        // Update worker's balance
        for (uint i = 0; i < workers.length; i++) {
            if (workers[i].addr == msg.sender) {
                workers[i].balance = deposits[msg.sender].amount;
                break;
            }
        }
    }

    function unregister() public {
        require(deposits[msg.sender].amount > 0, "You have no deposit to withdraw");

        // Call withdraw function
        withdraw(msg.sender);

        // Find and remove the worker from the workers array
        for (uint i = 0; i < workers.length; i++) {
            if (workers[i].addr == msg.sender) {
                workers[i] = workers[workers.length - 1];
                workers.pop();
                break;
            }
        }
    }

    function isWorker(address _address) public view returns (bool) {
        for (uint i = 0; i < workers.length; i++) {
            if (workers[i].addr == _address) {
                return true;
            }
        }
        return false;
    }

    function setMinimumRegister(uint _minimumRegister) public onlyOwner {
        minimumRegister = _minimumRegister;
    }

}
