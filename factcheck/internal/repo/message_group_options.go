package repo

type OptionMessageGroup func(*OptionsMessageGroup)

type OptionsMessageGroup struct {
	Options
	LikeMessageText string
}

func MessageGroupLikeMessageText(text string) OptionMessageGroup {
	return func(opts *OptionsMessageGroup) {
		opts.LikeMessageText = substringAuto(text)
	}
}
