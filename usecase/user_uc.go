package usecase

import (
	"errors"
	"majoo-backend/helper"
	"majoo-backend/model"
	"majoo-backend/pkg/bcrypt"
	"majoo-backend/pkg/logruslogger"
	"majoo-backend/pkg/str"
	"majoo-backend/server/request"
	"majoo-backend/usecase/viewmodel"
	"strings"
	"time"
)

// UserUC ...
type UserUC struct {
	*ContractUC
}

// BuildBody ...
func (uc UserUC) BuildBody(data *model.UserEntity, res *viewmodel.UserVM, isShowPassword bool) {

	res.ID = data.ID
	res.Name = data.Name.String
	res.UserName = data.UserName
	res.Password = str.ShowString(isShowPassword, data.Password)
	if data.ImagePath.String != "" {
		res.ImagePath = uc.EnvConfig["APP_IMAGE_URL"] + uc.EnvConfig["FILE_PATH"] + "/" + data.ImagePath.String
	}
	res.CreatedAt = data.CreatedAt
	res.UpdatedAt = data.UpdatedAt
	res.DeletedAt = data.DeletedAt.String
}

// Login ...
func (uc UserUC) Login(data request.UserLoginRequest) (res viewmodel.JwtVM, err error) {
	ctx := "UserUC.Login"

	if len(data.Password) < 8 {
		logruslogger.Log(logruslogger.WarnLevel, "", ctx, "password_length", uc.ReqID)
		return res, errors.New(helper.InvalidCredentials)
	}

	user, err := uc.FindByUserName(data.UserName, true)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_by_email", uc.ReqID)
		return res, errors.New(helper.InvalidCredentials)
	}

	isMatch := bcrypt.CheckPasswordHash(data.Password, user.Password)
	if !isMatch {
		logruslogger.Log(logruslogger.WarnLevel, "", ctx, "invalid_password", uc.ReqID)
		return res, errors.New(helper.InvalidCredentials)
	}

	// Jwe the payload & Generate jwt token
	payload := map[string]interface{}{
		"id": user.ID,
	}
	jwtUc := JwtUC{ContractUC: uc.ContractUC}
	err = jwtUc.GenerateToken(payload, &res)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "jwt", uc.ReqID)
		return res, errors.New(helper.InternalServer)
	}

	return res, err
}

// FindAll ...
func (uc UserUC) FindAll(search string, page, limit int, by, sort string) (res []viewmodel.UserVM, pagination viewmodel.PaginationVM, err error) {
	ctx := "UserUC.FindAll"

	if !str.Contains(model.UserBy, by) {
		by = model.DefaultUserBy
	}
	if !str.Contains(SortWhitelist, strings.ToLower(sort)) {
		sort = DescSort
	}

	limit = uc.LimitMax(limit)
	limit, offset := uc.PaginationPageOffset(page, limit)

	m := model.NewUserModel(uc.DB)
	data, count, err := m.FindAll(search, offset, limit, by, sort)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, pagination, err
	}
	pagination = PaginationRes(page, count, limit)

	for _, r := range data {
		temp := viewmodel.UserVM{}
		uc.BuildBody(&r, &temp, false)
		res = append(res, temp)
	}

	return res, pagination, err
}

// FindByID ...
func (uc UserUC) FindByID(id string, isShowPassword bool) (res viewmodel.UserVM, err error) {
	ctx := "UserUC.FindByID"

	m := model.NewUserModel(uc.DB)
	data, err := m.FindByID(id)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}
	uc.BuildBody(&data, &res, isShowPassword)

	return res, err
}

// FindByUserName ...
func (uc UserUC) FindByUserName(userName string, isShowPassword bool) (res viewmodel.UserVM, err error) {
	ctx := "UserUC.FindByEmail"

	m := model.NewUserModel(uc.DB)
	data, err := m.FindByUserName(userName)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	uc.BuildBody(&data, &res, isShowPassword)

	return res, err
}

// CheckDetails ...
func (uc UserUC) CheckDetails(data *request.UserRequest, oldData *viewmodel.UserVM) (err error) {
	ctx := "UserUC.CheckDetails"

	user, _ := uc.FindByUserName(data.UserName, false)
	if user.ID != "" && user.ID != oldData.ID {
		logruslogger.Log(logruslogger.WarnLevel, data.UserName, ctx, "duplicate_username", uc.ReqID)
		return errors.New(helper.DuplicateUserName)
	}

	if data.Password == "" && oldData.Password == "" {
		logruslogger.Log(logruslogger.WarnLevel, data.UserName, ctx, "empty_password", uc.ReqID)
		return errors.New(helper.InvalidPassword)
	}

	// Decrypt password input
	if data.Password == "" {
		data.Password = oldData.Password
	}

	// Encrypt password
	data.Password, err = bcrypt.HashPassword(data.Password)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "encrypt_password", uc.ReqID)
		return err
	}

	return err
}

// Create ...
func (uc UserUC) Create(data *request.UserRequest) (res viewmodel.UserVM, err error) {
	ctx := "UserUC.Create"

	err = uc.CheckDetails(data, &res)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "check_details", uc.ReqID)
		return res, err
	}

	now := time.Now().UTC()
	res = viewmodel.UserVM{
		Name:      data.Name,
		UserName:  data.UserName,
		Password:  data.Password,
		CreatedAt: now.Format(time.RFC3339),
		UpdatedAt: now.Format(time.RFC3339),
	}
	m := model.NewUserModel(uc.DB)
	res.ID, err = m.Store(res, now)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	return res, err
}

// Update ...
func (uc UserUC) Update(id string, data *request.UserRequest) (res viewmodel.UserVM, err error) {
	ctx := "UserUC.Update"

	oldData, err := uc.FindByID(id, true)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_user", uc.ReqID)
		return res, err
	}

	err = uc.CheckDetails(data, &oldData)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "check_details", uc.ReqID)
		return res, err
	}

	now := time.Now().UTC()
	res = viewmodel.UserVM{
		Name:      data.Name,
		UserName:  data.UserName,
		UpdatedAt: now.Format(time.RFC3339),
	}
	m := model.NewUserModel(uc.DB)
	res.ID, err = m.Update(id, res, now)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	return res, err
}

// UpdateImage ...
func (uc UserUC) UpdateImage(id string, data *request.UserUploadImageRequest) (res viewmodel.UserUploadImageVM, err error) {
	ctx := "UserUC.Update"

	now := time.Now().UTC()
	res = viewmodel.UserUploadImageVM{
		ID:        id,
		Path:      uc.EnvConfig["APP_IMAGE_URL"] + uc.EnvConfig["FILE_PATH"] + "/" + data.Path,
		CreatedAt: now.Format(time.RFC3339),
	}
	m := model.NewUserModel(uc.DB)
	res.ID, err = m.UpdateImage(id, data.Path, now)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	return res, err
}

// Delete ...
func (uc UserUC) Delete(id string) (res viewmodel.UserVM, err error) {
	ctx := "UserUC.Delete"

	now := time.Now().UTC()
	m := model.NewUserModel(uc.DB)
	res.ID, err = m.Destroy(id, now)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	return res, err
}
