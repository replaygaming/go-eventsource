#!/bin/ash

echo "Waiting for services"

sleep 21

echo "Starting test"

echo "Spawning a couple of curl clients with different settings"
./curl_client.sh 100 &
./curl_client.sh 200 &
./curl_client.sh 300 &

echo "Sending some messages to the rabbitmq server"
# Note that curl is a bit finicky with how it flushes data to the result files.
# If the sent test data isn't just over a flush limit some data will be missing
# from the result files even if it has been sent by eventsource and recived by
# curl. So take care if/when modifying this payload.
./rabbitmq_client.sh "Lorem ipsum dolor sit amet, consectetur adipiscing elit."
./rabbitmq_client.sh "Duis consectetur eros sit amet justo condimentum, quis pellentesque enim rutrum."
./rabbitmq_client.sh "Donec tempor blandit orci, vel facilisis libero."
./rabbitmq_client.sh "Integer porta nulla quis fermentum semper."
./rabbitmq_client.sh "Curabitur tempus feugiat fermentum."
./rabbitmq_client.sh "Mauris leo urna, maximus sed diam eu, tristique tincidunt dolor."
./rabbitmq_client.sh "Donec elementum purus in est elementum, vitae dignissim ligula dictum."
./rabbitmq_client.sh "Donec nec mauris vitae lectus feugiat vestibulum eu ultrices lacus."
./rabbitmq_client.sh "Mauris a velit vitae nunc ullamcorper fringilla."

expected_count=9

echo "Verifying the output"
for file in result*; do
  echo "  Verifying $file:"
  echo "*********************"
  cat $file
  echo "*********************"
  count=$(grep -c "data" $file)
  if [[ $count -eq $expected_count ]]; then
    echo "    Found $expected_count recived messages."
  else
    echo "    Missing messages! Only $count of $expected_count found."
    exit 1
  fi
  echo "  File good."
done
echo "All good."
