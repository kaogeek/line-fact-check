package repo

type OptionMessageGroup func(*OptionsMessageGroup)

type OptionsMessageGroup struct {
	Options
	LikeMessageText string
	IDIn            []string
	IDNotIn         []string
}

func MessageGroupLikeMessageText(text string, idIn []string, idNotIn []string) OptionMessageGroup {
	return func(opts *OptionsMessageGroup) {
		opts.LikeMessageText = substringAuto(text)
		opts.IDIn = idIn
		opts.IDNotIn = idNotIn
	}
}
