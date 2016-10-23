package uidplus

import (
	"time"

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

// Same as Client.Append, but can also return the UID of the appended message
// and the UID validity of the destination mailbox. The server can choose not
// to return these values, in this case uid and validity will be equal to zero.
func (c *Client) Append(mbox string, flags []string, date time.Time, msg imap.Literal) (validity, uid uint32, err error) {
	if c.c.State & imap.AuthenticatedState == 0 {
		err = client.ErrNotLoggedIn
		return
	}

	cmd := &commands.Append{
		Mailbox: mbox,
		Flags:   flags,
		Date:    date,
		Message: msg,
	}

	status, err := c.c.Execute(cmd, nil)
	if err != nil {
		return
	}
	if err = status.Err(); err != nil {
		return
	}

	if status.Code == CodeAppendUid && len(status.Arguments) >= 2 {
		validity, _ = imap.ParseNumber(status.Arguments[0])
		uid, _ = imap.ParseNumber(status.Arguments[1])
	}
	return
}
