package app

import (
  "errors"
  "os"
  "strings"
)

type Config struct { Port, MongoURI, MongoDatabase, RedisURL, JWTSecret, AllowedOrigins string; CookieSecure bool }
func LoadConfig() Config { return Config{Port: env("PORT","8080"), MongoURI: env("MONGO_URI","mongodb://localhost:27017"), MongoDatabase: env("MONGO_DATABASE","polling_tool"), RedisURL: env("REDIS_URL","redis://localhost:6379"), JWTSecret: os.Getenv("JWT_SECRET"), AllowedOrigins: env("ALLOWED_ORIGINS","http://localhost:5173"), CookieSecure: env("COOKIE_SECURE","false")=="true"} }
func env(k, fallback string) string { if v:=os.Getenv(k); v!="" { return v }; return fallback }
func (c Config) Validate() error { if len(c.JWTSecret)<32 || c.JWTSecret=="change-me" { return errors.New("JWT_SECRET must be set to a random value of at least 32 characters") }; if strings.TrimSpace(c.MongoURI)=="" || strings.TrimSpace(c.RedisURL)=="" { return errors.New("MONGO_URI and REDIS_URL are required") }; return nil }
func origins(v string) []string { out:=[]string{}; for _, item:=range strings.Split(v,",") { if s:=strings.TrimSpace(item); s!="" { out=append(out,s) } }; return out }
