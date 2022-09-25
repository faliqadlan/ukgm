package response

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"time"

	"ukgm/internal/app/config/env"
)

/*
============================
			NOTE
============================
This is responses template
*/
type PageRequest struct {
	Page           int64
	Paginate       int64
	PerPage        int64
	Search         string
	Filters        []Filter
	UnionFilters   []Filter
	Projections    []string
	Negations      []string
	Sorts          []Sort
	Date           time.Time
	IsRemoveDetail bool
}

type Filter struct {
	Option    string `json:"option"`
	Operator  string `json:"operator"`
	Value     string `json:"value"`
	ValueType string `json:"type"`
}
type Sort struct {
	Option string `json:"option"`
	Value  string `json:"value"`
}

type MessageNew []string
type Message map[string][]string

func NewResponse(data interface{}) *SingleResponse {
	response := &SingleResponse{
		Data: data,
		Meta: SimpleMeta{
			StatusCode: 200,
			Message: Message{
				"success": []string{"Success"},
			},
		},
	}
	return response
}

type SingleResponse struct {
	Data interface{} `json:"data"`
	Meta SimpleMeta  `json:"meta"`
}

type SingleResponseNew struct {
	Data interface{}   `json:"data"`
	Meta NewSimpleMeta `json:"meta"`
}

type SimpleMeta struct {
	StatusCode int     `json:"status_code"`
	Message    Message `json:"message"`
}

type NewSimpleMeta struct {
	StatusCode int        `json:"status_code"`
	Message    MessageNew `json:"message"`
}

type Pagination struct {
	CurrentPage int64       `json:"current_page"`
	From        int64       `json:"from"`
	LastPage    int64       `json:"last_page"`
	NextPageUrl string      `json:"next_page_url,omitempty"`
	Path        string      `json:"path"`
	PerPage     int64       `json:"per_page"`
	PrevPageUrl string      `json:"prev_page_url,omitempty"`
	To          int64       `json:"to"`
	Total       int64       `json:"total"`
	Data        interface{} `json:"data"`
	Meta        SimpleMeta  `json:"meta"`
}

func (singleRespNew *SingleResponseNew) ConvertToSingleResponse() *SingleResponse {
	result := new(SingleResponse)

	result.Data = singleRespNew.Data
	message := make(Message)
	message["message"] = singleRespNew.Meta.Message
	result.Meta = SimpleMeta{
		StatusCode: singleRespNew.Meta.StatusCode,
		Message:    message,
	}

	return result
}

func SendErrorJSON(res *SingleResponse, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Meta.StatusCode)
	json.NewEncoder(w).Encode(res)
	return
}

// SendJSON : Send Response to JSON
func SendJSON(data interface{}, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func SetResponsePagination(pageRequest *PageRequest, data interface{}, count int, pageUrl string, params string, isFirstParam bool) interface{} {
	response := &Pagination{
		Data: data,
		Meta: SimpleMeta{
			StatusCode: 200,
			Message: Message{
				"success": []string{"Success"},
			},
		},
		CurrentPage: pageRequest.Page,
		From:        (pageRequest.Page-1)*10 + 1,
		LastPage:    int64(math.Ceil(float64(count) / float64(pageRequest.PerPage))),
		PerPage:     pageRequest.PerPage,
		To:          int64(int(pageRequest.Page) * int(pageRequest.PerPage)),
		Total:       int64(count),
		Path:        env.Get("APP_HOST") + pageUrl + "?" + params,
	}

	if pageRequest.Page >= 1 && pageRequest.Page < response.LastPage {
		nextPage := "page=" + strconv.Itoa(int(pageRequest.Page)+1)
		if !isFirstParam {
			nextPage = "&" + nextPage
		}
		response.NextPageUrl = env.Get("APP_HOST") + pageUrl + "?" + params + nextPage
	}

	if pageRequest.Page > 1 && pageRequest.Page <= response.LastPage {
		prevPage := "page=" + strconv.Itoa(int(pageRequest.Page)-1)
		if !isFirstParam {
			prevPage = "&" + prevPage
		}
		response.PrevPageUrl = env.Get("APP_HOST") + pageUrl + "?" + params + prevPage
	}

	return response
}

func SetPaginate(r *http.Request, pageRequest *PageRequest, data interface{}, count int, pageUrl string) interface{} {
	var res interface{}
	if pageRequest == nil {
		pageRequest = &PageRequest{}
	}
	buffer := new(bytes.Buffer)
	isFirstParam := true
	for k, v := range r.URL.Query() {
		if k != "page" {
			if !isFirstParam {
				buffer.WriteString("&")
			} else {
				isFirstParam = false
			}
			buffer.WriteString(k)
			buffer.WriteString("=")
			buffer.WriteString(v[0])
		}
	}

	if pageRequest.Paginate == 1 {
		res = SetResponsePagination(pageRequest, data, count, pageUrl, buffer.String(), isFirstParam)
	} else {
		res = NewResponse(data)
		return res
	}

	return res
}

func PrintErrorMessage(err error) *SingleResponse {
	var ret *SingleResponse

	switch err.Error() {
	case "400":
		ret = &SingleResponse{
			Meta: SimpleMeta{
				StatusCode: http.StatusBadRequest,
				Message: Message{
					"errors": []string{"Bad Request"},
				},
			},
		}
		break
	case "400-Token-Parse":
		ret = &SingleResponse{
			Meta: SimpleMeta{
				StatusCode: http.StatusBadRequest,
				Message: Message{
					"errors": []string{"Cannot Read JWT Payload"},
				},
			},
		}
		break
	case "400-Signature":
		ret = &SingleResponse{
			Meta: SimpleMeta{
				StatusCode: http.StatusBadRequest,
				Message: Message{
					"errors": []string{"Invalid JWT Signature"},
				},
			},
		}
		break
	case "401-Expired":
		ret = &SingleResponse{
			Meta: SimpleMeta{
				StatusCode: http.StatusBadRequest,
				Message: Message{
					"errors": []string{"Invalid JWT Signature"},
				},
			},
		}
		break
	case "401":
		ret = &SingleResponse{
			Meta: SimpleMeta{
				StatusCode: http.StatusUnauthorized,
				Message: Message{
					"errors": []string{"JWT Token Is Required"},
				},
			},
		}
		break
	case "404":
		ret = &SingleResponse{
			Meta: SimpleMeta{
				StatusCode: http.StatusNotFound,
				Message: Message{
					"errors": []string{"Not Found"},
				},
			},
		}
		break
	case "405":
		ret = &SingleResponse{
			Meta: SimpleMeta{
				StatusCode: http.StatusMethodNotAllowed,
				Message: Message{
					"errors": []string{"Method Not Allowed"},
				},
			},
		}
		break
	case "500":
		ret = &SingleResponse{
			Meta: SimpleMeta{
				StatusCode: http.StatusInternalServerError,
				Message: Message{
					"errors": []string{err.Error()},
				},
			},
		}
		break
	default:
		ret = &SingleResponse{
			Meta: SimpleMeta{
				StatusCode: http.StatusBadRequest,
				Message: Message{
					"errors": []string{err.Error()},
				},
			},
		}
		break
	}

	return ret
}

func SendFile(downloadBytes []byte, fileName string, w http.ResponseWriter, r *http.Request) {
	mime := http.DetectContentType(downloadBytes)
	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Expires", "0")
	w.WriteHeader(http.StatusOK)

	// fileSize := len(string(downloadBytes))
	// Generate the server headers
	// w.Header().Set("Content-Transfer-Encoding", "binary")
	// w.Header().Set("Content-Length", strconv.Itoa(fileSize))
	// w.Header().Set("Content-Control", "private, no-transform, no-store, must-revalidate")

	// force it down the client's.....
	// http.ServeContent(w, r, fileName, time.Now(), bytes.NewReader(downloadBytes))
	w.Write(downloadBytes)
}
