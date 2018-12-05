package controllerV1

import (
	"Server/app/constants"
	"Server/app/model"
	"Server/app/model/database"
	"Server/app/response"
	"Server/app/utility"
	"github.com/arangodb/go-driver"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//Get phone number from client and send active code by SMS
func GetPhoneNumber(c *gin.Context) {

	//Validation
	var phone model.Phone
	if err := c.ShouldBind(&phone); err != nil {
		c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: err.Error()}})
		return
	}

	t := time.Now().Unix()
	phone.CreatedAt = t

	activeCode := utility.GenerateRandomCode()

	//Search PHONE NUMBER of c.phonePost in phones collection
	if phoneDoc, found := database.Builder().Collection(constants.Phones).Filter("phone", "==", phone.Phone).Build().FindItem(); !found {

		//Create new phonePost
		phoneId, err := model.Create(&phone)

		if err != nil {
			c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: err.Error()}})
			return
		}

		//Create new active code
		if codeId, err := model.Create(&model.Code{Code: activeCode, CreatedAt: t}); err != nil {
			c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: err.Error()}})
			return
		} else {

			//Insert new phone_to_code record in users_edge collection
			if _, err := model.Create(&model.UsersEdge{From: phoneId, To: codeId, Type: constants.PhoneToCode}); err != nil {
				c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: err.Error()}})
				return
			} else {

				//TODO send active code SMS to user phone.Phone
				c.JSON(http.StatusOK, &response.Data{
					Data: &response.SendActiveCode{
						PhoneKey: phoneId.Key(),
					},
					State: true,
				})
			}
		}

	} else {
		phoneId := driver.NewDocumentID(constants.Phones, phoneDoc["_key"].(string))
		//Search code in users GRAPH by PHONE ID with phone_to_code TYPE
		if code, state := database.Builder().From(phoneId).Graph(constants.UsersGraph).Type(constants.PhoneToCode).Build().FindItemInGraph(); !state {
			c.JSON(http.StatusNotFound, &response.Data{Data: &response.DatabaseError{Message: "رکوردی با این شماره تماس یافت نشد"}})
			return
		} else {
			//Update CODE with new active code
			if err := model.Update(code["_key"].(string), &model.Code{Code: activeCode, CreatedAt: t}); err != nil {
				c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: err.Error()}})
				return
			}

			//TODO send active code SMS to user phone.Phone
			c.JSON(http.StatusOK, &response.Data{
				Data: &response.SendActiveCode{
					PhoneKey: phoneDoc["_key"].(string),
				},
				State: true,
			})
		}
	}
}

//Check active code from client and validate
func CheckCode(c *gin.Context) {

	phoneKeyQuery := strings.TrimSpace(c.Query("phone"))
	activeCodeQuery := strings.TrimSpace(c.Query("code"))

	//Validation Queries
	if phoneKeyQuery == "" {
		c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: "کلید شماره تماس الزامی است"}})
		return
	} else if activeCodeQuery == "" {
		c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: "کد فعالسازی الزامی است"}})
		return
	}

	code, _ := strconv.Atoi(activeCodeQuery)

	phoneId := driver.NewDocumentID(constants.Phones, phoneKeyQuery)

	//Search PHONE ID in users GRAPH by phone_to_code TYPE
	if codeDoc, found := database.Builder().Graph(constants.UsersGraph).From(phoneId).Type(constants.PhoneToCode).Build().FindItemInGraph(); !found {
		c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ServerError{Message: "کدی برای این شماره تماس یافت نشد"}})
		return
	} else {

		//Check active code
		if int(codeDoc["code"].(float64)) != code {
			c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: "کد فعالسازی صحیح نمیباشد"}})
		} else {

			//Check active code TimeOut
			if !utility.CheckActiveTime(int64(codeDoc["created_at"].(float64)), time.Now().Unix()) {
				c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: "مدت زمان اعتبار کد فعالسازی تمام شده است"}})
				return
			}

			//Search PHONE ID in users GRAPH by phone_to_user TYPE
			builder := database.Builder().Graph(constants.UsersGraph).From(phoneId).Type(constants.PhoneToUser).Build()
			if userDoc, found := builder.FindItemInGraph(); found {

				token := utility.Token(userDoc["_key"].(string), constants.UserRole)

				//This user already registered
				c.JSON(http.StatusOK, &response.Data{
					Data: &response.NewToken{
						Token:     token,
						IsNewUser: false,
					},
					State: true,
				})
			} else {

				token := utility.RegisterToken(phoneId.Key())

				//This is new user
				c.JSON(http.StatusOK, &response.Data{
					Data: &response.NewToken{
						Token:     token,
						IsNewUser: true,
					},
					State: true,
				})
			}
		}
	}
}

//Create new active code and send by SMS
func RenewCode(c *gin.Context) {

	phoneKeyPost := strings.TrimSpace(c.PostForm("phone_key"))

	if phoneKeyPost == "" {
		c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: "کلید شماره تماس الزامی است"}})
		return
	}

	phoneId := driver.NewDocumentID(constants.Phones, phoneKeyPost)

	//Search code  by PHONE ID in users GRAPH with phone_to_code TYPE
	if codeDoc, found := database.Builder().Graph(constants.UsersGraph).From(phoneId).Type(constants.PhoneToCode).Build().FindItemInGraph(); !found {
		c.JSON(http.StatusNotFound, &response.Data{Data: &response.DatabaseError{Message: "اطلاعاتی برای این شماره تماس یافت نشد"}})
		return
	} else {

		activeCode := utility.GenerateRandomCode()

		//Update CODE with new active code
		if err := model.Update(codeDoc["_key"].(string), &model.Code{Code: activeCode, CreatedAt: time.Now().Unix()}); err != nil {
			c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: err.Error()}})
			return
		}

		//TODO Send active code with SMS
		c.JSON(http.StatusOK, &response.Data{
			Data: &response.SendActiveCode{
				PhoneKey: phoneKeyPost,
			},
			State: true,
		})
	}
}

//Create new user
func CreateUser(c *gin.Context) {

	claims := c.MustGet(constants.TokenClaims).(jwt.MapClaims)
	phoneKey := claims["phone"].(string)
	imageFile, _ := c.FormFile("image")

	//Validation phone key
	if phoneKey == "" {
		c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: "کلید شماره تماس الزامی است"}})
		return
	}

	//Search phoneKey exist or no
	if _, err := model.Find(phoneKey, &model.Phone{}); driver.IsNotFound(err) {
		c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: "چنین کلید شماره تماسی یافت نشد"}})
		return
	}

	phoneId := driver.NewDocumentID(constants.Phones, phoneKey)

	//Check user is new or not
	if _, found := database.Builder().Graph(constants.UsersGraph).From(phoneId).Type(constants.PhoneToUser).Build().FindItemInGraph(); found {
		c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: "این کاربر قبلا ثبت نام کرده است"}})
		return
	}

	//Validation
	var user model.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: err.Error()}})
		return
	}

	//Check exist username
	if _, found := database.Builder().Collection(constants.Users).Filter("username", "==", user.Username).Build().FindItem(); found {
		c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: "نام کاربری قبلا ثبت شده است"}})
		return
	}

	if imageFile != nil {
		isUpload, path, errMsg := utility.UploadImage(imageFile, constants.ImageUser)
		if !isUpload {
			c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: errMsg}})
			return
		}

		//Upload image
		err := c.SaveUploadedFile(imageFile, path)
		if err != nil {
			c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ServerError{Message: "مشکل در آپلود تصویر"}})
			return
		}

		user.Image = path
	} else {
		user.Image = constants.DefaultAvatarPath
	}

	t := time.Now().Unix()
	user.Approved = true
	user.CreatedAt = t
	user.UpdatedAt = t

	//Create new user
	userId, err := model.Create(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: err.Error()}})
		return
	}

	//Create new record in users_edge from phoneId to userId by phone_to_user TYPE
	if _, err := model.Create(&model.UsersEdge{From: phoneId, To: userId, Type: constants.PhoneToUser}); err != nil {
		c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: err.Error()}})
		return
	}

	//Create new record in users_edge from userId to phoneId by user_to_phone TYPE
	if _, err := model.Create(&model.UsersEdge{From: userId, To: phoneId, Type: constants.UserToPhone}); err != nil {
		c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: err.Error()}})
		return
	}

	//Create new Token
	token := utility.Token(userId.Key(), constants.UserRole)

	c.JSON(http.StatusOK, &response.Data{
		Data: &response.Token{
			Token: token,
		},
		State: true,
	})
}

//Check username
func CheckUsername(c *gin.Context) {

	usernamePost := c.PostForm("username")

	if usernamePost == "" {
		c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: "نام کاربری نمیتواند خالی باشد"}})
		return
	}

	if _, found := database.Builder().Collection(constants.Users).Filter("username", "==", usernamePost).Build().FindItem(); found {
		c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: "نام کاربری قبلا ثبت شده است"}})
		return
	}

	c.JSON(http.StatusOK, &response.Data{
		Data:  &response.CheckedUsername{Message: "نام کاربری قابل ثبت است"},
		State: true,
	})
}

//Get user account information
func GetUserInfo(c *gin.Context) {
	claims := c.MustGet(constants.TokenClaims).(jwt.MapClaims)
	userKey := claims["key"].(string)

	user, err := model.Find(userKey, &model.User{})

	if driver.IsNotFound(err) || err != nil {
		c.JSON(http.StatusUnauthorized, &response.Data{Data: &response.ServerError{Message: "کاربری با این مشخصات یافت نشد"}})
		return
	}

	c.JSON(http.StatusOK, &response.Data{
		Data: &response.FindDocument{
			Document: map[string]interface{}{
				"key": user["_key"].(string),
				"name": user["name"].(string),
				"family": user["family"].(string),
				"username": user["username"].(string),
				"image": user["image"].(string),
				"created_at": user["created_at"].(float64),
			},
		},
		State: true,
	})
}
