// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Task3RomanToInteger {
    function romanToInt(string memory s) public pure returns (uint256) {
        bytes memory b = bytes(s);
        uint256 result = 0;

        for (uint256 i = 0; i < b.length; i++) {
            uint256 curr = value(b[i]);

            // 如果当前值小于下一个值，说明是减法规则
            if (i + 1 < b.length && curr < value(b[i + 1])) {
                result -= curr;
            } else {
                result += curr;
            }
        }

        return result;
    }

    function value(bytes1 c) internal pure returns (uint256) {
        if (c == "I") return 1;
        if (c == "V") return 5;
        if (c == "X") return 10;
        if (c == "L") return 50;
        if (c == "C") return 100;
        if (c == "D") return 500;
        if (c == "M") return 1000;
        revert("Invalid Roman character");
    }
}
