package app

import ("errors"; "regexp"; "strings"; "time"; "github.com/google/uuid")
var emailRE=regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
func validateCredentials(email,password string) error { if !emailRE.MatchString(strings.TrimSpace(email)) || len(password)<8 || len(password)>72 { return errors.New("use a valid email and an 8-72 character password") }; return nil }
func validatePoll(question string, raw []string, expires *time.Time) ([]Option,error) { question=strings.TrimSpace(question); if len(question)<5||len(question)>240||len(raw)<2||len(raw)>10{return nil,errors.New("question must be 5-240 characters and include 2-10 options")}; seen:=map[string]bool{}; out:=make([]Option,0,len(raw)); for _,text:=range raw { text=strings.TrimSpace(text); key:=strings.ToLower(text); if len(text)<1||len(text)>100||seen[key]{return nil,errors.New("options must be unique and 1-100 characters")}; seen[key]=true; out=append(out,Option{ID:uuid.NewString(),Text:text}) }; if expires!=nil&&!expires.After(time.Now()){return nil,errors.New("expiration must be in the future")}; return out,nil }
