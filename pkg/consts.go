package pkg

const (
	Deposit           EventSig = "Deposit(uint8,bytes32,uint64)"
	ProposalCreated   EventSig = "ProposalCreated(uint8,uint8,uint64,bytes32,bytes32)"
	ProposalVote      EventSig = "ProposalVote(uint8,uint8,uint64,bytes32,uint8)"
	ProposalFinalized EventSig = "ProposalFinalized(uint8,uint8,uint64,bytes32)"
	ProposalExecuted  EventSig = "ProposalExecuted(uint8,uint8,uint64)"
	ProposalEvent     EventSig = "ProposalEvent(uint8,uint64,uint8,bytes32,bytes32)"
)

type ProposalStatus int

const (
	Inactive ProposalStatus = iota
	Active
	Passed
	Executed
	Cancelled
)
