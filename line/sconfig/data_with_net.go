package sconfig

import (
	"io"

	"github.com/scryinfo/dot/dot"
)

var _ dot.ConfigGetter = (*DataWithNet)(nil)

type DataWithNet struct {
	ReadCloser io.ReadCloser
	fileType   string
}

// Close implements [dot.ConfigGetter].
func (p *DataWithNet) Close() error {
	return p.ReadCloser.Close()
}

// FileType implements [dot.ConfigGetter].
func (p *DataWithNet) FileType() string {
	return p.fileType
}

// Read implements [dot.ConfigGetter].
func (d *DataWithNet) Read(p []byte) (n int, err error) {
	return d.ReadCloser.Read(p)
}

type DataWithNetConfig struct {
	Url      string
	ApiToken string
}

func NewDataWithNet(config *DataWithNetConfig) (*DataWithNet, error) {
	//todo
	return nil, nil
}
