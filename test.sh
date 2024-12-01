#!/bin/bash
casetype=${1:-"1"}
caseduration=${2:-"600"}

basedir=$(pwd)
casedir="${basedir}/case"
export BASEDIR="$basedir/"

updategenesis() {
	docker run -it --rm -v "${basedir}/config:/root/config" --name generate --entrypoint /usr/bin/prysmctl tscel/ethnettools:0627 \
		testnet \
		generate-genesis \
		--fork=deneb \
		--num-validators=256 \
		--genesis-time-delay=15 \
		--output-ssz=/root/config/genesis.ssz \
		--chain-config-file=/root/config/config.yml \
		--geth-genesis-json-in=/root/config/genesis.json \
		--geth-genesis-json-out=/root/config/genesis.json
}

testcase() {
	subdir=$1
	resultdir="${basedir}/results/${subdir}"

	# if resultdir exist, delete it.
	if [ -d $resultdir ]; then
		mv $resultdir $subdir_bak
	fi
	mkdir -p $resultdir

	echo "Running testcase $subdir"
	updategenesis
	file=$targetdir/docker-compose.yml
	TEST_STRATEGY=$1 docker compose -f docker-compose.yml up -d 
	echo "wait $epochsToWait epochs" && sleep $caseduration
	docker compose -f docker-compose.yml down
	sudo mv data $resultdir/data

	echo "test done and result in $resultdir"
}

fullstrategy=( "all" "confuse"  "exante"  "random"  "sandwich"  "staircase"  "unrealized"  "withholding" "ext_exante"  "ext_sandwich"  "ext_staircase"  "ext_unrealized"  "ext_withholding")

found=false
for item in "${fullstrategy[@]}"; do
	if [[ "$item" == "$casetype" ]]; then
		found=true
		break
	fi
done

if $found; then
	testcase $casetype
else
	echo "$casetype strategy is not support"
fi
