package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cafeList = map[string][]string{
	"moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	req := httptest.NewRequest(http.MethodGet, "/cafe?count=10&city=moscow", nil) // здесь нужно создать запрос к сервису

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// здесь нужно добавить необходимые проверки
	require.Equal(t, responseRecorder.Code, http.StatusOK) //проверка, запрос сформирован корректно, сервис возвращает код ответа 200

	body := responseRecorder.Body.String()
	list := strings.Split(body, ",")

	require.NotEmpty(t, body)            // проверка, тело ответа не пустое
	assert.Len(t, len(list), totalCount) // проверка, длина тела соответствует ожидаемой

}

func TestMainHandlerWhenCityStatusBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/cafe?count=10&city=omsk", nil) //город omsk, который передаётся в параметре city, не поддерживается

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	body := responseRecorder.Body.String()
	assert.Equal(t, responseRecorder.Code, http.StatusBadRequest) //проверка, сервис возвращает код ответа 400 и ошибку wrong city value в теле ответа.
	assert.Equal(t, body, "wrong city value")                     //проверка, сервис возвращает ошибку wrong city value в теле ответа.
}

func TestMainHandlerWhenCorrectRequest(t *testing.T) {
	totalCount := 4
	req := httptest.NewRequest(http.MethodGet, "/cafe?count=3&city=moscow", nil) //

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, responseRecorder.Code, http.StatusOK) //проверка, запрос сформирован корректно, сервис возвращает код ответа 200

	body := responseRecorder.Body.String()
	list := strings.Split(body, ",")

	require.NotEmpty(t, body)            // проверка, тело ответа не пустое
	assert.Len(t, len(list), totalCount) // проверка, длина тела соответствует ожидаемой
}
