pragma solidity 0.5.16;
pragma experimental ABIEncoderV2;

import "./Owned.sol";
import "./MixinResolver.sol";
import "./interfaces/IERC20.sol";
import "./interfaces/ILock.sol";
import "./Lock.sol";
import "openzeppelin-solidity-2.3.0/contracts/math/SafeMath.sol";

contract BaseCoordinator is Lock, Owned, MixinResolver {

    using SafeMath for uint;
    using SafeMath for uint256;
    bytes32 public constant CONTRACT_MAI_PROXY = "MAI_PROXY";
    bytes32 public constant CONTRACT_LOCK = "LockCoordinatorsV1";
    bool public staking_initialized;


    struct CoordinatorInfo {
        string name;
        address ownerAddress;
        address[] functionAddresses;
        uint stakedAmount;}

    CoordinatorInfo[] public coordinators;

    constructor(
        address _lockOwner,
        uint _unlockTime,
        address _owner,
        address _resolver
    )
    Lock( _lockOwner, _unlockTime) // Initialize Lock
    Owned(_owner) // Initialize Owned
    MixinResolver(_resolver) // Initialize MixinResolver
    public
    {
        staking_initialized = false;
    }

    function _initialize() public onlyOwner {
        require(!staking_initialized, "Already initialized");
        token = maiTokenContract();
        staking_initialized = true;
    }

    function resolverAddressesRequired() public view returns (bytes32[] memory addresses) {
        addresses = new bytes32[](1);
        addresses[0] = CONTRACT_MAI_PROXY;
        addresses[1] = CONTRACT_LOCK;
    }

    function isCoordinatorAddress(address _address) public view returns (bool) {
        for (uint i = 0; i < coordinators.length; i++) {
            // Check if the address is the owner address of a coordinator
            if (coordinators[i].ownerAddress == _address) {
                return true;
            }

            // Check if the address is one of the function addresses of a coordinator
            for (uint j = 0; j < coordinators[i].functionAddresses.length; j++) {
                if (coordinators[i].functionAddresses[j] == _address) {
                    return true;
                }
            }
        }
        return false; // Address not found in any coordinators
    }

    function maiTokenContract() internal view returns (IERC20) {
        return IERC20(requireAndGetAddress(CONTRACT_MAI_PROXY));
    }

    function coordinatorLockContract() internal view returns (ILock) {
        return ILock(requireAndGetAddress(CONTRACT_LOCK));
    }

    function isInFunctionAddresses(address _address, address[] memory functionAddresses) internal pure returns (bool) {
        for (uint j = 0; j < functionAddresses.length; j++) {
            if (functionAddresses[j] == _address) {
                return true;
            }
        }
        return false;
    }

    function getCoordinatorIndex() internal view returns (uint) {
        for (uint i = 0; i < coordinators.length; i++) {
            if (coordinators[i].ownerAddress == msg.sender) {
                return i;
            }
        }
        revert("Caller is not the owner of any coordinator");
    }


}