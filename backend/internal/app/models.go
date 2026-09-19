package app

import ("net/http"; "time"; "github.com/gin-gonic/gin"; "github.com/golang-jwt/jwt/v5"; "github.com/redis/go-redis/v9"; "go.mongodb.org/mongo-driver/mongo")
type User struct { ID string `json:"id" bson:"_id"`; Email string `json:"email" bson:"email"`; PasswordHash string `json:"-" bson:"passwordHash"`; CreatedAt time.Time `json:"createdAt" bson:"createdAt"` }
type Option struct { ID string `json:"id" bson:"id"`; Text string `json:"text" bson:"text"` }
type Poll struct { ID string `json:"id" bson:"_id"`; OwnerID string `json:"-" bson:"ownerId"`; Question string `json:"question" bson:"question"`; Options []Option `json:"options" bson:"options"`; Status string `json:"status" bson:"status"`; CreatedAt time.Time `json:"createdAt" bson:"createdAt"`; UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`; ExpiresAt *time.Time `json:"expiresAt,omitempty" bson:"expiresAt,omitempty"` }
type Vote struct { ID string `bson:"_id"`; PollID string `bson:"pollId"`; OptionID string `bson:"optionId"`; VoterID string `bson:"voterId"`; CreatedAt time.Time `bson:"createdAt"` }
type Claims struct { UserID string `json:"userId"`; jwt.RegisteredClaims }
type ResultOption struct { ID string `json:"id"`; Text string `json:"text"`; Votes int64 `json:"votes"`; Percentage float64 `json:"percentage"` }
type Results struct { Poll Poll `json:"poll"`; Options []ResultOption `json:"options"`; TotalVotes int64 `json:"totalVotes"` }
type App struct { cfg Config; mongoClient *mongo.Client; db *mongo.Database; redis *redis.Client; router *gin.Engine; httpServer *http.Server }
