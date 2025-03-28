package core

type Comics struct {
	ID     string
	URL    string
	ImgUrl string
	Words  []string
}

// Структура, хранящая в себе полностью готовый комикс и счетчик слов,
// который подходит к поисковому запросу
type ComicMatch struct {
	Comic Comics
	Count int
}
