package port

type AvatarProvider interface {
	Next() (name, avatarURL string, err error)
}
