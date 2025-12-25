// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;


import "./NFTAuctionUUPS.sol";


contract NFTAuctionUUPSV2 is NFTAuctionUUPS {
    function getVersion() public pure override returns (string memory) {
        return "V2";
    }
}
