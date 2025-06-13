package internal

//
//import (
//	"crypto/ecdsa"
//	"crypto/elliptic"
//	"crypto/rand"
//	"local-chain/internal/types"
//	"log"
//	"testing"
//
//	"github.com/google/uuid"
//)
//
//func TestTransaction_VerifySignature(t *testing.T) {
//	type fields struct {
//		SenderID   uuid.UUID
//		ReceiverID uuid.UUID
//		Amount     types.Amount
//		Key        *ecdsa.PrivateKey
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		want    bool
//		wantErr bool
//	}{
//		{
//			name: "",
//			fields: fields{
//				SenderID:   uuid.New(),
//				ReceiverID: uuid.New(),
//				Amount:     types.Amount{Value: 10, Unit: types.CurrencyUnit},
//				Key: func() *ecdsa.PrivateKey {
//					key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
//					if err != nil {
//						log.Fatal(err)
//					}
//					return key
//				}(),
//			},
//			want:    true,
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			tx := types.NewTransaction(tt.fields.SenderID, tt.fields.ReceiverID, tt.fields.Amount)
//			err := tx.Sign(tt.fields.Key)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("VerifySignature() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			got := tx.VerifySignature(&tt.fields.Key.PublicKey)
//			if got != tt.want {
//				t.Errorf("VerifySignature() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
