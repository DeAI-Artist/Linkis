pragma solidity ^0.5.16;

contract Owned {
    address public owner;
    address public nominatedOwner;

    constructor(address _owner) public {
        require(_owner != address(0), "Owner address cannot be zero");
        owner = _owner;
        emit OwnerChanged(address(0), _owner);
    }

    function nominateNewOwner(address _owner) external onlyOwner {
        nominatedOwner = _owner;
        emit OwnerNominated(_owner);
    }

    function acceptOwnership() external {
        require(msg.sender == nominatedOwner, "You must be nominated before you can accept ownership");
        emit OwnerChanged(owner, nominatedOwner);
        owner = nominatedOwner;
        nominatedOwner = address(0);
    }

    modifier onlyOwner {
        _onlyOwner();
        _;
    }

    function _onlyOwner() private view {
        require(msg.sender == owner, "Only the contract owner may perform this action");
    }

    event OwnerNominated(address newOwner);
    event OwnerChanged(address oldOwner, address newOwner);
}

contract Validator is Owned {
    // Array to store validator addresses
    address[] public validators;
    mapping(address => bool) private isValidatorMap;

    // Event emitted when a validator is added
    event ValidatorRegistered(address indexed validator);

    // Constructor that sets the initial owner
    constructor(address _owner) public Owned(_owner) {}

    // Function to register a validator, restricted to onlyOwner
    function registerValidator(address _validator) external onlyOwner {
        require(_validator != address(0), "Validator address cannot be zero");
        require(!isValidatorMap[_validator], "Address is already a validator");

        validators.push(_validator);
        isValidatorMap[_validator] = true;

        emit ValidatorRegistered(_validator);
    }

    // Public view function to check if an address is a validator
    function isValidator(address _address) external view returns (bool) {
        return isValidatorMap[_address];
    }

    // Function to get the total count of validators
    function getValidatorCount() external view returns (uint) {
        return validators.length;
    }

    // Function to get a validator by index
    function getValidator(uint index) external view returns (address) {
        require(index < validators.length, "Index out of bounds");
        return validators[index];
    }
}
