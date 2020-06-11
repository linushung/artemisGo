1. Install the Golang tools: [***Installer***](https://golang.org/dl/)

2. Check that Go is installed correctly
    - Go distribution in **`/usr/local/go`**
    - **`/usr/local/go/bin`** directory in your PATH

3. Set up a workspace and .bash_profile | .profile for GOPATH
    - *export GOPATH=$HOME/go*
    - *export GOBIN=$GOPATH/bin*
    - *export PATH=$PATH:$GOBIN*

4.  Install the Protocol Compiler & Protobuf Runtime of Golang
    - **Compiler**
      - Download a pre-built binary from release page [**Release**](https://github.com/protocolbuffers/protobuf/releases)
      - Copy protoc to `/usr/local/bin`
      - If want to use well known types, then copy the contents of the 'include' directory to `/usr/local/include`
    - **Runtime**
      ```
      go get -u github.com/golang/protobuf/protoc-gen-go
      ```
5. Check IDE or code editor's go-plugin uses **goimports** for formatting code

6. Install missing dependency packages of artemis
   ```
   make clean
   ```

7. Start project
   ```
   make artemis
   ```
