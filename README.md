# qbittorrent-tool

一个用于管理 qBittorrent 的工具，支持自动分类、标签设置、做种限制等功能。

## 功能特性

- 根据 tracker 域名自动设置标签（如 baidu.com）
- 根据种子保存路径自动设置分类
- 基于标签、分类、tracker、关键词等条件精确控制做种时间、分享率、上传速度，并支持删种或例外保种
- 根据 tracker 工作状态（尤其是错误状态）自动添加标签以便识别问题
- 批量替换 tracker（待实现）

## Python 版本

本项目现在同时提供 Python 版本，位于 `python-version` 分支。Python 版本具有以下特点：

1. 与 Go 版本功能一致
2. 支持打包为可执行文件
3. 更容易扩展和维护

### 安装 Python 版本

```bash
# 克隆项目并切换到 python-version 分支
git clone https://github.com/fengqi/qbittorrent-tool.git
cd qbittorrent-tool
git checkout python-version

# 安装依赖
pip install -r requirements.txt
```

### 使用 Python 版本

```bash
# 直接运行
python main.py

# 或者安装后运行
pip install -e .
qbittorrent-tool
```

### 打包为可执行文件

```bash
# 安装打包工具
pip install pyinstaller

# 打包
make build

# 或者直接使用 pyinstaller
pyinstaller pyinstaller.spec
```

打包后的可执行文件位于 `dist/` 目录中。

## 配置文件

配置文件格式与 Go 版本一致，请参考 `example.config.json`。

## 使用方法

配置完成后，可以通过 qBittorrent 的"运行外部程序"功能或定时任务触发执行。

```bash
# Go 版本
./qbittorrent-tool -c ./config.json

# Python 版本
python main.py -c ./config.json
```

## 技术栈

- Go 版本: Go 1.19+
- Python 版本: Python 3.6+

## 许可证

MIT
