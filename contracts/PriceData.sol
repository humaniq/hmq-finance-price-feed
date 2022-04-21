pragma solidity ^0.6.10;

contract PriceData {

    struct Record {
        uint64 timestamp;
        uint64 value;
    }

    mapping(address => mapping(string => Record)) private data;


    function put(bytes memory message, uint64 timestamp, string memory key, uint64 value, bytes calldata signature) external returns (string memory) {

    }
}
