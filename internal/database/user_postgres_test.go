package database

import (
	"database/sql"
	"database/sql/driver"
	errors "errors"
	"fmt"
	appErrors "go-video-hosting/internal/errors"
	"go-video-hosting/internal/model"
	"os"
	"testing"
	"time"

	"cnb.cool/ordermap/ordermap"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var (
	sqlxDB       *sqlx.DB
	userPostgres *UserPostgres
	mock         sqlmock.Sqlmock
)

func TestMain(m *testing.M) {
	db, mk, err := sqlmock.New()
	if err != nil {
		fmt.Printf("an error %s occurred when opening a stub database connection", err)
		return
	}
	defer db.Close()

	mock = mk
	sqlxDB = sqlx.NewDb(db, "postgres")
	userPostgres = NewUserPostgres(sqlxDB)

	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestCreateUser(t *testing.T) {
	user := model.Users{
		NickName:     "testuser",
		Email:        "test@example.com",
		Password:     "password",
		ActivateLink: "link",
	}

	mockFuncQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery("INSERT INTO USER").
			WithArgs(user.NickName, user.Email, user.Password, user.ActivateLink)
	}

	for _, test := range []struct {
		name     string
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "succesful create user",
			mockFunc: func() {
				mock.ExpectBegin()
				mockFuncQuery().WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "unique violation error",
			mockFunc: func() {
				mock.ExpectBegin()
				mockFuncQuery().WillReturnError(&pq.Error{Code: UniqueViolation})
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "unknown error",
			mockFunc: func() {
				mock.ExpectBegin()
				mockFuncQuery().WillReturnError(errors.New("some unknown error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mockFunc()

			tx, err := sqlxDB.Begin()
			if err != nil {
				t.Fatalf("failed to begin transaction: %s", err.Error())
			}

			id, appErr := userPostgres.CreateUser(tx, user)

			if test.wantErr {
				if appErr == nil {
					t.Error("expected an error, got none")
				}
				if id != 0 {
					t.Errorf("expected id to be 0, got %d", id)
				}
				if err := tx.Rollback(); err != nil {
					t.Error("transaction not rolled back")
				}
			} else {
				if appErr != nil {
					t.Errorf("expected no error, got %v", appErr)
				}
				if id != 1 {
					t.Errorf("expected id to be 1, got %d", id)
				}
				if err := tx.Commit(); err != nil {
					t.Error("transaction not commited")
				}
			}
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	email := "test@example.com"
	createTimestamp, _ := time.Parse(time.RFC3339, "2023-10-01T10:00:00Z")
	mockFuncQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery("SELECT \\* FROM USERS WHERE email=\\$1").
			WithArgs(email)
	}

	for _, test := range []struct {
		name     string
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "successful get user by email",
			mockFunc: func() {
				mockFuncQuery().
					WillReturnRows(sqlmock.NewRows([]string{"id", "nickname", "email", "password", "firstname", "lastname", "birthdate", "avatar", "role", "activatelink", "isactivate", "isbanned", "channelscount", "createtimestamp"}).
						AddRow(1, "testuser", email, "password", "First", "Last", "1965-01-24", "avatar.png", "user", "link", true, false, 3, createTimestamp))

			},
			wantErr: false,
		},
		{
			name: "error executing query",
			mockFunc: func() {
				mockFuncQuery().WillReturnError(errors.New("some database error"))
			},
			wantErr: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mockFunc()

			result, err := userPostgres.GetUserByEmail(email)

			if test.wantErr {
				if err == nil {
					t.Error("expected an error, got none")
				}
				if result != nil {
					t.Errorf("expected result to be nil, got %+v", result)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if result == nil {
					t.Fatalf("expected user to be returned, got nil")
				}
			}
		})
	}
}

func TestGetUserToRefreshById(t *testing.T) {
	id := 1
	mockFuncQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery("SELECT id, nickName, email, role FROM USERS WHERE id=\\$1").
			WithArgs(id)
	}

	for _, test := range []struct {
		name     string
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "successful get user for refreshToken by id",
			mockFunc: func() {
				mockFuncQuery().
					WillReturnRows(sqlmock.NewRows([]string{"id", "nickName", "email", "role"}).
						AddRow(id, "testuser", "test@example.com", "user"))

			},
			wantErr: false,
		},
		{
			name: "error executing query",
			mockFunc: func() {
				mockFuncQuery().WillReturnError(errors.New("some database error"))
			},
			wantErr: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mockFunc()

			result, err := userPostgres.GetUserForRefreshById(id)

			if test.wantErr {
				if err == nil {
					t.Error("expected an error, got none")
				}
				if result != nil {
					t.Errorf("expected result to be nil, got %+v", result)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if result == nil {
					t.Fatalf("expected user to be returned, got nil")
				}
			}
		})
	}
}

func TestGetAvatarByUserId(t *testing.T) {
	userId := 1
	mockFuncQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery("SELECT avatar FROM USERS WHERE id=\\$1").
			WithArgs(userId)
	}

	for _, test := range []struct {
		name     string
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "successful get avatarFileName by userId",
			mockFunc: func() {
				mockFuncQuery().
					WillReturnRows(sqlmock.NewRows([]string{"avatar"}).
						AddRow(userId))

			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("user with id = %d not found", userId),
			mockFunc: func() {
				mockFuncQuery().WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
		{
			name: "error executing query",
			mockFunc: func() {
				mockFuncQuery().WillReturnError(errors.New("some database error"))
			},
			wantErr: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mockFunc()

			result, err := userPostgres.GetAvatarByUserId(userId)

			if test.wantErr {
				if err == nil {
					t.Error("expected an error, got none")
				}
				if result != "" {
					t.Errorf("expected result to be empty string, got %+v", result)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if result == "" {
					t.Fatalf("expected user to be returned, got nil")
				}
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	mockFuncQuery := func(mockQuery string, mockArgs []driver.Value) *sqlmock.ExpectedExec {
		return mock.ExpectExec(mockQuery).WithArgs(mockArgs...)
	}
	for _, test := range []struct {
		name          string
		setData       func() *ordermap.OrderMap
		id            int
		mockFunc      func()
		expectedError *appErrors.AppError
	}{
		{
			name: "successful update",
			setData: func() *ordermap.OrderMap {
				om := ordermap.New()
				om.Store("nickName", "testname")
				om.Store("email", "test@example.com")
				return om
			},
			id: 1,
			mockFunc: func() {
				mockFuncQuery("UPDATE USERS SET nickName = \\$1, email = \\$2 WHERE id = \\$3", []driver.Value{"testname", "test@example.com", 1}).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "user not found",
			setData: func() *ordermap.OrderMap {
				om := ordermap.New()
				om.Store("nickName", "testname")
				return om
			},
			id: 2,
			mockFunc: func() {
				mockFuncQuery("UPDATE USERS SET nickName = \\$1 WHERE id = \\$2", []driver.Value{"testname", 2}).
					WillReturnError(sql.ErrNoRows)
			},
			expectedError: appErrors.New(appErrors.NotFound, sql.ErrNoRows.Error()),
		},
		{
			name: "unique violation error",
			setData: func() *ordermap.OrderMap {
				om := ordermap.New()
				om.Store("email", "duplicate-email@example.com")
				return om
			},
			id: 3,
			mockFunc: func() {
				mockFuncQuery("UPDATE USERS SET email = \\$1 WHERE id = \\$2", []driver.Value{"duplicate-email@example.com", 3}).
					WillReturnError(&pq.Error{Code: UniqueViolation, Message: "unique_violation_error"})
			},
			expectedError: appErrors.New(appErrors.NotUnique, "unique violation: pq: unique_violation_error"),
		},
		{
			name: "other database error",
			setData: func() *ordermap.OrderMap {
				om := ordermap.New()
				om.Store("firstName", "First")
				om.Store("lastName", "Last")
				om.Store("isActivate", true)
				return om
			},
			id: 4,
			mockFunc: func() {
				mockFuncQuery("UPDATE USERS SET firstName = \\$1, lastName = \\$2, isActivate = \\$3 WHERE id = \\$4",
					[]driver.Value{"First", "Last", true, 4}).
					WillReturnError(errors.New("other database error"))

			},
			expectedError: appErrors.New(appErrors.UnknownError, "other database error"),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mockFunc()

			err := userPostgres.UpdateUser(test.id, test.setData())

			if test.expectedError == nil {
				if err != nil {
					t.Errorf("expected no error, got %v", *err)
				}
			} else {
				if err == nil || err.Message != test.expectedError.Message {
					t.Errorf("expected error %v, got %v", *test.expectedError, *err)
				}
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	id := 1
	mockFuncQuery := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec("DELETE FROM USERS WHERE id = \\$1").
			WithArgs(id)
	}

	for _, test := range []struct {
		name     string
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "successful deleting user",
			mockFunc: func() {
				mockFuncQuery().WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "user not exist error",
			mockFunc: func() {
				mockFuncQuery().WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
		{
			name: "other database error",
			mockFunc: func() {
				mockFuncQuery().WillReturnError(errors.New("other database error"))
			},
			wantErr: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mockFunc()

			err := userPostgres.DeleteUser(id)

			if test.wantErr {
				if err == nil {
					t.Error("expected an error, got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", *err)
				}
			}
		})
	}
}

func TestGetUserByActivateLink(t *testing.T) {
	activateLink := "activateLink"
	mockFuncQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery("SELECT id FROM USERS WHERE activateLink = \\$1").
			WithArgs(activateLink)
	}

	for _, test := range []struct {
		name     string
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "successful getting userId by activateLink",
			mockFunc: func() {
				mockFuncQuery().WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantErr: false,
		},
		{
			name: "user not exist error",
			mockFunc: func() {
				mockFuncQuery().WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
		{
			name: "other database error",
			mockFunc: func() {
				mockFuncQuery().WillReturnError(errors.New("other database error"))
			},
			wantErr: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mockFunc()

			result, err := userPostgres.GetUserByActivateLink(activateLink)

			if test.wantErr {
				if err == nil {
					t.Error("expected an error, got none")
				}
				if result != 0 {
					t.Errorf("expected result to be 0, got %d", result)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", *err)
				}
				if result == 0 {
					t.Error("expected result to be not 0, got 0")
				}
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	mockFuncQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery("SELECT id, nickName, email, firstName, lastName, birthDate, role, isBanned, channelsCount, createTimestamp FROM USERS").
			WithoutArgs()
	}

	createTimestamp, _ := time.Parse(time.RFC3339, "2023-10-01T10:00:00Z")
	for _, test := range []struct {
		name     string
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "successful getting all users",
			mockFunc: func() {
				mockFuncQuery().
					WillReturnRows(sqlmock.NewRows([]string{"id", "nickname", "email", "firstname", "lastname", "birthdate", "role", "isbanned", "channelscount", "createtimestamp"}).
						AddRow(1, "testuser", "test@example.com", "First", "Last", "1965-01-24", "user", false, 3, createTimestamp))
			},
			wantErr: false,
		},
		{
			name: "other database error",
			mockFunc: func() {
				mockFuncQuery().WillReturnError(errors.New("other database error"))
			},
			wantErr: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mockFunc()

			result, err := userPostgres.GetAll()

			if test.wantErr {
				if err == nil {
					t.Error("expected an error, got none")
				}
				if result != nil {
					t.Errorf("expected result to be nil, got %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if result == nil {
					t.Error("expected result to be array, got nil")
				}
			}
		})
	}
}

func TestGetById(t *testing.T) {
	id := 1
	mockFuncQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery("SELECT id, nickName, email, firstName, lastName, birthDate, role, isBanned, channelsCount, createTimestamp FROM USERS WHERE id = \\$1").
			WithArgs(id)
	}

	createTimestamp, _ := time.Parse(time.RFC3339, "2023-10-01T10:00:00Z")
	for _, test := range []struct {
		name     string
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "successful getting user by id",
			mockFunc: func() {
				mockFuncQuery().
					WillReturnRows(sqlmock.NewRows([]string{"id", "nickname", "email", "firstname", "lastname", "birthdate", "role", "isbanned", "channelscount", "createtimestamp"}).
						AddRow(1, "testuser", "test@example.com", "First", "Last", "1965-01-24", "user", false, 3, createTimestamp))
			},
			wantErr: false,
		},
		{
			name: "user not exist error",
			mockFunc: func() {
				mockFuncQuery().WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
		{
			name: "other database error",
			mockFunc: func() {
				mockFuncQuery().WillReturnError(errors.New("other database error"))
			},
			wantErr: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mockFunc()

			result, err := userPostgres.GetById(id)

			if test.wantErr {
				if err == nil {
					t.Error("expected an error, got none")
				}
				if result != nil {
					t.Errorf("expected result to be nil, got %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if result == nil {
					t.Error("expected result to be struct FindUser, got nil")
				}
			}
		})
	}
}

func TestGetNickNameById(t *testing.T) {
	id := 1
	mockFuncQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery("SELECT nickName FROM USERS WHERE id = \\$1").
			WithArgs(id)
	}

	for _, test := range []struct {
		name     string
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "successful getting nickName by id",
			mockFunc: func() {
				mockFuncQuery().
					WillReturnRows(sqlmock.NewRows([]string{"nickname"}).AddRow("testuser"))
			},
			wantErr: false,
		},
		{
			name: "user not exist error",
			mockFunc: func() {
				mockFuncQuery().WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
		{
			name: "other database error",
			mockFunc: func() {
				mockFuncQuery().WillReturnError(errors.New("other database error"))
			},
			wantErr: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mockFunc()

			result, err := userPostgres.GetNickNameById(id)

			if test.wantErr {
				if err == nil {
					t.Error("expected an error, got none")
				}
				if result != "" {
					t.Errorf("expected result to be empty string, got %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if result == "" {
					t.Error("expected result to be nickName, got empty string")
				}
			}
		})
	}
}

func TestCheckIsUnique(t *testing.T) {
	key, value := "password", "secret"
	mockFuncQuery := func() *sqlmock.ExpectedQuery {
		return mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM USERS WHERE password = \\$1\\)").
			WithArgs(value)
	}

	for _, test := range []struct {
		name       string
		mockFunc   func()
		waitResult bool
		wantErr    bool
	}{
		{
			name: "value of key is unique",
			mockFunc: func() {
				mockFuncQuery().
					WillReturnRows(sqlmock.NewRows([]string{"exist"}).AddRow(false))
			},
			waitResult: true,
		},
		{
			name: "value of key is not unique",
			mockFunc: func() {
				mockFuncQuery().
					WillReturnRows(sqlmock.NewRows([]string{"exist"}).AddRow(true))
			},
			waitResult: false,
		},
		{
			name: "database error",
			mockFunc: func() {
				mockFuncQuery().WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mockFunc()

			result, err := userPostgres.CheckIsUnique(key, value)

			if test.wantErr {
				if err == nil {
					t.Error("expected an error, got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if result != test.waitResult {
					t.Errorf("expected result to be %v, got %v", test.waitResult, result)
				}
			}
		})
	}
}

func TestChangeChannelsCountOfUser(t *testing.T) {
	userId := 1
	mockFuncQuery := func() *sqlmock.ExpectedExec {
		return mock.ExpectExec("UPDATE USERS SET channelsCount = channelsCount \\+ \\$1 WHERE id = \\$2").
			WithArgs(1, userId)
	}

	for _, test := range []struct {
		name     string
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "succesful change counts of channel",
			mockFunc: func() {
				mock.ExpectBegin()
				mockFuncQuery().WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "user not found error",
			mockFunc: func() {
				mock.ExpectBegin()
				mockFuncQuery().WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "unknown error",
			mockFunc: func() {
				mock.ExpectBegin()
				mockFuncQuery().WillReturnError(errors.New("some unknown error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.mockFunc()

			tx, err := sqlxDB.Begin()
			if err != nil {
				t.Fatalf("failed to begin transaction: %s", err.Error())
			}

			appErr := userPostgres.ChangeChannelsCountOfUser(tx, userId, true)

			if test.wantErr {
				if appErr == nil {
					t.Error("expected an error, got none")
				}
				if err := tx.Rollback(); err != nil {
					t.Error("transaction not rolled back")
				}
			} else {
				if appErr != nil {
					t.Errorf("expected no error, got %v", *appErr)
				}
				if err := tx.Commit(); err != nil {
					t.Error("transaction not commited")
				}
			}
		})
	}
}
