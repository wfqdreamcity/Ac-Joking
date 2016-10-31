/*
token 处理模块
*/
package token

import (
	"lib"
	"net/http"
	"errors"
	"time"
	"math/rand"
	"crypto/md5"
	"strconv"
	"encoding/hex"
)

//获取token
func GetToken(rw http.ResponseWriter , r *http.Request){
	r.ParseForm()

	//var token []map[string]string
	var appKey string
	var secretkey string

	if len(r.Form["appKey"]) == 0 {
		lib.Error(rw,"参数appKey缺失")
		return
	}
	if len(r.Form["secretkey"]) ==0 {
		lib.Error(rw,"参数secretkey缺失")
		return
	}

	appKey = r.Form["appKey"][0]
	secretkey = r.Form["secretkey"][0]


	token , err := getTokenInfo(appKey,secretkey)
	if err != nil {
		lib.Error(rw,err.Error())
		return
	}

	lib.Success(rw,token)
}

//获取token （更新token）
func getTokenInfo(appKey , secretkey string) (string,error){

	var token string
	var secret string

	key := "app_token:"+appKey
	tokenArray , err  :=  lib.Rclient.HGetAll(key).Result()

	if err != nil  || tokenArray["token"] == ""{
		//账号检测
		sqlCheck := "select secret_key,user_token,update_time from app_token where app_key=? and secret_key=? limit 1"
		CheckRestlt ,err := lib.DB.Query(sqlCheck , appKey,secretkey)
		defer  CheckRestlt.Close()

		if err != nil {
			panic(err)
		}

		if !CheckRestlt.Next() {
			return token ,errors.New("opp,get the wrong secret or appkey")
		}


		//token check
		times :=(time.Now().Second()-259200)
		sql := "select secret_key,user_token,update_time from app_token where app_key=? and secret_key=? and update_time >? limit 1"
		rows ,err := lib.DB.Query(sql , appKey,secretkey,times)
		defer  rows.Close()

		if err != nil {
			panic(err)
		}

		var secret string
		var user_token string
		var update_time string
		 redisToken := make(map[string]string)
		if rows.Next() {
			//reset the redis token
			if err := rows.Scan(&secret,&user_token,&update_time); err != nil {
				return token ,errors.New("opp,get token fail (form mysql)!")
			}

			token = user_token

			redisToken["token"] = user_token
			redisToken["secret_key"] = secret
			redisToken["update_time"] = update_time
			lib.Rclient.HMSet("app_token:"+appKey,redisToken)

			courceTime , _ := time.Parse("20060102150405","19700101000000")
			update_time , _ := time.ParseDuration(update_time+"s")
			lastUpdate := courceTime.Add(update_time)
			expireTime , _ :=time.ParseDuration("259200s")
			valibeTime := lastUpdate.Add(expireTime)

			nowtime := time.Now()

			expireTime = valibeTime.Sub(nowtime)

			lib.Rclient.Expire("app_token:"+appKey,expireTime)
			//设置验证apptoken
			lib.Rclient.Set("app_token_check:"+token,token,expireTime)


		}else{
			//reset the mysql and redis token
			randNum := strconv.Itoa(GetRand(111111,999999))

			update_time := time.Now().Unix()
			update_time_string := strconv.FormatInt(update_time,10)
			token = GetMD5Hash(appKey+randNum+update_time_string)

			sql := "update app_token set user_token=? ,update_time = ? where app_key =? and secret_key=?"
			result , err := lib.DB.Query(sql,token,update_time,appKey,secretkey)
			if err !=nil {
				return token,errors.New("reset token fail (from mysql) :"+err.Error())
			}
			result.Close()

			redisToken["token"] = token
			redisToken["secret_key"] = secretkey
			redisToken["update_time"] = update_time_string
			lib.Rclient.HMSet("app_token:"+appKey,redisToken)

			expireTime , _ :=time.ParseDuration("259200s")
			lib.Rclient.Expire("app_token:"+appKey,expireTime)
			//设置验证apptoken
			lib.Rclient.Set("app_token_check:"+token,token,expireTime)


		}

		return token , nil

	}else {
		secret = tokenArray["secret_key"]
		if secretkey == secret {
			token = tokenArray["token"]
		}else{
			return token ,errors.New("oop ,get wrong secret!")
		}
	}

	return token , nil

}

//获取指定范围随机数
func GetRand(min ,max int) int{

	subNum := max -min

	return min+rand.Intn(subNum)

}

//字符串md5 加密
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}