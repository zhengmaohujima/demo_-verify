const dotenv = require("dotenv");
const dotenvExpand = require("dotenv-expand");

// 👇 关键：只加载一次，用 expand 包住
dotenvExpand.expand(dotenv.config());

require("@nomicfoundation/hardhat-toolbox");
require("@openzeppelin/hardhat-upgrades");

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
    solidity:{
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
    networks: {
        hardhat: {},
        sepolia: {
            url: process.env.SEPOLIA_RPC,
            accounts: [process.env.PRIVATE_KEY],
        },
        // mainnet: {
        //     url: process.env.MAINNET_RPC,
        //     accounts: [process.env.PRIVATE_KEY],
        // },
    },
};
