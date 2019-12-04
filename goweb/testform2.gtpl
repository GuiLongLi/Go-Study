<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>修改资料</title>
</head>
<body>
<form action="/info" method="post" enctype="multipart/form-data">
    头像：<input type="file" name="userheader">
    <br>
    用户名：<input type="text" name="username">
    <br>
    密码：<input type="password" name="password">
    <br>
    年龄：<input type="text" name="age">
    <br>
    邮箱：<input type="text" name="email">
    <br>
    手机号码：<input type="text" name="mobile">
    <br>
    身份证号码：<input type="text" name="usercard">
    <br>
    测试xss：<input type="text" name="xss">
    <br>
    性别：
    <select name="sex" >
        <option value="0">无</option>
        <option value="1">男</option>
        <option value="2">女</option>
    </select>
    <br>
    兴趣：
    <br>
    <input type="checkbox" name="interest" value="no">无
    <input type="checkbox" name="interest" value="football">足球
    <input type="checkbox" name="interest" value="basketball">篮球
    <input type="checkbox" name="interest" value="tennis">网球
    <br>
    <input type="hidden" name="token" value="{{.}}">
    <input type="submit" value="提交">
</form>
</body>
</html>