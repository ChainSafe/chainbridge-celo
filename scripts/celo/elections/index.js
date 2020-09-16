const { privateKeyToAddress, privateKeyToPublicKey, trimLeading0x } = require('@celo/utils/lib/address');
const { getBlsPoP, getBlsPublicKey } = require('@celo/utils/lib/bls');
const Web3 = require('web3');
const Signatures = require('./contracts/Signatures.json');
const AddressLinkedList = require('./contracts/AddressLinkedList.json');
const AddressSortedLinkedList = require('./contracts/AddressSortedLinkedList.json');
const Registry = require('./contracts/Registry.json');
const Freezer = require('./contracts/Freezer.json');
const Accounts = require('./contracts/Accounts.json');
const Election = require('./contracts/Election.json');
const LockedGold = require('./contracts/LockedGold.json');
const Validators = require('./contracts/Validators.json');
const Random = require('./contracts/Random.json');
const BlockchainParameters = require('./contracts/BlockchainParameters.json');
const GoldToken = require('./contracts/GoldToken.json');
const EpochRewards = require('./contracts/EpochRewards.json');
const Proxy = require('./contracts/Proxy.json');
const Promise = require('bluebird');
const newKit = require('@celo/contractkit').newKit;
const kit = newKit('http://localhost:8546')
const web3 = kit.web3;
const nonceTracker = {};
const NULL_ADDRESS = '0x0000000000000000000000000000000000000000';
const keys = [
  '0x487e418d6566e9b3c6a85545be986e571b09d836d3b2d0bf38889cfc7dbc717a',
  '0x000000000000000000000000000000000000000000000000000000616c696365',
  '0x0000000000000000000000000000000000000000000000000000000000626f62',
];
const groupKeys = [
  '0x00000000000000000000000000000000000000000000000000636861726c6965',
  '0x0000000000000000000000000000000000000000000000000000000064617665',
  '0x0000000000000000000000000000000000000000000000000000000000657665',
];

async function sendTransaction(tx, privateKey, waitFor = 'transactionHash') {
  nonceTracker[privateKey] = nonceTracker[privateKey] || 0;
  tx.nonce = nonceTracker[privateKey]++;
  tx.gas = 20000000;
  //tx.gasPrice = 20000000000;
  tx.from = privateKeyToAddress(privateKey);
  const result = new Promise((resolve, reject) => {
    web3.eth.sendTransaction(tx)
    .once(waitFor, resolve)
    .on('error', reject);
  });
  return result;
}

async function deploy(contract, privateKey) {
  const receipt = await sendTransaction({
    data: contract.bytecode,
  }, privateKey, 'receipt');
  console.log(`${contract.contractName} deployed: ${receipt.contractAddress}`);
  return new web3.eth.Contract(contract.abi, receipt.contractAddress);
}

async function lockGold(accounts, lockedGold, value, privateKey) {
  const createAccountData = accounts.methods.createAccount().encodeABI();
  await sendTransaction({
    to: accounts._address,
    data: createAccountData,
  }, privateKey);

  const lockData = lockedGold.methods.lock().encodeABI();

  await sendTransaction({
    to: lockedGold._address,
    data: lockData,
    value: value.toString(10),
  }, privateKey);
}

async function registerValidatorGroup(
  name,
  accounts,
  lockedGold,
  validators,
  privateKey,
  lockedGoldValue
) {
  console.info(`    - ${name}`);
  console.info(`      - lock gold`);
  await lockGold(accounts, lockedGold, lockedGoldValue, privateKey);

  console.info(`      - setName`);
  const setNameData = accounts.methods.setName(name).encodeABI();
  await sendTransaction({
    data: setNameData,
    to: accounts._address,
  }, privateKey);

  console.info(`      - registerValidatorGroup`);
  const txData = validators.methods.registerValidatorGroup(1).encodeABI();

  await sendTransaction({
    data: txData,
    to: validators._address,
  }, privateKey, 'receipt');
}

async function registerValidator(
  accounts,
  lockedGold,
  validators,
  validatorPrivateKey,
  groupAddress,
  index
) {
  const valName = `Validator #${index}`;

  console.info(`    - lockGold ${valName}`);
  await lockGold(
    accounts,
    lockedGold,
    100,
    validatorPrivateKey
  );

  console.info(`    - setName ${valName}`);

  const setNameData = accounts.methods.setName(valName).encodeABI();
  await sendTransaction({
    data: setNameData,
    to: accounts._address,
  }, validatorPrivateKey);

  console.info(`    - registerValidator ${valName}`);
  const publicKey = privateKeyToPublicKey(validatorPrivateKey);
  const blsPublicKey = getBlsPublicKey(validatorPrivateKey);
  const blsPoP = getBlsPoP(privateKeyToAddress(validatorPrivateKey), validatorPrivateKey);

  const registerData = validators.methods.registerValidator(publicKey, blsPublicKey, blsPoP).encodeABI();

  await sendTransaction({
    data: registerData,
    to: validators._address,
  }, validatorPrivateKey);

  console.info(`    - affiliate ${valName}`);

  const affiliateData = validators.methods.affiliate(groupAddress).encodeABI();

  await sendTransaction({
    data: affiliateData,
    to: validators._address,
  }, validatorPrivateKey);

  console.info(`    - setAccountDataEncryptionKey ${valName}`);

  const registerDataEncryptionKeyData = accounts.methods.setAccountDataEncryptionKey(
    privateKeyToPublicKey(validatorPrivateKey)
  ).encodeABI();

  await sendTransaction({
    data: registerDataEncryptionKeyData,
    to: accounts._address,
  }, validatorPrivateKey, 'receipt');

  console.info(`    - done ${valName}`);
}

function replaceAll(where, what, withWhat) {
  return where.split(what).join(withWhat);
}

function not(bool) {
  return !bool;
}

async function wait(condition, sleepMs = 1000) {
  while (not(await condition())) {
    await Promise.delay(sleepMs);
  }
}

(async () => {
  const EPOCH_SIZE_PRECOMPILE = '0x00000000000000000000000000000000000000f8';
  const ownerPrivateKey = keys[0];
  const epochSize = web3.utils.hexToNumber(await web3.eth.call({to: EPOCH_SIZE_PRECOMPILE, data: '0x'}));

  console.log(`Epoch Size: ${epochSize}`);

  const freezer = await deploy(Freezer, ownerPrivateKey);

  const celoRandom = await deploy(Random, ownerPrivateKey);
  const randomInit = celoRandom.methods.initialize(120).encodeABI();
  await sendTransaction({
    data: randomInit,
    to: celoRandom._address,
  }, ownerPrivateKey);
  const blockchainParams = await deploy(BlockchainParameters, ownerPrivateKey);
  const blockchainParamsInit = blockchainParams.methods.initialize(
    0,
    0,
    0,
    1000000,
    30000000,
  ).encodeABI();
  await sendTransaction({
    data: blockchainParamsInit,
    to: blockchainParams._address,
  }, ownerPrivateKey);

  const addressLinkedList = await deploy(AddressLinkedList, ownerPrivateKey);
  const addressSortedLinkedList = await deploy(AddressSortedLinkedList, ownerPrivateKey);
  const signatures = await deploy(Signatures, ownerPrivateKey);

  Accounts.bytecode = replaceAll(
    Accounts.bytecode,
    '__Signatures____________________________',
    trimLeading0x(signatures._address)
  );

  Election.bytecode = replaceAll(
    Election.bytecode,
    '__AddressSortedLinkedList_______________',
    trimLeading0x(addressSortedLinkedList._address)
  );

  Validators.bytecode = replaceAll(
    Validators.bytecode,
    '__AddressLinkedList_____________________',
    trimLeading0x(addressLinkedList._address)
  );

  // Deploy Freezer
  // Deploy Random, initialize(...)
  // Deploy Registry, initialize()
  // Deploy Accounts, initialize(registry)
  // Deploy Validators, initialize(registry, ...)
  // Deploy EpochRewards, initialize(registry, ...)
  // Deploy LockedGold, initialize(registry, ...)
  // Deploy Election, initialize(registry, ...)
  // Deploy GoldToken, initialize(registry)
  // Registry.setAddressFor('Freezer', Freezer._address)
  // Registry.setAddressFor('Random', Election._address)
  // Registry.setAddressFor('Accounts', Accounts._address)
  // Registry.setAddressFor('Validators', Validators._address)
  // Registry.setAddressFor('LockedGold', LockedGold._address)
  // Registry.setAddressFor('Election', Election._address)
  // Registry.setAddressFor('EpochRewards', EpochRewards._address)
  // Registry.setAddressFor('GoldToken', GoldToken._address)

  const registryPrototype = await deploy(Registry, ownerPrivateKey);
  let tempRegistry;
  if ((await web3.eth.getCode('0x000000000000000000000000000000000000ce10')).length > 2) {
    const registryProxy = new web3.eth.Contract(Proxy.abi, '0x000000000000000000000000000000000000ce10');
    const registryInit = registryPrototype.methods.initialize().encodeABI();
    const setImplementationData = registryProxy.methods._setAndInitializeImplementation(
      registryPrototype._address,
      registryInit
    ).encodeABI();
    await sendTransaction({
      data: setImplementationData,
      to: registryProxy._address,
    }, ownerPrivateKey, 'receipt');
    tempRegistry = new web3.eth.Contract(Registry.abi, '0x000000000000000000000000000000000000ce10');
  } else {
    tempRegistry = registryPrototype;
    const registryInit = tempRegistry.methods.initialize().encodeABI();
    await sendTransaction({
      data: registryInit,
      to: tempRegistry._address,
    }, ownerPrivateKey);
  }
  const registry = tempRegistry;

  const goldToken = await deploy(GoldToken, ownerPrivateKey);
  const goldTokenInit = goldToken.methods.initialize(registry._address).encodeABI();
  await sendTransaction({
    data: goldTokenInit,
    to: goldToken._address,
  }, ownerPrivateKey);

  const accounts = await deploy(Accounts, ownerPrivateKey);
  const accountsInit = accounts.methods.initialize(registry._address).encodeABI();
  await sendTransaction({
    data: accountsInit,
    to: accounts._address,
  }, ownerPrivateKey);

  const validators = await deploy(Validators, ownerPrivateKey);
  const validatorsInit = validators.methods.initialize(
    registry._address,
    1, // uint256 groupRequirementValue
    1, // uint256 groupRequirementDuration
    1, // uint256 validatorRequirementValue
    1, // uint256 validatorRequirementDuration
    1, // uint256 validatorScoreExponent
    1, // uint256 validatorScoreAdjustmentSpeed
    10, // uint256 _membershipHistoryLength
    1, // uint256 _slashingMultiplierResetPeriod
    10, // uint256 _maxGroupSize
    10, // uint256 _commissionUpdateDelay
  ).encodeABI();
  await sendTransaction({
    data: validatorsInit,
    to: validators._address,
  }, ownerPrivateKey);

  const epochRewards = await deploy(EpochRewards, ownerPrivateKey);
  const epochRewardsInit = epochRewards.methods.initialize(
    registry._address,
    1, // targetVotingYieldInitial
    2, // targetVotingYieldMax
    1, // targetVotingYieldAdjustmentFactor
    1, // rewardsMultiplierMax
    1, // rewardsMultiplierUnderspendAdjustmentFactor
    1, // rewardsMultiplierOverspendAdjustmentFactor
    1, // _targetVotingGoldFraction
    1, // _targetValidatorEpochPayment
    1, // _communityRewardFraction
    privateKeyToAddress(ownerPrivateKey), // _carbonOffsettingPartner
    1, // _carbonOffsettingFraction
  ).encodeABI();
  await sendTransaction({
    data: epochRewardsInit,
    to: epochRewards._address,
  }, ownerPrivateKey);

  const lockedGold = await deploy(LockedGold, ownerPrivateKey);
  const lockedGoldInit = lockedGold.methods.initialize(
    registry._address,
    1,
  ).encodeABI();
  await sendTransaction({
    data: lockedGoldInit,
    to: lockedGold._address,
  }, ownerPrivateKey);

  const election = await deploy(Election, ownerPrivateKey);
  const electionInit = election.methods.initialize(
    registry._address,
    1, // minElectableValidators
    3, // maxElectableValidators
    3, // _maxNumGroupsVotedFor
    1, // _electabilityThreshold
  ).encodeABI();
  await sendTransaction({
    data: electionInit,
    to: election._address,
  }, ownerPrivateKey);

  await sendTransaction({
    data: registry.methods.setAddressFor('Freezer', freezer._address).encodeABI(),
    to: registry._address,
  }, ownerPrivateKey);
  await sendTransaction({
    data: registry.methods.setAddressFor('Random', celoRandom._address).encodeABI(),
    to: registry._address,
  }, ownerPrivateKey);
  await sendTransaction({
    data: registry.methods.setAddressFor('BlockchainParameters', blockchainParams._address).encodeABI(),
    to: registry._address,
  }, ownerPrivateKey);
  await sendTransaction({
    data: registry.methods.setAddressFor('Accounts', accounts._address).encodeABI(),
    to: registry._address,
  }, ownerPrivateKey);
  await sendTransaction({
    data: registry.methods.setAddressFor('Validators', validators._address).encodeABI(),
    to: registry._address,
  }, ownerPrivateKey);
  await sendTransaction({
    data: registry.methods.setAddressFor('LockedGold', lockedGold._address).encodeABI(),
    to: registry._address,
  }, ownerPrivateKey);
  await sendTransaction({
    data: registry.methods.setAddressFor('GoldToken', goldToken._address).encodeABI(),
    to: registry._address,
  }, ownerPrivateKey);
  await sendTransaction({
    data: registry.methods.setAddressFor('EpochRewards', epochRewards._address).encodeABI(),
    to: registry._address,
  }, ownerPrivateKey);
  await sendTransaction({
    data: registry.methods.setAddressFor('Election', election._address).encodeABI(),
    to: registry._address,
  }, ownerPrivateKey, 'receipt');

  // Create 3 groups
  // Lock gold in every group
  // Add 1 validator in every group
  // Vote every group for itself
  // Next election 3rd group revokes her votes
  // Next election 3rd group votes for itself again
  // Repeat elections

  console.info(`  * Registering validator groups.`);
  await registerValidatorGroup(
    'Group 0',
    accounts,
    lockedGold,
    validators,
    groupKeys[0],
    1000,
  );
  await registerValidatorGroup(
    'Group 1',
    accounts,
    lockedGold,
    validators,
    groupKeys[1],
    999,
  );
  await registerValidatorGroup(
    'Group 2',
    accounts,
    lockedGold,
    validators,
    groupKeys[2],
    998,
  );

  console.info(`  * Registering validators ...`)
  await registerValidator(
    accounts,
    lockedGold,
    validators,
    keys[0],
    privateKeyToAddress(groupKeys[0]),
    0
  );
  await registerValidator(
    accounts,
    lockedGold,
    validators,
    keys[1],
    privateKeyToAddress(groupKeys[1]),
    1
  );
  await registerValidator(
    accounts,
    lockedGold,
    validators,
    keys[2],
    privateKeyToAddress(groupKeys[2]),
    2
  );

  const add0Data = validators.methods.addFirstMember(
    privateKeyToAddress(keys[0]),
    NULL_ADDRESS,
    NULL_ADDRESS
  ).encodeABI();
  await sendTransaction({
    data: add0Data,
    to: validators._address,
  }, groupKeys[0], 'receipt');
  const add1Data = validators.methods.addFirstMember(
    privateKeyToAddress(keys[1]),
    NULL_ADDRESS,
    privateKeyToAddress(groupKeys[0])
  ).encodeABI();
  await sendTransaction({
    data: add1Data,
    to: validators._address,
  }, groupKeys[1], 'receipt');
  const add2Data = validators.methods.addFirstMember(
    privateKeyToAddress(keys[2]),
    NULL_ADDRESS,
    privateKeyToAddress(groupKeys[1])
  ).encodeABI();
  await sendTransaction({
    data: add2Data,
    to: validators._address,
  }, groupKeys[2], 'receipt');

  console.info(`  * Voting ...`)
  
  const vote0Data = election.methods.vote(
    privateKeyToAddress(groupKeys[0]),
    1000,
    privateKeyToAddress(groupKeys[1]),
    NULL_ADDRESS
  ).encodeABI();
  await sendTransaction({
    data: vote0Data,
    to: election._address,
  }, groupKeys[0], 'receipt');
  const vote1Data = election.methods.vote(
    privateKeyToAddress(groupKeys[1]),
    999,
    privateKeyToAddress(groupKeys[2]),
    privateKeyToAddress(groupKeys[0])
  ).encodeABI();
  await sendTransaction({
    data: vote1Data,
    to: election._address,
  }, groupKeys[1], 'receipt');

  let thirdVoting = true;
  while (true) {
    if (thirdVoting) {
      console.info('Third group voting.');
      const vote2Data = election.methods.vote(
        privateKeyToAddress(groupKeys[2]),
        998,
        NULL_ADDRESS,
        privateKeyToAddress(groupKeys[1])
      ).encodeABI();
      await sendTransaction({
        data: vote2Data,
        to: election._address,
      }, groupKeys[2], 'receipt');
    } else {
      console.info('Third group revoking votes.');
      const revokeVote2Data = election.methods.revokePending(
        privateKeyToAddress(groupKeys[2]),
        998,
        NULL_ADDRESS,
        privateKeyToAddress(groupKeys[1]),
        0
      ).encodeABI();
      await sendTransaction({
        data: revokeVote2Data,
        to: election._address,
      }, groupKeys[2], 'receipt');
    }
    thirdVoting = not(thirdVoting);

    console.info('Waiting for new epoch.');
    await wait(async () => {
      const block = await web3.eth.getBlockNumber();
      // console.log('Current block:', block);
      // Waiting till the first block of next epoch.
      return (block % epochSize) === 1;
    });
    const numberOfValidators = await election.methods.electValidatorSigners().call();
    console.info(`New epoch. Number of validators: ${numberOfValidators.length}`);
  }
})();

