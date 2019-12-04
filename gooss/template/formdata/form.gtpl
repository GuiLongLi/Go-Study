<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>oss上传文件</title>
</head>
<body>
<form action="/ossupload" method="post" enctype="multipart/form-data">
    文件：<input type="file" name="files">
    <input type="hidden" name="token" value="{{ .token }}">
    <input type="submit" value="提交">
</form>
</body>
</html>