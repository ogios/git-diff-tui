> 请在开始前确保checkout到所需分支

使用方法:

运行`go install`安装为可执行文件

> 删除可直接`where merg-repo`/`which merg-repo`然后直接删除掉就行

```
merge-repo <commit_hash> <commit_hash> <path_regex>
# or
merge-repo -r <path_regex>
```

`-r`的意思是使用`reflog`获取当前分支的创建点的`commit_hash`到当前的`hash`

进入TUI后使用方法与vim类似（窗口与光标的移动都使用`j`/`k`/`ctrl+d`/`ctrl+u`），分为三个窗口，每个窗口之间使用`tab`/`shift+tab`切换：

- 左侧目录可以使用`a`切换全选，或`space`切换选择文件夹或文件。左侧目录view面板实现了滚动，如果文件长度超过了当前面板，可以使用`h`/`l`左右滚动
- 中间展示文件内容，需要注意的是并非当前文件系统中的文件，而是传入的最后一个commit记录的文件内容
- 右侧显示与当前文件相关的传入的`起始commit`与`终止commit`之间的所有commit

在左侧选择完成后按下`c`即可开始复制，但需要注意，复制的内容是文件系统里记录的文件内容，并非`终止commit`中记录的文件内容(关于这点，后续可以修改)

按下`ctrl+c`/`q`退出
