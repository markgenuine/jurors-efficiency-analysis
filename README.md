# jurors-efficiency-analysis
[Contest proposal: Jurors efficiency analysis](https://forum.freeton.org/t/contest-proposal-jurors-efficiency-analysis/5295)

## Installation

```sh
 $ go get -u github.com/markgenuine/jurors-efficiency-analysis
 ```
 
 ## Installation dependency
 ```
export CGO_LDFLAGS="-L/path-to-installation/ -lton_client"
```
#### Linux:
```
export LD_LIBRARY_PATH=/path-to-installation/TON-SDK/target/release/deps/
```
#### MacOS:
```
export DYLD_LIBRARY_PATH=/path-to-installation/TON-SDK/target/release/deps/
```

path-to-installation

From source code Ton Labs:
/TON-SDK/target/release/deps/

From binaray:
https://github.com/move-ton/ton-client-go/lib
