// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;


import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {IERC721} from "@openzeppelin/contracts/token/ERC721/IERC721.sol";


using SafeERC20 for IERC20;

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

contract NFTAuctionUUPS is
    Initializable,
    OwnableUpgradeable,
    UUPSUpgradeable,
    ReentrancyGuardUpgradeable
{

    struct Auction {
        address seller;
        address nftContract;
        uint256 tokenId;
        uint256 minPriceUSD;
        uint256 endTime;

        address highestBidder;
        uint256 highestBid;
        address highestBidToken;

        bool ended;
        bool nftClaimed;
        bool fundsClaimed;
    }

    uint256 public auctionCount;
    address public constant ETH = address(0);
    mapping(uint256 => Auction) public auctions;
    mapping(address => AggregatorV3Interface) public priceFeeds;

    event AuctionCreated(uint256 indexed auctionId);
    event BidPlaced(uint256 indexed auctionId, address bidder, uint256 amount, uint256 usdValue);
    event AuctionFinalized(uint256 indexed auctionId);
    event NFTClaimed(uint256 indexed auctionId, address winner);
    event FundsClaimed(uint256 indexed auctionId, address seller, uint256 amount);

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize() public initializer {
        __Ownable_init(msg.sender);
        __ReentrancyGuard_init();
    }

    function _authorizeUpgrade(address) internal override onlyOwner {}

    /* ========== Price ========== */
    function setETHPriceFeed(address feed) external onlyOwner {
        priceFeeds[ETH] = AggregatorV3Interface(feed);
    }

    function setPriceFeed(address token, address feed) external onlyOwner {
        priceFeeds[token] = AggregatorV3Interface(feed);
    }

    function getPriceInUSD(address token, uint256 amount) public view returns (uint256) {
        AggregatorV3Interface feed = priceFeeds[token];
        require(address(feed) != address(0), "Price feed not set");
        (, int256 price,,,) = feed.latestRoundData();
        require(price > 0, "Invalid price");
        return (amount * uint256(price)) / (10 ** feed.decimals());
    }

    /* ========== Auction ========== */

    function createAuction(
        address nftContract,
        uint256 tokenId,
        uint256 minPriceUSD,
        uint256 duration
    ) external {
        IERC721 nft = IERC721(nftContract);
        require(nft.ownerOf(tokenId) == msg.sender, "Not owner");
        require(duration > 0, "Invalid duration");

        auctionCount++;
        auctions[auctionCount] = Auction({
            seller: msg.sender,
            nftContract: nftContract,
            tokenId: tokenId,
            minPriceUSD: minPriceUSD,
            endTime: block.timestamp + duration,
            highestBidder: address(0),
            highestBid: 0,
            highestBidToken: address(0),
            ended: false,
            nftClaimed: false,
            fundsClaimed: false
        });

        nft.transferFrom(msg.sender, address(this), tokenId);
        emit AuctionCreated(auctionCount);
    }

    function bid(uint256 auctionId, address bidToken, uint256 amount) external payable nonReentrant {
        Auction storage a = auctions[auctionId];
        require(block.timestamp < a.endTime, "Auction ended");

        uint256 usdValue;
        uint256 currentHighestUSD = a.highestBid == 0
            ? 0
            : getPriceInUSD(a.highestBidToken, a.highestBid);

        if (bidToken == address(0)) {
            require(msg.value == amount, "ETH mismatch");
            usdValue = getPriceInUSD(address(0), amount);
        } else {
            require(msg.value == 0, "No ETH");
            IERC20(bidToken).safeTransferFrom(msg.sender, address(this), amount);
            usdValue = getPriceInUSD(bidToken, amount);
        }

        require(usdValue >= a.minPriceUSD, "Below min price");
        require(usdValue > currentHighestUSD, "Bid too low");

        // refund previous bidder
        if (a.highestBidder != address(0)) {
            if (a.highestBidToken == address(0)) {
                payable(a.highestBidder).transfer(a.highestBid);
            } else {
                IERC20(a.highestBidToken).safeTransfer(a.highestBidder, a.highestBid);
            }
        }

        a.highestBid = amount;
        a.highestBidToken = bidToken;
        a.highestBidder = msg.sender;

        emit BidPlaced(auctionId, msg.sender, amount, usdValue);
    }

    /* ========== Finalize & Claim ========== */

    function finalizeAuction(uint256 auctionId) external {
        Auction storage a = auctions[auctionId];
        require(block.timestamp >= a.endTime, "Not ended");
        require(!a.ended, "Already finalized");

        a.ended = true;
        emit AuctionFinalized(auctionId);
    }

    function claimNFT(uint256 auctionId) external nonReentrant {
        Auction storage a = auctions[auctionId];
        require(a.ended, "Not finalized");
        require(msg.sender == a.highestBidder, "Not winner");
        require(!a.nftClaimed, "NFT claimed");

        a.nftClaimed = true;
        IERC721(a.nftContract).transferFrom(address(this), msg.sender, a.tokenId);

        emit NFTClaimed(auctionId, msg.sender);
    }

    function claimFunds(uint256 auctionId) external nonReentrant {
        Auction storage a = auctions[auctionId];
        require(a.ended, "Not finalized");
        require(msg.sender == a.seller, "Not seller");
        require(!a.fundsClaimed, "Funds claimed");
        require(a.highestBidder != address(0), "No bids");

        a.fundsClaimed = true;

        if (a.highestBidToken == address(0)) {
            payable(msg.sender).transfer(a.highestBid);
        } else {
            IERC20(a.highestBidToken).transfer(msg.sender, a.highestBid);
        }

        emit FundsClaimed(auctionId, msg.sender, a.highestBid);
    }

    function getVersion() public pure virtual returns (string memory) {
        return "V1";
    }

}
