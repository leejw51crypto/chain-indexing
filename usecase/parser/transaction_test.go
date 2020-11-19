package parser_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/crypto-com/chainindex/entity/command"
	"github.com/crypto-com/chainindex/infrastructure/tendermint"
	"github.com/crypto-com/chainindex/usecase/coin"
	command_usecase "github.com/crypto-com/chainindex/usecase/command"
	"github.com/crypto-com/chainindex/usecase/model"
	"github.com/crypto-com/chainindex/usecase/parser"
	usecase_parser_test "github.com/crypto-com/chainindex/usecase/parser/test"
)

var _ = Describe("TransactionParser", func() {
	Describe("TxHash", func() {
		It("should return transaction hash from hex encouded tx data", func() {
			txHex := "Cp4CCowBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEmwKK3Rjcm8xNjV0emNyaDJ5bDgzZzhxZXF4dWVnMmc1Z3pndTU3eTNmZTNrYzMSK3Rjcm8xODRsdGEybHN5dTQ3dnd5cDJlOHptdGNhM2s1eXE4NXA2YzR2cDMaEAoIYmFzZXRjcm8SBDEwMDAKjAEKHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQSbAordGNybzE4NGx0YTJsc3l1NDd2d3lwMmU4em10Y2EzazV5cTg1cDZjNHZwMxIrdGNybzE2NXR6Y3JoMnlsODNnOHFlcXh1ZWcyZzVnemd1NTd5M2ZlM2tjMxoQCghiYXNldGNybxIEMjAwMBKsAQpRCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAgiLen9uwpvsreYibwgnQtzupil7kyNJl4oTG3Wl6oIEEgQKAggBGLdPClEKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDw9KBooWSrc6BvuMJTwDq4mkyy8aC+6I5uQ9H2sn+cDYSBAoCCAEYyk8SBBDAmgwaQMtPcJacL5aryCBZz7bL4vKrOLFi07rejX0nMvBRA7BSd09ywefL+VMSkC/UwqhHC28pRTHhEDiNApbxrIYBVvIaQE0+gltCOfawUGDJU9nXJJkFLPmjMKJMYvt3UtTMjPR2bws7l78EzaUfrjtbmrkIokoxAW8GBgTuhEkC2Frr6Q0="
			Expect(parser.TxHash(txHex)).To(Equal(
				"4936522F7391D425F2A93AD47576F8AEC3947DC907113BE8A2FBCFF8E9F2A416",
			))
		})
	})

	Describe("ParseTransactionCommands", func() {
		It("should parse Transaction commands when there is two Msg in one transaction", func() {
			txFeeParser := parser.NewTxDecoder("basetcro")
			block, _, _ := tendermint.ParseBlockResp(strings.NewReader(usecase_parser_test.ONE_TX_TWO_MSG_BLOCK_RESP))
			blockResults, _ := tendermint.ParseBlockResultsResp(strings.NewReader(usecase_parser_test.ONE_TX_TWO_MSG_BLOCK_RESULTS_RESP))

			cmds, err := parser.ParseTransactionCommands(
				txFeeParser,
				block,
				blockResults,
			)
			Expect(err).To(BeNil())
			Expect(cmds).To(HaveLen(1))
			expectedBlockHeight := int64(343358)
			Expect(cmds).To(Equal([]command.Command{command_usecase.NewCreateTransaction(
				expectedBlockHeight,
				model.CreateTransactionParams{
					TxHash:    "4936522F7391D425F2A93AD47576F8AEC3947DC907113BE8A2FBCFF8E9F2A416",
					Code:      0,
					Log:       "[{\"msgIndex\":0,\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"send\"},{\"key\":\"sender\",\"value\":\"tcro165tzcrh2yl83g8qeqxueg2g5gzgu57y3fe3kc3\"},{\"key\":\"module\",\"value\":\"bank\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"tcro184lta2lsyu47vwyp2e8zmtca3k5yq85p6c4vp3\"},{\"key\":\"sender\",\"value\":\"tcro165tzcrh2yl83g8qeqxueg2g5gzgu57y3fe3kc3\"},{\"key\":\"amount\",\"value\":\"1000basetcro\"}]}]},{\"msgIndex\":1,\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"send\"},{\"key\":\"sender\",\"value\":\"tcro184lta2lsyu47vwyp2e8zmtca3k5yq85p6c4vp3\"},{\"key\":\"module\",\"value\":\"bank\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"tcro165tzcrh2yl83g8qeqxueg2g5gzgu57y3fe3kc3\"},{\"key\":\"sender\",\"value\":\"tcro184lta2lsyu47vwyp2e8zmtca3k5yq85p6c4vp3\"},{\"key\":\"amount\",\"value\":\"2000basetcro\"}]}]}]",
					MsgCount:  2,
					Fee:       coin.MustNewCoinFromInt(int64(0)),
					GasWanted: "200000",
					GasUsed:   "80148",
				},
			)}))
		})

		It("should parse Transaction commands when there is transaction fee", func() {
			txFeeParser := parser.NewTxDecoder("basetcro")
			block, _, _ := tendermint.ParseBlockResp(strings.NewReader(usecase_parser_test.TX_WITH_FEE_BLOCK_RESP))
			blockResults, _ := tendermint.ParseBlockResultsResp(strings.NewReader(usecase_parser_test.TX_WITH_FEE_BLOCK_RESULTS_RESP))

			cmds, err := parser.ParseTransactionCommands(
				txFeeParser,
				block,
				blockResults,
			)
			Expect(err).To(BeNil())
			Expect(cmds).To(HaveLen(1))
			expectedBlockHeight := int64(377673)
			Expect(cmds).To(Equal([]command.Command{command_usecase.NewCreateTransaction(
				expectedBlockHeight,
				model.CreateTransactionParams{
					TxHash:    "2A2A64A310B3D0E84C9831F4353E188A6E63BF451975C859DF40C54047AC6324",
					Code:      0,
					Log:       "[{\"msgIndex\":0,\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"send\"},{\"key\":\"sender\",\"value\":\"tcro1feqh6ad9ytjkr79kjk5nhnl4un3wez0ynurrwv\"},{\"key\":\"module\",\"value\":\"bank\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"tcro1feqh6ad9ytjkr79kjk5nhnl4un3wez0ynurrwv\"},{\"key\":\"sender\",\"value\":\"tcro1feqh6ad9ytjkr79kjk5nhnl4un3wez0ynurrwv\"},{\"key\":\"amount\",\"value\":\"1000000000basetcro\"}]}]}]",
					MsgCount:  1,
					Fee:       coin.MustNewCoinFromString("8000000"),
					GasWanted: "80000000",
					GasUsed:   "62582",
				},
			)}))
		})

		It("should parse Transaction commands when transaction failed with fee", func() {
			txFeeParser := parser.NewTxDecoder("basetcro")
			block, _, _ := tendermint.ParseBlockResp(strings.NewReader(usecase_parser_test.FAILED_TX_WITH_FEE_BLOCK_RESP))
			blockResults, _ := tendermint.ParseBlockResultsResp(strings.NewReader(usecase_parser_test.FAILED_TX_WITH_FEE_BLOCK_RESULTS_RESP))

			cmds, err := parser.ParseTransactionCommands(
				txFeeParser,
				block,
				blockResults,
			)
			Expect(err).To(BeNil())
			Expect(cmds).To(HaveLen(1))
			expectedBlockHeight := int64(420301)
			Expect(cmds).To(Equal([]command.Command{command_usecase.NewCreateTransaction(
				expectedBlockHeight,
				model.CreateTransactionParams{
					TxHash:    "2A2A64A310B3D0E84C9831F4353E188A6E63BF451975C859DF40C54047AC6324",
					Code:      11,
					Log:       "out of gas in location: WriteFlat; gasWanted: 80000000, gasUsed: 80150021: out of gas",
					MsgCount:  1,
					Fee:       coin.MustNewCoinFromString("8000000"),
					GasWanted: "80000000",
					GasUsed:   "80150021",
				},
			)}))
		})

		It("should parse Transaction commands when transaction failed without fee", func() {
			txFeeParser := parser.NewTxDecoder("basetcro")
			block, _, _ := tendermint.ParseBlockResp(strings.NewReader(usecase_parser_test.FAILED_TX_WITHOUT_FEE_BLOCK_RESP))
			blockResults, _ := tendermint.ParseBlockResultsResp(strings.NewReader(usecase_parser_test.FAILED_TX_WITHOUT_FEE_BLOCK_RESULTS_RESP))

			cmds, err := parser.ParseTransactionCommands(
				txFeeParser,
				block,
				blockResults,
			)
			Expect(err).To(BeNil())
			Expect(cmds).To(HaveLen(1))
			expectedBlockHeight := int64(3245)
			Expect(cmds).To(Equal([]command.Command{command_usecase.NewCreateTransaction(
				expectedBlockHeight,
				model.CreateTransactionParams{
					TxHash:    "CDBA166168176BF7ECA2EAC9E9B49054F1BF4C8799B8C26CC0B9EE85CB93AF27",
					Code:      11,
					Log:       "out of gas in location: WriteFlat; gasWanted: 200000, gasUsed: 201420: out of gas",
					MsgCount:  5,
					Fee:       coin.Zero(),
					GasWanted: "200000",
					GasUsed:   "201420",
				},
			)}))
		})
	})
})
