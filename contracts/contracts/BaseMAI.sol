pragma solidity ^0.5.16;

// Inheritance
import "./interfaces/IERC20.sol";
import "./ExternStateToken.sol";
import "./MixinResolver.sol";

// Internal references
import "./TokenState.sol";
import "./interfaces/IExchanger.sol";
import "./interfaces/IRewardsDistribution.sol";
import "./interfaces/IRewardEscrowV2.sol";

contract BaseMAI is IERC20, ExternStateToken, MixinResolver {
    // ========== STATE VARIABLES ==========

    // Available Synths which can be used with the system
    string public constant TOKEN_NAME = "MintAI Network Token";
    string public constant TOKEN_SYMBOL = "MAI";
    uint8 public constant DECIMALS = 18;

    // ========== ADDRESS RESOLVER CONFIGURATION ==========

    bytes32 private constant CONTRACT_EXCHANGER = "Exchanger";
    bytes32 private constant CONTRACT_REWARDSDISTRIBUTION = "RewardsDistribution";
    bytes32 private constant CONTRACT_REWARDESCROW_V2 = "RewardEscrowV2";

    // ========== CONSTRUCTOR ==========

    constructor(
        address payable _proxy,
        TokenState _tokenState,
        address _owner,
        uint _totalSupply,
        address _resolver
    )
    public
    ExternStateToken(_proxy, _tokenState, TOKEN_NAME, TOKEN_SYMBOL, _totalSupply, DECIMALS, _owner)
    MixinResolver(_resolver)
    {}

    // ========== VIEWS ==========

    // Note: use public visibility so that it can be invoked in a subclass
    function resolverAddressesRequired() public view returns (bytes32[] memory addresses) {
        addresses = new bytes32[](3);
        addresses[0] = CONTRACT_EXCHANGER;
        addresses[1] = CONTRACT_REWARDSDISTRIBUTION;
        addresses[2] = CONTRACT_REWARDESCROW_V2;
    }

    function exchanger() internal view returns (IExchanger) {
        return IExchanger(requireAndGetAddress(CONTRACT_EXCHANGER));
    }

    function rewardsDistribution() internal view returns (IRewardsDistribution) {
        return IRewardsDistribution(requireAndGetAddress(CONTRACT_REWARDSDISTRIBUTION));
    }

    function rewardEscrowV2() internal view returns (IRewardEscrowV2) {
        return IRewardEscrowV2(requireAndGetAddress(CONTRACT_REWARDESCROW_V2));
    }

    // ========== MUTATIVE FUNCTIONS ==========

    function transfer(address to, uint value) external onlyProxyOrInternal returns (bool) {
        // Perform the transfer: if there is a problem an exception will be thrown in this call.
        _transferByProxy(messageSender, to, value);

        return true;
    }

    function transferFrom(
        address from,
        address to,
        uint value
    ) external onlyProxyOrInternal returns (bool) {

        // Perform the transfer: if there is a problem,
        // an exception will be thrown in this call.
        return _transferFromByProxy(messageSender, from, to, value);
    }

    function mint() external returns (bool) {
        _notImplemented();
    }

    function _notImplemented() internal pure {
        revert("Cannot be run on this layer");
    }

    // ========== MODIFIERS ==========

    modifier onlyExchanger() {
        _onlyExchanger();
        _;
    }

    function _onlyExchanger() private view {
        require(msg.sender == address(exchanger()), "Only Exchanger can invoke this");
    }

    modifier onlyProxyOrInternal {
        _onlyProxyOrInternal();
        _;
    }

    function _onlyProxyOrInternal() internal {
        if (msg.sender == address(proxy)) {
            // allow proxy through, messageSender should be already set correctly
            return;
        } else if (_isInternalTransferCaller(msg.sender)) {
            // optionalProxy behaviour only for the internal legacy contracts
            messageSender = msg.sender;
        } else {
            revert("Only the proxy can call");
        }
    }

    /// some legacy internal contracts use transfer methods directly on implementation
    /// which isn't supported due to SIP-238 for other callers
    function _isInternalTransferCaller(address caller) internal view returns (bool) {
        // These entries are not required or cached in order to allow them to not exist (==address(0))
        // e.g. due to not being available on L2 or at some future point in time.
        return
            caller == resolver.getAddress("RewardEscrowV2");
    }
}