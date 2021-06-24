package handler

import (
	"fmt"
	"net/http"
	"sort"
	"stefma.guru/valueStorage/apicommon"
	"stefma.guru/valueStorage/storage"
	"strconv"
	"strings"
	"time"
)

func HandleGet(w http.ResponseWriter, r *http.Request) {
	err := apicommon.CheckToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key := r.URL.Query().Get("key")

	storage, err := storage.CreateStorage()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer storage.Close()

	data, err := storage.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create time slice and sort it
	// to display the map correctly...
	sortedKeys := make([]time.Time, 0, len(data))
	for keyInData, _ := range data {
		sortedKeys = append(sortedKeys, keyInData)
	}
	sort.Slice(sortedKeys, func(i, j int) bool {
		return sortedKeys[i].Before(sortedKeys[j])
	})

	url := buildChartUrl(data, sortedKeys)
	fmt.Fprint(w, url)
	http.Redirect(w, r, url, 301)
}

func buildChartUrl(data map[time.Time]string, sortedKeys []time.Time) string {
	var builder strings.Builder
	builder.WriteString("https://quickchart.io/chart?c={type:'line',data:{labels:[")
	for _, key := range sortedKeys {
		builder.WriteString("'" + timeToString(key) + "'")
		builder.WriteString(",")
	}
	builder.WriteString("],")
	builder.WriteString("datasets:[{label:'Stuff',data:[")
	for _, key := range sortedKeys {
		builder.WriteString(data[key])
		builder.WriteString(",")
	}
	builder.WriteString("]}]}}")
	return builder.String()
}

func timeToString(time time.Time) string {
	var builder strings.Builder
	builder.WriteString(strconv.Itoa(time.Day()))
	builder.WriteString(".")
	builder.WriteString(strconv.Itoa(int(time.Month())))
	builder.WriteString(".")
	builder.WriteString(strconv.Itoa(time.Year()))
	builder.WriteString(" / ")
	builder.WriteString(strconv.Itoa(time.Hour()))
	builder.WriteString(":")
	builder.WriteString(strconv.Itoa(time.Minute()))
	builder.WriteString(":")
	builder.WriteString(strconv.Itoa(time.Second()))
	return builder.String()
}
