const { ethers, upgrades } = require("hardhat");

async function main() {
    const Auction = await ethers.getContractFactory("NFTAuctionUUPS");

    const auction = await upgrades.deployProxy(Auction, [], {
        initializer: "initialize",
        kind: "uups",
    });

    await auction.waitForDeployment();

    const address = await auction.getAddress();
    console.log("UUPS NFT Auction deployed to:", address);
}

main().catch((error) => {
    console.error(error);
    process.exit(1);
});