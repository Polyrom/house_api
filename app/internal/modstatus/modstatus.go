package modstatus

type ModerationStatus struct {
	s string
}

func (ms ModerationStatus) String() string {
	return ms.s
}

var (
	Created      = ModerationStatus{"created"}
	Approved     = ModerationStatus{"approved"}
	Declined     = ModerationStatus{"declined"}
	OnModeration = ModerationStatus{"on moderation"}
)
