ps=$(docker ps | grep $1)
lines=$(echo "$ps" | wc -l)

if [ $lines = 1 ]; then
  id=$(echo "$ps" | awk '{print$1}')
  image=$(echo "$ps" | awk '{print$2}')
  name=$(echo "$ps" | awk '{print$NF}')
  echo found "$id" "$image" "$name"
  docker exec -ti "$id" /bin/sh
else
  echo "multiple instances found:"
  line_number=1
  echo "$ps" | awk '{print $1, $2, $NF}' | while read -r id image name; do
    echo "$line_number:\t$id \t$image \t$name"
    ((line_number++))
  done

  read -p "Choose one instance:" line_number

  selected=$(echo "$ps" | sed -n "${line_number}p")
  id=$(echo "$selected" | awk '{print$1}')
  image=$(echo "$selected" | awk '{print$2}')
  name=$(echo "$selected" | awk '{print$NF}')
  echo selected "$id" "$image" "$name"
  docker exec -ti "$id" /bin/sh
fi
