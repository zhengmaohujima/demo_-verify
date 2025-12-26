const dotenv = require("dotenv");
const dotenvExpand = require("dotenv-expand");

dotenvExpand.expand(dotenv.config());

require("@nomicfoundation/hardhat-toolbox");
require("@openzeppelin/hardhat-upgrades");

// 添加这个 subtask 来包含 test 目录下的 sol 文件
const { subtask } = require("hardhat/config");
const { TASK_COMPILE_SOLIDITY_GET_SOURCE_PATHS } = require("hardhat/builtin-tasks/task-names");

subtask(TASK_COMPILE_SOLIDITY_GET_SOURCE_PATHS)
    .setAction(async (_, __, runSuper) => {
        const paths = await runSuper();
        const glob = require("glob");
        const testContracts = glob.sync("test/**/*.sol");
        return [...paths, ...testContracts];
    });

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
    solidity: {
        compilers: [
            {
                version: "0.8.28",
                settings: {
                    optimizer: {
                        enabled: true,
                        runs: 200,
                    },
                },
            },
            {
                version: "0.8.22",
                settings: {
                    optimizer: {
                        enabled: true,
                        runs: 200,
                    },
                },
            },
            {
                version: "0.8.21",
                settings: {
                    optimizer: {
                        enabled: true,
                        runs: 200,
                    },
                },
            },
        ],
    },
    paths: {
        sources: "./contracts",
        tests: "./test",
        cache: "./cache",
        artifacts: "./artifacts",
    },
    networks: {
        hardhat: {},
        sepolia: {
            url: process.env.SEPOLIA_RPC,
            accounts: [process.env.PRIVATE_KEY],
        },
    },
};