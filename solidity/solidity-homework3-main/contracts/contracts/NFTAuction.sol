// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

// Chainlink 正确路径
interface AggregatorV3Interface {
    function decimals() external view returns (uint8);
    function latestRoundData()
    external
    view
    returns (
        uint80 roundId,
        int256 answer,
        uint256 startedAt,
        uint256 updatedAt,
        uint80 answeredInRound
    );
}

contract NFTAuction is Ownable {
    struct Auction {
        address seller;
        address nftContract;
        uint256 tokenId;
        uint256 minPrice;
        uint256 endTime;
        address highestBidder;
        uint256 highestBid;
        address bidToken;
        bool ended;
    }

    uint256 public auctionCount;
    mapping(uint256 => Auction) public auctions;
    mapping(address => AggregatorV3Interface) public priceFeeds;

    event AuctionCreated(
        uint256 indexed auctionId,
        address indexed seller,
        address nftContract,
        uint256 tokenId,
        uint256 minPrice,
        uint256 endTime,
        address bidToken
    );
    event BidPlaced(uint256 indexed auctionId, address indexed bidder, uint256 amount);
    event AuctionEnded(uint256 indexed auctionId, address winner, uint256 amount);

    constructor() Ownable(msg.sender) {}

    function setPriceFeed(address token, address feed) external onlyOwner {
        priceFeeds[token] = AggregatorV3Interface(feed);
    }

    function createAuction(
        address nftContract,
        uint256 tokenId,
        uint256 minPrice,
        uint256 duration,
        address bidToken
    ) external {
        IERC721 nft = IERC721(nftContract);
        require(nft.ownerOf(tokenId) == msg.sender, "Not NFT owner");
        require(nft.isApprovedForAll(msg.sender, address(this)) || nft.getApproved(tokenId) == address(this), "Auction not approved");

        auctionCount += 1;
        auctions[auctionCount] = Auction({
            seller: msg.sender,
            nftContract: nftContract,
            tokenId: tokenId,
            minPrice: minPrice,
            endTime: block.timestamp + duration,
            highestBidder: address(0),
            highestBid: 0,
            bidToken: bidToken,
            ended: false
        });

        nft.transferFrom(msg.sender, address(this), tokenId);
        emit AuctionCreated(auctionCount, msg.sender, nftContract, tokenId, minPrice, block.timestamp + duration, bidToken);
    }

    function bidUSD(uint256 auctionId) external payable {
        Auction storage a = auctions[auctionId];
        require(block.timestamp < a.endTime, "Auction ended");
        require(a.bidToken == address(0), "ETH only");

        uint256 usd = getPriceInUSD(address(0), msg.value);
        require(usd >= a.minPrice, "Bid below min USD");
        require(usd > getPriceInUSD(address(0), a.highestBid), "Bid not higher than current highest");

        if (a.highestBidder != address(0)) {
            payable(a.highestBidder).transfer(a.highestBid);
        }

        a.highestBid = msg.value;
        a.highestBidder = msg.sender;
        emit BidPlaced(auctionId, msg.sender, msg.value);
    }


    function getPriceInUSD(address token, uint256 amount) public view returns (uint256) {
        AggregatorV3Interface feed = priceFeeds[token];
        require(address(feed) != address(0), "Price feed not set");

        (, int price,,,) = feed.latestRoundData();
        require(price > 0, "Invalid price");

        uint8 feedDecimals = feed.decimals();
        return (amount * uint256(price)) / (10 ** feedDecimals);
    }

    function endAuction(uint256 auctionId) external {
        Auction storage a = auctions[auctionId];
        require(block.timestamp >= a.endTime, "Auction not ended");
        require(!a.ended, "Auction already ended");

        a.ended = true;
        IERC721 nft = IERC721(a.nftContract);

        if (a.highestBidder != address(0)) {
            nft.transferFrom(address(this), a.highestBidder, a.tokenId);
            if (a.bidToken == address(0)) {
                payable(a.seller).transfer(a.highestBid);
            } else {
                IERC20(a.bidToken).transfer(a.seller, a.highestBid);
            }
            emit AuctionEnded(auctionId, a.highestBidder, a.highestBid);
        } else {
            nft.transferFrom(address(this), a.seller, a.tokenId);
            emit AuctionEnded(auctionId, address(0), 0);
        }
    }
}