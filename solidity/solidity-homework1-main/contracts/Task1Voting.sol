// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract Task1Voting {

    // 存储候选人 => 得票数
    mapping(string => uint256) private votes;

    // 记录所有候选人，方便重置
    string[] private candidates;

    // 判断候选人是否已存在
    mapping(string => bool) private candidateExists;

    /**
     * @dev 给某个候选人投票
     */
    function vote(string memory candidate) public {
        // 如果是新候选人，加入候选人列表
        if (!candidateExists[candidate]) {
            candidateExists[candidate] = true;
            candidates.push(candidate);
        }

        votes[candidate] += 1;
    }

    /**
     * @dev 获取某个候选人的得票数
     */
    function getVotes(string memory candidate) public view returns (uint256) {
        return votes[candidate];
    }

    /**
     * @dev 重置所有候选人的得票数
     */
    function resetVotes() public {
        for (uint256 i = 0; i < candidates.length; i++) {
            votes[candidates[i]] = 0;
        }
    }
}
