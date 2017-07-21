package client

import (
	"flag"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/ybbus/jsonrpc"
)

// Credential holds a credential to connect to FairCoin2 RPC service
type Credential struct {
	user     string
	password string
}

// Client is a client of FairCoin RPC service
type Client struct {
	c *jsonrpc.RPCClient
}

var (
	cookiePath = flag.String("cookie.path", "~/.faircoin2/.cookie", "Path to the cookie for connecting to rpc server")
)

func expandUser(path string) string {
	expanded := path

	if strings.HasPrefix(path, "~/") {
		usr, _ := user.Current()
		dir := usr.HomeDir
		expanded = filepath.Join(dir, path[2:])
	}

	return expanded
}

// CookieCredential reads credential from default cookie
func CookieCredential() Credential {
	path := expandUser(*cookiePath)

	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("error opening %s: %s", path, err)
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if err != nil {
		log.Fatalf("error reading: %s", err)
	}

	content := string(buf[:n])
	lines := strings.Split(content, "\n")
	fields := strings.Split(lines[0], ":")

	cred := Credential{
		user:     fields[0],
		password: fields[1],
	}

	return cred
}

// NewClient returns a FairCoin RPC client
func NewClient(url string, cred Credential) *Client {
	rpcClient := jsonrpc.NewRPCClient(url)
	rpcClient.SetBasicAuth(cred.user, cred.password)

	return &Client{
		c: rpcClient,
	}
}
