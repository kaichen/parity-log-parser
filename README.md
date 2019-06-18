# parity log parser

A parity log file parser, and analysize rpc requests.

Requirement turn on the parity log options and set to trace level, as "-l rpc=trace".

run this
```shell
make build
./parity-log-parser -logfile=/path/to/parity.log
```

sample output
```
start process log file: /Users/kaichen/parity-rpc-debug/parity.0.log
from:	 2019-06-17 04:08:10 +0800 CST
to:	 2019-06-17 04:21:46 +0800 CST
requests:	 93243
qps:	 114.26838235294117
+---------------------------+--------+
|          METHOD           | COUNT  |
+---------------------------+--------+
| eth_call                  | 116118 |
| eth_blockNumber           |  22699 |
| eth_getBalance            |  15057 |
| eth_getTransactionReceipt |  14274 |
| eth_gasPrice              |   4995 |
| eth_getTransactionCount   |   2581 |
| eth_estimateGas           |   1233 |
| eth_getBlockByNumber      |    896 |
| trace_transaction         |    100 |
| eth_getTransactionByHash  |     90 |
| eth_getCode               |     28 |
+---------------------------+--------+
```
