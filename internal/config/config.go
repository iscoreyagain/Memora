package config

var Host = "0.0.0.0"
var Port = "6379"
var Protocol = "tcp"

const MAXIMUM_CONNECTION = 20000

var MAX_KEYS int = 10
var EVICTION_RATIO = 0.1

var EVICTION_POLICY string = "allkeys-random"

var POOL_MAXIMUM_SIZE = 16
var POOL_LRU_SIZE = 5 //Sample
