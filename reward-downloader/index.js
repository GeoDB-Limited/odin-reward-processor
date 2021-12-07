const axios = require('axios');
const fs = require('fs');

(async function () {
    const nodeAddresses = ["35.195.4.110", "34.78.138.110", "34.78.239.23", "34.77.171.169"];

    let nodes = [];
    for (let i = 0; i < nodeAddresses.length; i++) {
        const result = await getNodeData(nodeAddresses[i]);
        nodes.push({[nodeAddresses[i]]: result});
    }

    var result = JSON.stringify({nodes: nodes});
    fs.writeFileSync("nodes_data.json", result);
}());

async function getNodeData(nodeAddress) {
    let result = {};
    const nodeData = await getNodeInfo(nodeAddress);
    result.node_data = nodeData;

    result.validators_delegations = {};
    result.outstanding_rewards = {};
    let addresses = {};
    // Getting validators data
    for (let i = 0; i < nodeData.validators.length; i++) {
        const validatorData = await getValidatorData(nodeAddress, nodeData.validators[i].operator_address);
        Object.assign(result.validators_delegations, {[nodeData.validators[i].operator_address]: validatorData});
        for (let j = 0; j < validatorData.delegation_responses.length; j++) {
            Object.assign(addresses, {[validatorData.delegation_responses[j].delegation.delegator_address]: true});
        }

        const outstandingRewardsData = await getOutstandingRewardsData(nodeAddress, nodeData.validators[i].operator_address);
        Object.assign(result.outstanding_rewards, {[nodeData.validators[i].operator_address]: outstandingRewardsData});
    }

    result.rewards = {};
    // Getting rewards
    for (const [key, value] of Object.entries(addresses)) {
        const accountData = await getDistributionData(nodeAddress, key);
        Object.assign(result.rewards, {[key]: accountData});
    }
    return result;
}

function execNodeValidators(response) {
    return response.data.validators.map(item => item.operator_address);
}

async function getNodeInfo(nodeAddress) {
    const url = `http://${nodeAddress}:1317/cosmos/staking/v1beta1/validators`;
    const response = await axios.get(url);
    if (response.status != 200) {
        console.error("Failed to make request: " + url);
    }
    // data.data.validators[i].operator_address
    return response.data;
}

async function getValidatorData(nodeAddress, validatorAddress) {
    const url = `http://${nodeAddress}:1317/cosmos/staking/v1beta1/validators/${validatorAddress}/delegations`;
    const response = await axios.get(url);
    if (response.status != 200) {
        console.error("Failed to make request: " + url);
    }
    return response.data;
}

async function getOutstandingRewardsData(nodeAddress, validatorAddress) {
    const url = `http://${nodeAddress}:1317/cosmos/distribution/v1beta1/validators/${validatorAddress}/outstanding_rewards`;
    const response = await axios.get(url);
    if (response.status != 200) {
        console.error("Failed to make request: " + url);
    }
    return response.data;
}

async function getDistributionData(nodeAddress, accountAddress) {
    const url = `http://${nodeAddress}:1317/cosmos/distribution/v1beta1/delegators/${accountAddress}/rewards`;
    const response = await axios.get(url);
    if (response.status != 200) {
        console.error("Failed to make request: " + url);
    }
    return response.data;
}