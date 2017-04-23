package updates

import "bytes"

func BuildResponse() string {
    return compileUpdate().String()
}

func compileUpdate() *update {
    return &update{
        matters: []*matter{
            {text: "Good afternoon Toby."},
            {text: "It's looking sunny outside - you can leave the umbrella at home."},
            {text: "You have the day off tomorrow with no appointments."},
            {text: "Remember to send Jess some sweet texts."},
            {text: "That's all for today. See you tomorrow."},
        },
    }
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
