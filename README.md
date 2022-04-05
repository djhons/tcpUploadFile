# tcpUploadFile
通过gzip一边压缩一边使用tcp上传文件夹。

# 前言：
以前需要从服务器上下载文件到本地电脑，需要先压缩再下载，如果文件夹特别大剩余的磁盘空间不够，就无法压缩，下载文件夹就是一个特别麻烦的事情。
<br>以前会开一个服务器，然后通过sftp把文件夹上传到服务器再压缩再下载，这样会需要输入账号密码，十分不安全。

# 实现：
使用filepath.Walk获取文件夹以及文件夹中的文件再使用gzip压缩，gzip压缩直接输出到远程服务器。所以需要在远程服务器上开启一个nc用于接收文件。

# 使用方法：
在服务器上使用`nc -lvvp 1221 >file.tar.gz`接收文件。
<br>将程序编译好后上传到要下载文件的服务器上，使用`main.exe -dir D:/web/log/ -host 8.8.8.8:1221`上传文件。
<br>可以使用`-skip "D:/web/log/tmp,D:/web/log/www"`的形式排除文件夹,遇到文件夹下有不可读权限文件时可以使用`-err False`跳过该文件

# 关于bug：
写这个程序时，是抱着能用就行的态度。所以有很多bug，欢迎提出来。
<br>已知bug：
  1. Linux下打包上传的文件夹在windows下解压会在file目录下多出来许多数字文件夹。
