package request

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"ukgm/internal/app/utils"

	responses "ukgm/internal/app/ukgm/model/response"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type PageRequest struct {
	Page           int64
	Paginate       int64
	PerPage        int64
	Search         string
	Filters        []Filter
	UnionFilters   []Filter
	GeoFilter      *GeoFilter
	Projections    []string
	Negations      []string
	Sorts          []Sort
	Date           time.Time
	IsRemoveDetail bool
}

type GeoFilter struct {
	Field       string
	MaxDistance float64
	Longitude   float64
	Latitude    float64
}

type Filter struct {
	Option    string `json:"option"`
	Operator  string `json:"operator"`
	Value     string `json:"value"`
	ValueType string `json:"type"`
}

func (f *Filter) ToBson() bson.D {
	var filter bson.D
	switch strings.ToUpper(f.Operator) {
	case ">":
		filter = bson.D{{Key: f.Option, Value: bson.D{{Key: "$gt", Value: f.getValue()}}}}
	case "<":
		filter = bson.D{{Key: f.Option, Value: bson.D{{Key: "$lt", Value: f.getValue()}}}}
	case ">=":
		filter = bson.D{{Key: f.Option, Value: bson.D{{Key: "$gte", Value: f.getValue()}}}}
	case "<=":
		filter = bson.D{{Key: f.Option, Value: bson.D{{Key: "$lte", Value: f.getValue()}}}}
	case "SIMILAR":
		filter = bson.D{{Key: f.Option, Value: primitive.Regex{Pattern: ".*" + f.Value + ".*", Options: "i"}}}
	case "=":
		filter = bson.D{{Key: f.Option, Value: f.getValue()}}
	case "!=":
		filter = bson.D{{Key: f.Option, Value: bson.M{"$ne": f.getValue()}}}
	case "IN":
		filter = bson.D{{Key: f.Option, Value: bson.M{"$in": f.getValue()}}}
	case "NIN":
		filter = bson.D{{Key: f.Option, Value: bson.M{"$nin": f.getValue()}}}
	case "DATE":
		if date, err := time.Parse("2006-01-02", f.Value); err == nil {
			filter = bson.D{{Key: f.Option, Value: bson.M{"$gte": utils.StartOfDay(date), "$lte": utils.EndOfDay(date)}}}
		}
	case "MONTH":
		if date, err := time.Parse("2006-01-02", f.Value); err == nil {
			filter = bson.D{{Key: f.Option, Value: bson.M{"$gte": utils.StartOfMonth(date), "$lte": utils.EndOfMonth(date)}}}
		}
	case "YEAR":
		if date, err := time.Parse("2006-01-02", f.Value); err == nil {
			filter = bson.D{{Key: f.Option, Value: bson.M{"$gte": utils.StartOfYear(date), "$lte": utils.EndOfYear(date)}}}
		}
	}

	return filter
}

func (f *Filter) getValue() interface{} {
	var datatype string

	f.ValueType = strings.ToUpper(f.ValueType)

	if f.ValueType == "" {
		switch strings.ToUpper(f.Operator) {
		case ">", "<", ">=", "<=":
			datatype = "NUM"
		case "IN", "NIN":
			datatype = "ARRAY"
		case "DATE", "MONTH", "YEAR":
			datatype = "DATE"
		case "=", "!=":
			datatype = "STRING"
		default:
			datatype = "STRING"
		}
	} else {
		if f.ValueType == "DATE" && f.Operator == ">=" || f.ValueType == "DATE" && f.Operator == ">" {
			datatype = "START DATE"
		} else if f.ValueType == "DATE" && f.Operator == "<=" || f.ValueType == "DATE" && f.Operator == "<" {
			datatype = "END DATE"
		} else {
			datatype = f.ValueType
		}
	}

	return convValue(f.Value, datatype)
}

func convValue(val, datatype string) interface{} {
	var result interface{}
	timeLoc, _ := time.LoadLocation("Asia/Jakarta")

	if val != "" {
		switch strings.ToUpper(datatype) {
		case "NUM":
			if num, err := strconv.Atoi(val); err == nil {
				result = num
			}
		case "UINT":
			if num, err := strconv.Atoi(val); err == nil {
				result = uint(num)
			}
		case "STRING":
			result = val
		case "DATE":
			if date, err := time.Parse("2006-01-02", val); err == nil {
				result = date
			}
		case "BOOL":
			result, _ = strconv.ParseBool(val)
		case "ARRAY":
			var temp []interface{}
			arr := strings.Split(val, ",")
			for _, v := range arr {
				if v == "" {
					temp = append(temp, nil)
				} else {
					temp = append(temp, v)
				}
			}

			result = temp
		case "NULL":
			result = nil
		case "START DATE":
			if date, err := utils.UniversalTimeParser(val); err == nil {
				result = utils.StartOfDay(date.Local().In(timeLoc))
			}
		case "END DATE":
			if date, err := utils.UniversalTimeParser(val); err == nil {
				result = utils.EndOfDay(date.Local().In(timeLoc))
			}
		default:
			result = val
		}
	} else {
		result = nil
	}

	return result
}

type Sort struct {
	Option string `json:"option"`
	Value  string `json:"value"`
}

func (s *Sort) ToBson() (sort bson.D) {
	if s.Option != "" {
		option := s.Option
		value := 1

		if s.Option == "created_at" {
			option = "timestamp.created_at"
		}

		if strings.ToUpper(s.Value) == "DESC" {
			value = -1
		}

		sort = bson.D{primitive.E{Key: option, Value: value}}
	} else {
		sort = nil
	}

	return sort
}

func (pr *PageRequest) QueryPipeline(searchFields []string) (filter bson.D, findOption *options.FindOptions) {
	filter = primitive.D{}
	var (
		allFilter       primitive.A   = primitive.A{}
		exclusiveFilter []primitive.D = []primitive.D{}
		unionFilter     []primitive.D = []primitive.D{}
		searchFilter    []primitive.D
	)
	filter = append(filter, primitive.E{Key: "timestamp.deleted_at", Value: nil})

	if len(pr.Filters) > 0 || pr.GeoFilter != nil {

		exclusiveFilter = append(exclusiveFilter, BuildExclusiveFilter(&pr.Filters)...)

		if pr.GeoFilter != nil {
			exclusiveFilter = append(exclusiveFilter, pr.GeoFilter.BuildGeoNearFilter())
		}
	}

	if pr.Search != "" && len(searchFields) > 0 {
		searchFilter = BuildSearchFilter(pr.Search, searchFields)
	}

	if len(pr.UnionFilters) > 0 {
		unionFilter = append(unionFilter, BuildExclusiveFilter(&pr.UnionFilters)...)
	}

	if len(exclusiveFilter) > 0 {
		allFilter = append(allFilter, bson.D{primitive.E{Key: "$and", Value: exclusiveFilter}})
	}
	if len(searchFilter) > 0 {
		allFilter = append(allFilter, bson.D{primitive.E{Key: "$or", Value: searchFilter}})
	}
	if len(unionFilter) > 0 {
		allFilter = append(allFilter, bson.D{primitive.E{Key: "$or", Value: unionFilter}})
	}

	if len(allFilter) > 0 {
		filter = append(filter, primitive.E{Key: "$and", Value: allFilter})
	}

	options := options.Find()
	if len(pr.Sorts) > 0 {
		for _, v := range pr.Sorts {
			if sor := v.ToBson(); sor != nil {
				options.SetSort(sor)
			}
		}
	}

	if pr.Paginate == 1 {
		options.SetLimit(pr.PerPage)
		if pr.Page >= 1 {
			options.SetSkip(pr.PerPage * (pr.Page - 1))
		}
	}

	if len(pr.Projections) > 0 {
		filter := bson.M{}
		for _, field := range pr.Projections {
			filter[field] = 1
		}
		options.SetProjection(filter)
	}

	if len(pr.Negations) > 0 {
		filter := bson.M{}
		for _, field := range pr.Negations {
			filter[field] = 0
		}
		options.SetProjection(filter)
	}
	findOption = options

	return
}

func (f *GeoFilter) BuildGeoNearFilter() (filter bson.D) {
	if f.Field == "" || f.MaxDistance < 0 || f.Latitude == 0 || f.Longitude == 0 {
		return
	}
	geoWithin := bson.D{
		{Key: "$centerSphere", Value: bson.A{bson.A{f.Longitude, f.Latitude}, (f.MaxDistance / 1000 / 6378.1)}},
	}
	filter = append(filter, bson.E{Key: f.Field, Value: bson.D{bson.E{Key: "$geoWithin", Value: geoWithin}}})
	return
}

// BuildSearchFilter :
func BuildSearchFilter(search string, fields []string) (filters []bson.D) {
	// MAIN VARIABLES
	var (
		filter bson.D
	)

	// BUILD SEARCH
	search = ".*" + search + ".*"

	// APPEND FILTER
	for _, field := range fields {
		// INIT FILTER
		filter = bson.D{
			{Key: field, Value: primitive.Regex{Pattern: search, Options: "i"}}}
		filters = append(filters, filter)
	}

	return
}

// BuildExclusiveFilter return bson that has been builded pre-execute
func BuildExclusiveFilter(f *[]Filter) (filters []bson.D) {
	// INITATION FILTER

	if len(*f) > 0 {
		for _, v := range *f {
			if fil := v.ToBson(); fil != nil {
				filters = append(filters, fil)
			}
		}
	}

	return filters
}

func PageRequestNewCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageRequest := &PageRequest{
			Filters:  make([]Filter, 0),
			Sorts:    make([]Sort, 0),
			Page:     1,
			PerPage:  10,
			Paginate: 0,
			Search:   "",
		}
		query := r.URL.Query()

		if query.Get("page") != "" {
			pageRequest.Page, _ = strconv.ParseInt(query.Get("page"), 0, 64)
		}

		if query.Get("per_page") != "" {
			pageRequest.PerPage, _ = strconv.ParseInt(query.Get("per_page"), 0, 64)
		}

		if query.Get("paginate") != "" {
			pageRequest.Paginate = 1
		}

		if query.Get("search") != "" {
			pageRequest.Search = query.Get("search")
		}

		if query.Get("sort") != "" {
			srts := strings.Split(query.Get("sort"), "|")
			if len(srts) >= 2 {
				pageRequest.Sorts = append(pageRequest.Sorts, Sort{Option: srts[0], Value: srts[1]})
			}
		}

		//Fetching filters data
		if len(query["filter[]"]) > 0 {
			for _, v := range query["filter[]"] {
				filter := Filter{}
				json.Unmarshal([]byte(v), &filter)
				pageRequest.Filters = append(pageRequest.Filters, filter)
			}
		}

		//Fetching sorts data
		if len(query["sorts[]"]) > 0 {
			for _, v := range query["sorts[]"] {
				sort := Sort{}
				json.Unmarshal([]byte(v), &sort)
				pageRequest.Sorts = append(pageRequest.Sorts, sort)
			}
		}

		ctx := context.WithValue(r.Context(), GetPageRequestCtxKey(), pageRequest)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (pr *PageRequest) ToResponsePR() (result *responses.PageRequest, err error) {
	temp, errTemp := json.Marshal(pr)
	if err != nil {
		err = errTemp
		zap.S().Error()
		return
	}
	json.Unmarshal(temp, &result)
	return
}
