#!/bin/bash
docker build -t attacker-scheduler:latest -f dockerfiles/attack.Dockerfile .
docker build -t generator:latest -f dockerfiles/strategy.Dockerfile .
docker build -t geth:latest -f dockerfiles/geth.Dockerfile .
docker build -t bs_beacon:latest -f dockerfiles/modified.beacon.Dockerfile .
docker build -t beacon:latest -f dockerfiles/normal.beacon.Dockerfile .
docker build -t bs_validator:latest -f dockerfiles/modified.validator.Dockerfile .
docker build -t validator:latest -f dockerfiles/normal.validator.Dockerfile .
