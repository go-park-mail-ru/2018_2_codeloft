package models

import (
	"github.com/go-park-mail-ru/2018_2_codeloft/database"

	"go.uber.org/zap"
	"gopkg.in/mgo.v2/bson"
)

type Message struct {
	Id            bson.ObjectId `bson:"_id,omitempty" json:"-"`
	ChatId        int           `bson:"chat_id" json:"chat_id"`
	SenderLogin   string        `bson:"sender_login" json:"sender_login,omitempty"`
	ReceiverLogin string        `bson:"receiver_login" json:"receiver_login,omitempty"`
	Message       string        `bson:"message" json:"message"`
	Date          string        `bson:"data" json:"date"`
	Type          string        `bson:"-" json:"type,omitempty"`
}

type Messages []Message

func (m *Message) Add(db *database.MongoDB) error {
	collection := db.Database.C("messages")
	err := collection.Insert(m)
	if err != nil {
		zap.S().Infow("Can not insert", "err", err)
	}
	return err
}

func GetMessageByChatId(chatId int, db *database.MongoDB) (res Messages, err error) {
	collection := db.Database.C("messages")
	err = collection.Find(bson.M{"chat_id": chatId}).All(&res)
	if err != nil {
		zap.S().Infow("Can not get by chat id", "err", err)
	}
	return
}
