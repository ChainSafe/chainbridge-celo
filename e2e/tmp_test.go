package e2e

//func (s *IntegrationTestSuite) TestCheckBalancesBalance() {
//erc20Contract2, err := erc20.NewERC20PresetMinterPauser(s.erc20ContractAddr, s.sender2.Sender)
//s.Nil(err)
//bobBalance, err := erc20Contract2.BalanceOf(s.sender.CallOpts, BobKp.CommonAddress())
//s.Nil(err)
//log.Debug().Msg(bobBalance.String())

//rec, err := s.sender2.Sender.TransactionReceipt(context.TODO(), common.HexToHash("0x23d14271e56ff495a9afdca7b91d2d4b250dbcdb783e1ea78647e600842f1fe0"))
//s.Nil(err)
//log.Debug().Msgf("%v", rec.Logs)

//bridgeContract, err := Bridge.NewBridge(s.bridgeAddr, s.sender2.Sender)
//s.Nil(err)
//dataHash := CreateProposalDataHash(data, handlerContract, m.MPParams, m.SVParams)
//
//res, err := bridgeContract.GetProposal(s.sender2.CallOpts, uint8(srcId), uint64(nonce), dataHash)
//s.Nil(err)
//log.Debug().Msgf("%v", res.Status)
//
//log.Debug().Msg(hexutils.BytesToHex(res.DataHash[:]))
//}

//func (s *IntegrationTestSuite) TestProposalEventSearch() {
//	query := buildQuery(s.bridgeAddr, pkg.ProposalEvent, big.NewInt(13300), big.NewInt(13600))
//	evts, err := s.sender2.Sender.FilterLogs(context.Background(), query)
//	if err != nil {
//		log.Error().Err(err).Msg("Failed to fetch logs")
//		return
//	}
//
//	// execute the proposal once we find the matching finalized event
//	for _, evt := range evts {
//		sourceId := evt.Topics[1].Big().Uint64()
//		depositNonce := evt.Topics[2].Big().Uint64()
//		status := evt.Topics[3].Big().Uint64()
//		log.Trace().Interface("src", sourceId).Interface("nonce", depositNonce).Uint64("status", status).Uint64("block", evt.BlockNumber).Msg("event")
//	}
//}
//func (s *IntegrationTestSuite) TestSimulate() {
//	_, err := simulate(s.sender2, big.NewInt(101), common.HexToHash("0xe0c7276f084e260b269b3722930af3a91f068a27121f0bcf6f1fc63f2d09d965"), AliceKp.CommonAddress())
//	s.Nil(err)
//}
//func (s *IntegrationTestSuite) TestCheckBalancesBalance() {
//erc20Contract2, err := erc20.NewERC20PresetMinterPauser(s.erc20ContractAddr, s.sender2.Sender)
//s.Nil(err)
//bobBalance, err := erc20Contract2.BalanceOf(s.sender.CallOpts, BobKp.CommonAddress())
//s.Nil(err)
//log.Debug().Msg(bobBalance.String())

//rec, err := s.sender2.Sender.TransactionReceipt(context.TODO(), common.HexToHash("0x23d14271e56ff495a9afdca7b91d2d4b250dbcdb783e1ea78647e600842f1fe0"))
//s.Nil(err)
//log.Debug().Msgf("%v", rec.Logs)

//bridgeContract, err := Bridge.NewBridge(s.bridgeAddr, s.sender2.Sender)
//s.Nil(err)
//dataHash := CreateProposalDataHash(data, handlerContract, m.MPParams, m.SVParams)
//
//res, err := bridgeContract.GetProposal(s.sender2.CallOpts, uint8(srcId), uint64(nonce), dataHash)
//s.Nil(err)
//log.Debug().Msgf("%v", res.Status)
//
//log.Debug().Msg(hexutils.BytesToHex(res.DataHash[:]))
//}

//func (s *IntegrationTestSuite) TestProposalEventSearch() {
//	query := buildQuery(s.bridgeAddr, pkg.ProposalEvent, big.NewInt(13300), big.NewInt(13600))
//	evts, err := s.sender2.Sender.FilterLogs(context.Background(), query)
//	if err != nil {
//		log.Error().Err(err).Msg("Failed to fetch logs")
//		return
//	}
//
//	// execute the proposal once we find the matching finalized event
//	for _, evt := range evts {
//		sourceId := evt.Topics[1].Big().Uint64()
//		depositNonce := evt.Topics[2].Big().Uint64()
//		status := evt.Topics[3].Big().Uint64()
//		log.Trace().Interface("src", sourceId).Interface("nonce", depositNonce).Uint64("status", status).Uint64("block", evt.BlockNumber).Msg("event")
//	}
//}
