# map_put map_name key value
map_put() {
  alias "${1}$2"="$3"
}

# map_get map_name key
# @return value
map_get() {
  alias "${1}$2" | awk -F"'" '{ print $2; }'
}

# map_keys map_name
# @return map keys
map_keys() {
  alias -p | grep $1 | cut -d'=' -f1 | awk -F"$1" '{print $2; }'
}
