package uidplus

import (
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-imap/commands"
	"github.com/emersion/go-imap/responses"
)

// Client is a UIDPLUS client.
type Client struct {
	c *client.Client
}

// NewClient creates a new UIDPLUS client.
func NewClient(c *client.Client) *Client {
	return &Client{c}
}

// SupportUidPlus checks if the server supports the UIDPLUS extension.
func (c *Client) SupportUidPlus() (bool, error) {
	return c.c.Support(Capability)
}

// UidExpunge permanently removes all messages that both have the \Deleted flag
// set and have a UID that is included in the specified sequence set from the
// currently selected mailbox.
func (c *Client) UidExpunge(seqSet *imap.SeqSet, ch chan uint32) error {
	defer close(ch)

	if c.c.State() != imap.SelectedState {
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

// Append is the same as Client.Append, but can also return the UID of the
// appended message and the UID validity of the destination mailbox. The server
// can choose not to return these values, in this case uid and validity will be
// equal to zero.
func (c *Client) Append(mbox string, flags []string, date time.Time, msg imap.Literal) (validity, uid uint32, err error) {
	if c.c.State()&imap.AuthenticatedState == 0 {
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

func (c *Client) copy(uid bool, seqSet *imap.SeqSet, dest string) (validity uint32, srcUids, dstUids *imap.SeqSet, err error) {
	if c.c.State()&imap.SelectedState == 0 {
		err = client.ErrNoMailboxSelected
		return
	}

	var cmd imap.Commander
	cmd = &commands.Copy{
		SeqSet:  seqSet,
		Mailbox: dest,
	}
	if uid {
		cmd = &commands.Uid{Cmd: cmd}
	}

	status, err := c.c.Execute(cmd, nil)
	if err != nil {
		return
	}
	if err = status.Err(); err != nil {
		return
	}

	if status.Code == CodeCopyUid && len(status.Arguments) >= 3 {
		validity, _ = imap.ParseNumber(status.Arguments[0])
		if seqSet, ok := status.Arguments[1].(string); ok {
			srcUids, _ = imap.ParseSeqSet(seqSet)
		}
		if seqSet, ok := status.Arguments[2].(string); ok {
			dstUids, _ = imap.ParseSeqSet(seqSet)
		}
	}
	return
}

// Copy is the same as Client.Copy, but can also return the source and
// destination UIDs of the copied messages.
func (c *Client) Copy(seqset *imap.SeqSet, dest string) (validity uint32, srcUids, dstUids *imap.SeqSet, err error) {
	return c.copy(false, seqset, dest)
}

// UidCopy is the same as Client.UidCopy, but can also return the source and
// destination UIDs of the copied messages.
func (c *Client) UidCopy(seqset *imap.SeqSet, dest string) (validity uint32, srcUids, dstUids *imap.SeqSet, err error) {
	return c.copy(true, seqset, dest)
}

// UidPlusClient is an alias used to compose multiple client extensions.
type UidPlusClient = Client
