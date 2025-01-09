package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"go-api-newspaper/api"
	"go-api-newspaper/app/models"
	"go-api-newspaper/pkg/tester"
)

// NewspaperControllersSuite はテスト用の構造体。suite を使用してテストを一元管理する。
type NewspaperControllersSuite struct {
	tester.DBSQLiteSuite
	newspaperHandler NewspaperHandler // テスト対象の NewspaperHandler
	originalDB   *gorm.DB             // モック前のデータベースの参照を保持
}

// TestNewspaperControllersTestSuite はテストスイートを実行するエントリーポイント。
func TestNewspaperControllersTestSuite(t *testing.T) {
	suite.Run(t, new(NewspaperControllersSuite))
}

func (suite *NewspaperControllersSuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite()
	suite.newspaperHandler = NewspaperHandler{} // NewspaperHandler の初期化
	suite.originalDB = models.DB
}

func (suite *NewspaperControllersSuite) MockDB() sqlmock.Sqlmock {
	mock, mockGormDB := tester.MockDB()
	models.DB = mockGormDB // models.DB をモックデータベースに置き換える
	return mock
}

func (suite *NewspaperControllersSuite) AfterTest(suiteName, testName string) {
	models.DB = suite.originalDB // テスト終了後、元のデータベースに戻す
}

// TestCreate は CreateNewspaper メソッドの正常系テスト。
func (suite *NewspaperControllersSuite) TestCreate() {
	// リクエストの準備
	request, _ := api.NewCreateNewspaperRequest("/api/v1", api.CreateNewspaperJSONRequestBody{
		Title:       "test",
		ColumnName:  "sports",
	})
	w := httptest.NewRecorder() // レスポンス記録用の HTTP テストレコーダー
	ginContext, _ := gin.CreateTestContext(w) // Gin のテストコンテキスト作成
	ginContext.Request = request // テスト用リクエストを Gin のコンテキストに設定

	// メソッド実行
	suite.newspaperHandler.CreateNewspaper(ginContext)

	suite.Assert().Equal(http.StatusCreated, w.Code)
	bodyBytes, _ := io.ReadAll(w.Body)
	var newspaperGetResponse api.NewspaperResponse
	err := json.Unmarshal(bodyBytes, &newspaperGetResponse)
	suite.Assert().Nil(err) // JSON のパースが成功していることを確認
	suite.Assert().Equal(http.StatusCreated, w.Code)
	suite.Assert().Equal("test", newspaperGetResponse.Title)
	suite.Assert().Equal("sports", newspaperGetResponse.ColumnName)
}

// TestCreateRequestBodyFailure はリクエストボディが不正な場合のテスト。
func (suite *NewspaperControllersSuite) TestCreateRequestBodyFailure() {
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)

	// リクエストボディなしのリクエストを作成
	req, _ := http.NewRequest("POST", "/api/v1/newspaper", nil)
	req.Header.Add("Content-Type", "application/json")
	ginContext.Request = req

	suite.newspaperHandler.CreateNewspaper(ginContext)
	suite.Assert().Equal(http.StatusBadRequest, w.Code)
	suite.Assert().JSONEq(`{"message": "invalid request"}`, w.Body.String())
}

// TestCreateFailure はデータベースエラー時のテスト。
func (suite *NewspaperControllersSuite) TestCreateFailure() {
	mockDB := suite.MockDB()
	// INSERT クエリを実行した際にエラーを返すよう設定
	mockDB.ExpectExec("INSERT INTO `newspapers`").WithArgs("Test", "sports").WillReturnError(errors.New("create error"))

	request, _ := api.NewCreateNewspaperRequest("/api/v1", api.CreateNewspaperJSONRequestBody{
		Title:       "test",
		ColumnName:    "sports",
	})
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.newspaperHandler.CreateNewspaper(ginContext)

	suite.Assert().Equal(http.StatusInternalServerError, w.Code)
	suite.Assert().True(strings.Contains(w.Body.String(), "create error"))
}

func (suite *NewspaperControllersSuite) TestGet() {
	// 新聞データを作成
	createdNewspaper, _ := models.CreateNewspaper("test", "sports")

	// HTTPリクエストを作成
	request, _ := api.NewGetNewspaperByIdRequest("/api/v1", createdNewspaper.ID)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request
	// GetNewspaperById メソッドを呼び出し
	suite.newspaperHandler.GetNewspaperById(ginContext, createdNewspaper.ID)
	bodyBytes, _ := io.ReadAll(w.Body)
	var newspaperGetResponse api.NewspaperResponse
	err := json.Unmarshal(bodyBytes, &newspaperGetResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusOK, w.Code)
	suite.Assert().Equal("test", newspaperGetResponse.Title)
	suite.Assert().Equal("sports", newspaperGetResponse.ColumnName)
}

func (suite *NewspaperControllersSuite) TestGetNoNewspaperFailure() {
	// 存在しないIDを設定
	doesNotExistNewspaperID := 1111
	// 存在しないIDでGetNewspaperを呼び出し、エラーが発生することを確認
	deletedNewspaper, err := models.GetNewspaper(doesNotExistNewspaperID)
	suite.Assert().NotNil(err)
	suite.Assert().Nil(deletedNewspaper)

	request, _ := api.NewGetNewspaperByIdRequest("/api/v1", doesNotExistNewspaperID)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request
	suite.newspaperHandler.GetNewspaperById(ginContext, doesNotExistNewspaperID)
	bodyBytes, _ := io.ReadAll(w.Body)
	var newspaperGetResponse api.NewspaperResponse
	err = json.Unmarshal(bodyBytes, &newspaperGetResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusInternalServerError, w.Code)
}

func (suite *NewspaperControllersSuite) TestUpdate() {
	createdNewspaper, _ := models.CreateNewspaper("test", "sports")

	// 更新データを設定
	title := "updated"
	request, _ := api.NewUpdateNewspaperByIdRequest("/api/v1", createdNewspaper.ID,
		api.UpdateNewspaperByIdJSONRequestBody{
			Title:    &title,
		},
	)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request
	suite.newspaperHandler.UpdateNewspaperById(ginContext, createdNewspaper.ID)
	bodyBytes, _ := io.ReadAll(w.Body)
	var newspaperGetResponse api.NewspaperResponse
	err := json.Unmarshal(bodyBytes, &newspaperGetResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusOK, w.Code)
	suite.Assert().Equal("updated", newspaperGetResponse.Title)
	suite.Assert().Equal("sports", newspaperGetResponse.ColumnName)
}

func (suite *NewspaperControllersSuite) TestUpdateRequestBodyFailure() {
	// 不正なリクエストを作成
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("PATCH", "/api/v1/newspaper", nil)
	req.Header.Add("Content-Type", "application/json")
	ginContext.Request = req

	suite.newspaperHandler.CreateNewspaper(ginContext)
	suite.Assert().Equal(http.StatusBadRequest, w.Code)
	suite.Assert().JSONEq(`{"message": "invalid request"}`, w.Body.String())
}

func (suite *NewspaperControllersSuite) TestUpdateNoNewspaperFailure() {
	doesNotExistNewspaperID := 1111
	deletedNewspaper, err := models.GetNewspaper(doesNotExistNewspaperID)
	suite.Assert().NotNil(err)
	suite.Assert().Nil(deletedNewspaper)

	title := "updated"
	request, _ := api.NewUpdateNewspaperByIdRequest("/api/v1", doesNotExistNewspaperID,
		api.UpdateNewspaperByIdJSONRequestBody{
			Title:    &title,
		},
	)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request
	suite.newspaperHandler.UpdateNewspaperById(ginContext, doesNotExistNewspaperID)
	bodyBytes, _ := io.ReadAll(w.Body)
	var newspaperGetResponse api.NewspaperResponse
	err = json.Unmarshal(bodyBytes, &newspaperGetResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusInternalServerError, w.Code)
}

func (suite *NewspaperControllersSuite) TestUpdateFailure() {
	mockDB := suite.MockDB()

	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `newspapers` WHERE `newspapers`.`id` = ? ORDER BY `newspapers`.`id` LIMIT ?")).WithArgs(1, 1).WillReturnError(errors.New("update error"))

	title := "updated"
	request, _ := api.NewUpdateNewspaperByIdRequest("/api/v1", 1,
		api.UpdateNewspaperByIdJSONRequestBody{
			Title:    &title,
		},
	)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.newspaperHandler.UpdateNewspaperById(ginContext, 1)

	suite.Assert().Equal(http.StatusInternalServerError, w.Code)
	suite.Assert().True(strings.Contains(w.Body.String(), "update error"))
}

func (suite *NewspaperControllersSuite) TestDelete() {
	createdNewspaper, _ := models.CreateNewspaper("test", "sports")

	request, _ := api.NewDeleteNewspaperByIdRequest("/api/v1", createdNewspaper.ID)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request
	suite.newspaperHandler.DeleteNewspaperById(ginContext, createdNewspaper.ID)
	suite.Assert().Equal(http.StatusNoContent, w.Code)

	// 削除後に存在しないことを確認
	deletedNewspaper, err := models.GetNewspaper(createdNewspaper.ID)
	suite.Assert().NotNil(err)
	suite.Assert().Nil(deletedNewspaper)
}

func (suite *NewspaperControllersSuite) TestDeleteNoNewspaperFailure() {
	doesNotExistNewspaperID := 1111
	// 削除対象が存在しないことを確認
	deletedNewspaper, err := models.GetNewspaper(doesNotExistNewspaperID)
	suite.Assert().NotNil(err)
	suite.Assert().Nil(deletedNewspaper)

	request, _ := api.NewDeleteNewspaperByIdRequest("/api/v1", doesNotExistNewspaperID)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request
	suite.newspaperHandler.DeleteNewspaperById(ginContext, doesNotExistNewspaperID)
	suite.Assert().Equal(http.StatusNoContent, w.Code)
}

func (suite *NewspaperControllersSuite) TestDeleteNewspaperFailure() {
	mockDB := suite.MockDB()

	// モックの期待動作を定義
	mockDB.ExpectBegin()
	mockDB.ExpectExec("DELETE FROM `newspapers`").WillReturnError(errors.New("delete error"))
	mockDB.ExpectCommit()

	request, _ := api.NewDeleteNewspaperByIdRequest("/api/v1", 1)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request
	suite.newspaperHandler.DeleteNewspaperById(ginContext, 1)
	suite.Assert().Equal(http.StatusInternalServerError, w.Code)
	suite.Assert().True(strings.Contains(w.Body.String(), "delete error"))
}
