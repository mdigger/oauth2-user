package user

import (
	"encoding/json"
	"io"
	"net/http"
)

// Provider описывает информацию о провайдере, который умеет запрашивать,
// разбирать и возвращать информацию о пользователе.
type Provider struct {
	URL    string // адрес для запроса информации
	Prefix string // префикс к идентификатору пользователя
	// функция для чтения, разбора и конвертации информации о пользователе
	// если не определена, то используется стандартный разбор из формата JSON
	Decoder func(dec *json.Decoder) (*User, error)
}

// Get возвращает информацию о пользователе.
func (p Provider) Get(accessToken string) (user *User, err error) {
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
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		// описание ошибки в формате JSON, как правило, находится внутри error
		var jerr struct {
			*Error `json:"error"`
		}
		if err = json.Unmarshal(body, &jerr); err != nil {
			// если ответ не формате JSON или мы его не смогли разобрать,
			// то подставляем сам возвращенный текст как сообщение
			jerr.Code = resp.StatusCode
			jerr.Message = string(body)
		}

		return nil, jerr.Error
	}

	// разбираем ответ с информацией о пользователе
	dec := json.NewDecoder(resp.Body)
	if p.Decoder != nil {
		user, err = p.Decoder(dec)
	} else {
		user = new(User)
		err = dec.Decode(&user)
	}
	if err != nil {
		return nil, err
	}

	if p.Prefix != "" {
		// добавляем префикс к идентификатору пользователя
		user.ID = p.Prefix + user.ID
	}

	return user, nil
}
