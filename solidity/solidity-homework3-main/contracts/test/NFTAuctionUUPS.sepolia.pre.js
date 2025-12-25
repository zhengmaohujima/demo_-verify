// test/NFTAuctionUUPS.sepolia.js
const { expect } = require("chai");
const { ethers } = require("hardhat");
const { time } = require("@nomicfoundation/hardhat-network-helpers");

describe("NFTAuctionUUPS", function () {

    let owner, seller, bidder1, bidder2;

    this.timeout(300000); // 5分钟超时

    before(async function () {
        // 如果只配置了一个 PRIVATE_KEY，只能获取一个 signer
        [owner] = await ethers.getSigners();
        seller = new ethers.Wallet(process.env.SELLER_PRIVATE_KEY, ethers.provider);
        bidder1 = new ethers.Wallet(process.env.BIDDER1_PRIVATE_KEY, ethers.provider);
        bidder2 = new ethers.Wallet(process.env.BIDDER2_PRIVATE_KEY, ethers.provider);

    });


    describe("版本检查", function () {
        it("应该返回正确的版本", async function () {
            console.log("ERC20 deployed:",  owner.address);
            console.log("ERC20 seller:",  seller.address);
        });
    });
});