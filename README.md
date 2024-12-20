# Sudoku Game

一个使用Go语言和Fyne框架开发的数独游戏。

## 效果图

![效果图](https://raw.githubusercontent.com/wangle201210/sudoku/main/image.png)

## 功能特点

- 标准9x9数独游戏
- 三种难度级别：
  - Easy: 保留41个数字
  - Medium: 保留31个数字
  - Hard: 保留21个数字
- 实时游戏计时器
- 实时输入验证
- 自动检测胜利条件

## 技术栈

- 语言: Go 1.22.4
- UI框架: Fyne v2

## 安装说明

1. 确保已安装Go 1.22.4或更高版本
2. 克隆仓库：
```bash
git clone https://github.com/wangle201210/sudoku.git
cd sudoku
```

3. 安装依赖：
```bash
go mod tidy
```

4. 运行游戏：
```bash
go run main.go
```

## 游戏规则

1. 在每个空格中填入1-9的数字
2. 每行必须包含1-9的数字，不能重复
3. 每列必须包含1-9的数字，不能重复
4. 每个3x3宫格必须包含1-9的数字，不能重复
5. 每个数独谜题只有一个唯一解

## 操作说明

- 选择难度：使用下拉菜单选择游戏难度
- 填写数字：点击格子输入1-9的数字
- 新游戏：点击"New Game"按钮开始新游戏
- 计时器：自动记录游戏时间

