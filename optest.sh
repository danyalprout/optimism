
print_block_info() {
  echo "$1";
  echo $2 | jq ".$1.number"
}

# OP Node
syncStatus=$(curl -k "0.0.0.0:7545" -s \
  -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"optimism_syncStatus","params":[],"id":1}' | jq ".result")

print_block_info "current_l1" "$syncStatus"
print_block_info "current_l1_finalized" "$syncStatus"
print_block_info "unsafe_l2" "$syncStatus"
print_block_info "safe_l2" "$syncStatus"

echo "L2 Latest Block"
curl -k "0.0.0.0:9545" -s \
  -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["latest", false],"id":1}' | jq ".result.number"

echo "L1 Latest Block"
curl -k "0.0.0.0:8545" -s \
  -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["latest", false],"id":1}' | jq ".result.number"

