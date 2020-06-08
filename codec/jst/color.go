package jst

type Color struct {
	Enabled      bool
	KeyHighlight []byte
	PlainValue   []byte
}

func (c *Color) initDefaults() {
	if !c.Enabled {
		return
	}
	if c.KeyHighlight == nil {
		c.KeyHighlight = []byte("\033[32m")
	}
	if c.PlainValue == nil {
		c.PlainValue = []byte("\033[1;34m")
	}
}
