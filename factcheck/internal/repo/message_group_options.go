package repo

import "github.com/kaogeek/line-fact-check/factcheck"

type OptionMessageGroup func(*OptionsMessageGroup)

type OptionsMessageGroup struct {
	Options
	LikeMessageText string
	IDIn            []string
	IDNotIn         []string
	Statuses        []factcheck.StatusMGroup
}

func MessageGroupLikeMessageText(text string) OptionMessageGroup {
	return func(opts *OptionsMessageGroup) {
		opts.LikeMessageText = substringAuto(text)
	}
}

func MessageGroupIDIn(idIn []string) OptionMessageGroup {
	return func(opts *OptionsMessageGroup) {
		opts.IDIn = idIn
	}
}

func MessageGroupIDNotIn(idNotIn []string) OptionMessageGroup {
	return func(opts *OptionsMessageGroup) {
		opts.IDNotIn = idNotIn
	}
}

func MessageGroupStatusesIn(statuses []factcheck.StatusMGroup) OptionMessageGroup {
	return func(opts *OptionsMessageGroup) {
		opts.Statuses = statuses
	}
}
