package service

import (
	"closer-api-go/dbclient"
	"closer-api-go/model"
	"fmt"
	"strings"
)

func CreateChatWithoutAnyInfo(userIds ...int) {
	chatId := dbclient.GetDbInstance().Insert(insertChatQuery)
	query := insertUserChatQuery
	var params []int
	for i, userId := range userIds {
		if i == 0 {
			query += "VALUES (?, ?)"
		} else {
			query += ", (?, ?)"
		}
		params = append(params, chatId)
		params = append(params, userId)
	}
	dbclient.GetDbInstance().Insert(query, params)
}

func GetUserChats(userId int) ([]model.Chat, error) {
	chats := make(map[int]model.Chat)
	rows, err := dbclient.GetDbInstance().Query(getChatsQuery, userId)
	if err != nil {
		fmt.Println(err)
		return []model.Chat{}, err
	}
	for rows.Next() {
		var chat model.Chat
		err = rows.Scan(&chat.Id)
		if err != nil {
			fmt.Println(err)
			return []model.Chat{}, err
		}
		chats[chat.Id] = chat
	}
	length := len(chats)
	slice := make([]string, length)
	for i := range slice {
		slice[i] = "?"
	}
	questionMarks := strings.Join(slice, ", ")
	query := strings.ReplaceAll(getUsersInChatBaseQuery, "%placeholder%", questionMarks)
	keys := make([]any, len(chats))
	i := 0
	for k := range chats {
		keys[i] = k
		i++
	}
	var params []any
	params = append(params, userId)
	params = append(params, keys...)
	fmt.Println(params)
	fmt.Println(query)
	rows, err = dbclient.GetDbInstance().Query(query, params...)

	if err != nil {
		fmt.Println(err)
		return []model.Chat{}, err
	}
	for rows.Next() {
		var user model.User
		var chatId int
		err = rows.Scan(&chatId, &user.Id, &user.Name, &user.PhoneNumber)
		chat := chats[chatId]
		chat.SetNewUser(user)
		chats[chatId] = chat
	}
	i, values := 0, make([]model.Chat, len(chats))
	for _, val := range chats {
		values[i] = val
		i++
	}
	return values, nil
}

const insertChatQuery = "INSERT INTO chats_go (name) VALUES (null)"
const insertUserChatQuery = "INSERT INTO user_chat_go (chat_id, user_id)"

const getChatsQuery = "select chat_id from user_chat_go where user_id = ? ORDER BY updated_at DESC"
const getUsersInChatBaseQuery = "select ucg.chat_id, u.id, u.name, u.phone_number " +
	"from users_go u " +
	"join user_chat_go ucg on u.id = ucg.user_id " +
	"where ucg.user_id != ? and ucg.chat_id in (%placeholder%)"
