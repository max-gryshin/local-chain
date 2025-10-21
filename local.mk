NODE_ID=1
RAFT_ADDR=127.0.0.1:8001
GRPC_ADDR=127.0.0.1:9001
DATA_DIR=./db
BOOTSTRAP=true

.PHONY: add-peer
add-peer:
	grpcurl -v -plaintext -max-time 30 -d '{"id":"00000000-0000-0000-0000-000000000002","address":"172.25.0.12:8001"}' 127.0.0.1:9001 LocalChain.AddPeer

.PHONY: add-voter
add-voter:
	grpcurl -v -plaintext -max-time 30 -d '{"id":"00000000-0000-0000-0000-000000000002","address":"172.25.0.12:8001"}' 127.0.0.1:9001 LocalChain.AddVoter

.PHONY: add-peer
remove-peer:
	grpcurl -v -plaintext -max-time 30 -d '{"id":"00000000-0000-0000-0000-000000000002","address":"172.25.0.12:8001"}' 127.0.0.1:9001 LocalChain.RemovePeer

.PHONY: docker-up
docker-up:
	docker-compose down -v
	docker rmi local-chain -f
	docker-compose up -V --build

.PHONY: oracle-connect
oracle-connect:
	ssh -i ~/.ssh/ssh-key-2025-10-01.key ubuntu@158.180.32.108
