pragma solidity ^0.6.10;

contract PriceData {

    event Written(address source, string symbol, string currency, uint64 value);
    event NotWritten(uint256 blockTimestamp, uint64 messageTimestamp, uint64 currentValueTimestamp);

    struct Record {
        uint64 timestamp;
        uint64 value;
    }

    mapping(address => mapping(string => mapping(string => Record))) private data;


    function putPrice(string memory symbol, string memory currency, uint64 value, uint64 timestamp) external {
        address source = msg.sender;
        Record storage current = data[source][symbol][currency];
        if (timestamp > current.timestamp && timestamp < block.timestamp + 60 minutes) {
            data[source][symbol][currency] = Record(timestamp, value);
            emit Written(source, symbol, currency, value);
        } else {
            emit NotWritten(block.timestamp, timestamp, current.timestamp);
        }
    }
    function putEthPrice(string memory symbol, uint64 value, uint64 timestamp) external {
        this.putPrice(symbol, "ETH", value, timestamp);
    }
    function putUsdPrice(string memory symbol, uint64 value, uint64 timestamp) external {
        this.putPrice(symbol, "USD", value, timestamp);
    }


    function getPrice(address memory source, string memory symbol, string memory currency) external view returns (uint64, uint64) {
        address source = msg.sender;
        Record storage value = data[source][symbol][currency];
        return (value.value, value.timestamp);
    }
    function getEthPrice(address memory source, string memory symbol) external view returns (uint64, uint64) {
        return getPrice(source, source, "ETH");
    }
}
