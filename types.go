package gitcomm

import (
	"bytes"
	"fmt"
)

// Message type holds all commit message fields
type Message struct {
	// Type message field
	Type string
	// Subject message field
	Subject string
	// Body message field
	//Body string
	// tapd id field
	TapdId int
	// TapdType tapd type: story/bug
	TapdType string
}

func (m Message) String() string {
	buf := bytes.NewBuffer(nil)
	fmt.Fprintf(buf, "%s: %s\n\n--%s=%d", m.Type, m.Subject, m.TapdType, m.TapdId)
	return buf.String()
}
