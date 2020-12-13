# jurors-efficiency-analysis
[Contest proposal: Jurors efficiency analysis](https://forum.freeton.org/t/contest-proposal-jurors-efficiency-analysis/5295)

You can use scrypty 2 ways:
## 1 Use binary:
1. Go to the link [https://github.com/markgenuine/jurors-efficiency-analysis/tree/main/bin] and downloads archive for you operationing system.
2. Unpacked archive.
3. Open console (terminal).
4. Run command and copy value with way files: 
```sh
pwd
```
5. Run command and change "past_value_from_pwd" value from 4:
```sh
export CGO_LDFLAGS="-L/past_value_from_pwd/ -lton_client"
```
4. Run command and change "past_value_from_pwd" value from:
```sh 
export LD_LIBRARY_PATH=/past_value_from_pwd/
```
5. Use scrypt:
```sh
./ContestsResults 0:0618a45b9fd55533e4108b2ee8a63d07c775550e5362f124342ae94a0d6158ec 0:f276d7d294a19db415359a466d248b43202edf3a81f91c6dcc017eaff9be308c 0:bed32b8670fec398973d97e2bd2f4c8125ed599182662b4383bdb3c1e996f55c 0:59ebb6b7f0bcc13fb0d239017ad1485930d08c8d97c9456675df7087c54e7971 0:824a244a2483873a43abf3d24d0637a0cfeccb6311e40ab0628dd5f96a41df84 0:099a8a476c5b85fe4271438ff9588a3d104d65233da1ef572e4b4d1c2e9a90f4
```

## 2 Use source code.
### Installation
1. Install golang [www.golang.org]
2. Run command:
```sh
go get -u github.com/markgenuine/jurors-efficiency-analysis
```

### Installation dependency
3. Run command:
##### Linux:
```
export CGO_LDFLAGS="-L/$GOPATH/github.com/move-ton/ton-client-go/lib/linux/ -lton_client"
export LD_LIBRARY_PATH=/$GOPATH/github.com/move-ton/ton-client-go/lib/linux/
```
or

##### MacOS:
```
export CGO_LDFLAGS="-L/$GOPATH/github.com/move-ton/ton-client-go/lib/darwin/ -lton_client"
export DYLD_LIBRARY_PATH=/$GOPATH/github.com/move-ton/ton-client-go/lib/darwin/
```
4. Run command
```sh
go build -o ContestsResults
```
5. Use scrypt. 

Example:
```sh
./ContestsResults 0:0618a45b9fd55533e4108b2ee8a63d07c775550e5362f124342ae94a0d6158ec 0:f276d7d294a19db415359a466d248b43202edf3a81f91c6dcc017eaff9be308c 0:bed32b8670fec398973d97e2bd2f4c8125ed599182662b4383bdb3c1e996f55c 0:59ebb6b7f0bcc13fb0d239017ad1485930d08c8d97c9456675df7087c54e7971 0:824a244a2483873a43abf3d24d0637a0cfeccb6311e40ab0628dd5f96a41df84 0:099a8a476c5b85fe4271438ff9588a3d104d65233da1ef572e4b4d1c2e9a90f4
```
