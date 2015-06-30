package controllers

import "github.com/revel/revel"
import  "time"
import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
	"strings"
	"github.com/nu7hatch/gouuid"

)

const KKDEV_DEV_MODE = true

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

type CommentItem struct{
	Id string  `bson:"_id,omitempty"`
	Dt string  `bson:"Dt,omitempty"`
	Email string `bson:"Email,omitempty"`
	Time int64 `bson:"Time,omitempty"`
	Ip string `bson:"Ip,omitempty"`
	CmChunk string `bson:"CmChunk,omitempty"`
	Name string `bson:"Name,omitempty"`

}


func (c App) CommentBasic(chunk string) revel.Result {

	if len(c.Params.Form["Name"])==1&&len(c.Params.Form["Email"])==1&&len(c.Params.Form["Dt"])==1{



		ret:=genComment(c.Params.Form["Dt"][0],c.Params.Form["Email"][0],getIp(c),c.Params.Form["Name"][0])
		ret.CmChunk=chunk
		postComment(ret,KKDEV_DEV_MODE)



	}


	CommentAr:=listCommentByChunk(chunk,KKDEV_DEV_MODE)
	return c.Render(CommentAr)


}


func (c App) API_GetCommentByChunk()revel.Result{

	//check request

	if(len(c.Params.Form["Chunk"])!=1){
		return nil //TODO: show err reason
	}

	chunk := c.Params.Form["Chunk"][0]


	res:=listCommentByChunk(chunk,KKDEV_DEV_MODE) //TODO:filter result

	return c.RenderJson(res)
}


func (c App)API_PostComment()revel.Result{

	if(len(c.Params.Form["Name"])!=1){
		return nil
	}

	if(len(c.Params.Form["Email"])!=1){
		return nil
	}

	if(len(c.Params.Form["Dt"])!=1){
		return nil
	}

	if(len(c.Params.Form["Chunk"])!=1){
		return nil
	}

	ret:=genComment(c.Params.Form["Dt"][0],c.Params.Form["Email"][0],getIp(c),c.Params.Form["Name"][0])
	ret.CmChunk=c.Params.Form["Chunk"][0]
	postComment(ret,KKDEV_DEV_MODE)

	return c.RenderJson(ret)




}


func getIp(c App) string{

	var uip string

	if ip := c.Request.Header.Get("X-Forwarded-For"); ip != "" {
	ips := strings.Split(ip, ",")
	if len(ips) > 0 && ips[0] != "" {
		rip := strings.Split(ips[0], ":")
		uip		= rip[0]
	}
} else {
	ip := strings.Split(c.Request.RemoteAddr, ":")
	if len(ip) > 0 {
		if ip[0] != "[" {
			uip	 = ip[0]
		}
	}
}
return uip
}


func getMongoDbUrl() string {

  mongodburl := ""

  if os.Getenv("KKDEV_MONGO_DB_URL") == "" {

		if os.Getenv("MONGO_PORT_27017_TCP_ADDR")!=""{
			//Running inside docker, we will construct mongodburl ourself
			mongodburl = os.Getenv("MONGO_PORT_27017_TCP_ADDR") +
			":" + os.Getenv("MONGO_PORT_27017_TCP_PORT") +
			"/" + os.Getenv("KKDEV_MONGO_DB_DBNAME") // not required if KKDEV_MONGO_DB_URL defined
		}else{
			//under development
			if os.Getenv("KKDEV_MONGO_DB_DBNAME")!=""{
				mongodburl="127.0.0.1:27017/"+os.Getenv("KKDEV_MONGO_DB_DBNAME")
			}else{
				mongodburl="127.0.0.1:27017/kkdev_kkcommentboxtest"
			}

		}

  } else {
    mongodburl = os.Getenv("KKDEV_MONGO_DB_URL")
  }

	return mongodburl

}

func getCommentbyid(id string, extout bool) CommentItem{

	mongodburl:=getMongoDbUrl()


	session,dialerr:=mgo.Dial(mongodburl)

	if dialerr!=nil{
		log.Println(dialerr)
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	db:=session.DB("")

	var CommentItemId CommentItem

	if (extout){
		m := bson.M{}
		prequeryerr:=db.C("Comments").FindId(id).Explain(m)

		log.Println("Explain: %#v\n", m)

		if prequeryerr!=nil{
			log.Println(prequeryerr)
		}

	}

	queryerr:=db.C("Comments").FindId(id).One(CommentItemId)

	if queryerr!=nil{
		log.Println(queryerr)
	}

	return CommentItemId

}

func postComment(target CommentItem, extout bool) string{

	mongodburl:=getMongoDbUrl()


	session,dialerr:=mgo.Dial(mongodburl)

	if dialerr!=nil{
		log.Println(dialerr)
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	db:=session.DB("")

	/*
	if (extout){
		m := bson.M{}
		prequeryerr:=db.C("Comments").Insert(target).Explain(m)

		log.Println("Explain: %#v\n", m)

		if prequeryerr!=nil{
			log.Println(prequeryerr)
		}

	}
*/
	queryerr:=db.C("Comments").Insert(target)

	if queryerr!=nil{
		log.Println(queryerr)
	}

	return target.Id


}

func listCommentByChunk(target string, extout bool)[]CommentItem{

	mongodburl:=getMongoDbUrl()


	session,dialerr:=mgo.Dial(mongodburl)

	if dialerr!=nil{
		log.Println(dialerr)
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	db:=session.DB("")

	var CommentItemIds []CommentItem

	if (extout){
		m := bson.M{}
		prequeryerr:=db.C("Comments").Find(bson.M{"CmChunk": target}).Sort("-Time").Explain(m)

		log.Println("Explain: %#v\n", m)

		if prequeryerr!=nil{
			log.Println(prequeryerr)
		}

	}

	queryerr:=db.C("Comments").Find(bson.M{"CmChunk": target}).Sort("-Time").All(&CommentItemIds)

	if queryerr!=nil{
		log.Println(queryerr)
	}

	return CommentItemIds


}


func genComment(Dt,Email,Ip,Name string)CommentItem{
	var CommentItema CommentItem
	id,_:=uuid.NewV4()
	CommentItema.Id=id.String()
	CommentItema.Dt=Dt
	CommentItema.Email=Email
	CommentItema.Ip=Ip
	CommentItema.Name=Name
	CommentItema.Time=time.Now().UTC().Unix()
	return CommentItema

}
