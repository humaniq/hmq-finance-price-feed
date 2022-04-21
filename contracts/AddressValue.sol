pragma solidity >=0.4.25 <0.7.0;
pragma experimental ABIEncoderV2;

contract AddressValue {
    address value;

    function set(bytes memory message, bytes memory signature) external returns (address) {
        address source = source(message, signature);
        value = source;
        return source;
    }

    function get() public view returns (address) {
        return value;
    }

    function source(bytes memory message, bytes memory signature) public pure returns (address) {
        bytes32 hash = keccak256(message);
        return ecrecover(hash, v, r, s);
    }
}
