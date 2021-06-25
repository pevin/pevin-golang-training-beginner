package repository_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/pevin/pevin-golang-training-beginner/model"
	postgresTest "github.com/pevin/pevin-golang-training-beginner/postgres"
	repository "github.com/pevin/pevin-golang-training-beginner/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type paymentCodeRepositoryTestSuite struct {
	postgresTest.Suite
}

func TestSuitePaymentCodeRepository(t *testing.T) {
	dsn := os.Getenv("POSTGRES_TEST_URL")
	if dsn == "" {
		dsn = postgresTest.DefaultTestDsn
	}

	paymentCodeRepoSuite := &paymentCodeRepositoryTestSuite{
		postgresTest.Suite{
			DSN:                     dsn,
			MigrationLocationFolder: "../db/migrations",
		},
	}

	suite.Run(t, paymentCodeRepoSuite)
}

func (s paymentCodeRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	ok, err := s.Migration.Up()
	s.Require().NoError(err)
	s.Require().True(ok)
}

func (s paymentCodeRepositoryTestSuite) AfterTest(suiteName, testName string) {
	ok, err := s.Migration.Down()
	s.Require().NoError(err)
	s.Require().True(ok)
}

func CreatePaymentCodePayload() model.PaymentCode {
	id, _ := uuid.NewRandom()
	model := model.PaymentCode{
		Id:          id.String(),
		PaymentCode: "test-payment-code-" + id.String(),
		Name:        "test name",
		Status:      "ACTIVE",
	}
	return model
}
func (s paymentCodeRepositoryTestSuite) TestCreatePaymentCode() {
	mockPaymentCode := CreatePaymentCodePayload()

	testCases := []struct {
		desc        string
		repo        repository.IPaymentCodeRepository
		expectedErr error
		ctx         context.Context
		reqBody     *model.PaymentCode
	}{
		{
			desc: "insert-success",
			repo: func() repository.IPaymentCodeRepository {
				repo := repository.PaymentCodeRepository{Db: s.DBConn}
				return repo
			}(),
			expectedErr: nil,
			ctx:         context.TODO(),
			reqBody:     &mockPaymentCode,
		},
	}

	for _, tC := range testCases {
		s.T().Run(tC.desc, func(t *testing.T) {
			err := tC.repo.Create(tC.ctx, tC.reqBody)
			s.Require().Equal(tC.expectedErr, err)
		})
	}
}

func (s paymentCodeRepositoryTestSuite) TestGetPaymentCodeById() {
	mockPaymentCode := CreatePaymentCodePayload()

	repo := repository.PaymentCodeRepository{Db: s.DBConn}
	err := repo.Create(context.TODO(), &mockPaymentCode)
	if err != nil {
		s.Fail("Error in creating seed settings", err)
	}

	id := mockPaymentCode.Id

	testCases := []struct {
		desc             string
		repo             repository.IPaymentCodeRepository
		expectedError    error
		expectedResponse model.PaymentCode
		id               string
		ctx              context.Context
	}{
		{
			desc:             "get-success",
			repo:             repo,
			expectedError:    nil,
			expectedResponse: mockPaymentCode,
			id:               id,
			ctx:              context.TODO(),
		},
		{
			desc: "not-found",
			repo: func() repository.IPaymentCodeRepository {
				repo := repository.PaymentCodeRepository{Db: s.DBConn}
				return repo
			}(),
			expectedError:    nil,
			expectedResponse: model.PaymentCode{},
			id:               "invalid-id",
			ctx:              context.TODO(),
		},
	}

	for _, tC := range testCases {
		// Run tests
		s.T().Run(tC.desc, func(t *testing.T) {
			res, err := tC.repo.Get(tC.ctx, tC.id)
			s.Require().Equal(tC.expectedError, err)

			if err == nil {
				s.Require().Equal(tC.expectedResponse.Id, res.Id)
				s.Require().Equal(tC.expectedResponse.PaymentCode, res.PaymentCode)
				s.Require().Equal(tC.expectedResponse.Name, res.Name)
			}
		})
	}
}

func (s paymentCodeRepositoryTestSuite) TestGetIdsToBeExpired() {
	var ids []string

	type seed func() []string
	repo := repository.PaymentCodeRepository{Db: s.DBConn}

	testCases := []struct {
		desc          string
		repo          repository.IPaymentCodeRepository
		seed          seed
		expectedError error
		ctx           context.Context
	}{
		{
			desc:          "no-ids-found",
			repo:          repo,
			expectedError: nil,
			seed: func() []string {
				return nil
			},
			ctx: context.TODO(),
		},
		{
			desc: "get-success",
			repo: repo,
			seed: func() []string {
				yesterday := time.Now().AddDate(0, 0, -5).UTC()
				mockPaymentCode := CreatePaymentCodePayload()
				mockPaymentCode.ExpirationDate = yesterday
				err := repo.Create(context.TODO(), &mockPaymentCode)
				if err != nil {
					s.Fail("Error in creating seed settings", err)
				}
				ids = append(ids, mockPaymentCode.Id)

				mockPaymentCode = CreatePaymentCodePayload()
				mockPaymentCode.ExpirationDate = yesterday
				err = repo.Create(context.TODO(), &mockPaymentCode)
				if err != nil {
					s.Fail("Error in creating seed settings", err)
				}
				ids = append(ids, mockPaymentCode.Id)
				return ids
			},
			expectedError: nil,
			ctx:           context.TODO(),
		},
	}

	for _, tC := range testCases {
		// Run tests
		s.T().Run(tC.desc, func(t *testing.T) {
			seededIds := tC.seed()
			res, err := tC.repo.GetIdsToBeExpired(tC.ctx)
			s.Require().Equal(tC.expectedError, err)

			if err == nil {
				s.Require().Equal(seededIds, res)
			}
		})
	}
}
func (s paymentCodeRepositoryTestSuite) TestUpdatePaymentCodeStatusById() {
	mockPaymentCode := CreatePaymentCodePayload()

	repo := repository.PaymentCodeRepository{Db: s.DBConn}
	err := repo.Create(context.TODO(), &mockPaymentCode)
	if err != nil {
		s.Fail("Error in creating seed settings", err)
	}

	testCases := []struct {
		desc        string
		repo        repository.IPaymentCodeRepository
		expectedErr error
		ctx         context.Context
		paymentCode *model.PaymentCode
	}{
		{
			desc:        "update-success",
			repo:        repo,
			expectedErr: nil,
			ctx:         context.TODO(),
			paymentCode: &mockPaymentCode,
		},
	}

	for _, tC := range testCases {
		s.T().Run(tC.desc, func(t *testing.T) {
			err := tC.repo.UpdateStatusById(tC.ctx, tC.paymentCode.Id, "EXPIRED")
			s.Require().Equal(tC.expectedErr, err)
		})
	}
}
