#!/bin/sh -e

echo "Usage: fryabunch.sh <file>"

for i in {200..1000}
do
	figlet "doing $i"

	tmp=`mktemp`
	echo "tempfile: $tmp"
	
	./deepfry.sh $1 "$tmp"
	enc=`base64 -w 0 < "$tmp"`
	#base64 -w 0 < "$tmp"

	echo "writing to file"

	echo "<img style=\"width:50%\" src=\"data:image/jpeg;base64,$enc\">" > "out/$i.html"
done
