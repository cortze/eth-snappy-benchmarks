#!/bin/bash
declare -i var=113893
declare -i blockid=0
for i in {1..20}
do
	echo $i;
        blockid=$(($var - $i))	
	echo $blockid;
	curl "localhost:5052/beacon/block?slot=$blockid" > ./data/block$i.json 
done
