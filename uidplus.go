// Implements the IMAP UIDPLUS extension, as defined in RFC 4315.
package uidplus

const Capability = "UIDPLUS"

// Additional response codes, defined in RFC 4315 section 3.
const (
	CodeAppendUid = "APPENDUID"
	CodeCopyUid = "COPYUID"
	CodeUidNotSticky = "UIDNOTSTICKY"
)
