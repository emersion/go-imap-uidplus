package uidplus

import (
	"errors"

	"github.com/emersion/go-imap"
)

// An UID EXPUNGE command, as defined in RFC
type ExpungeCommand struct {
	SeqSet *imap.SeqSet
}

func (cmd *ExpungeCommand) Command() *imap.Command {
	return &imap.Command{
		Name: imap.Expunge,
		Arguments: []interface{}{cmd.SeqSet},
	}
}

func (cmd *ExpungeCommand) Parse(fields []interface{}) error {
	if len(fields) < 1 {
		return errors.New("Not enough arguments")
	}

	if seqSet, ok := fields[0].(string); !ok {
		return errors.New("Invalid sequence set")
	} else if seqSet, err := imap.NewSeqSet(seqSet); err != nil {
		return err
	} else {
		cmd.SeqSet = seqSet
	}

	return nil
}