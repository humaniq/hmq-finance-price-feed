pragma solidity >=0.4.25 <0.7.0;
pragma experimental ABIEncoderV2;

contract AddressValue {
    function source(bytes memory message, bytes memory signature) public pure returns (address) {
        (bytes32 r, bytes32 s, uint8 v) = abi.decode(signature, (bytes32, bytes32, uint8));
        bytes32 hash = keccak256(message);
        return ecrecover(hash, v, r, s);
    }
}
