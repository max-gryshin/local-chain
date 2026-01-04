package service_test

import (
	"testing"

	"local-chain/internal/pkg/crypto"

	"local-chain/internal/service"

	"local-chain/internal/types"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestTransactor_CreateTx(t1 *testing.T) {
	type args struct {
		txReq     *types.TransactionRequest
		txStore   service.TransactionStore
		utxoStore service.UTXOStore
		txPool    service.TxPool
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
				fromPubKey := crypto.PublicKeyToBytes(&from.PublicKey)
				to := crypto.GenerateKeyEllipticP256()

				tx1 := types.NewTransaction().WithOutput(types.NewAmount(30), &from.PublicKey)
				tx2 := types.NewTransaction().WithOutput(types.NewAmount(50), &from.PublicKey)
				tx3 := types.NewTransaction().WithOutput(types.NewAmount(20), &from.PublicKey)

				txStore := NewMockTransactionStore(ctrl)
				txStore.EXPECT().Get(tx1.GetHash()).Return(tx1, nil).Times(1)
				txStore.EXPECT().Get(tx2.GetHash()).Return(tx2, nil).Times(1)
				txStore.EXPECT().Get(tx3.GetHash()).Return(tx3, nil).Times(1)

				txPool := NewMockTxPool(ctrl)
				txPool.EXPECT().GetUTXOs(fromPubKey).Return(nil).Times(1)
				txPool.EXPECT().GetPool().Return(nil).Times(3)
				txPool.EXPECT().AddUtxos(gomock.Any(), gomock.Any()).Times(1)
				txPool.EXPECT().AddUtxos(gomock.Any(), gomock.Any()).Times(1)
				txPool.EXPECT().AddTx(gomock.Any()).Return(nil).Times(1)
				utxoStore := NewMockUTXOStore(ctrl)
				utxoStore.EXPECT().Get(fromPubKey).Return([]*types.UTXO{
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
				}, nil).Times(1)
				return args{
					txReq: &types.TransactionRequest{
						Sender:   from,
						Receiver: &to.PublicKey,
						Amount:   *types.NewAmount(100),
					},
					txStore:   txStore,
					txPool:    txPool,
					utxoStore: utxoStore,
				}
			},
			transactor: func(args args) *service.Transactor {
				t := service.NewTransactor(args.txStore, args.utxoStore, args.txPool)

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
				from := crypto.GenerateKeyEllipticP256()
				// fromPubKey := crypto.PublicKeyToBytes(&from.PublicKey)
				fakeFrom := crypto.GenerateKeyEllipticP256()
				// make fakeFrom have the same public key as "from"
				fakeFrom.PublicKey = from.PublicKey
				fakeFromPubKey := crypto.PublicKeyToBytes(&fakeFrom.PublicKey)
				to := crypto.GenerateKeyEllipticP256()

				tx1 := types.NewTransaction().WithOutput(types.NewAmount(30), &from.PublicKey)
				tx2 := types.NewTransaction().WithOutput(types.NewAmount(50), &from.PublicKey)
				tx3 := types.NewTransaction().WithOutput(types.NewAmount(20), &from.PublicKey)

				txStore := NewMockTransactionStore(ctrl)
				txStore.EXPECT().Get(tx1.GetHash()).Return(tx1, nil).Times(1)

				txPool := NewMockTxPool(ctrl)
				txPool.EXPECT().GetUTXOs(fakeFromPubKey).Return(nil).Times(1)
				txPool.EXPECT().GetPool().Return(nil).Times(1)
				utxoStore := NewMockUTXOStore(ctrl)
				utxoStore.EXPECT().Get(fakeFromPubKey).Return([]*types.UTXO{
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
				}, nil).Times(1)

				return args{
					txReq: &types.TransactionRequest{
						Sender:   fakeFrom,
						Receiver: &to.PublicKey,
						Amount:   *types.NewAmount(100),
					},
					txStore:   txStore,
					txPool:    txPool,
					utxoStore: utxoStore,
				}
			},
			transactor: func(args args) *service.Transactor {
				t := service.NewTransactor(args.txStore, args.utxoStore, args.txPool)

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
