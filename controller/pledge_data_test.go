package controller

import (
	"fmt"
	"github.com/stretchr/testify/mock"
	"github.org/cclose/rsi-pledge-track/model"
	serviceMock "github.org/cclose/rsi-pledge-track/service/mock"
	"time"

	//"errors"
	"net/http"
	"net/http/httptest"

	echo "github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	_ "github.com/stretchr/testify/mock"
)

var _ = Describe("controller/pledge_data", func() {
	var (
		mockPledgeDataService *serviceMock.PledgeDataServiceMock
		e                     *echo.Echo
		w                     *httptest.ResponseRecorder
		//ctx                 context.Context
	)
	BeforeEach(func() {
		mockPledgeDataService = &serviceMock.PledgeDataServiceMock{}

		e = echo.New()
		w = httptest.NewRecorder()

		pdc := NewPledgeDataController(mockPledgeDataService)
		e.GET("/pledge-data", pdc.GetPledgeData)
		//ctx = context.Background()
	})
	Context("GetPledgeData", func() {
		It("Passing Call", func() {
			mockPledgeDataService.On("GetAll", mock.Anything, mock.Anything).Return([]*model.PledgeData{}, nil)

			req, _ := http.NewRequest("GET", "/pledge-data", nil)
			//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			e.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(200))
			mockPledgeDataService.AssertCalled(GinkgoT(), "GetAll", 0, 0)
		})
		It("Offset Query", func() {
			mockPledgeDataService.On("GetAll", mock.Anything, mock.Anything).Return([]*model.PledgeData{}, nil)

			req, _ := http.NewRequest("GET", "/pledge-data?offset=1", nil)
			//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			e.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(200))
			mockPledgeDataService.AssertCalled(GinkgoT(), "GetAll", 1, 0)
		})
		It("Limit Query", func() {
			mockPledgeDataService.On("GetAll", mock.Anything, mock.Anything).Return([]*model.PledgeData{}, nil)

			req, _ := http.NewRequest("GET", "/pledge-data?limit=100", nil)
			//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			e.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(200))
			mockPledgeDataService.AssertCalled(GinkgoT(), "GetAll", 0, 100)
		})
		It("Timestamp Query", func() {
			//mockPledgeDataService.On("GetAll", mock.Anything, mock.Anything).Return([]*model.PledgeData{}, nil)
			mockPledgeDataService.On("GetByTimestamp", mock.Anything, mock.Anything).Return(&model.PledgeData{}, nil)

			req, _ := http.NewRequest("GET", "/pledge-data?timestamp=2016-12-06T19:09:05Z", nil)

			e.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(200))
			mockPledgeDataService.AssertCalled(GinkgoT(), "GetByTimestamp", mock.Anything, mock.Anything)
		})
		It("AfterTimestamp Query", func() {
			//mockPledgeDataService.On("GetAll", mock.Anything, mock.Anything).Return([]*model.PledgeData{}, nil)
			mockPledgeDataService.On("GetAfterTimestamp", mock.Anything, mock.Anything, mock.Anything).Return([]*model.PledgeData{}, nil)
			timestamp := "2016-12-06T19:09:05Z"
			startingTime, _ := time.Parse(time.RFC3339, timestamp)

			req, _ := http.NewRequest("GET", "/pledge-data?startingDateTime="+timestamp, nil)

			e.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(200))
			mockPledgeDataService.AssertCalled(GinkgoT(), "GetAfterTimestamp", startingTime, 0, 0)
		})
		It("Timestamp Query Invalid format", func() {
			req, _ := http.NewRequest("GET", "/pledge-data?timestamp=IAmAPotat03", nil)

			e.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(400))
			eErr := "Error parsing request: \ntimestamp=\"[IAmAPotat03]\": parsing time \"IAmAPotat03\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"IAmAPotat03\" as \"2006\""
			Expect(w.Body.String()).To(Equal(eErr))
		})
		for _, goodFormat := range model.ValidRequestFormats {
			When("Format - "+string(goodFormat), func() {
				It("Format passes", func() {
					mockPledgeDataService.On("GetAll", mock.Anything, mock.Anything).Return([]*model.PledgeData{}, nil)

					req, _ := http.NewRequest("GET", "/pledge-data?format="+string(goodFormat), nil)

					e.ServeHTTP(w, req)
					Expect(w.Code).To(Equal(200))
					//Expect(w.Header().Get("Content-Type")).To(Equal(string(goodFormat)))
				})
			})
		}
		It("Format - invalid", func() {
			badFormat := "xml"
			mockPledgeDataService.On("GetAll", mock.Anything, mock.Anything).Return([]*model.PledgeData{}, nil)

			req, _ := http.NewRequest("GET", "/pledge-data?format="+badFormat, nil)

			e.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(400))
			eErr := fmt.Sprintf("Error parsing request: \nformat=\"[%s]\": Invalid format: %s", badFormat, badFormat)
			Expect(w.Body.String()).To(Equal(eErr))
		})
	})
})
