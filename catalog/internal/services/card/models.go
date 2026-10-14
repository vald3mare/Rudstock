package card

// Пагинация: 0 означает значение по умолчанию, limit больше максимума срезается до него
const (
	defaultCardsLimit = 20
	maxCardsLimit     = 100
)

// идем сверху вниз
// данные(CardInput) -> CardService.Create(input) -> cardService.Create(input) -> CardRepo.Create() -> postgres.CardRepo.Create()

// CreateCardInput входные данные для создания карточки.
// Контракт: category_id, title, description, price, photo_url.
type CreateCardInput struct {
	CategoryID  int64  // обязательно, категория должна существовать
	Title       string // обязательно, не пустое (пробелы по краям обрезаются)
	Description string // необязательно
	Price       int64  // обязательно, > 0, копейки
	PhotoURL    string // необязательно, ссылка на объект в minio
}

// ListCardsInput параметры выборки списка карточек.
type ListCardsInput struct {
	CategoryID int64 // необязательно, 0 означает все категории
	Page       int   // необязательно, с 1; 0 означает первую страницу
	Limit      int   // необязательно, 0 означает defaultCardsLimit, максимум maxCardsLimit
}

// UpdateCardInput входные данные для частичного обновления карточки.
// nil означает "не менять поле"; хотя бы одно поле должно быть задано.
type UpdateCardInput struct {
	CategoryID  *int64  // если задано: > 0, категория должна существовать
	Title       *string // если задано: не пустое (пробелы по краям обрезаются)
	Description *string // если задано: любое, в том числе пустое
	Price       *int64  // если задано: > 0, копейки
	PhotoURL    *string // если задано: любое, в том числе пустое
}
