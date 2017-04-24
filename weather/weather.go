package weather

import (
    "time"
    "net/http"
    "fmt"
    "encoding/json"
    "os"
    "errors"
)

const (
    ConditionClear = "Clear"
    ConditionClouds = "Clouds"
    ConditionRain = "Rain"
    ConditionDrizzle = "Drizzle"
    ConditionThunderstorm = "Thunderstorm"
    ConditionSnow = "Snow"
    ConditionExtreme = "Extreme"
    ConditionOther = "Other"

    cityVancouver = 6173331
)

var owmAppId string

func init() {
    owmAppId = os.Getenv("OPEN_WEATHER_MAP_APPID")
}

type Forecast struct {
    Items []*ForecastItem
}

type ForecastItem struct {
    Time time.Time
    Condition string
}

func VancouverForecast() (*Forecast, error) {
    url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast?id=%d&appid=%s", cityVancouver, owmAppId)
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode != http.StatusOK {
        return nil, errors.New("Received an unexpected response code.")
    }

    decoder := json.NewDecoder(resp.Body)
    var fc owmForecastResponse
    err = decoder.Decode(&fc)
    if err != nil {
        return nil, err
    }

    return mapOwmRespToForecast(&fc), nil
}

type owmForecastResponse struct {
    Code int `json:"cod"`
    List []*owmLine `json:"list"`
}

type owmLine struct {
    DateUnix int `json:"dt"`
    Weather []*owmWeather `json:"weather"`
    Rain *owmRain `json:"rain"`
}

type owmWeather struct {
    ConditionId int `json:"id"`
}

type owmRain struct {
    Amount3h float32 `json:"3h"`
}

func mapOwmRespToForecast(r *owmForecastResponse) *Forecast {
    items := make([]*ForecastItem, len(r.List))

    for i, li := range r.List {
        dt := time.Unix(int64(li.DateUnix), 0)
        cond := ConditionClear
        if len(li.Weather) > 0 {
            cond = mapOwmConditionCodeToCondition(li.Weather[0].ConditionId)
        }

        items[i] = &ForecastItem{
            Time: dt,
            Condition: cond,
        }
    }

    return &Forecast{
        Items: items,
    }
}

func mapOwmConditionCodeToCondition(code int) string {
    switch {
    case code >= 200 && code < 300:
        return ConditionThunderstorm
    case code >= 300 && code < 400:
        return ConditionDrizzle
    case code >= 500 && code < 600:
        return ConditionRain
    case code >= 600 && code < 700:
        return ConditionSnow
    case code == 800:
        return ConditionClear
    case code > 800 && code < 900:
        return ConditionClouds
    case code >= 900 && code < 910:
        return ConditionExtreme
    default:
        return ConditionOther
    }
}