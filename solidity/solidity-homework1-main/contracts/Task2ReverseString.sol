// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Task2ReverseString {
    function reverse(string memory input) public pure returns (string memory) {
        bytes memory strBytes = bytes(input);
        uint len = strBytes.length;

        bytes memory reversed = new bytes(len);

        for (uint i = 0; i < len; i++) {
            reversed[i] = strBytes[len - 1 - i];
        }

        return string(reversed);
    }
}
