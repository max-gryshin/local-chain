#!/bin/bash

geth account new --datadir ./var/data && \
geth -datadir ./var/data/ init /var/genesis.json && \
geth --http --http.corsdomain "http://localhost:8000" #--networkid 1999