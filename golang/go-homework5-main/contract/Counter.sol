// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;
//0xaC7d68AE475b14d0F0067E64e2862fBBb92D47d8
contract Counter {
    uint256 private count;

    // 获取当前计数
    function get() public view returns (uint256) {
        return count;
    }

    // 计数 +1
    function increment() public {
        count += 1;
    }

    // 计数 -1
    function decrement() public {
        require(count > 0, "count is already zero");
        count -= 1;
    }
}
