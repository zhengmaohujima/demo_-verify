// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract BeggingContract {
    // 合约所有者
    address public owner;

    // 记录每个地址的捐赠金额
    mapping(address => uint256) private donations;

    // 捐赠事件（额外挑战）
    event Donation(address indexed donor, uint256 amount);

    // 仅允许合约所有者调用
    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }

    // 构造函数，部署合约时执行
    constructor() {
        owner = msg.sender;
    }

    /**
     * @dev 捐赠函数
     * 用户向合约发送 ETH
     */
    function donate() external payable {
        require(msg.value > 0, "Donation must be greater than 0");

        donations[msg.sender] += msg.value;

        emit Donation(msg.sender, msg.value);
    }

    /**
     * @dev 查询某个地址的捐赠金额
     */
    function getDonation(address donor) external view returns (uint256) {
        return donations[donor];
    }

    /**
     * @dev 合约所有者提取所有捐赠资金
     */
    function withdraw() external onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No funds to withdraw");

        payable(owner).transfer(balance);
    }

    /**
     * @dev 查看合约当前余额（辅助函数）
     */
    function getContractBalance() external view returns (uint256) {
        return address(this).balance;
    }
}
