package model

/*********普通队列 使用先进先出FIFO *********/
type NormalQueue struct {
	//数据
	List []string `json:"list"`
}
//队列长度
func (queue *NormalQueue) Length() int{
	return len(queue.List)
}

//插入队列
func (queue *NormalQueue) Innormalqueue(value string) []string {
	if value == "" {
		return []string{}
	}
	//插入
	queue.List = append(queue.List,value)
	return queue.List
}

//排在第一的, 先出队列
func (queue *NormalQueue) Outnormalqueue() []string {
	lenth := queue.Length()
	if(lenth > 1){
		queue.List = queue.List[1:]
	}else{
		queue.List = []string{}
	}
	return queue.List
}
