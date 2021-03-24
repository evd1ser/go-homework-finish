package apiserver

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"

	"github.com/evd1ser/go-homework-finish/internal/app/middleware"
	"github.com/evd1ser/go-homework-finish/internal/app/models"
	"github.com/form3tech-oss/jwt-go"
)

type Message struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message,omitempty"`
	Error      string `json:"error,omitempty"`
	IsError    bool   `json:"is_error"`
}

func initHeaders(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
}

func (api *APIServer) PostUserRegister(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Post User Register POST /api/v1/user/register")
	var user models.User
	err := json.NewDecoder(req.Body).Decode(&user)

	if err != nil {
		api.logger.Info("Invalid json recieved from client")
		msg := Message{
			StatusCode: 400,
			Message:    "Provided json is invalid",
			IsError:    true,
		}
		writer.WriteHeader(msg.StatusCode)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	//Пытаемся найти пользователя с таким логином в бд
	_, ok, err := api.store.User().FindByLogin(user.Username)
	if err != nil {
		api.logger.Info("Troubles while accessing database table (users) with id. err:", err)
		msg := Message{
			StatusCode: 500,
			Message:    "We have some troubles to accessing database. Try again",
			IsError:    true,
		}
		writer.WriteHeader(msg.StatusCode)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	//Смотрим, если такой пользователь уже есть - то никакой регистрации мы не делаем!
	if ok {
		api.logger.Info("User with that ID already exists")
		msg := Message{
			StatusCode: 400,
			Message:    "User with that login already exists in database",
			IsError:    true,
		}
		writer.WriteHeader(msg.StatusCode)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	//Теперь пытаемся добавить в бд
	_, err = api.store.User().Create(&user)
	if err != nil {
		api.logger.Info("Troubles while accessing database table (users) with id. err:", err)
		msg := Message{
			StatusCode: 500,
			Message:    "We have some troubles to accessing database. Try again",
			IsError:    true,
		}
		writer.WriteHeader(msg.StatusCode)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	msg := Message{
		StatusCode: 201,
		Message:    "User created. Try to auth",
		IsError:    false,
	}

	writer.WriteHeader(msg.StatusCode)
	json.NewEncoder(writer).Encode(msg)
}

func (api *APIServer) PostToAuth(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Post to Auth POST /api/v1/user/auth")
	var userFromJson models.User
	err := json.NewDecoder(req.Body).Decode(&userFromJson)

	//Обрабатываем случай, если json - вовсе не json или в нем какие-либо пробелмы
	if err != nil {
		api.logger.Info("Invalid json recieved from client")
		msg := Message{
			StatusCode: 400,
			Message:    "Provided json is invalid",
			IsError:    true,
		}
		writer.WriteHeader(msg.StatusCode)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	//Необходимо попытаться обнаружить пользователя с таким login в бд
	userInDB, ok, err := api.store.User().FindByLogin(userFromJson.Username)
	// Проблема доступа к бд
	if err != nil {
		api.logger.Info("Can not make user search in database:", err)
		msg := Message{
			StatusCode: 500,
			Message:    "We have some troubles while accessing database",
			IsError:    true,
		}
		writer.WriteHeader(msg.StatusCode)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	//Если подключение удалось , но пользователя с таким логином нет
	if !ok {
		api.logger.Info("User with that login does not exists")
		msg := Message{
			StatusCode: 400,
			Message:    "User with that login does not exists in database. Try register first",
			IsError:    true,
		}
		writer.WriteHeader(msg.StatusCode)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	//Если пользователь с таким логином ест ьв бд - проверим, что у него пароль совпадает с фактическим

	err = bcrypt.CompareHashAndPassword([]byte(userInDB.Password), []byte(userFromJson.Password))

	if err != nil {
		api.logger.Info("Invalid credetials to auth")
		msg := Message{
			StatusCode: 404,
			Message:    "Your password is invalid",
			IsError:    true,
		}
		writer.WriteHeader(msg.StatusCode)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	//Теперь выбиваем токен как знак успешной аутентифкации
	token := jwt.New(jwt.SigningMethodHS256)             // Тот же метод подписания токена, что и в JwtMiddleware.go
	claims := token.Claims.(jwt.MapClaims)               // Дополнительные действия (в формате мапы) для шифрования
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix() //Время жизни токена
	claims["admin"] = true
	claims["name"] = userInDB.Username
	tokenString, err := token.SignedString(middleware.SecretKey)
	//В случае, если токен выбить не удалось!
	if err != nil {
		api.logger.Info("Can not claim jwt-token")
		msg := Message{
			StatusCode: 500,
			Message:    "We have some troubles. Try again",
			IsError:    true,
		}
		writer.WriteHeader(msg.StatusCode)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	//В случае, если токен успешно выбит - отдаем его клиенту
	msg := Message{
		StatusCode: 201,
		Message:    tokenString,
		IsError:    false,
	}
	writer.WriteHeader(msg.StatusCode)
	json.NewEncoder(writer).Encode(msg)

}

//GetStock...
func (api *APIServer) GetStock(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)

	cars, err := api.store.Auto().GetAll()

	if err != nil {
		api.logger.Info("Can not find cars in database:", err)
		msg := Message{
			StatusCode: 500,
			Message:    "We have some troubles while accessing database",
			IsError:    true,
		}
		writer.WriteHeader(msg.StatusCode)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	writer.WriteHeader(200)
	json.NewEncoder(writer).Encode(cars)
}

func (api *APIServer) PostAutoCreate(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)

	mark := mux.Vars(req)["mark"]

	var autoFromJson models.Auto
	err := json.NewDecoder(req.Body).Decode(&autoFromJson)

	//Обрабатываем случай, если json - вовсе не json или в нем какие-либо пробелмы
	if err != nil {
		api.logger.Info("Invalid json recieved from client")
		msg := Message{
			StatusCode: 400,
			Message:    "Provided json is invalid",
			IsError:    true,
		}
		writer.WriteHeader(msg.StatusCode)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	autoFromJson.Mark = mark

	_, exist, err := api.store.Auto().Create(&autoFromJson)

	if err != nil {
		api.logger.Info("Can not find cars in database:", err)
		msg := Message{
			StatusCode: 500,
			Message:    "We have some troubles while accessing database",
			IsError:    true,
		}
		writer.WriteHeader(msg.StatusCode)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	if !exist {
		msg := Message{
			StatusCode: 400,
			Error:      "Auto with that mark exists",
			IsError:    true,
		}
		writer.WriteHeader(msg.StatusCode)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	//В случае, если токен успешно выбит - отдаем его клиенту
	msg := Message{
		StatusCode: 201,
		Message:    "Auto created",
		IsError:    false,
	}
	writer.WriteHeader(msg.StatusCode)
	json.NewEncoder(writer).Encode(msg)
}

func (api *APIServer) GetAuto(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)

	mark := mux.Vars(req)["mark"]

	auto, found, err := api.store.Auto().GetByMark(mark)

	if err != nil {
		api.logger.Info("Can not find cars in database:", err)
		msg := Message{
			StatusCode: 500,
			Message:    "We have some troubles while accessing database",
			IsError:    true,
		}
		writer.WriteHeader(msg.StatusCode)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	if !found {
		msg := Message{
			StatusCode: 404,
			Error:      "Auto with that mark not found",
			IsError:    true,
		}
		writer.WriteHeader(msg.StatusCode)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	writer.WriteHeader(200)
	json.NewEncoder(writer).Encode(auto)
}

func (api *APIServer) PutAutoUpdate(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)

	mark := mux.Vars(req)["mark"]

	auto, found, err := api.store.Auto().GetByMark(mark)

	if err != nil {
		api.logger.Info("Can not find cars in database:", err)
		msg := Message{
			StatusCode: 500,
			Message:    "We have some troubles while accessing database",
			IsError:    true,
		}
		writer.WriteHeader(msg.StatusCode)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	if !found {
		msg := Message{
			StatusCode: 404,
			Error:      "Auto with that mark not found",
			IsError:    true,
		}
		writer.WriteHeader(msg.StatusCode)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	err = json.NewDecoder(req.Body).Decode(auto)

	//Обрабатываем случай, если json - вовсе не json или в нем какие-либо пробелмы
	if err != nil {
		api.logger.Info("Invalid json recieved from client")
		msg := Message{
			StatusCode: 400,
			Message:    "Provided json is invalid",
			IsError:    true,
		}
		writer.WriteHeader(msg.StatusCode)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	auto, err = api.store.Auto().Update(auto)

	if err != nil {
		api.logger.Info("Can not find cars in database:", err)
		msg := Message{
			StatusCode: 500,
			Message:    "We have some troubles while accessing database",
			IsError:    true,
		}
		writer.WriteHeader(msg.StatusCode)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	writer.WriteHeader(202)
	json.NewEncoder(writer).Encode(auto)
}

func (api *APIServer) DeleteAuto(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)

	mark := mux.Vars(req)["mark"]

	auto, found, err := api.store.Auto().GetByMark(mark)

	if err != nil {
		api.logger.Info("Can not find cars in database:", err)
		msg := Message{
			StatusCode: 500,
			Message:    "We have some troubles while accessing database",
			IsError:    true,
		}
		writer.WriteHeader(500)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	if !found {
		msg := Message{
			StatusCode: 404,
			Error:      "Auto with that mark not found",
			IsError:    true,
		}
		writer.WriteHeader(msg.StatusCode)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	_, err = api.store.Auto().Delete(auto)

	if err != nil {
		api.logger.Info("Can not find cars in database:", err)
		msg := Message{
			StatusCode: 500,
			Message:    "We have some troubles while accessing database",
			IsError:    true,
		}
		writer.WriteHeader(msg.StatusCode)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	msg := Message{
		StatusCode: 202,
		Message:    "Auto deleted",
		IsError:    false,
	}

	writer.WriteHeader(msg.StatusCode)
	json.NewEncoder(writer).Encode(msg)
}
