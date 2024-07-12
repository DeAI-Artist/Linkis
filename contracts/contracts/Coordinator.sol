pragma solidity 0.5.16;
pragma experimental ABIEncoderV2;

import "./BaseCoordinator.sol";

contract Coordinator is BaseCoordinator {

    // Constructor for the Coordinator contract
    constructor(
        address _lockOwner,
        uint _unlockTime,
        address _owner,
        address _resolver
    )
    BaseCoordinator(_lockOwner, _unlockTime, _owner, _resolver) // Initialize BaseCoordinator
    public
    {
        // Additional initialization specific to Coordinator contract (if any)
    }

    function initialize() public {
        _initialize();
    }

    function nominateSelfAsCoordinator(
        string memory _name,
        address[] memory _functionAddresses,
        uint _amount) public {
        uint256 totalSupply = token.totalSupply();
        uint256 requiredTokenAmount = totalSupply.mul(nominationPercentage).div(100);

        require(token.balanceOf(msg.sender) >= requiredTokenAmount, "Insufficient token balance to nominate");
        require(_amount >= requiredTokenAmount, "Staked amount does not meet the required threshold");

        // Directly call the deposit function
        deposit(_amount);

        // Create new coordinator information with the staked amount
        CoordinatorInfo memory newCoordinator = CoordinatorInfo({
            name: _name,
            ownerAddress: msg.sender,
            functionAddresses: _functionAddresses,
            stakedAmount: _amount  // Storing the staked amount
        });

        coordinators.push(newCoordinator);

        // Emit an event or other logic as necessary
    }

    function addCoordinator(
        string memory _name,
        address _ownerAddress,
        address[] memory _functionAddresses,
        uint _stakedAmount) public onlyOwner {
        // existing checks and logic...

        CoordinatorInfo memory newCoordinator = CoordinatorInfo({
            name: _name,
            ownerAddress: _ownerAddress,
            functionAddresses: _functionAddresses,
            stakedAmount: _stakedAmount // Optional: Include this if staking is required
        });

        coordinators.push(newCoordinator);
    }

    function resignAsCoordinator() public onlyCoordinator {
        // Call the internal withdraw function
        withdraw();

        // Get the remaining balance of tokens in the contract
        uint remainingAmount = token.balanceOf(address(this));

        // Remove the coordinator from the list
        for (uint i = 0; i < coordinators.length; i++) {
            if (coordinators[i].ownerAddress == msg.sender) {
                // Shift other elements to fill the gap left by the removed coordinator
                for (uint j = i; j < coordinators.length - 1; j++) {
                    coordinators[j] = coordinators[j + 1];
                }
                // Remove the last element
                coordinators.pop();
                break;
            }
        }

        // Distribute the remaining 0.1% evenly among remaining coordinators
        uint numberOfCoordinators = coordinators.length;
        if (numberOfCoordinators > 0) {
            uint amountPerCoordinator = remainingAmount.div(numberOfCoordinators);
            for (uint i = 0; i < numberOfCoordinators; i++) {
                // Transfer the amount to each coordinator's owner address
                require(token.transfer(coordinators[i].ownerAddress, amountPerCoordinator), "Failed to distribute remaining stake");
            }
        }

        // Emit an event or other logic as necessary
    }

    function removeCoordinator(address _address) public onlyOwner {
        require(isCoordinatorAddress(_address), "Address is not a coordinator");

        for (uint i = 0; i < coordinators.length; i++) {
            // Check if the address is the owner address or one of the function addresses of a coordinator
            if (coordinators[i].ownerAddress == _address || isInFunctionAddresses(_address, coordinators[i].functionAddresses)) {
                // Shift elements to fill the gap left by the removed coordinator
                for (uint j = i; j < coordinators.length - 1; j++) {
                    coordinators[j] = coordinators[j + 1];
                }
                coordinators.pop(); // Remove the last element which is now a duplicate

                // Emit an event or handle other logic as necessary
                return;
            }
        }

        revert("Coordinator with the given address not found");
    }

    function replaceFunctionAddress(address oldAddress, address newAddress) public onlyOwner {
        require(isCoordinatorAddress(oldAddress), "Old address is not a coordinator function address");

        for (uint i = 0; i < coordinators.length; i++) {
            for (uint j = 0; j < coordinators[i].functionAddresses.length; j++) {
                if (coordinators[i].functionAddresses[j] == oldAddress) {
                    coordinators[i].functionAddresses[j] = newAddress;
                    return;
                }
            }
        }
        revert("Function address not found");
    }

    function deleteFunctionAddress(address _address) public onlyOwner {
        require(isCoordinatorAddress(_address), "Address is not a coordinator function address");

        for (uint i = 0; i < coordinators.length; i++) {
            require(coordinators[i].functionAddresses.length > 1, "Cannot delete the only function address");
            for (uint j = 0; j < coordinators[i].functionAddresses.length; j++) {
                if (coordinators[i].functionAddresses[j] == _address) {
                    // Remove the address by shifting elements
                    for (uint k = j; k < coordinators[i].functionAddresses.length - 1; k++) {
                        coordinators[i].functionAddresses[k] = coordinators[i].functionAddresses[k + 1];
                    }
                    coordinators[i].functionAddresses.pop();
                    return;
                }
            }
        }
        revert("Function address not found");
    }

    function replaceFunctionAddress(address oldAddress, address newAddress) public {
        uint coordinatorIndex = getCoordinatorIndex();

        for (uint j = 0; j < coordinators[coordinatorIndex].functionAddresses.length; j++) {
            if (coordinators[coordinatorIndex].functionAddresses[j] == oldAddress) {
                coordinators[coordinatorIndex].functionAddresses[j] = newAddress;
                return;
            }
        }
        revert("Function address not found");
    }

    function deleteFunctionAddress(address _address) public {
        uint coordinatorIndex = getCoordinatorIndex();
        require(coordinators[coordinatorIndex].functionAddresses.length > 1, "Cannot delete the only function address");

        for (uint j = 0; j < coordinators[coordinatorIndex].functionAddresses.length; j++) {
            if (coordinators[coordinatorIndex].functionAddresses[j] == _address) {
                for (uint k = j; k < coordinators[coordinatorIndex].functionAddresses.length - 1; k++) {
                    coordinators[coordinatorIndex].functionAddresses[k] = coordinators[coordinatorIndex].functionAddresses[k + 1];
                }
                coordinators[coordinatorIndex].functionAddresses.pop();
                return;
            }
        }
        revert("Function address not found");
    }

    function addFunctionAddress(address _address) public {
        uint coordinatorIndex = getCoordinatorIndex();
        require(coordinators[coordinatorIndex].functionAddresses.length < 3, "Cannot have more than 3 function addresses");

        coordinators[coordinatorIndex].functionAddresses.push(_address);
    }

    function transferCoordinatorDeposit(address newAddress) public onlyCoordinator {
        // Transfer the deposit to the new address
        transferDeposit(newAddress);

        // Update coordinator information in the list
        uint coordinatorIndex = getCoordinatorIndex();
        coordinators[coordinatorIndex].ownerAddress = newAddress;
    }

    function withdrawForAddress(address recipient) public onlyOwner {
        withdrawByAddress(recipient);
    }


    modifier onlyCoordinator() {
        require(isCoordinatorAddress(msg.sender), "Caller is not a coordinator");
        _;
    }

}
