package uidplus

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-imap/commands"
	"github.com/emersion/go-imap/responses"
)

// A UIDPLUS client.
type Client struct {
	c *client.Client
}

// NewClient creates a new UIDPLUS client.
func NewClient(c *client.Client) *Client {
	return &Client{c}
}

// UidExpunge permanently removes all messages that both have the \Deleted flag
// set and have a UID that is included in the specified sequence set from the
// currently selected mailbox.
func (c *Client) UidExpunge(seqSet *imap.SeqSet, ch chan uint32) error {
	defer close(ch)

	if c.c.State != imap.SelectedState {
		return client.ErrNoMailboxSelected
	}

	cmd := &commands.Uid{
		Cmd: &ExpungeCommand{SeqSet: seqSet},
	}

	var res *responses.Expunge
	if ch != nil {
		res = &responses.Expunge{SeqNums: ch}
	}

	status, err := c.c.Execute(cmd, res)
	if err != nil {
		return err
	}

	return status.Err()
}
