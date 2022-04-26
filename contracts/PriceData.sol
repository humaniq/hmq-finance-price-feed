pragma solidity ^0.6.10;

contract PriceData {

    event Written(string symbol, string currency, uint64 value);
    event NotWritten(uint256 blockTimestamp, uint64 messageTimestamp, uint64 currentValueTimestamp);

    struct Record {
        uint64 timestamp;
        uint64 value;
    }

    mapping(string => mapping(string => Record)) private data;


    function putPrice(string memory symbol, string memory currency, uint64 value, uint64 timestamp) external {
        Record storage current = data[symbol][currency];
        if (timestamp > current.timestamp && timestamp < block.timestamp + 60 minutes) {
            data[symbol][currency] = Record(timestamp, value);
            emit Written(symbol, currency, value);
        } else {
            emit NotWritten(block.timestamp, timestamp, current.timestamp);
        }
    }
    function putEth(string memory symbol, uint64 value, uint64 timestamp) external {
        this.putPrice(symbol, "ETH", value, timestamp);
    }

    function getPrice(string memory symbol, string memory currency) external view returns (uint64, uint64) {
        Record storage value = data[symbol][currency];
        return (value.value, value.timestamp);
    }
    function getEth(string memory symbol) external view returns (uint64, uint64) {
        Record storage value = data[symbol]["ETH"];
        return (value.value, value.timestamp);
    }
}
