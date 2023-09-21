### list
#### group
|  Group func  | status |
|:------------:|:------:|
| createGroup  | finish |
|  joinGroup   | finish |
|  createUser  | finish |
|  kickMember  | finish |
| createMember | finish |
|  quitGroup   | finish |

#### friend

|   Friend func    | status |
|:----------------:|:------:|
|  applyFriendReq  | finish |
|  agreeFriendReq  | finish |
| refuseFriendReq  | finish |
| getSendFriendReq | finish |
| getRecvFriendReq | finish |
|   deleteFriend   | finish |
| setFriendRemark  | finish |
|  getFriendList   | finish |


### wait for solution
1. 添加好友的备注问题
2. 删除好友后的好友申请请求是否删除问题(handle_result=1)
3. redis的存储info后续再加
4. 群成员退出后删除聊天记录

### redis Design

#### friend
|         key          |     value     |  type  |
|:--------------------:|:-------------:|:------:|
|  friend_ids:ownerid  | [xxx,xxx,xxx] |  list  |
| friend_count:ownerid |       0       | string |


#### group
|           key            |     value     |  type  |
|:------------------------:|:-------------:|:------:|
| group_member_ids:groupid | [xxx,xxx,xxx] |  list  |
|   group_count:groupid    |       0       | string |