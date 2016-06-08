package pastebin

// privacy types
const (
	Public   = 0
	Unlisted = 1
	Private  = 2
)

type Paste struct {
	DevKey       string
	PasteCode    string
	PastePrivate int
	PasteName    string
	ExpireDate   string
	UserKey      string
}
