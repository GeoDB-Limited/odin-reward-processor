package bank

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	txClient "github.com/cosmos/cosmos-sdk/client/tx"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"google.golang.org/grpc"
)

type CosmosAccount struct {
	PrivateKey      cryptotypes.PrivKey
	AccountSequence uint64
	AccountNumber   uint64
	AccountAddress  types.AccAddress
}

type TxProcessor struct {
	TxBuilder client.TxBuilder
	EncCfg    simappparams.EncodingConfig
	ChainID   string
	Signers   []CosmosAccount
}

func NewTxProcessor(chainID string, signers []CosmosAccount) TxProcessor {
	var t TxProcessor

	t.EncCfg = simapp.MakeTestEncodingConfig()
	t.TxBuilder = t.EncCfg.TxConfig.NewTxBuilder()
	t.ChainID = chainID
	t.Signers = signers

	return t
}

func (t *TxProcessor) createMsg(msg *banktypes.MsgSend) error {
	err := t.TxBuilder.SetMsgs(msg)
	if err != nil {
		return err
	}

	// TODO add configuration for timeout
	// as t.TxBuilder.SetTimeoutHeight(6000)

	return nil
}

func (t *TxProcessor) signTx() error {
	var sigsV2 []signing.SignatureV2
	for _, signer := range t.Signers {
		sigV2 := signing.SignatureV2{
			PubKey: signer.PrivateKey.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode:  t.EncCfg.TxConfig.SignModeHandler().DefaultMode(),
				Signature: nil,
			},
			Sequence: signer.AccountSequence,
		}

		sigsV2 = append(sigsV2, sigV2)
	}
	err := t.TxBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return err
	}

	// Second round: all signer infos are set, so each signer can sign.
	sigsV2 = []signing.SignatureV2{}
	for _, signer := range t.Signers {
		signerData := xauthsigning.SignerData{
			ChainID:       t.ChainID,
			AccountNumber: signer.AccountNumber,
			Sequence:      signer.AccountSequence,
		}
		sigV2, err := txClient.SignWithPrivKey(
			t.EncCfg.TxConfig.SignModeHandler().DefaultMode(), signerData,
			t.TxBuilder, signer.PrivateKey, t.EncCfg.TxConfig, signer.AccountSequence)
		if err != nil {
			return err
		}

		sigsV2 = append(sigsV2, sigV2)
	}
	err = t.TxBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return err
	}
	return nil
}

func (t *TxProcessor) broadcastTx() error {
	// Generated Protobuf-encoded bytes.
	txBytes, err := t.EncCfg.TxConfig.TxEncoder()(t.TxBuilder.GetTx())
	if err != nil {
		return err
	}

	// Generate a JSON string.
	//txJSONBytes, err := b.EncCfg.TxConfig.TxJSONEncoder()(b.TxBuilder.GetTx())
	//if err != nil {
	//	return err
	//}
	//txJSON := string(txJSONBytes)

	// Create a connection to the gRPC server.
	grpcConn, _ := grpc.Dial(
		"127.0.0.1:9090",    // Or your gRPC server address.
		grpc.WithInsecure(), // The Cosmos SDK doesn't support any transport security mechanism.
	)
	defer grpcConn.Close()

	// Broadcast the tx via gRPC. We create a new client for the Protobuf Tx
	// service.
	txClient := tx.NewServiceClient(grpcConn)
	// We then call the BroadcastTx method on this client.
	grpcRes, err := txClient.BroadcastTx(
		context.TODO(),
		&tx.BroadcastTxRequest{
			Mode:    tx.BroadcastMode_BROADCAST_MODE_SYNC,
			TxBytes: txBytes, // Proto-binary of the signed transaction, see previous step.
		},
	)
	if err != nil {
		return err
	}

	fmt.Println(grpcRes.TxResponse.Code) // Should be `0` if the tx is successful

	return nil
}

func (t *TxProcessor) ProcessTx(msg *banktypes.MsgSend) error {
	if err := t.createMsg(msg); err != nil {
		return err
	}

	if err := t.signTx(); err != nil {
		return err
	}

	if err := t.broadcastTx(); err != nil {
		return err
	}

	return nil
}
