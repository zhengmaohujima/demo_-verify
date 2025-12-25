// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract BinarySearch {
    /**
     * @dev 在有序数组中查找 target
     * @param nums 升序数组
     * @param target 目标值
     * @return index 找到返回索引，找不到返回 -1
     */
    function binarySearch(uint[] memory nums, uint target) public pure returns (int) {
        if (nums.length == 0) {
            return -1;
        }

        uint left = 0;
        uint right = nums.length - 1;

        while (left <= right) {
            uint mid = left + (right - left) / 2;

            if (nums[mid] == target) {
                return int(mid);
            } else if (nums[mid] < target) {
                left = mid + 1;
            } else {
                // 防止 uint 下溢
                if (mid == 0) {
                    break;
                }
                right = mid - 1;
            }
        }

        return -1;
    }
}
