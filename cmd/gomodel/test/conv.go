package test

//gomodel conv = [
//  INSERT INTO {Follow}({Follow:UserId, FollowUserId})
//      SELECT ?, ? FROM DUAL
//      WHERE EXISTS(SELECT {User:Id} FROM User WHERE {User:Id}=?)
//]
func Conv() {

}
