// test/NFTAuctionUUPS.hardhat.js
const { expect } = require("chai");
const { ethers } = require("hardhat");
const { time } = require("@nomicfoundation/hardhat-network-helpers");

describe("NFTAuctionUUPS", function () {

    let auction, nft, mockERC20, ethPriceFeed, btcPriceFeed;
    let owner, seller, bidder1, bidder2;
    const ETH_ADDRESS = ethers.ZeroAddress;

    beforeEach(async function () {
        [owner, seller, bidder1, bidder2] = await ethers.getSigners();

        // 使用 upgrades 插件部署 UUPS 代理
        const NFTAuctionUUPS = await ethers.getContractFactory("NFTAuctionUUPS");
        auction = await upgrades.deployProxy(NFTAuctionUUPS, [], {
            kind: "uups"
        });

        const MyNFT = await ethers.getContractFactory("MyNFT");
        nft = await MyNFT.deploy();

        const MockERC20 = await ethers.getContractFactory("MockERC20");
        mockERC20 = await MockERC20.deploy("Mock BTC", "MBTC");

        const MockPriceFeed = await ethers.getContractFactory("MockPriceFeed");
        ethPriceFeed = await MockPriceFeed.deploy(8, 2000_00000000);
        btcPriceFeed = await MockPriceFeed.deploy(8, 40000_00000000);

        // 设置 ETH 价格预言机
        await auction.setETHPriceFeed(await ethPriceFeed.getAddress());

        await mockERC20.mint(bidder1.address, ethers.parseEther("100"));
        await mockERC20.mint(bidder2.address, ethers.parseEther("100"));
    });

    describe("价格预言机功能", function () {
        it("应该正确获取ETH的USD价格", async function () {
            const amount = ethers.parseEther("1");
            const price = await auction.getPriceInUSD(ETH_ADDRESS, amount);
            expect(price).to.equal(2000n * 10n**18n);
        });

        it("应该正确获取ERC20的USD价格", async function () {
            await auction.setPriceFeed(await mockERC20.getAddress(), await btcPriceFeed.getAddress());
            const amount = ethers.parseEther("1");
            const price = await auction.getPriceInUSD(await mockERC20.getAddress(), amount);
            expect(price).to.equal(40000n * 10n**18n);
        });

        it("未设置价格预言机应该失败", async function () {
            const randomToken = ethers.Wallet.createRandom().address;
            await expect(auction.getPriceInUSD(randomToken, ethers.parseEther("1")))
                .to.be.revertedWith("Price feed not set");
        });
    });

    describe("创建拍卖", function () {
        it("应该成功创建拍卖", async function () {
            await nft.mint(seller.address);
            const tokenId = 1;
            await nft.connect(seller).approve(await auction.getAddress(), tokenId);

            const minPriceUSD = ethers.parseEther("1000");
            const duration = 86400;

            await expect(auction.connect(seller).createAuction(
                await nft.getAddress(),
                tokenId,
                minPriceUSD,
                duration
            )).to.emit(auction, "AuctionCreated");

            const auctionCount = await auction.auctionCount();
            const auctionData = await auction.auctions(auctionCount);

            expect(auctionData.seller).to.equal(seller.address);
            expect(auctionData.nftContract).to.equal(await nft.getAddress());
            expect(auctionData.tokenId).to.equal(tokenId);
            expect(auctionData.minPriceUSD).to.equal(minPriceUSD);
            expect(auctionData.ended).to.be.false;
            expect(await nft.ownerOf(tokenId)).to.equal(await auction.getAddress());
        });

        it("非NFT所有者不能创建拍卖", async function () {
            await nft.mint(seller.address);
            const tokenId = 1;

            await expect(auction.connect(bidder1).createAuction(
                await nft.getAddress(),
                tokenId,
                ethers.parseEther("1000"),
                86400
            )).to.be.reverted;
        });

        it("持续时间为0应该失败", async function () {
            await nft.mint(seller.address);
            const tokenId = 1;
            await nft.connect(seller).approve(await auction.getAddress(), tokenId);

            await expect(auction.connect(seller).createAuction(
                await nft.getAddress(),
                tokenId,
                ethers.parseEther("1000"),
                0
            )).to.be.revertedWith("Invalid duration");
        });
    });

    describe("ETH出价", function () {
        let auctionId, tokenId;

        beforeEach(async function () {
            await nft.mint(seller.address);
            tokenId = 1;
            await nft.connect(seller).approve(await auction.getAddress(), tokenId);
            await auction.connect(seller).createAuction(
                await nft.getAddress(),
                tokenId,
                ethers.parseEther("1000"),
                86400
            );
            auctionId = await auction.auctionCount();
        });

        it("应该接受有效的ETH出价", async function () {
            const bidAmount = ethers.parseEther("1");

            await expect(auction.connect(bidder1).bid(auctionId, ETH_ADDRESS, bidAmount, { value: bidAmount }))
                .to.emit(auction, "BidPlaced")
                .withArgs(auctionId, bidder1.address, bidAmount, await auction.getPriceInUSD(ETH_ADDRESS, bidAmount));

            const auctionData = await auction.auctions(auctionId);
            expect(auctionData.highestBidder).to.equal(bidder1.address);
            expect(auctionData.highestBid).to.equal(bidAmount);
            expect(auctionData.highestBidToken).to.equal(ETH_ADDRESS);
        });

        it("ETH金额不匹配应该失败", async function () {
            const bidAmount = ethers.parseEther("1");
            await expect(auction.connect(bidder1).bid(auctionId, ETH_ADDRESS, bidAmount, { value: ethers.parseEther("0.5") }))
                .to.be.revertedWith("ETH mismatch");
        });

        it("低于最低价格应该失败", async function () {
            const bidAmount = ethers.parseEther("0.1");
            await expect(auction.connect(bidder1).bid(auctionId, ETH_ADDRESS, bidAmount, { value: bidAmount }))
                .to.be.revertedWith("Below min price");
        });

        it("应该退还前一个出价者的ETH", async function () {
            const bid1 = ethers.parseEther("1");
            await auction.connect(bidder1).bid(auctionId, ETH_ADDRESS, bid1, { value: bid1 });

            const balanceBefore = await ethers.provider.getBalance(bidder1.address);

            const bid2 = ethers.parseEther("1.5");
            await auction.connect(bidder2).bid(auctionId, ETH_ADDRESS, bid2, { value: bid2 });

            const balanceAfter = await ethers.provider.getBalance(bidder1.address);
            expect(balanceAfter - balanceBefore).to.equal(bid1);
        });

        it("出价不高于当前最高价应该失败", async function () {
            const bid1 = ethers.parseEther("2");
            await auction.connect(bidder1).bid(auctionId, ETH_ADDRESS, bid1, { value: bid1 });

            const bid2 = ethers.parseEther("1.5");
            await expect(auction.connect(bidder2).bid(auctionId, ETH_ADDRESS, bid2, { value: bid2 }))
                .to.be.revertedWith("Bid too low");
        });
    });

    describe("ERC20出价", function () {
        let auctionId, tokenId;

        beforeEach(async function () {
            await auction.setPriceFeed(await mockERC20.getAddress(), await btcPriceFeed.getAddress());

            await nft.mint(seller.address);
            tokenId = 1;
            await nft.connect(seller).approve(await auction.getAddress(), tokenId);
            await auction.connect(seller).createAuction(
                await nft.getAddress(),
                tokenId,
                ethers.parseEther("1000"),
                86400
            );
            auctionId = await auction.auctionCount();
        });

        it("应该接受有效的ERC20出价", async function () {
            const bidAmount = ethers.parseEther("0.1");
            await mockERC20.connect(bidder1).approve(await auction.getAddress(), bidAmount);

            await expect(auction.connect(bidder1).bid(auctionId, await mockERC20.getAddress(), bidAmount))
                .to.emit(auction, "BidPlaced");

            const auctionData = await auction.auctions(auctionId);
            expect(auctionData.highestBidder).to.equal(bidder1.address);
            expect(auctionData.highestBid).to.equal(bidAmount);
            expect(auctionData.highestBidToken).to.equal(await mockERC20.getAddress());
        });

        it("应该退还前一个出价者的ERC20", async function () {
            const bid1 = ethers.parseEther("0.1");
            await mockERC20.connect(bidder1).approve(await auction.getAddress(), bid1);
            await auction.connect(bidder1).bid(auctionId, await mockERC20.getAddress(), bid1);

            const balanceBefore = await mockERC20.balanceOf(bidder1.address);

            const bid2 = ethers.parseEther("0.15");
            await mockERC20.connect(bidder2).approve(await auction.getAddress(), bid2);
            await auction.connect(bidder2).bid(auctionId, await mockERC20.getAddress(), bid2);

            const balanceAfter = await mockERC20.balanceOf(bidder1.address);
            expect(balanceAfter - balanceBefore).to.equal(bid1);
        });

        it("ERC20出价时发送ETH应该失败", async function () {
            const bidAmount = ethers.parseEther("0.1");
            await mockERC20.connect(bidder1).approve(await auction.getAddress(), bidAmount);

            await expect(auction.connect(bidder1).bid(auctionId, await mockERC20.getAddress(), bidAmount, { value: ethers.parseEther("0.1") }))
                .to.be.revertedWith("No ETH");
        });
    });

    describe("混合出价场景", function () {
        let auctionId;

        beforeEach(async function () {
            await auction.setPriceFeed(await mockERC20.getAddress(), await btcPriceFeed.getAddress());

            await nft.mint(seller.address);
            await nft.connect(seller).approve(await auction.getAddress(), 1);
            await auction.connect(seller).createAuction(
                await nft.getAddress(),
                1,
                ethers.parseEther("1000"),
                86400
            );
            auctionId = await auction.auctionCount();
        });

        it("ETH出价可以被ERC20出价覆盖", async function () {
            const ethBid = ethers.parseEther("1");
            await auction.connect(bidder1).bid(auctionId, ETH_ADDRESS, ethBid, { value: ethBid });

            const erc20Bid = ethers.parseEther("0.06");
            await mockERC20.connect(bidder2).approve(await auction.getAddress(), erc20Bid);
            await auction.connect(bidder2).bid(auctionId, await mockERC20.getAddress(), erc20Bid);

            const auctionData = await auction.auctions(auctionId);
            expect(auctionData.highestBidder).to.equal(bidder2.address);
            expect(auctionData.highestBidToken).to.equal(await mockERC20.getAddress());
        });

        it("ERC20出价可以被ETH出价覆盖", async function () {
            const erc20Bid = ethers.parseEther("0.05");
            await mockERC20.connect(bidder1).approve(await auction.getAddress(), erc20Bid);
            await auction.connect(bidder1).bid(auctionId, await mockERC20.getAddress(), erc20Bid);

            const ethBid = ethers.parseEther("2.5");
            await auction.connect(bidder2).bid(auctionId, ETH_ADDRESS, ethBid, { value: ethBid });

            const auctionData = await auction.auctions(auctionId);
            expect(auctionData.highestBidder).to.equal(bidder2.address);
            expect(auctionData.highestBidToken).to.equal(ETH_ADDRESS);
        });
    });

    describe("拍卖结束", function () {
        let auctionId;

        beforeEach(async function () {
            await nft.mint(seller.address);
            await nft.connect(seller).approve(await auction.getAddress(), 1);
            await auction.connect(seller).createAuction(
                await nft.getAddress(),
                1,
                ethers.parseEther("1000"),
                86400
            );
            auctionId = await auction.auctionCount();

            const bidAmount = ethers.parseEther("1");
            await auction.connect(bidder1).bid(auctionId, ETH_ADDRESS, bidAmount, { value: bidAmount });
        });

        it("时间到期后应该能结束拍卖", async function () {
            await time.increase(86401);

            await expect(auction.finalizeAuction(auctionId))
                .to.emit(auction, "AuctionFinalized")
                .withArgs(auctionId);

            const auctionData = await auction.auctions(auctionId);
            expect(auctionData.ended).to.be.true;
        });

        it("时间未到不能结束拍卖", async function () {
            await expect(auction.finalizeAuction(auctionId))
                .to.be.revertedWith("Not ended");
        });

        it("已结束的拍卖不能再次结束", async function () {
            await time.increase(86401);
            await auction.finalizeAuction(auctionId);

            await expect(auction.finalizeAuction(auctionId))
                .to.be.revertedWith("Already finalized");
        });

        it("拍卖结束后不能继续出价", async function () {
            await time.increase(86401);
            await auction.finalizeAuction(auctionId);

            const bidAmount = ethers.parseEther("2");
            await expect(auction.connect(bidder2).bid(auctionId, ETH_ADDRESS, bidAmount, { value: bidAmount }))
                .to.be.revertedWith("Auction ended");
        });
    });

    describe("领取NFT", function () {
        let auctionId;

        beforeEach(async function () {
            await nft.mint(seller.address);
            await nft.connect(seller).approve(await auction.getAddress(), 1);
            await auction.connect(seller).createAuction(
                await nft.getAddress(),
                1,
                ethers.parseEther("1000"),
                86400
            );
            auctionId = await auction.auctionCount();

            const bidAmount = ethers.parseEther("1");
            await auction.connect(bidder1).bid(auctionId, ETH_ADDRESS, bidAmount, { value: bidAmount });

            await time.increase(86401);
            await auction.finalizeAuction(auctionId);
        });

        it("获胜者应该能领取NFT", async function () {
            await expect(auction.connect(bidder1).claimNFT(auctionId))
                .to.emit(auction, "NFTClaimed")
                .withArgs(auctionId, bidder1.address);

            expect(await nft.ownerOf(1)).to.equal(bidder1.address);

            const auctionData = await auction.auctions(auctionId);
            expect(auctionData.nftClaimed).to.be.true;
        });

        it("非获胜者不能领取NFT", async function () {
            await expect(auction.connect(bidder2).claimNFT(auctionId))
                .to.be.revertedWith("Not winner");
        });

        it("拍卖未结束不能领取NFT", async function () {
            await nft.mint(seller.address);
            await nft.connect(seller).approve(await auction.getAddress(), 2);
            await auction.connect(seller).createAuction(
                await nft.getAddress(),
                2,
                ethers.parseEther("1000"),
                86400
            );
            const newAuctionId = await auction.auctionCount();

            const bidAmount = ethers.parseEther("1");
            await auction.connect(bidder1).bid(newAuctionId, ETH_ADDRESS, bidAmount, { value: bidAmount });

            await expect(auction.connect(bidder1).claimNFT(newAuctionId))
                .to.be.revertedWith("Not finalized");
        });

        it("不能重复领取NFT", async function () {
            await auction.connect(bidder1).claimNFT(auctionId);

            await expect(auction.connect(bidder1).claimNFT(auctionId))
                .to.be.revertedWith("NFT claimed");
        });
    });

    describe("领取资金", function () {
        let auctionId;

        beforeEach(async function () {
            await nft.mint(seller.address);
            await nft.connect(seller).approve(await auction.getAddress(), 1);
            await auction.connect(seller).createAuction(
                await nft.getAddress(),
                1,
                ethers.parseEther("1000"),
                86400
            );
            auctionId = await auction.auctionCount();
        });

        it("卖家应该能领取ETH资金", async function () {
            const bidAmount = ethers.parseEther("1");
            await auction.connect(bidder1).bid(auctionId, ETH_ADDRESS, bidAmount, { value: bidAmount });

            await time.increase(86401);
            await auction.finalizeAuction(auctionId);

            const balanceBefore = await ethers.provider.getBalance(seller.address);
            const tx = await auction.connect(seller).claimFunds(auctionId);
            const receipt = await tx.wait();
            const gasUsed = receipt.gasUsed * receipt.gasPrice;
            const balanceAfter = await ethers.provider.getBalance(seller.address);

            expect(balanceAfter - balanceBefore + gasUsed).to.equal(bidAmount);

            const auctionData = await auction.auctions(auctionId);
            expect(auctionData.fundsClaimed).to.be.true;
        });

        it("卖家应该能领取ERC20资金", async function () {
            await auction.setPriceFeed(await mockERC20.getAddress(), await btcPriceFeed.getAddress());

            const bidAmount = ethers.parseEther("0.1");
            await mockERC20.connect(bidder1).approve(await auction.getAddress(), bidAmount);
            await auction.connect(bidder1).bid(auctionId, await mockERC20.getAddress(), bidAmount);

            await time.increase(86401);
            await auction.finalizeAuction(auctionId);

            const balanceBefore = await mockERC20.balanceOf(seller.address);
            await auction.connect(seller).claimFunds(auctionId);
            const balanceAfter = await mockERC20.balanceOf(seller.address);

            expect(balanceAfter - balanceBefore).to.equal(bidAmount);
        });

        it("非卖家不能领取资金", async function () {
            const bidAmount = ethers.parseEther("1");
            await auction.connect(bidder1).bid(auctionId, ETH_ADDRESS, bidAmount, { value: bidAmount });

            await time.increase(86401);
            await auction.finalizeAuction(auctionId);

            await expect(auction.connect(bidder1).claimFunds(auctionId))
                .to.be.revertedWith("Not seller");
        });

        it("拍卖未结束不能领取资金", async function () {
            const bidAmount = ethers.parseEther("1");
            await auction.connect(bidder1).bid(auctionId, ETH_ADDRESS, bidAmount, { value: bidAmount });

            await expect(auction.connect(seller).claimFunds(auctionId))
                .to.be.revertedWith("Not finalized");
        });

        it("没有出价不能领取资金", async function () {
            await time.increase(86401);
            await auction.finalizeAuction(auctionId);

            await expect(auction.connect(seller).claimFunds(auctionId))
                .to.be.revertedWith("No bids");
        });

        it("不能重复领取资金", async function () {
            const bidAmount = ethers.parseEther("1");
            await auction.connect(bidder1).bid(auctionId, ETH_ADDRESS, bidAmount, { value: bidAmount });

            await time.increase(86401);
            await auction.finalizeAuction(auctionId);

            await auction.connect(seller).claimFunds(auctionId);

            await expect(auction.connect(seller).claimFunds(auctionId))
                .to.be.revertedWith("Funds claimed");
        });
    });

    describe("完整拍卖流程", function () {
        it("应该完成完整的拍卖流程", async function () {
            await nft.mint(seller.address);
            await nft.connect(seller).approve(await auction.getAddress(), 1);

            await auction.connect(seller).createAuction(
                await nft.getAddress(),
                1,
                ethers.parseEther("1000"),
                86400
            );
            const auctionId = await auction.auctionCount();

            const bid1 = ethers.parseEther("1");
            await auction.connect(bidder1).bid(auctionId, ETH_ADDRESS, bid1, { value: bid1 });

            const bid2 = ethers.parseEther("1.5");
            await auction.connect(bidder2).bid(auctionId, ETH_ADDRESS, bid2, { value: bid2 });

            await time.increase(86401);
            await auction.finalizeAuction(auctionId);

            await auction.connect(bidder2).claimNFT(auctionId);
            expect(await nft.ownerOf(1)).to.equal(bidder2.address);

            const sellerBalanceBefore = await ethers.provider.getBalance(seller.address);
            const tx = await auction.connect(seller).claimFunds(auctionId);
            const receipt = await tx.wait();
            const gasUsed = receipt.gasUsed * receipt.gasPrice;
            const sellerBalanceAfter = await ethers.provider.getBalance(seller.address);

            expect(sellerBalanceAfter - sellerBalanceBefore + gasUsed).to.equal(bid2);
        });
    });

    describe("版本检查", function () {
        it("应该返回正确的版本", async function () {
            expect(await auction.getVersion()).to.equal("V1");
        });
    });
});