package updates

import (
    "bytes"
    "time"
    "github.com/tobyjsullivan/moneypenny/weather"
)

func BuildResponse() string {
    return compileUpdate().String() + "That's all for today. See you tomorrow."
}

func compileUpdate() *update {
    fGreeting := asyncMatter(greeting)
    fWeather := asyncMatter(weatherUpdate)
    fCalendar := asyncMatter(calendar)
    fReminder := asyncMatter(reminders)

    out := &update{
        matters: []*matter{
            <- fGreeting,
            <- fWeather,
            <- fCalendar,
            <- fReminder,
        },
    }

    return out
}

type update struct {
    matters []*matter
}

func (u *update) String() string {
    buf := new(bytes.Buffer)
    for _, m := range u.matters {
        buf.WriteString(m.text + "\n")
    }

    return buf.String()
}

type matter struct {
    text string
}

func asyncMatter(f func() *matter) chan *matter {
    var c = make(chan *matter)
    go func() {
        res := f()
        c <- res
    }()

    return c
}

func greeting() *matter {
    loc, _ := time.LoadLocation("America/Vancouver")
    hour := time.Now().In(loc).Hour()
    if hour < 12 {
        return &matter{text: "Good morning Toby"}
    } else if hour < 18 {
        return &matter{text: "Good afternoon Toby"}
    } else {
        return &matter{text: "Good evening Toby"}
    }
}

func weatherUpdate() *matter {
    fc, err := weather.VancouverForecast()
    if err != nil {
        println("ERROR: "+err.Error())
        return &matter{text: "I encountered an error while trying to check the weather."}
    }

    var drizzle bool
    var rain bool
    var snow bool
    var thunderstorm bool
    var extreme bool

    forecastHorizon := time.Now().Add(24 * time.Hour)

    for _, i := range fc.Items {
        if i.Time.After(forecastHorizon) {
            continue
        }

        switch i.Condition {
        case weather.ConditionDrizzle:
            drizzle = true
        case weather.ConditionRain:
            rain = true
        case weather.ConditionSnow:
            snow = true
        case weather.ConditionThunderstorm:
            thunderstorm = true
        case weather.ConditionExtreme:
            extreme = true
        }
    }

    if extreme {
        return &matter{text: "Be careful outside, there is extreme weather in the forecast. Check weather alerts for Vancouver."}
    } else if thunderstorm {
        return &matter{text: "It looks like there's a thunderstorm today - take your umbrella and plan to stay indoors."}
    } else if snow {
        return &matter{text: "It looks like there will be some snow today - plan accordingly."}
    } else if rain {
        return &matter{text: "It is going to be a rainy day today - take your umbrella."}
    } else if drizzle {
        return &matter{text: "It is going to drizzle a bit today - it might be smart to take your umbrella with you."}
    }

    return &matter{text: "It's looking clear outside today - you can leave the umbrella at home."}
}

func calendar() *matter {
    return &matter{text: "I don't have access to your calendar yet so you'll have to check it yourself."}
}

func reminders() *matter {
    return &matter{text: "Remember to send Jess some sweet texts."}
}