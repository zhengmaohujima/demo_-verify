// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Task5MergeSortedArray {

    function merge(
        uint[] memory a,
        uint[] memory b
    ) public pure returns (uint[] memory) {

        uint i = 0;
        uint j = 0;
        uint k = 0;

        uint[] memory result = new uint[](a.length + b.length);

        // 同时遍历两个数组
        while (i < a.length && j < b.length) {
            if (a[i] <= b[j]) {
                result[k] = a[i];
                i++;
            } else {
                result[k] = b[j];
                j++;
            }
            k++;
        }

        // a 有剩余
        while (i < a.length) {
            result[k] = a[i];
            i++;
            k++;
        }

        // b 有剩余
        while (j < b.length) {
            result[k] = b[j];
            j++;
            k++;
        }

        return result;
    }
}
