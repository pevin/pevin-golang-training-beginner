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

type paymentRepositoryTestSuite struct {
	postgresTest.Suite
}

func TestSuitePaymentRepository(t *testing.T) {
	dsn := os.Getenv("POSTGRES_TEST_URL")
	if dsn == "" {
		dsn = postgresTest.DefaultTestDsn
	}

	paymentRepoSuite := &paymentRepositoryTestSuite{
		postgresTest.Suite{
			DSN:                     dsn,
			MigrationLocationFolder: "../db/migrations",
		},
	}

	suite.Run(t, paymentRepoSuite)
}

func (s paymentRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	ok, err := s.Migration.Up()
	s.Require().NoError(err)
	s.Require().True(ok)
}

func (s paymentRepositoryTestSuite) AfterTest(suiteName, testName string) {
	ok, err := s.Migration.Down()
	s.Require().NoError(err)
	s.Require().True(ok)
}

func CreatePaymentPayload() model.Payment {
	id, _ := uuid.NewRandom()
	model := model.Payment{
		Id:            id.String(),
		PaymentCode:   "test-payment-code-" + id.String(),
		TransactionId: "test-trx-id",
		Name:          "Test Name",
		Amount:        "10000",
	}
	return model
}
func (s paymentRepositoryTestSuite) TestCreate() {
	mockInquiry := CreatePaymentPayload()

	testCases := []struct {
		desc        string
		repo        repository.IPaymentRepository
		expectedErr error
		ctx         context.Context
		reqBody     *model.Payment
	}{
		{
			desc: "insert-success",
			repo: func() repository.IPaymentRepository {
				repo := repository.PaymentRepository{Db: s.DBConn}
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
