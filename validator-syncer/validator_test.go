package validator_syncer

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type WriterTestSuite struct {
	suite.Suite
}

func TestRunTestSuite(t *testing.T) {
	suite.Run(t, new(WriterTestSuite))
}

func (s *WriterTestSuite) SetupSuite()    {}
func (s *WriterTestSuite) TearDownSuite() {}
func (s *WriterTestSuite) SetupTest()     {}
func (s *WriterTestSuite) TearDownTest()  {}

//
//func (s * WriterTestSuite) TestRemove() {
//	arr1 := []int{1,2,3,4,5,6}
//	remove(arr1[], 2)
//
//}
