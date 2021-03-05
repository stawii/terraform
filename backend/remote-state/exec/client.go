package exec

import (
	"os/exec"
	"runtime"
	"fmt"
  "os"
	"crypto/md5"

	"github.com/hashicorp/terraform/states/remote"
	"github.com/hashicorp/terraform/states/statemgr"
)

type execClient struct {
  loadCommand string
  saveCommand string
  lockCommand string
  unlockCommand string
}

func (c *execClient) execute(command string, data *[]byte) ([]byte, error) {
	var cmdargs []string
	if runtime.GOOS == "windows" {
		cmdargs = []string{"cmd", "/C"}
	} else {
		cmdargs = []string{"/bin/sh", "-c"}
	}
	cmdargs = append(cmdargs, command)
  fmt.Printf("cmdargs: %v\n", cmdargs)
	cmd := exec.Command(cmdargs[0], cmdargs[1:]...)
  cmd.Stderr = os.Stderr
  if data != nil {
    input, err := cmd.StdinPipe()
  	if err != nil {
  		return nil, err
  	}
    input.Write(*data)
    input.Close()
  }
	return cmd.Output()
}

func (c *execClient) Get() (*remote.Payload, error) {
	resp, err := c.execute(c.loadCommand, nil)
  if err != nil {
    return nil, err
  }
	if len(resp) == 0 {
		return nil, nil
	}
	payload := &remote.Payload{
		Data: resp,
	}
  hash := md5.Sum(payload.Data)
  payload.MD5 = hash[:]
	return payload, nil
}

func (c *execClient) Put(data []byte) error {
	_, err := c.execute(c.saveCommand, &data)
  if err != nil {
    return err
  }
  return nil
}

func (c *execClient) Delete() error {
  panic("Delete?!")
  return nil
}

func (c *execClient) Lock(info *statemgr.LockInfo) (string, error) {
  if c.lockCommand == "" {
    return "", nil
  }
  json := info.Marshal()
	_, err := c.execute(c.lockCommand, &json)
  if err != nil {
    return "", err
  }
  return info.ID, nil
}

func (c *execClient) Unlock(id string) error {
  if c.unlockCommand == "" {
    return nil
  }
  bid := []byte(id)
	_, err := c.execute(c.unlockCommand, &bid)
  if err != nil {
    return err
  }
  return nil
}
