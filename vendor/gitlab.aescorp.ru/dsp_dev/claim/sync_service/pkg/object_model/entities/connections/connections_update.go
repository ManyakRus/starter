package connections

import (
	"gitlab.aescorp.ru/dsp_dev/claim/sync_service/pkg/db/constants"
)

// Update_BranchID - изменяет объект в БД по ID, присваивает BranchID
func (m *Connection) Update_BranchID() error {
	if Crud_Connection == nil {
		return constants.ErrorCrudIsNotInit
	}

	err := Crud_Connection.Update_BranchID(m)

	return err
}

// Update_DbName - изменяет объект в БД по ID, присваивает DbName
func (m *Connection) Update_DbName() error {
	if Crud_Connection == nil {
		return constants.ErrorCrudIsNotInit
	}

	err := Crud_Connection.Update_DbName(m)

	return err
}

// Update_DbScheme - изменяет объект в БД по ID, присваивает DbScheme
func (m *Connection) Update_DbScheme() error {
	if Crud_Connection == nil {
		return constants.ErrorCrudIsNotInit
	}

	err := Crud_Connection.Update_DbScheme(m)

	return err
}

// Update_IsLegal - изменяет объект в БД по ID, присваивает IsLegal
func (m *Connection) Update_IsLegal() error {
	if Crud_Connection == nil {
		return constants.ErrorCrudIsNotInit
	}

	err := Crud_Connection.Update_IsLegal(m)

	return err
}

// Update_Login - изменяет объект в БД по ID, присваивает Login
func (m *Connection) Update_Login() error {
	if Crud_Connection == nil {
		return constants.ErrorCrudIsNotInit
	}

	err := Crud_Connection.Update_Login(m)

	return err
}

// Update_Name - изменяет объект в БД по ID, присваивает Name
func (m *Connection) Update_Name() error {
	if Crud_Connection == nil {
		return constants.ErrorCrudIsNotInit
	}

	err := Crud_Connection.Update_Name(m)

	return err
}

// Update_Password - изменяет объект в БД по ID, присваивает Password
func (m *Connection) Update_Password() error {
	if Crud_Connection == nil {
		return constants.ErrorCrudIsNotInit
	}

	err := Crud_Connection.Update_Password(m)

	return err
}

// Update_Port - изменяет объект в БД по ID, присваивает Port
func (m *Connection) Update_Port() error {
	if Crud_Connection == nil {
		return constants.ErrorCrudIsNotInit
	}

	err := Crud_Connection.Update_Port(m)

	return err
}

// Update_Server - изменяет объект в БД по ID, присваивает Server
func (m *Connection) Update_Server() error {
	if Crud_Connection == nil {
		return constants.ErrorCrudIsNotInit
	}

	err := Crud_Connection.Update_Server(m)

	return err
}
