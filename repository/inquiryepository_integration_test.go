package repository_test

import (
	"context"
	"os"
	"testing"

	"github.com/pevin/pevin-golang-training-beginner/model"
	postgresTest "github.com/pevin/pevin-golang-training-beginner/postgres"
	repository "github.com/pevin/pevin-golang-training-beginner/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type inquiryRepositoryTestSuite struct {
	postgresTest.Suite
}

func TestSuiteInquiryRepository(t *testing.T) {
	dsn := os.Getenv("POSTGRES_TEST_URL")
	if dsn == "" {
		dsn = postgresTest.DefaultTestDsn
	}

	inquiryRepoSuite := &inquiryRepositoryTestSuite{
		postgresTest.Suite{
			DSN:                     dsn,
			MigrationLocationFolder: "../db/migrations",
		},
	}

	suite.Run(t, inquiryRepoSuite)
}

func (s inquiryRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	ok, err := s.Migration.Up()
	s.Require().NoError(err)
	s.Require().True(ok)
}

func (s inquiryRepositoryTestSuite) AfterTest(suiteName, testName string) {
	ok, err := s.Migration.Down()
	s.Require().NoError(err)
	s.Require().True(ok)
}

func CreateInquiryPayload() model.Inquiry {
	id, _ := uuid.NewRandom()
	model := model.Inquiry{
		Id:            id.String(),
		PaymentCode:   "test-payment-code-" + id.String(),
		TransactionId: "test-trx-id",
	}
	return model
}
func (s inquiryRepositoryTestSuite) TestCreate() {
	mockInquiry := CreateInquiryPayload()

	testCases := []struct {
		desc        string
		repo        repository.IInquiryRepository
		expectedErr error
		ctx         context.Context
		reqBody     *model.Inquiry
	}{
		{
			desc: "insert-success",
			repo: func() repository.IInquiryRepository {
				repo := repository.InquiryRepository{Db: s.DBConn}
				return repo
			}(),
			expectedErr: nil,
			ctx:         context.TODO(),
			reqBody:     &mockInquiry,
		},
	}

	for _, tC := range testCases {
		s.T().Run(tC.desc, func(t *testing.T) {
			err := tC.repo.Create(tC.ctx, tC.reqBody)
			s.Require().Equal(tC.expectedErr, err)
		})
	}
}

func (s inquiryRepositoryTestSuite) TestGetByTransactionId() {
	mockInquiry := CreateInquiryPayload()

	repo := repository.InquiryRepository{Db: s.DBConn}
	err := repo.Create(context.TODO(), &mockInquiry)
	if err != nil {
		s.Fail("Error in creating seed settings", err)
	}

	mockTrxId := mockInquiry.TransactionId

	testCases := []struct {
		desc             string
		repo             repository.IInquiryRepository
		expectedError    error
		expectedResponse model.Inquiry
		trxId            string
		ctx              context.Context
	}{
		{
			desc:             "get-success",
			repo:             repo,
			expectedError:    nil,
			expectedResponse: mockInquiry,
			trxId:            mockTrxId,
			ctx:              context.TODO(),
		},
		{
			desc: "not-found",
			repo: func() repository.IInquiryRepository {
				repo := repository.InquiryRepository{Db: s.DBConn}
				return repo
			}(),
			expectedError:    nil,
			expectedResponse: model.Inquiry{},
			trxId:            "invalid-id",
			ctx:              context.TODO(),
		},
	}

	for _, tC := range testCases {
		// Run tests
		s.T().Run(tC.desc, func(t *testing.T) {
			res, err := tC.repo.GetByTransactionId(tC.ctx, tC.trxId)
			s.Require().Equal(tC.expectedError, err)

			if err == nil {
				s.Require().Equal(tC.expectedResponse.Id, res.Id)
				s.Require().Equal(tC.expectedResponse.PaymentCode, res.PaymentCode)
				s.Require().Equal(tC.expectedResponse.TransactionId, res.TransactionId)
				s.Require().Equal(tC.expectedResponse.Amount, res.Amount)
			}
		})
	}
}
