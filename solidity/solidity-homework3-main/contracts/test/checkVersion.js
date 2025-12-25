const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("Check deployed NFTAuctionUUPS version", function () {
    it("Reads version from deployed proxy contract", async function () {
        // 这里填你已部署的代理合约地址
        const proxyAddress = "0xd962100895F7f02a459191Fc306147AbA6F7a104";

        // 获取已部署合约实例
        const auction = await ethers.getContractAt("NFTAuctionUUPS", proxyAddress);

        // 调用 getVersion 方法
        const version = await auction.getVersion();
        console.log("Deployed contract version:", version);

        // 可加断言
        expect(version).to.be.a("string");
    });
});
