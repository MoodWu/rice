/*
平衡二叉树实现
MoodWu
*/
package avl
 
import (
	 "fmt"
	 "math"
)

type Node struct {
	Depth int //节点的高度
	bf int // 节点的平衡因子
	Key interface{}
	Data interface{}
	LNode *Node //左树
	RNode *Node //右树
	Compare CompareDataFunc
}

type CompareDataFunc func(a,b interface{}) int

type AVL interface {
	//增加一个元素
	Put(data interface{}) *Node
	//删除一个元素
	Del(data interface{}) *Node
	//查找一个元素
	Search(data interface{}) *Node
	//比较两个数据 1：a > b 0：a = b -1：a < b
	CompareData(a,b interface{}) int
}

//需要后续类型重载
func CompareData(a,b interface{}) int {
	switch a.(type) {
	case int:
		return CompareInt(int64(a.(int)),int64(b.(int)))
	case float32,float64:
		return CompareFloat(a.(float64),b.(float64))
	case string:
		return CompareString(a.(string),b.(string))
	default:
		return 0
	}
}

func CompareInt(a,b int64) int {
	if a > b {
		return 1
	} else if a < b {
		return -1
	} else {
		return 0
	}
}

func CompareFloat(a,b float64) int{
	if a > b {
		return 1
	} else if a < b {
		return -1
	} else {
		return 0
	}
}

func CompareString(a,b string) int {
	if a > b {
		return 1
	} else if a < b {
		return -1
	} else {
		return 0
	}
}

func (node *Node)Search(data interface{}) *Node  {
	if node == nil {
		return node
	}

	switch node.Compare(node.Key,data){
	case 0:
		return node
	case -1:
		return node.RNode.Search(data)
	case 1:
		return node.LNode.Search(data)
	}
	return nil 
}

//需要后续类型重载
func (node *Node)Del(data interface{}) *Node {

	return nil 
}


func InitTree(key,data interface{},f CompareDataFunc) *Node{
	if f == nil {
		return &Node {		
			Depth :1,
			bf :0,
			Key: key,
			Data: data,
			LNode :nil,		
			RNode :nil,
			Compare : CompareData,
		}
	} else {
		return &Node {				
			Depth :1,
			bf :0,
			Key: key,
			Data: data,
			LNode :nil,		
			RNode :nil,
			Compare : f,
		}
	}
}

func (node *Node) Put(key,data interface{},f CompareDataFunc) *Node {
	if f == nil {
		f = CompareData
	}
	if node == nil {
		return &Node {
			Depth :1,
			bf :0,
			Key: key,
			Data: data,
			Compare: f,
			LNode :nil,
			RNode :nil,
		}
	}

	switch node.Compare(node.Key,key) {
	case -1,0:
		//新元素大于等于当前节点，放到右树
		newNode := node.RNode.Put(key,data,node.Compare)
		node.RNode = newNode
	case 1:
		//新元素小与当前节点，放到左树
		newNode := node.LNode.Put(key,data,node.Compare)
		node.LNode = newNode
	}

	node.bf = GetNodeHeight(node.LNode) - GetNodeHeight(node.RNode)
	node.Depth = Max(GetNodeHeight(node.LNode),GetNodeHeight(node.RNode)) + 1
	node = node.ReBalance()

	return node
}

func (node *Node) RotateLeft() *Node {	
	r := node.RNode;  // 取得node的右树	
    node.RNode = r.LNode // 将node右树的左树（"拖油瓶"结点）链接到旋转后的node的右树中
    r.LNode = node; // 调转node和它右树的父子关系，使node成为它原右树的左树
    RefreshNode(node) // 更新并维护受影响结点的高度
	RefreshNode(r) // 更新并维护受影响结点的高度
	
    return r; // 将r返回
}



func (node *Node) RotateRight() *Node {	
	r := node.LNode;  // 取得node的左树	
    node.LNode = r.RNode // 将node左树的右树（"拖油瓶"结点）链接到旋转后的node的左树中
    r.RNode = node; // 调转node和它左树的父子关系，使node成为它原左树的右树
    RefreshNode(node) // 更新并维护受影响结点的高度
	RefreshNode(r) // 更新并维护受影响结点的高度
	
    return r; // 将r返回
}

func (node *Node) ReBalance() *Node {
    bf := GetBalance(node)
    if(bf > 1 && GetBalance(node.LNode)>0) { // LL型，进行单次右旋
     	return node.RotateRight();
    }
	if(bf > 1 && GetBalance(node.LNode)<=0) { //LR型 先左旋再右旋
		//先旋转为LL形
		r := node.LNode		
		y := r.RNode
		r.RNode = y.LNode
		y.LNode = r	
		node.LNode = y
		
		RefreshNode(node.LNode.LNode)
		RefreshNode(node.LNode)
		RefreshNode(node)
      	return node.RotateRight()      	
    }
    if(bf < -1 && GetBalance(node.RNode)<=0) {//RR型， 进行单次左旋
      return node.RotateLeft();
    }
	if(bf < -1 && GetBalance(node.RNode)>0) {// RL型，先右旋再左旋
		//先转成RR形
		r := node.RNode
		y := r.LNode
		r.LNode = y.RNode
		y.RNode = r
		node.RNode = y
		
		RefreshNode(node.RNode.RNode)
		RefreshNode(node.RNode)
		RefreshNode(node)
      	return node.RotateLeft();      
    }
    return node; 	
}

func RefreshNode(node *Node) {
	node.Depth = Max(GetNodeHeight(node.LNode),GetNodeHeight(node.RNode)) + 1; // 更新并维护受影响结点的高度
	node.bf = GetNodeHeight(node.LNode) - GetNodeHeight(node.RNode)
}

func Max(a int ,b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func GetNodeHeight(node *Node) int {
	if node == nil {
		return 0
	} else {
		return node.Depth
	}
}

func GetBalance(node *Node) int {
	if node == nil {
		return 0
	} else {
		return node.bf
	}
}


func (node *Node)Print() {
	lines := make(map[int]string,node.Depth)

	node.printnode(node.Depth,lines)
	for i:=node.Depth;i>0;i-- {
		fmt.Println(lines[i])
	}
	fmt.Println("--------------")
}

func (node *Node)printnode(depth int,lines map[int]string){
	max := int(math.Exp2(float64(depth)))
	for i:=1;i<max;i++ {
		lines[depth] += "  "
	}
	lines[depth] += fmt.Sprintf("%v",node.Key) 
	
	if node.LNode != nil {
		node.LNode.printnode(depth-1,lines)
		  lines[depth] += "[" + fmt.Sprintf("%v",node.LNode.Key) + ","
	} else {
		 lines[depth] += "[ ,"
		for i:=1;i<max;i++ {
			lines[depth-1] += "  "
		}	
	}

	if node.RNode != nil {
		node.RNode.printnode(depth-1,lines)
		lines[depth] += fmt.Sprintf("%v",node.RNode.Key) + "]"
	} else {
		lines[depth] += " ]"
		for i:=1;i<max;i++ {
			lines[depth-1] += "  "
		}	
	}
}


