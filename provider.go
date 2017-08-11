package user

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// Provider описывает информацию о провайдере, который умеет запрашивать,
// разбирать и возвращать информацию о пользователе.
type Provider struct {
	URL    string // адрес для запроса информации
	Prefix string // префикс к идентификатору пользователя
	// функция для чтения, разбора и конвертации информации о пользователе
	reader func(r io.Reader) (*User, error)
}

// Get возвращает информацию о пользователе.
func (p *Provider) Get(accessToken string) (*User, error) {
	req, err := http.NewRequest("GET", p.URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// разбираем ответ с описанием ошибки
	if resp.StatusCode >= 400 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		// описание ошибки в формате JSON, как правило, находится внутри error
		var jerr = new(struct {
			*Error `json:"error"`
		})
		if err = json.Unmarshal(body, jerr); err != nil {
			// если ответ не формате JSON или мы его не смогли разобрать,
			// то подставляем сам возвращенный текст как сообщение
			jerr.Code = resp.StatusCode
			jerr.Message = string(body)
		}
		return nil, jerr.Error
	}
	// разбираем ответ с информацией о пользователе
	user, err := p.reader(resp.Body)
	// добавляем префикс к идентификатору пользователя
	if err == nil && p.Prefix != "" {
		user.ID = p.Prefix + user.ID
	}
	return user, err
}

// Error описывает ошибку, возвращаемую сервером провайдера при запросе
// информации о пользователе.
type Error struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
}

// Error возвращает строку с описанием ошибки.
func (e *Error) Error() string {
	return e.Message
}
