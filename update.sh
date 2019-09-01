#!/bin/bash

for i in `cat follows`
do
	echo "Updating $i..."
	wget https://www.twitter.com/$i
	./twitter-sanity $i
	cat tweets.csv >> $i.csv
	rm tweets.csv
	rm $i
done

