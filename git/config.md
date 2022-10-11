# git配置信息

### 尾部换行符

提交时转换为LF，检出时转换为CRLF  
git config --global core.autocrlf true  
提交时转换为LF，检出时不转换  
git config --global core.autocrlf input  
提交检出均不转换 git config --global core.autocrlf false  
