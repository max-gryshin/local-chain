package service_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"local-chain/internal/pkg/crypto"
	"testing"

	"local-chain/internal/service"

	"local-chain/internal/types"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestTransactor_CreateTx(t1 *testing.T) {
	type args struct {
		txReq  *types.TransactionRequest
		txPool service.TransactionStore
	}
	tests := []struct {
		name       string
		args       func(ctrl *gomock.Controller) args
		transactor func(args args) *service.Transactor
		want       func(args args) *types.Transaction
		wantErr    bool
	}{
		{
			name: "ok full amount",
			args: func(ctrl *gomock.Controller) args {
				from := crypto.GenerateKeyEllipticP256()
				to := crypto.GenerateKeyEllipticP256()
				txStore := NewMockTransactionStore(ctrl)
				tx1 := types.NewTransaction().WithOutput(types.NewAmount(30), &from.PublicKey)
				txStore.EXPECT().Get(tx1.GetHash()).Return(tx1, nil).Times(1)
				tx2 := types.NewTransaction().WithOutput(types.NewAmount(50), &from.PublicKey)
				txStore.EXPECT().Get(tx2.GetHash()).Return(tx2, nil).Times(1)
				tx3 := types.NewTransaction().WithOutput(types.NewAmount(20), &from.PublicKey)
				txStore.EXPECT().Get(tx3.GetHash()).Return(tx3, nil).Times(1)

				return args{
					txReq: &types.TransactionRequest{
						Sender:   from,
						Receiver: &to.PublicKey,
						Amount:   *types.NewAmount(100),
						Utxos: []*types.UTXO{
							{
								TxHash: tx1.GetHash(),
								Index:  0,
							},
							{
								TxHash: tx2.GetHash(),
								Index:  0,
							},
							{
								TxHash: tx3.GetHash(),
								Index:  0,
							},
						},
					},
					txPool: txStore,
				}
			},
			transactor: func(args args) *service.Transactor {
				t := service.NewTransactor(args.txPool)

				return t
			},
			want: func(args args) *types.Transaction {
				return types.NewTransaction().WithOutput(types.NewAmount(100), args.txReq.Receiver)
			},
			wantErr: false,
		},
		{
			name: "err sender is not owner of the output",
			args: func(ctrl *gomock.Controller) args {
				from, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
				require.NoError(t1, err)
				fakeFrom, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
				require.NoError(t1, err)
				fakeFrom.PublicKey = from.PublicKey // make fakeFrom have the same public key as "from"
				to, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
				require.NoError(t1, err)
				pool := NewMockTransactionStore(ctrl)
				utxoTx1 := &types.UTXO{
					TxHash: types.NewTransaction().GetHash(),
					Index:  0,
				}
				r, s, err := utxoTx1.Sign(from)
				require.NoError(t1, err)
				tx1 := types.NewTransaction().WithOutput(types.NewAmount(100), &from.PublicKey)
				tx1.AddInput(&types.TxIn{
					Prev:       utxoTx1,
					PubKey:     crypto.PublicKeyToBytes(&from.PublicKey),
					SignatureR: r,
					SignatureS: s,
					NSequence:  0,
				})
				pool.EXPECT().Get(tx1.GetHash()).Return(tx1, nil).Times(1)

				return args{
					txReq: &types.TransactionRequest{
						Sender:   fakeFrom,
						Receiver: &to.PublicKey,
						Amount:   *types.NewAmount(100),
						Utxos: []*types.UTXO{
							{
								TxHash: tx1.GetHash(),
								Index:  0,
							},
						},
					},
					txPool: pool,
				}
			},
			transactor: func(args args) *service.Transactor {
				t := service.NewTransactor(args.txPool)

				return t
			},
			want: func(args args) *types.Transaction {
				return types.NewTransaction().WithOutput(types.NewAmount(100), args.txReq.Receiver)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			ctrl := gomock.NewController(t1)
			tArgs := tt.args(ctrl)
			transactor := tt.transactor(tArgs)
			newTx, err := transactor.CreateTx(tArgs.txReq)
			if (err != nil) != tt.wantErr {
				t1.Errorf("CreateTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				require.Error(t1, err)
				return
			}
			want := tt.want(tArgs)
			require.Equal(t1, want.Outputs[0].Amount, newTx.Outputs[0].Amount)
			require.Equal(t1, want.Outputs[0].PubKey, newTx.Outputs[0].PubKey)
		})
	}
}
