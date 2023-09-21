package group

import "server/pkg/common/db"

func User2GroupMember(user *db.UserModel) *db.GroupMemberModel {
	return &db.GroupMemberModel{
		UserID:   user.UserID,
		NickName: user.NickName,
		FaceURL:  user.FaceURL,
	}
}
