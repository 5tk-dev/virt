package virt

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/google/uuid"
	yaml "gopkg.in/yaml.v3"
)

var (
	SocketPath    = "sock/"  // fique a vontade para mudar
	VmDataPath    = "data/"  // fique a vontade para mudar
	VmStoragePath = "disks/" // fique a vontade para mudar
)

func createGuest(g *Guest) error {
	fPath := path.Join(VmDataPath, fmt.Sprintf("%s.yaml", g.Name))
	_, err := os.Stat(fPath)
	if err == nil {
		return os.ErrExist
	}

	if g.UUID == "" {
		g.UUID = uuid.NewString()
	}

	f, err := os.Create(fPath)
	if err != nil {
		return err
	}
	defer f.Close()

	gData, err := yaml.Marshal(g)
	if err != nil {
		return nil
	}
	f.Write(gData)
	return nil
}

func loadGuest(guestPath string) (*Guest, error) {
	f, err := os.Open(guestPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var data bytes.Buffer
	data.ReadFrom(f)

	g := &Guest{}
	err = yaml.Unmarshal(data.Bytes(), g)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func startGuest(g *Guest) error {
	var stdout, stderr bytes.Buffer

	a := g.ToArgs()
	cmd := exec.Command(a[0], a[1:]...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return errors.New(stderr.String())
	}
	return nil
}

/*
usage:

	yalmPath := path.Join("/custom/path/","guestName.yaml")
	guest, err := qmp.LoadGuestFromPath(yalmPath)

load a guest from custom path
*/
func LoadGuestFromPath(yalmPath string) (*Guest, error) { return loadGuest(yalmPath) }

/*
usage:

	guest,err := qmp.LoadGuest("guestName")

load a guest from default path
*/
func LoadGuest(guestName string) (*Guest, error) {
	return loadGuest(path.Join(VmDataPath, fmt.Sprintf("%s.yaml", guestName)))
}

func CreateGuest(g *Guest) error { return createGuest(g) }
func StartGuest(g *Guest) error  { return startGuest(g) }
