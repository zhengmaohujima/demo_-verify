// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract MyNFT is ERC721, Ownable {
    uint256 private _tokenIds;

    constructor() ERC721("MyNFT", "MNFT") Ownable(msg.sender) {}

    function mint(address to) external onlyOwner returns (uint256) {
        _tokenIds += 1;
        uint256 newItemId = _tokenIds;
        _mint(to, newItemId);
        return newItemId;
    }
}