pragma solidity ^0.5.16;

import "./RegisterLock.sol";
import "./Owned.sol";

contract Client is RegisterLock, Owned {

    struct ClientStruct {
        address addr;
        uint balance;
    }

    ClientStruct[] public clients;

    uint public minimumRegister = 50 * 1e18;
    // Constructor
    constructor(address _tokenAddress, address _owner, uint _unlockTime)
    RegisterLock(_tokenAddress, _owner, _unlockTime)
    Owned(_owner)
    public {}

    function register(uint _amount) public {
        require(_amount >= minimumRegister, "Deposit is below the minimum required amount");
        deposit(_amount);

        // Add new client to the list
        clients.push(ClientStruct({
            addr: msg.sender,
            balance: deposits[msg.sender].amount // Assuming deposit updates this balance
        }));
    }

    function increaseDeposit(uint _additionalAmount) public {
        require(_additionalAmount > 0, "Additional amount must be greater than 0");
        deposit(_additionalAmount);
        // Update client's balance
        for (uint i = 0; i < clients.length; i++) {
            if (clients[i].addr == msg.sender) {
                clients[i].balance = deposits[msg.sender].amount;
                break;
            }
        }
    }

    function unregister(address recipient) public {
        // Call withdraw function
        withdraw(recipient);

        // Find and remove the client from the clients array
        for (uint i = 0; i < clients.length; i++) {
            if (clients[i].addr == recipient) {
                // Move the last element into the place to delete
                clients[i] = clients[clients.length - 1];
                // Remove the last element
                clients.pop();
                break;
            }
        }
    }

    function unregister() public {
        // Ensure only the client can unregister themselves
        require(deposits[msg.sender].amount > 0, "You have no deposit to withdraw");

        // Call withdraw function
        withdraw(msg.sender);

        // Find and remove the client from the clients array
        for (uint i = 0; i < clients.length; i++) {
            if (clients[i].addr == msg.sender) {
                // Move the last element into the place to delete
                clients[i] = clients[clients.length - 1];
                // Remove the last element
                clients.pop();
                break;
            }
        }
    }

    function isClient(address _address) public view returns (bool) {
        for (uint i = 0; i < clients.length; i++) {
            if (clients[i].addr == _address) {
                return true;
            }
        }
        return false;
    }

    function setMinimumRegister(uint _minimumRegister) public onlyOwner {
        minimumRegister = _minimumRegister;
    }

}
