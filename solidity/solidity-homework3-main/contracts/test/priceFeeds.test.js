const { ethers } = require("hardhat");
const { expect } = require("chai");

describe("NFTAuctionUUPS Sepolia Price Feed Test", function () {
    let nftAuction;
    let deployer;

    const ETH_USD_FEED = "0x694AA1769357215DE4FAC081bf1f309aDC325306";
    const BTC_USD_FEED = "0x1b44F3514812d835EB1BDB0acB33d3fA3351Ee43";
    const NFT_AUCTION_ADDRESS = process.env.NFT_AUCTION_ADDRESS;

    before(async () => {
        [deployer] = await ethers.getSigners();
        nftAuction = await ethers.getContractAt("NFTAuctionUUPS", NFT_AUCTION_ADDRESS);
    });

    it("Set ETH/USD feed and read price", async () => {
        await nftAuction.setETHPriceFeed(ETH_USD_FEED);

        const amount = ethers.parseEther("1");
        const price = await nftAuction.getPriceInUSD(ethers.ZeroAddress, amount);
        console.log("1 ETH in USD (Sepolia):", price.toString());

        expect(price).to.be.gt(0);
    });

    it("Set BTC/USD feed and read price", async () => {
        const BTC_ADDRESS = "0x0000000000000000000000000000000000000002";

        const tx = await nftAuction.setPriceFeed(BTC_ADDRESS, BTC_USD_FEED);
        await tx.wait();

        const amount = ethers.parseUnits("1", 8);
        const price = await nftAuction.getPriceInUSD(BTC_ADDRESS, amount);
        console.log("1 BTC in USD (Sepolia):", price.toString());

        expect(price).to.be.gt(0);
    });
});