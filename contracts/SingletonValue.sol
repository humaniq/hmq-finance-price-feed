pragma solidity >=0.4.25 <0.7.0;

contract SingletonValue {
    uint value;

    function set(uint x) public {
        value = x;
    }

    function get() public view returns (uint) {
        return value;
    }
}
