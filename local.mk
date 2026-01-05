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


# List of names to generate keys for (100 total)
NAMES := \
	alice bob charlie dave eve frank grace hank irene jack \
	karen leo mia nick olivia paul quinn rachel steve tina \
	uma victor wendy xavier yvonne zack abby brad chris \
	diana elena felix gina harry isabel jim kevin lara \
	mike nora oscar peter queen ron sara tom ursula \
	vince will xena yasmin zane aaron beth cody dana \
	edgar fiona gary holly ivan jill kyle liam maggie \
	neil opal priya quincy rose sean troy una vera \
	walter xia yang zoe albert bella carl denise \
	eric faith gabriel heidi ian jen ken luis mandy \
	nate owen paula qadir rita sam tyler ugo val \
	wayne xiao yara ziad

.PHONY: add-users
add-users:
	@echo Generating users for $$(echo $(NAMES) | wc -w) names...
	@for name in $(NAMES); do \
  		bin/debug add-user --name $$name; \
	done

.PHONY: send-money
send-money:
	@echo Give money to users for $$(echo $(NAMES) | wc -w) names...
	@for name in $(NAMES); do \
  		bin/debug send -s admin -r $$name -a 1000; \
	done

.PHONY: add-peers
add-peers:
	@echo Adding peers...
	./bin/debug add-peer --address 172.25.0.12:8001 --id 00000000-0000-0000-0000-000000000002 && \
	./bin/debug add-peer --address 172.25.0.13:8001 --id 00000000-0000-0000-0000-000000000003

.PHONY: add-voters
add-voters:
	@echo Adding voters...
	./bin/debug add-voter --address 172.25.0.12:8001 --id 00000000-0000-0000-0000-000000000002 && \
	./bin/debug add-voter --address 172.25.0.13:8001 --id 00000000-0000-0000-0000-000000000003

.PHONY: balance-check
balance-check:
	@echo Checking balances for users...
	@for name in $(NAMES); do \
  		bin/debug --server 127.0.0.1:9002 balance --name $$name; \
	done

.PHONY: full-emission
full-emission:
	@echo "Full emission complete."
	./bin/debug full-emission

.PHONY: load-test-data
load-test-data: add-users send-money add-peers add-voters

.PHONY: load-test-data
debug: send-money balance-check full-emission
	@echo "Test data loading complete."