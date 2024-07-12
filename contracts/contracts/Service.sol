pragma solidity ^0.5.16;

import "./Owner.sol";

contract Service is Owner {
    // Define a struct to represent each service
    struct ServiceStruct {
        string name;
        uint serviceType;
        uint weight;
        uint threshold;
        address beneficiary;
        uint beneficiaryPercentage;
        bytes32[] metadata;
    }

    // Authority address with exclusive permission
    address public authority;

    // Dynamic array to store a list of services
    ServiceStruct[] public services;

    // Constructor to set the initial authority
    constructor(address _authority, address _owner) public Owner(_owner) {
        require(_authority != address(0), "Authority address cannot be zero");
        authority = _authority;
    }

    // Modifier to restrict access to the authority
    modifier onlyAuthority() {
        require(msg.sender == authority, "Caller is not the authority");
        _;
    }

    // Function to change the authority (onlyOwner)
    function setAuthority(address _newAuthority) public onlyOwner {
        require(_newAuthority != address(0), "New authority address cannot be zero");
        authority = _newAuthority;
    }

    // Function to add a new service to the list (restricted to authority)
    function addService(
        string memory _name,
        uint _serviceType,
        uint _weight,
        uint _threshold,
        address _beneficiary,
        uint _beneficiaryPercentage,
        bytes32[] memory _metadata
    ) public onlyAuthority {
        require(_beneficiaryPercentage >= 1 && _beneficiaryPercentage <= 100, "Percentage must be between 1 and 100");
        services.push(ServiceStruct({
            name: _name,
            serviceType: _serviceType,
            weight: _weight,
            threshold: _threshold,
            beneficiary: _beneficiary,
            beneficiaryPercentage: _beneficiaryPercentage,
            metadata: _metadata
        }));
    }

    function removeService(uint index) public onlyAuthority {
        require(index < services.length, "Index out of bounds");

        // Move the last element to the place of the one to be removed
        services[index] = services[services.length - 1];

        // Remove the last element
        services.pop();
    }

    function removeServiceByName(string memory serviceName) public onlyAuthority {
        uint index = findServiceIndexByName(serviceName);
        require(index != uint(-1), "Service not found");

        // Call removeService with the found index
        removeService(index);
    }

    // Helper function to find the index of a service by name
    function findServiceIndexByName(string memory serviceName) internal view returns (uint) {
        for (uint i = 0; i < services.length; i++) {
            if (keccak256(abi.encodePacked(services[i].name)) == keccak256(abi.encodePacked(serviceName))) {
                return i;
            }
        }
        return uint(-1); // Return an invalid index if not found
    }
}