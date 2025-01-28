DAY=$1
cd $1
go build -o ./cmd/$1
cat input.txt | ./cmd/$1
