package controller

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/mock"
	"github.org/cclose/rsi-pledge-track/model"
	serviceMock "github.org/cclose/rsi-pledge-track/service/mock"
	"time"

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

		sampleData = []*model.PledgeData{
			{
				ID:        1,
				TimeStamp: "2023-10-01 12:00:00",
				Funds:     10,
				Citizens:  5,
				Fleet:     0,
			},
			{
				ID:        2,
				TimeStamp: "2023-10-02 12:00:00",
				Funds:     20,
				Citizens:  10,
				Fleet:     0,
			},
			{
				ID:        3,
				TimeStamp: "2023-10-03 12:00:00",
				Funds:     30,
				Citizens:  15,
				Fleet:     0,
			},
		}
	)
	BeforeEach(func() {
		mockPledgeDataService = &serviceMock.PledgeDataServiceMock{}

		e = echo.New()
		w = httptest.NewRecorder()

		pdc := NewPledgeDataController(mockPledgeDataService)
		e.Renderer = pdc
		e.GET("/pledge-data", pdc.GetPledgeData)
	})
	Context("GetPledgeData", func() {
		It("Passing Call", func() {
			mockPledgeDataService.On("GetAll", mock.Anything, mock.Anything).Return([]*model.PledgeData{}, nil)

			req, _ := http.NewRequest("GET", "/pledge-data", nil)

			e.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(200))
			mockPledgeDataService.AssertCalled(GinkgoT(), "GetAll", 0, 0)
		})
		It("Id Query", func() {
			mockPledgeDataService.On("Get", mock.Anything).Return(&model.PledgeData{}, nil)

			req, _ := http.NewRequest("GET", "/pledge-data?id=42", nil)

			e.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(200))
			mockPledgeDataService.AssertCalled(GinkgoT(), "Get", 42)
		})
		It("Offset Query", func() {
			mockPledgeDataService.On("GetAll", mock.Anything, mock.Anything).Return([]*model.PledgeData{}, nil)

			req, _ := http.NewRequest("GET", "/pledge-data?offset=1", nil)

			e.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(200))
			mockPledgeDataService.AssertCalled(GinkgoT(), "GetAll", 1, 0)
		})
		It("Limit Query", func() {
			mockPledgeDataService.On("GetAll", mock.Anything, mock.Anything).Return([]*model.PledgeData{}, nil)

			req, _ := http.NewRequest("GET", "/pledge-data?limit=100", nil)

			e.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(200))
			mockPledgeDataService.AssertCalled(GinkgoT(), "GetAll", 0, 100)
		})
		It("Timestamp Query", func() {
			mockPledgeDataService.On("GetByTimestamp", mock.Anything, mock.Anything).Return(&model.PledgeData{}, nil)

			req, _ := http.NewRequest("GET", "/pledge-data?timestamp=2016-12-06T19:09:05Z", nil)

			e.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(200))
			mockPledgeDataService.AssertCalled(GinkgoT(), "GetByTimestamp", mock.Anything, mock.Anything)
			mockPledgeDataService.AssertNotCalled(GinkgoT(), "GetAll", mock.Anything, mock.Anything)
		})
		It("AfterTimestamp Query", func() {
			mockPledgeDataService.On("GetAfterTimestamp", mock.Anything, mock.Anything, mock.Anything).Return([]*model.PledgeData{}, nil)
			timestamp := "2016-12-06T19:09:05Z"
			startingTime, _ := time.Parse(time.RFC3339, timestamp)

			req, _ := http.NewRequest("GET", "/pledge-data?startingDateTime="+timestamp, nil)

			e.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(200))
			mockPledgeDataService.AssertCalled(GinkgoT(), "GetAfterTimestamp", startingTime, 0, 0)
			mockPledgeDataService.AssertNotCalled(GinkgoT(), "GetAll", mock.Anything, mock.Anything)
		})
		It("Timestamp Query Invalid format", func() {
			req, _ := http.NewRequest("GET", "/pledge-data?timestamp=IAmAPotat03", nil)

			e.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(400))
			eErr := "Error parsing request: \ntimestamp=\"[IAmAPotat03]\": parsing time \"IAmAPotat03\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"IAmAPotat03\" as \"2006\""
			Expect(w.Body.String()).To(Equal(eErr))
		})
		Context("Format Param", func() {
			var entries []TableEntry
			for _, format := range model.ValidRequestFormats {
				entries = append(entries, Entry("Good Format "+string(format), string(format), 200, ""))
			}
			entries = append(entries, Entry("Bad Format xml", "xml", 400,
				"Error parsing request: \nformat=\"[xml]\": Invalid format: xml"),
			)
			DescribeTable("Validate Formats ", func(goodFormat string, expectedCode int, expectedBody string) {
				mockPledgeDataService.On("GetAll", mock.Anything, mock.Anything).Return([]*model.PledgeData{}, nil)

				req, _ := http.NewRequest("GET", "/pledge-data?format="+goodFormat, nil)

				e.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(expectedCode))
				if expectedBody != "" {
					Expect(w.Body.String()).To(Equal(expectedBody))
				}
			},
				entries,
			)
			DescribeTable("Format Response",
				func(ContentType, ContentDispo, Extension string, prepareBody func(body *bytes.Buffer) interface{}, ExpectedBody interface{}) {
					mockPledgeDataService.On("GetAll", mock.Anything, mock.Anything).Return(sampleData, nil)

					req, _ := http.NewRequest("GET", "/pledge-data?format="+Extension, nil)

					w := httptest.NewRecorder()
					e.ServeHTTP(w, req)

					Expect(w.Code).To(Equal(200))
					Expect(w.Header().Get("Content-Type")).To(Equal(ContentType))
					Expect(w.Header().Get("Content-Disposition")).To(Equal(ContentDispo))
					Expect(prepareBody(w.Body)).To(Equal(ExpectedBody))
				},
				Entry("CSV", "text/csv", "attachment; filename=pledgeData.csv", "csv",
					func(body *bytes.Buffer) interface{} {
						return body.String()
					},
					"ID,TimeStamp,Funds,Citizens,Fleet\n"+
						"1,2023-10-01 12:00:00,10,5,0\n"+
						"2,2023-10-02 12:00:00,20,10,0\n"+
						"3,2023-10-03 12:00:00,30,15,0\n",
				),
				Entry("JSON", "application/json; charset=UTF-8", "", "json",
					func(body *bytes.Buffer) interface{} {
						var payload []model.PledgeData
						err := json.Unmarshal([]byte(body.String()), &payload)
						Expect(err).NotTo(HaveOccurred())
						return payload
					},
					[]model.PledgeData{
						{
							ID:        1,
							TimeStamp: "2023-10-01 12:00:00",
							Funds:     10,
							Citizens:  5,
							Fleet:     0,
						},
						{
							ID:        2,
							TimeStamp: "2023-10-02 12:00:00",
							Funds:     20,
							Citizens:  10,
							Fleet:     0,
						},
						{
							ID:        3,
							TimeStamp: "2023-10-03 12:00:00",
							Funds:     30,
							Citizens:  15,
							Fleet:     0,
						},
					},
				),
			)
		})
	})
})
