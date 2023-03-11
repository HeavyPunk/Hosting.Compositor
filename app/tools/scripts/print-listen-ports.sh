raw_out=$(lsof -i -P -n | grep LISTEN)
echo $raw_out