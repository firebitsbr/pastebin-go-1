package pastebin

// privacy types
const (
	Public   = 0
	Unlisted = 1
	Private  = 2
)

type Credentials struct {
	DevKey  string
	UserKey string
}

type Paste struct {
	Code       string
	Name       string
	Format     string
	Expiration string
	Privacy    int
}

type Response struct {
	ResCode int `json:""`
}

func NewPastebin(devKey, userKey string) Credentials {
	return Credentials{
		DevKey:  devKey,
		UserKey: userKey,
	}
}

func (c Credentials) NewPaste(code, name, format, expiration string, privacy int) Paste {
	return Paste{
		Code:       code,
		Name:       name,
		Format:     format,
		Expiration: experiation,
		Privacy, privacy,
	}
}

func (c Credentials) SendPaste() (Response, error) {
}
