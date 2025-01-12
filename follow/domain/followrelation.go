package domain

type FollowRelation struct {
	//被关注的人
	Followee int64
	//关注的人
	Follower int64
}

type FollowStatics struct {
	//被多少人关注（粉丝）
	Followers int64
	//自己关注了多少人
	Followees int64
}
