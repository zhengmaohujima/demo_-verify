const { ethers, upgrades } = require("hardhat");

async function main() {
    const NFT_AUCTION_ADDRESS = process.env.NFT_AUCTION_ADDRESS;
    const AuctionV2 = await ethers.getContractFactory("NFTAuctionUUPS_V2"); // 新版本合约
    const upgraded = await upgrades.upgradeProxy(NFT_AUCTION_ADDRESS, AuctionV2);
    const address = await upgraded.getAddress();
    console.log("Upgraded contract at:", address);
}

main();
