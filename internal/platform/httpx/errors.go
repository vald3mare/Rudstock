package httpx

// ответы с ошибками для http сервера

type Error struct {
	Slug string `json:"slug"`
	Msg  string `json:"msg"`
}
