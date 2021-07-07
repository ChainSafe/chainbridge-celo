### Glossary

1. [Zk-Snark Based Syncing](#zk-snark-based-syncing)
2. [Epochs](#epochs)
3. [Proving Mechanism](#proving-mechanism)
4. [Seals](#seals)
5. [Validator Sync](#validator-sync)
6. [Block Submission](#block-submission)
7. [Validation](#block-validation)
8. [Finality](#finality)
9. [Elections](#elections)
10. [Accounts](#accounts)
11. [IBFT](#ibft)

---


### Zk-Snark Based Syncing
_Cerifies validity of 1 year's worth of epoch block headers_

- Instead of syncing 1 block header per day, it downloads 1 snark per year

### Epochs
_A fixed number of blocks, configured in the network's genesis block, during which the same validator set is used for consensus._

- Epoch = 1/day
- MaxValidators = 100
- At each epoch, every elected Validator must be re-elected to continue. Validators are selected in proportion to votes received for each Validator Group.

*Epoch Based Sycning: Light client only has to download one block header per epoch.

In each epoch block header there is the validator set diff (`ValidatorSetDiff`). 
- This allows for calculation of the current validator set of the current head block, explained further in [Validator Sync](#validator-sync).

### Proving Mechanisms
_A two-part proving mechanism used to verify transfers between Celo chains._

#### Merkle Proofs
_Also known as an "inclusion proof", serves as part-one of the proving mechanism by ensuring that the transaction that initiated a transfer was included within a block._

Similar to Ethereum, blocks within the Celo blockchain contain a bundle of transactions. These transactions are then used to construct a merkle trie with the resulting root of this trie (merkle root) becoming part of the block's header.

A Merkle Proof consists of:

- `TxRootHash`: the root hash of the merkle trie
- `Key`: the transaction hash we are trying to prove
- `Nodes`: the nodes of the tree that, when hashed together, produce the rootHash

Merkle Proof data structure found [here](https://github.com/ChainSafe/chainbridge-celo/blob/main/utils/msg.go#L43-L47):
```go
type MerkleProof struct {
	TxRootHash [32]byte // Expected root of trie, in our case should be transactionsRoot from block
	Key        []byte   // RLP encoding of tx index, for the tx we want to prove
	Nodes      []byte   // The actual proof, all the nodes of the trie that between leaf value and root
}
```

We can then pass this data along with the [Validator Signature Verification](#validator-signature-verification) data, shown below, as parameters into the `executeProposal` contract method. With this data, the contract is able to verify that the block header was signed by the provided APK and ensures the transaction was included in the block.

_Two snippets below detail how this data is passed in to the `executeProposal` contract call._

Code snippet from [proposal data](https://github.com/ChainSafe/chainbridge-celo/blob/main/chain/writer/proposal_data.go#L51-L61) detailing what parameters are required to construct the proposal data hash:
```go
func CreateProposalDataHash(data []byte, handler common.Address, mp *utils.MerkleProof, sv *utils.SignatureVerification) common.Hash
```


Code snippet from [the writer](https://github.com/ChainSafe/chainbridge-celo/blob/main/chain/writer/writer.go#L93-L103):
```go
dataHash := CreateProposalDataHash(data, handlerContract, m.MPParams, m.SVParams)

if !w.shouldVote(m, dataHash) {
    if w.proposalIsPassed(m.Source, m.DepositNonce, dataHash) {
        // We should not vote for this proposal but it is ready to be executed
        w.executeProposal(m, data, dataHash)
        return true
    } else {
        return false
    }
}
```

#### Validator Signature Verification
_Part-two of the proving mechanism works by verifying that the block containing a transaction is valid; this is necessary to ensure that the merkle proof from part-one can be trusted._

The Celo protocol uses [BLS](https://github.com/celo-org/celo-bls-go) signatures in its consensus to ultimately determine whether or not a particular block is valid.

Many BLS signatures over the same content can be combined into a single "aggregated signature." This allows several kilobytes of signatures to be compressed into fewer than 100 bytes, ensuring that the block headers remain compact and light client friendly.

- Compressing keys in this way drastically reduces size of block header.
- During verification, the verifier will be able to know quickly which validator keys signed and which did not.

To verify the signature we need:

- `AggregatePublicKey`: APK; combined public keys of the validators
- `BlockHash`: Hash of the block being proved
- `Signature`: The signature of the block being proved

Signature Verification data structure found [here](https://github.com/ChainSafe/chainbridge-celo/blob/main/utils/msg.go#L49-53):

```go
type SignatureVerification struct {
	AggregatePublicKey []byte      // Aggregated public key of block validators
	BlockHash          common.Hash // Hash of block we are proving
	Signature          []byte      // Signature of block we are proving
}
```

#### Istanbul Epoch Validator Definitions:
Bitmap: Memory organization structure which, in this case, allows for quick, efficient determination of which validators are included in this signature and which are not.

Signature: This is an aggregated BLS signature resulting from signatures by each validator that signed the block.

```go
type IstanbulEpochValidatorSetSeal struct {
    // Bitmap is a bitmap having an active bit for each validator that signed this epoch data
    Bitmap *big.Int
    // Signature is an aggregated BLS signature resulting from signatures by each validator that signed this block
    // 64 bytes long
    Signature []byte
}
```

* When a validator has committed to a block, then they will add their signature to the block header.
* The IBFT consensus calls for 2/3rds of the validators to sign these blocks in order to ensure finality.

### Seals
_Data produced by the consensus engine and proving the authorship of the block producer._

### Validator Sync (package: validatorsync)
_A service that stores validator data, performs operations on this data and retrieves it for use._

Some detailed features:
- Holds data structure for new `db` operations.
- Provides `db` methods that return latest block that was parsed wtihin the [epoch](#epochs); determines which validators are to assist in validating a block; constructs the aggregated public keys of the validators chosen for the validation of a block (this rotates); and then some functionality to synchronize the validators who were chosen to validate a block from the pool of eligible validators for that [epoch](#epochs).
- Determines the "validator diff", explained further below.

#### Validator Diff
_This method encodes the index of the validators that were removed and the public keys of the validators that were added._

The "validator diff" can be extracted from the first block of the epoch.
    - The exact format of the data encoded in the block header can be seen [here](https://github.com/celo-org/celo-blockchain/blob/b815edcefa2649bc9b73720d19b313dd02b10aeb/core/types/istanbul.go#L84-L97). 

The Listener will need to include the current APK (Aggregated Public Key) in each message it produces. 
Thus, it should always ensure the Validator Syncer is up to date, and will need to retrieve the APK from the Syncer on-demand.

### Block Submission
_Process by which the client asserts that a block is fit to be added to the blockchain._

- This means that the block is consistent with the world state and transitions from the state of the system to a new valid state.
- Celo processes **1 block every 5 seconds.**

### Block Validation
- A block is deemed valid if the block author had the authorship right for the slot during which the slot was built as well as if the transactions in the block constitute a valid transition of states.

### Finality
_The process of finalizing blocks is obtained by consecutive rounds of voting by validator nodes._

### Elections
_A validator election is carried out after the last block of an epoch, and any resulting changes to the validator set are written into that block's header._

- The active validator set is updated by running an election in the final block of each epoch, after processing transactions and Epoch Rewards.

### Accounts
- Any private key generated for use in the Celo protocol has a corresponding address. 
- The account address is the last 20 bytes of the hash of the corresponding public key, just as in Ethereum. 
- Celo account keys can be used to sign and send transactions on the Celo network.

### IBFT 
_Istanbul Byzantine Fault Tolerance_

A type of Byzantine Fault Tolerance (BFT) consensus algorithm in a Proof of Stake network in which a defined set of validator nodes broadcast signed messages between themselves in a sequence of steps to reach agreement even when up to a third of the total nodes are offline, faulty or malicious. When a quorum of validators have reached agreement, that decision is final.

Validator sets are chosen per [epoch](#epochs) and then a small subgroup of this pool is chosen randomly to sign and verify each block.

- In the original IBFT implementation, the block headers contained an array of signatures of the validators.
    - This was determined to be bulky and therefore made more efficient through Celo's implementation.

Resources:
- https://www.youtube.com/watch?v=FhXCenm1Yok
- https://docs.celo.org/overview
- https://docs.celo.org/validator-guide/summary/detailed
- https://docs.celo.org/validator-guide/summary/key-rotation
- https://docs.celo.org/celo-owner-guide/voting-validators#validator-elections
- https://docs.celo.org/celo-codebase/protocol/proof-of-stake/validator-elections
- https://github.com/celo-org/celo-blockchain/core/types/istanbul.go
