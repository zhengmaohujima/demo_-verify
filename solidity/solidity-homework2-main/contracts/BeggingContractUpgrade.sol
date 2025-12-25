// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract BeggingContractUpGrade {
    address public owner;

    // 捐赠记录
    mapping(address => uint256) private donations;

    // 前 3 名捐赠者
    address[3] public topDonors;

    // 捐赠时间窗口
    uint256 public startTime;
    uint256 public endTime;

    event Donation(address indexed donor, uint256 amount);

    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner");
        _;
    }

    constructor(uint256 _startTime, uint256 _endTime) {
        require(_startTime < _endTime, "Invalid time range");

        owner = msg.sender;
        startTime = _startTime;
        endTime = _endTime;
    }

    /**
     * 捐赠（仅在时间范围内）
     */
    function donate() external payable {
        require(msg.value > 0, "Donation must be > 0");
        require(
            block.timestamp >= startTime && block.timestamp <= endTime,
            "Donation not allowed now"
        );

        donations[msg.sender] += msg.value;

        _updateTopDonors(msg.sender);

        emit Donation(msg.sender, msg.value);
    }

    /**
     * 更新前 3 名捐赠者
     */
    function _updateTopDonors(address donor) internal {
        for (uint256 i = 0; i < 3; i++) {
            if (topDonors[i] == donor) {
                break;
            }

            if (
                topDonors[i] == address(0) ||
                donations[donor] > donations[topDonors[i]]
            ) {
                // 向后挪一位
                for (uint256 j = 2; j > i; j--) {
                    topDonors[j] = topDonors[j - 1];
                }
                topDonors[i] = donor;
                break;
            }
        }
    }

    /**
     * 查询某地址的捐赠金额
     */
    function getDonation(address donor) external view returns (uint256) {
        return donations[donor];
    }

    /**
     * 获取 Top 3 捐赠者及其金额
     */
    function getTopDonors()
    external
    view
    returns (address[3] memory donors, uint256[3] memory amounts)
    {
        for (uint256 i = 0; i < 3; i++) {
            donors[i] = topDonors[i];
            amounts[i] = donations[topDonors[i]];
        }
    }

    /**
     * 合约所有者提取所有资金
     */
    function withdraw() external onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No balance");

        payable(owner).transfer(balance);
    }

    /**
     * 查看合约余额
     */
    function getContractBalance() external view returns (uint256) {
        return address(this).balance;
    }
}
