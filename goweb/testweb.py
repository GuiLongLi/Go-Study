#!/usr/bin/python
# -*- coding: UTF-8 -*-

import urllib
import urllib2
import time
import hashlib
import json

def post(url, data):
  req = urllib2.Request(url)
  data = urllib.urlencode(data)
  #enable cookie
  opener = urllib2.build_opener(urllib2.HTTPCookieProcessor())
  response = opener.open(req, data)
  return response.read()
def main(time,hashlib,json):
  posturl = "https://dev-gateway.xuanliyi.com"
  #测试秘钥
  secret = '51ac9aefa0772f6519e955b8202b9002'
  #测试appver
  appver = 'pc'
  #测试time
  time = int(time.time())
  strings = str(time)
  print strings
  strings_reverse = '{0}{1}{2}{3}{4}'.format(strings[3:4],strings[1:3],strings[0:1],strings[4:7],"000")
  print strings
  #测试版本
  version = '1.0.0'
  action = "Article.getArticleList"
  api_param = {"category_id":1001}
  api_param = json.dumps(api_param)

  #测试签名
  #sign = 'action={0}&api_param={1}&appver={2}&scaret={3}&timestamp={4}&version={5}'.format(action,api_param,appver,secret,strings_reverse,version)
  sign = 'version={0}&timestamp={1}&scaret={2}&appver={3}&api_param={4}&action={5}'.format(version,strings_reverse,secret,appver,api_param,action)

  print sign

  sign = hashlib.md5(sign.encode(encoding='UTF-8')).hexdigest()

  print sign

  data = {'action':action, 'api_param':api_param, 'appver':appver, 'version':version,'timestamp':strings, 'sign':sign}
  print post(posturl, data)
if __name__ == '__main__':
  main(time,hashlib,json)
