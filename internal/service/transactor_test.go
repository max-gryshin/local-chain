package service_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"local-chain/internal/service"

	"local-chain/internal/types"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestTransactor_CreateTx(t1 *testing.T) {
	type args struct {
		privKey  *ecdsa.PrivateKey
		toPubKey *ecdsa.PublicKey
		amount   *types.Amount
		utxos    []*types.UTXO
		txPool   service.TxPool
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
				from, err := ecdsa.GenerateKey(elliptic.P256().Params(), rand.Reader)
				require.NoError(t1, err)
				to, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
				require.NoError(t1, err)
				pool := NewMockTxPool(ctrl)
				tx1 := types.NewTransaction().WithOutput(types.NewAmount(30), &from.PublicKey)
				pool.EXPECT().Get(tx1.GetHash()).Return(tx1).Times(1)
				tx2 := types.NewTransaction().WithOutput(types.NewAmount(50), &from.PublicKey)
				pool.EXPECT().Get(tx2.GetHash()).Return(tx2).Times(1)
				tx3 := types.NewTransaction().WithOutput(types.NewAmount(20), &from.PublicKey)
				pool.EXPECT().Get(tx3.GetHash()).Return(tx3).Times(1)

				return args{
					privKey:  from,
					toPubKey: &to.PublicKey,
					amount:   types.NewAmount(100),
					utxos: []*types.UTXO{
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
					txPool: pool,
				}
			},
			transactor: func(args args) *service.Transactor {
				t := service.NewTransactor(args.txPool)

				return t
			},
			want: func(args args) *types.Transaction {
				return types.NewTransaction().WithOutput(types.NewAmount(100), args.toPubKey)
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
				pool := NewMockTxPool(ctrl)
				utxoTx1 := &types.UTXO{
					TxHash: types.NewTransaction().GetHash(),
					Index:  0,
				}
				r, s, err := utxoTx1.Sign(from)
				require.NoError(t1, err)
				tx1 := types.NewTransaction().WithOutput(types.NewAmount(100), &from.PublicKey).WithInput(&types.TxIn{
					Prev:       utxoTx1,
					PubKey:     &from.PublicKey,
					SignatureR: r,
					SignatureS: s,
					NSequence:  0,
				})
				pool.EXPECT().Get(tx1.GetHash()).Return(tx1).Times(1)

				return args{
					privKey:  fakeFrom,
					toPubKey: &to.PublicKey,
					amount:   types.NewAmount(100),
					utxos: []*types.UTXO{
						{
							TxHash: tx1.GetHash(),
							Index:  0,
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
				return types.NewTransaction().WithOutput(types.NewAmount(100), args.toPubKey)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			ctrl := gomock.NewController(t1)
			tArgs := tt.args(ctrl)
			transactor := tt.transactor(tArgs)
			newTx, err := transactor.CreateTx(tArgs.privKey, tArgs.toPubKey, *tArgs.amount, tArgs.utxos)
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
