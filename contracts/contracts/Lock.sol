pragma solidity ^0.5.16;

import "./interfaces/IERC20.sol";

contract Lock {
    IERC20 public token;
    uint public unlockTime; // Duration in seconds
    address public owner;
    bool internal locked = false;
    uint public rate = 999; // Variable to store the rate

    struct Deposit {
        uint amount;
        uint unlockTime; // Specific unlock time for each deposit
    }

    mapping(address => Deposit) public deposits;

    event DepositMade(address indexed depositor, uint amount, uint unlockTime);
    event Withdrawal(address indexed withdrawer, uint amount);
    event DepositTransferred(address indexed from, address indexed to, uint amount);

    constructor(address _owner, uint _unlockTime) public {
        owner = _owner;
        unlockTime = _unlockTime;
    }


    // Set unlock time in seconds (e.g., 1 day = 86400 seconds)
    function setUnlockTime(uint _durationInSeconds) public {
        require(msg.sender == owner, "Only owner can set unlock time");
        unlockTime = _durationInSeconds;
    }

    // Function to set the rate
    function setRate(uint _rate) public {
        require(msg.sender == owner, "Only owner can set rate");
        require(_rate >= 1 && _rate <= 1000, "Rate must be between 1 and 1000");
        rate = _rate;
    }


    function deposit(uint _amount) public {
        require(_amount > 0, "Amount must be greater than 0");
        require(token.transferFrom(msg.sender, address(this), _amount), "Transfer failed");

        uint depositUnlockTime = block.timestamp + unlockTime;

        if (deposits[msg.sender].amount > 0) {
            deposits[msg.sender].amount += _amount;
            deposits[msg.sender].unlockTime = depositUnlockTime;
        } else {
            deposits[msg.sender] = Deposit(_amount, depositUnlockTime);
        }

        emit DepositMade(msg.sender, _amount, depositUnlockTime);
    }

    function withdraw() internal noReentrant{
        require(block.timestamp >= deposits[msg.sender].unlockTime, "You can't withdraw yet");
        require(deposits[msg.sender].amount > 0, "No deposit to withdraw");

        uint totalAmount = deposits[msg.sender].amount;
        uint amountToWithdraw = totalAmount.mul(rate).div(1000);

        deposits[msg.sender].amount = 0;

        require(token.transfer(msg.sender, amountToWithdraw), "Transfer failed");

        emit Withdrawal(msg.sender, amountToWithdraw);
    }

    function withdrawByAddress(address recipient) internal noReentrant {
        require(recipient != address(0), "Recipient address cannot be zero");
        require(block.timestamp >= deposits[recipient].unlockTime, "You can't withdraw yet");
        require(deposits[recipient].amount > 0, "No deposit to withdraw");


        uint totalAmount = deposits[recipient].amount;
        uint amountToWithdraw = totalAmount

        deposits[recipient].amount = 0;

        require(token.transfer(recipient, amountToWithdraw), "Transfer failed");

        emit Withdrawal(recipient, amountToWithdraw);
    }


    function transferDeposit(address newAddress) internal noReentrant{
        require(newAddress != address(0), "New address cannot be the zero address");
        require(deposits[msg.sender].amount > 0, "No deposit to transfer");

        // Store the deposit amount before resetting
        uint depositAmount = deposits[msg.sender].amount;

        // Reset the deposit of the original depositor
        deposits[msg.sender] = Deposit(0, 0);

        // Transfer the deposit details to the new address
        deposits[newAddress] = Deposit(depositAmount, block.timestamp + unlockTime);

        // Emit an event to record the transfer of deposit
        emit DepositTransferred(msg.sender, newAddress, depositAmount);
    }

    // Function to get the staked balance of an address
    function getStakedBalance(address _user) public view returns (uint) {
        return deposits[_user].amount;
    }

    modifier noReentrant() {
        require(!locked, "Race condition detected");
        locked = true;
        _;
        locked = false;
    }

}
