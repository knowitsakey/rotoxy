package core

import (
	"golang.org/x/crypto/ed25519"
	"os/exec"
	"strconv"
	"strings"
)

//i user, hidden service, ssh private key, torsocks5port, doubleproxyport,
type SshProxy struct {
	user            *string
	hiddenservice   *string
	privatekey      ed25519.PrivateKey
	torsocks5port   *string
	doubleproxyport *string
}

//takes in a string of user@onion, appends @8439@8349 sends to ssh
func CreateSshProxy(txtline string, tp *TorProxy1) (*SshProxy, error) {
	var err error
	sp := &SshProxy{}
	sshArg := txtline + "@" + strconv.Itoa(*tp.ProxyPort) + "@" + strconv.Itoa(*tp.DoubleProxyPort)

	inputstr := strings.SplitN(sshArg, "@", 4)
	sp.user = &inputstr[0]
	sp.hiddenservice = &inputstr[1]
	sp.torsocks5port = &inputstr[2]
	sp.doubleproxyport = &inputstr[3]
	mc := exec.Command("./darkssh", sshArg)
	mc.Start()
	return sp, err
}
