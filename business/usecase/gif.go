package usecase

type gifGetter interface {
	GetGIF() (string, error)
}

type ProcessGIF struct {
	getter gifGetter
}

func NewProcessGIF(gif gifGetter) ProcessGIF {
	return ProcessGIF{getter: gif}
}

func (p ProcessGIF) GetGif() (string, error) {
	return p.getter.GetGIF()
}
